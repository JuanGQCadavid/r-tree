package rtree

import (
	"log"
	"testing"

	"github.com/JuanGQCadavid/r-tree/app/core/domain"
)

// func TestInsertOnRoot(t *testing.T) {
// 	tree := NewRTree[string](3)

// 	tree.InsertLocation(&domain.LatLon{Lat: 1, Lon: 1}, "A", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 50, Lon: 50}, "B", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 100, Lon: 100}, "C", tree.Root)

// 	tree.TraverseAndPrint()
// }

// func TestInsertOnRoot5(t *testing.T) {
// 	tree := NewRTree[string](3)

// 	tree.InsertLocation(&domain.LatLon{Lat: 1, Lon: 1}, "A", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 50, Lon: 50}, "B", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 100, Lon: 100}, "C", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 200, Lon: 200}, "D", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 300, Lon: 300}, "E", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 400, Lon: 400}, "F", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 500, Lon: 500}, "G", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 600, Lon: 600}, "H", tree.Root)

// 	tree.TraverseAndPrint()
// }

// func TestInsertSearch(t *testing.T) {
// 	tree := NewRTree[string](3)

// 	tree.InsertLocation(&domain.LatLon{Lat: 1, Lon: 1}, "A", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 50, Lon: 50}, "B", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 100, Lon: 100}, "C", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 200, Lon: 200}, "D", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 300, Lon: 300}, "E", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 400, Lon: 400}, "F", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 500, Lon: 500}, "G", tree.Root)
// 	tree.InsertLocation(&domain.LatLon{Lat: 600, Lon: 600}, "H", tree.Root)

// 	tree.TraverseAndPrint()

// 	resp := tree.Search(
// 		&domain.LatLon{
// 			Lat: 2,
// 			Lon: 40,
// 		},
// 		tree.Root)

// 	for _, val := range resp {
// 		log.Println(val.Value)
// 	}

// }

func TestInsertSearchAreaGrande(t *testing.T) {
	tree := NewRTree[string](3)

	tree.InsertLocation(&domain.LatLon{Lat: 1, Lon: 1}, "A", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 50, Lon: 50}, "B", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 100, Lon: 100}, "C", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 200, Lon: 200}, "D", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 300, Lon: 300}, "E", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 400, Lon: 400}, "F", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 500, Lon: 500}, "G", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 600, Lon: 600}, "H", tree.Root)

	tree.TraverseAndPrint()

	resp := tree.Search(
		&domain.LatLon{
			Lat: 249,
			Lon: 249,
		},
		tree.Root)

	for _, val := range resp {
		log.Println(val.Value)
	}

}

// func TestSplit(t *testing.T) {
// 	tree := NewRTree[string](3)

// 	l := &domain.Node[string]{
// 		Parent: nil,
// 		Locations: []*domain.Location[string]{
// 			{Value: "1", LimitA: &domain.LatLon{0, 0}, LimitB: &domain.LatLon{1, 1}},     // Area = 1
// 			{Value: "2", LimitA: &domain.LatLon{2, 2}, LimitB: &domain.LatLon{4, 4}},     // Area = 4
// 			{Value: "3", LimitA: &domain.LatLon{5, 5}, LimitB: &domain.LatLon{6, 6}},     // Area = 1
// 			{Value: "5", LimitA: &domain.LatLon{10, 10}, LimitB: &domain.LatLon{14, 14}}, // Area = 16
// 			{Value: "6", LimitA: &domain.LatLon{15, 15}, LimitB: &domain.LatLon{20, 20}}, // Area = 25,
// 			{Value: "4", LimitA: &domain.LatLon{7, 7}, LimitB: &domain.LatLon{9, 9}},     // Area = 4
// 		},
// 	}

// 	l_1, l_2 := tree.SplitNodeQuadraticCost(l)

// 	log.Println("L1")
// 	for _, v := range l_1.Locations {
// 		log.Println(v.Value, v.LimitA)
// 	}

// 	log.Println("l_2")
// 	for _, v := range l_2.Locations {
// 		log.Println(v.Value, v.LimitA)
// 	}

// }
