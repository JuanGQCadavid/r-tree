package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/JuanGQCadavid/r-tree/app/core/domain"
	"github.com/JuanGQCadavid/r-tree/app/core/mathstuff"
	"github.com/JuanGQCadavid/r-tree/app/core/utils"
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

func (rtree *RTree[T]) AdjustTree(l *domain.Node[T], ll *domain.Node[T]) {
	// Are we in the root?

	// if l.Parent == nil {
	// 	newRoot := &domain.Node[T]{
	// 		Parent:    nil,
	// 		Locations: make([]*domain.Location[T], 0, rtree.MaxValues),
	// 	}
	// 	ll_a, ll_b := mathstuff.MinSquare()
	// }
}

func (rtree *RTree[T]) PickSeeds(entries []*domain.Location[T]) (*domain.Location[T], int, *domain.Location[T], int) {
	var (
		LimitA, LimitB *domain.Location[T]
		aI             int
		bI             int
		maxArea        = math.Inf(-1)
	)

	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
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

func (rtree *RTree[T]) PickNext(origin, l_1, l_2 []*domain.Location[T]) ([]*domain.Location[T], []*domain.Location[T], []*domain.Location[T]) {

	var (
		nextIndex = 0
		minArea   = math.Inf(0)
		to_l_1    = true

		l_1_coords = l_1[0]
		l_2_coords = l_1[0]
	)

	for i := 1; i < len(l_1); i++ {
		a, b := mathstuff.MinSquare(l_1[i], l_1_coords)
		l_1_coords = &domain.Location[T]{
			LimitA: a,
			LimitB: b,
		}
	}

	for i := 1; i < len(l_2); i++ {
		a, b := mathstuff.MinSquare(l_2[i], l_2_coords)
		l_2_coords = &domain.Location[T]{
			LimitA: a,
			LimitB: b,
		}
	}

	for i, v := range origin {
		l_1_area := mathstuff.CalculateArea(mathstuff.MinSquare(l_1_coords, v))
		l_2_area := mathstuff.CalculateArea(mathstuff.MinSquare(l_1_coords, v))

		if l_1_area < minArea {
			nextIndex = i
			minArea = l_1_area
			to_l_1 = true
		}

		if l_2_area < minArea {
			nextIndex = i
			minArea = l_2_area
			to_l_1 = false
		}
	}

	if to_l_1 {
		l_1 = append(l_1, origin[nextIndex])
	} else {
		l_2 = append(l_2, origin[nextIndex])
	}

	origin = utils.DeleteElements(origin, nextIndex)

	return origin, l_1, l_2
}

func (rtree *RTree[T]) SplitNodeQuadraticCost(newLocation *domain.Location[T], l *domain.Node[T]) (*domain.Node[T], *domain.Node[T]) {

	l_1 := make([]*domain.Location[T], 0, rtree.MaxValues)
	l_2 := make([]*domain.Location[T], 0, rtree.MaxValues)

	totalEntries := make([]*domain.Location[T], 0)
	totalEntries = append(totalEntries, l.Locations...)
	totalEntries = append(totalEntries, newLocation)

	// Selecting seeds
	a, aI, b, bI := rtree.PickSeeds(totalEntries)
	l_1 = append(l_1, a)
	l_2 = append(l_2, b)

	// Removing seeds
	totalEntries = utils.DeleteElements(totalEntries, aI, bI)

	// for _, v := range totalEntries {
	// 	fmt.Print(v.Value, ", ")
	// }
	// fmt.Println()

	// log.Println("----- start -----")
	for len(totalEntries) > 0 {
		totalEntries, l_1, l_2 = rtree.PickNext(totalEntries, l_1, l_2)
		// log.Println("a: ", a.Value)
		// log.Println("b: ", b.Value)

		// for _, v := range totalEntries {
		// 	fmt.Print(v.Value, ", ")
		// }
		// fmt.Println()
		// log.Println("----")
	}

	// Balancing, the last three elements are the ones with bigger areas.
	if len(l_1) > 3 {
		l_2 = append(l_2, l_1[3:]...)
		l_1 = l_1[0:3]
	}

	if len(l_2) > 3 {
		l_1 = append(l_1, l_2[3:]...)
		l_2 = l_2[0:3]
	}

	return &domain.Node[T]{
			Parent:    l.Parent,
			Locations: l_1,
		}, &domain.Node[T]{
			Parent:    l.Parent,
			Locations: l_2,
		}
}

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
		l, ll := rtree.SplitNodeQuadraticCost(&domain.Location[T]{
			Value:  value,
			LimitA: latLon,
			LimitB: latLon,
		}, node)
		rtree.AdjustTree(l, ll)
	}
}

func main() {

}
