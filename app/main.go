package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/JuanGQCadavid/r-tree/app/core/domain"
	"github.com/JuanGQCadavid/r-tree/app/core/mathstuff"
)

type Location[T any] struct {
	Value        T
	ChildPointer *Node[T]
	LimitA       *domain.LatLon
	LimitB       *domain.LatLon
}

type Node[T any] struct {
	Parent    *Node[T]
	Locations []*Location[T]
}

type RTree[T any] struct {
	Root      *Node[T]
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
func traverseNode[T any](node *Node[T], level int) {
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
func (node *Node[T]) IsLeaf() bool {
	for _, val := range node.Locations {
		if val != nil && val.ChildPointer != nil {
			return false
		}
	}
	return true
}

func NewRTree[T any](valuesPerNode int) *RTree[T] {
	return &RTree[T]{
		Root: &Node[T]{
			Parent:    nil,
			Locations: make([]*Location[T], 0, valuesPerNode),
		},
		MaxValues: valuesPerNode,
		MinValues: valuesPerNode / 2,
	}
}

func (rtree *RTree[T]) ChooseLeaf(latLon *domain.LatLon, node *Node[T]) *Node[T] {
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

func (rtree *RTree[T]) adjustTree(l *Node[T], ll *Node[T]) {

}

func (rtree *RTree[T]) pickSeeds(entries []*Location[T]) (*Location[T], *Location[T]) {

	var (
		LimitA, LimitB *Location[T]
		maxArea        = math.Inf(-1)
	)

	for i := 0; i < len(entries); i++ {
		for j := i; j < len(entries); j++ {
			areaI := mathstuff.CalculateArea(entries[i].LimitA, entries[i].LimitB)
			areaJ := mathstuff.CalculateArea(entries[j].LimitA, entries[j].LimitB)
		}
	}

	return nil, nil
}

func (rtree *RTree[T]) splitNodeQuadraticCost(newLocation *Location[T], l *Node[T]) (*Node[T], *Node[T]) {

	ll := &Node[T]{
		Parent:    nil,
		Locations: make([]*Location[T], 0, rtree.MaxValues),
	}

	totalEntries := make([]*Location[T], 0)
	totalEntries = append(totalEntries, l.Locations...)
	totalEntries = append(totalEntries, newLocation)

	return nil, nil
}

// func (rtree *RTree[T]) splitNode(latLon *domain.LatLon, value T, l *Node[T]) (*Node[T], *Node[T]) {
// 	ll := &Node[T]{
// 		Parent:    nil,
// 		Locations: make([]*Location[T], 0, rtree.MaxValues),
// 	}

// 	return nil, nil
// }

func (rtree *RTree[T]) InsertLocation(latLon *domain.LatLon, value T, node *Node[T]) {
	l := rtree.ChooseLeaf(latLon, node)

	if len(l.Locations) < cap(l.Locations) {
		l.Locations = append(l.Locations, &Location[T]{
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
