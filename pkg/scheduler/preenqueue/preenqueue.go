package preenqueue

import (
	"strconv"

	"github.com/weyg0/hyperactivity-disorder/pkg/scheduler/policy"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

func GetPodWeight(pod *v1.Pod) float64 {
	if w, ok := pod.GetLabels()["weight"]; ok {
		weight, _ := strconv.ParseFloat(w, 64)
		return weight
	}
	klog.Errorf("[PreEnqueue] GetPodWeight Error. Labels: %v", pod.Labels)
	return policy.PodWeight
}

func GetPodMinSelectFreq(pod *v1.Pod) float64 {
	if f, ok := pod.GetLabels()["select/freq"]; ok {
		freq, _ := strconv.ParseFloat(f, 64)
		return freq
	}
	klog.Errorf("[PreEnqueue] GetPodMinSelectFreq Error. Labels: %v", pod.Labels)
	return policy.PodMinSelectFreq
}
