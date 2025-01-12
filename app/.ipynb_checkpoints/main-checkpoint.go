package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/JuanGQCadavid/r-tree/app/core/domain"
	"github.com/JuanGQCadavid/r-tree/app/core/mathstuff"
)

type RTree[T any] struct {
	Root      *domain.Node[T]
	MaxValues int
	MinValues int
}

// TraverseAndPrint traverses the R-tree and prints its hierarchy.
func (tree *RTree[T]) TraverseAndPrint() {
	if tree.Root == nil {
		fmt.Println("The R-tree is empty.")
		return
	}
	traverseNode(tree.Root, 0)
}

// Helper function to traverse a node and print its hierarchy.
func traverseNode[T any](node *domain.Node[T], level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%sNode:\n", indent)

	for i, loc := range node.Locations {
		fmt.Printf("%s  Location %d: Value: %v, Limits: [%v, %v]\n",
			indent, i+1, loc.Value, loc.LimitA, loc.LimitB)
		if loc.ChildPointer != nil {
			fmt.Printf("%s  -> Child Node:\n", indent)
			traverseNode(loc.ChildPointer, level+1)
		}
	}
}

// TODO - What about making it as a variable and chaning it when splitting?
func (node *domain.Node[T]) IsLeaf() bool {
	for _, val := range node.Locations {
		if val != nil && val.ChildPointer != nil {
			return false
		}
	}
	return true
}

func NewRTree[T any](valuesPerNode int) *RTree[T] {
	return &RTree[T]{
		Root: &domain.Node[T]{
			Parent:    nil,
			Locations: make([]*domain.Location[T], 0, valuesPerNode),
		},
		MaxValues: valuesPerNode,
		MinValues: valuesPerNode / 2,
	}
}

func (rtree *RTree[T]) ChooseLeaf(latLon *domain.LatLon, node *domain.Node[T]) *domain.Node[T] {
	if node.IsLeaf() {
		return node
	}
	var (
		PreArea   = math.Inf(0)
		PreDelta  = math.Inf(0)
		nodeIndex = 0
	)
	for i, loc := range node.Locations {
		_, _, delta, newArea := mathstuff.NewCoords(loc.LimitA, loc.LimitB, latLon)

		if delta < PreDelta {
			PreDelta = delta
			nodeIndex = i
		}

		if delta == PreDelta {
			if newArea < PreArea {
				PreArea = newArea
				nodeIndex = i
			}
		}
	}

	return rtree.ChooseLeaf(latLon, node.Locations[nodeIndex].ChildPointer)
}

func (rtree *RTree[T]) adjustTree(l *domain.Node[T], ll *domain.Node[T]) {

}

func (rtree *RTree[T]) pickSeeds(entries []*domain.Location[T]) (*domain.Location[T], int, *domain.Location[T], int) {
	var (
		LimitA, LimitB *domain.Location[T]
		aI             int
		bI             int
		maxArea        = math.Inf(-1)
	)

	for i := 0; i < len(entries); i++ {
		for j := i; j < len(entries); j++ {
			areaE1 := mathstuff.CalculateAreaV2(entries[i], entries[i])
			areaE2 := mathstuff.CalculateAreaV2(entries[j], entries[j])
			areaJ := mathstuff.CalculateAreaV2(entries[i], entries[j])

			d := areaJ - areaE1 - areaE2

			if d >= maxArea {
				LimitA = entries[i]
				aI = i
				bI = j
				LimitB = entries[j]
				maxArea = d
			}
		}
	}

	return LimitA, aI, LimitB, bI
}

func (rtree *RTree[T]) splitNodeQuadraticCost(newLocation *domain.Location[T], l *domain.Node[T]) (*domain.Node[T], *domain.Node[T]) {

	l_1 := make([]*domain.Location[T], 0, rtree.MaxValues)
	l_2 := make([]*domain.Location[T], 0, rtree.MaxValues)

	totalEntries := make([]*domain.Location[T], 0)
	totalEntries = append(totalEntries, l.Locations...)
	totalEntries = append(totalEntries, newLocation)

	for len(totalEntries) > 0 {
		a, aI, b, bI := rtree.pickSeeds(totalEntries)
		l_1 = append(l_1, a)
		l_2 = append(l_1, b)

		totalEntries[len(totalEntries)-1], totalEntries[aI] = totalEntries[aI], totalEntries[len(totalEntries)-1]
		totalEntries[len(totalEntries)-2], totalEntries[bI] = totalEntries[bI], totalEntries[len(totalEntries)-2]

		totalEntries = totalEntries[0 : len(totalEntries)-2] // Cut!
	}

	return nil, nil
}

// func (rtree *RTree[T]) splitNode(latLon *domain.LatLon, value T, l *domain.Node[T]) (*domain.Node[T], *domain.Node[T]) {
// 	ll := &domain.Node[T]{
// 		Parent:    nil,
// 		Locations: make([]*domain.Location[T], 0, rtree.MaxValues),
// 	}

// 	return nil, nil
// }

func (rtree *RTree[T]) InsertLocation(latLon *domain.LatLon, value T, node *domain.Node[T]) {
	l := rtree.ChooseLeaf(latLon, node)

	if len(l.Locations) < cap(l.Locations) {
		l.Locations = append(l.Locations, &domain.Location[T]{
			Value:        value,
			ChildPointer: nil,
			LimitA:       latLon,
			LimitB:       latLon,
		})
	} else {
		l, ll := rtree.splitNodeQuadraticCost(latLon, value, node)
		rtree.adjustTree(l, ll)
	}
}

func main() {

}
