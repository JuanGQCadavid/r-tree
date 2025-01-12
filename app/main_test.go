package main

import (
	"log"
	"testing"

	"github.com/JuanGQCadavid/r-tree/app/core/domain"
)

func TestInsertOnRoot(t *testing.T) {
	tree := NewRTree[string](3)

	tree.InsertLocation(&domain.LatLon{Lat: 1, Lon: 1}, "A", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 50, Lon: 50}, "B", tree.Root)
	tree.InsertLocation(&domain.LatLon{Lat: 100, Lon: 100}, "C", tree.Root)

	tree.TraverseAndPrint()
}

func TestSplit(t *testing.T) {
	tree := NewRTree[string](3)

	// SplitNodeQuadraticCost(newLocation *domain.Location[T], l *domain.Node[T]) (*domain.Node[T], *domain.Node[T])
	newLocation := &domain.Location[string]{Value: "4", LimitA: &domain.LatLon{7, 7}, LimitB: &domain.LatLon{9, 9}} // Area = 4

	l := &domain.Node[string]{
		Parent: nil,
		Locations: []*domain.Location[string]{
			{Value: "1", LimitA: &domain.LatLon{0, 0}, LimitB: &domain.LatLon{1, 1}},     // Area = 1
			{Value: "2", LimitA: &domain.LatLon{2, 2}, LimitB: &domain.LatLon{4, 4}},     // Area = 4
			{Value: "3", LimitA: &domain.LatLon{5, 5}, LimitB: &domain.LatLon{6, 6}},     // Area = 1
			{Value: "5", LimitA: &domain.LatLon{10, 10}, LimitB: &domain.LatLon{14, 14}}, // Area = 16
			{Value: "6", LimitA: &domain.LatLon{15, 15}, LimitB: &domain.LatLon{20, 20}}, // Area = 25
		},
	}

	l_1, l_2 := tree.SplitNodeQuadraticCost(newLocation, l)

	log.Println("L1")
	for _, v := range l_1.Locations {
		log.Println(v.Value, v.LimitA)
	}

	log.Println("l_2")
	for _, v := range l_2.Locations {
		log.Println(v.Value, v.LimitA)
	}

}
