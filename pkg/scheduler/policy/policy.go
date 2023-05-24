package policy

import "k8s.io/apimachinery/pkg/types"

var (
	PodSet = map[types.UID]Pod{}
	Time   = 0.0
)

// Max-Weight policy
type Pod struct {
	Priority      float64
	Weight        float64
	AoI           float64 // Age of Information
	Debt          float64 // Pod selection
	SelectedTimes float64
	MinSelectFreq float64
}

const (
	V = 1.0 // Trade-off between AoI and Debt
)
