package domain

type LatLon struct {
	Lat float64
	Lon float64
}

type Location[T any] struct {
	Value        T
	ChildPointer *Node[T]
	LimitA       *LatLon
	LimitB       *LatLon
}

type Node[T any] struct {
	Parent    *Node[T]
	Locations []*Location[T]
}

// TODO - What about making it as a variable and chaning it when splitting?
func (node *Node[T]) IsLeaf() bool {
	for _, val := range node.Locations {
		if val != nil && val.ChildPointer != nil {
			return false
		}
	}
	return true
}
