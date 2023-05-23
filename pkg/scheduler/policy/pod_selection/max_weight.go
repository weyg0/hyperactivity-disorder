package podselection

import "k8s.io/apimachinery/pkg/types"

type Pod struct {
	weight float64
	aoi    int // Age of Information
	debt   float64
}

var PodSet map[types.UID]Pod
