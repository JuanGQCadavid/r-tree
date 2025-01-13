package domain

import "math"

type LatLon struct {
	Lat float64 `json:"lat,omitempty"`
	Lon float64 `json:"lon,omitempty"`
}

type Location[T any] struct {
	Value        T        `json:"id,omitempty"`
	ChildPointer *Node[T] `json:"-"`
	LimitA       *LatLon  `json:"coords,omitempty"`
	LimitB       *LatLon  `json:"-"`
}

type Node[T any] struct {
	Parent    *Node[T]
	Locations []*Location[T]
}

func (node *Node[T]) UpdateMBR() {
	for _, locP := range node.Locations {
		if locP.ChildPointer != nil {
			a, b := locP.ChildPointer.GetMRB()
			locP.LimitA = a
			locP.LimitB = b
		}
	}
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

func (node *Node[T]) GetMRB() (*LatLon, *LatLon) {
	if len(node.Locations) == 0 {
		return nil, nil
	}

	var (
		l_1_coords = node.Locations[0]
	)

	for i := 1; i < len(node.Locations); i++ {
		a, b := minSquare[T](node.Locations[i], l_1_coords)
		l_1_coords = &Location[T]{
			LimitA: a,
			LimitB: b,
		}
	}

	return l_1_coords.LimitA, l_1_coords.LimitB
}

func minSquare[T any](limitA, limitB *Location[T]) (*LatLon, *LatLon) {
	minALat := math.Min(limitA.LimitA.Lat, limitA.LimitB.Lat)
	minBLat := math.Min(limitB.LimitA.Lat, limitB.LimitB.Lat)
	minABLat := math.Min(minALat, minBLat)

	minALon := math.Min(limitA.LimitA.Lon, limitA.LimitB.Lon)
	minBLon := math.Min(limitB.LimitA.Lon, limitB.LimitB.Lon)
	minABLon := math.Min(minALon, minBLon)

	maxALon := math.Max(limitA.LimitA.Lon, limitA.LimitB.Lon)
	maxBLon := math.Max(limitB.LimitA.Lon, limitB.LimitB.Lon)
	maxABLon := math.Max(maxALon, maxBLon)

	maxALat := math.Max(limitA.LimitA.Lat, limitA.LimitB.Lat)
	maxBLat := math.Max(limitB.LimitA.Lat, limitB.LimitB.Lat)
	maxABLat := math.Max(maxALat, maxBLat)

	return &LatLon{
			Lat: minABLat,
			Lon: minABLon,
		}, &LatLon{
			Lat: maxABLon,
			Lon: maxABLat,
		}
}
