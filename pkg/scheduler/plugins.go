package scheduler

import (
	"context"
	"fmt"
	"math"

	"github.com/weyg0/hyperactivity-disorder/pkg/scheduler/policy"
	"github.com/weyg0/hyperactivity-disorder/pkg/scheduler/preenqueue"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// Name is the name of the plugin used in the Registry and configurations.
const (
	Name = "active-defense"
)

// 定义 ActiveDefense 结构类型
type ActiveDefense struct {
	handle framework.Handle
}

// 确保 ActiveDefense 实现调度框架接口
var (
	_ framework.PreEnqueuePlugin = &ActiveDefense{}
	_ framework.QueueSortPlugin  = &ActiveDefense{}
	_ framework.ScorePlugin      = &ActiveDefense{}
	_ framework.ScoreExtensions  = &ActiveDefense{}
)

func (ad *ActiveDefense) Name() string {
	return Name
}

func (ad *ActiveDefense) PreEnqueue(ctx context.Context, pod *v1.Pod) *framework.Status {
	uid := pod.UID
	if p, ok := policy.PodSet[uid]; ok { // 更新 Debt 和 Priority
		p.Debt = policy.Time*preenqueue.GetPodMinSelectFreq(pod) - p.SelectedTimes
		p.Priority = preenqueue.GetPodWeight(pod)/2*p.AoI*(p.AoI+2) + policy.V*math.Max(p.Debt, 0)
	} else { // 初始化 Pod
		policy.PodSet[uid] = policy.Pod{
			AoI: 1,
		}
	}
	return framework.NewStatus(framework.Success, "")
}

func (ad *ActiveDefense) Less(pInfo1, pInfo2 *framework.QueuedPodInfo) bool {
	// 根据 Priority 选择最佳 Pod
	p1 := policy.PodSet[pInfo1.Pod.UID]
	p2 := policy.PodSet[pInfo2.Pod.UID]
	return p1.Priority > p2.Priority
}

// Score invoked at the score extension point.
func (ad *ActiveDefense) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := ad.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}

	// pe.score favors nodes with terminating pods instead of nominated pods
	// It calculates the sum of the node's terminating pods and nominated pods
	return ad.score(nodeInfo)
}

// ScoreExtensions of the Score plugin.
func (ad *ActiveDefense) ScoreExtensions() framework.ScoreExtensions {
	return ad
}

func (ad *ActiveDefense) score(nodeInfo *framework.NodeInfo) (int64, *framework.Status) {
	var terminatingPodNum, nominatedPodNum int64
	// get nominated Pods for node from nominatedPodMap
	nominatedPodNum = int64(len(ad.handle.NominatedPodsForNode(nodeInfo.Node().Name)))
	for _, p := range nodeInfo.Pods {
		// Pod is terminating if DeletionTimestamp has been set
		if p.Pod.DeletionTimestamp != nil {
			terminatingPodNum++
		}
	}
	return terminatingPodNum - nominatedPodNum, nil
}

func (ad *ActiveDefense) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	// Find highest and lowest scores.
	var highest int64 = -math.MaxInt64
	var lowest int64 = math.MaxInt64
	for _, nodeScore := range scores {
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}

	// Transform the highest to lowest score range to fit the framework's min to max node score range.
	oldRange := highest - lowest
	newRange := framework.MaxNodeScore - framework.MinNodeScore
	for i, nodeScore := range scores {
		if oldRange == 0 {
			scores[i].Score = framework.MinNodeScore
		} else {
			scores[i].Score = ((nodeScore.Score - lowest) * newRange / oldRange) + framework.MinNodeScore
		}
	}

	return nil
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {

	return &ActiveDefense{handle: h}, nil
}
