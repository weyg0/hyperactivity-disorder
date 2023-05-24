package policy

import "k8s.io/apimachinery/pkg/types"

var (
	PodSet = map[types.UID]Pod{}
	Time   = 0.0
)

// Max-Weight policy
type Pod struct {
	Priority      float64
	AoI           float64 // Age of Information
	Debt          float64 // Pod selection
	SelectedTimes float64
}

// Pod 常量
const (
	V                = 1.0            // Trade-off between AoI and Debt
	PodWeight        = 1.0            // 默认权重
	PodNumbers       = 20             // 总数
	PodMinSelectFreq = 1 / PodNumbers // 默认最小选择频率

)
