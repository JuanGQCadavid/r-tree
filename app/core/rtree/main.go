package rtree

import (
	"fmt"
	"log"
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
		MinValues: int(math.Ceil(float64(valuesPerNode) / 2)),
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

	if l.Parent == nil {
		newRoot := &domain.Node[T]{
			Parent:    nil,
			Locations: make([]*domain.Location[T], 0),
		}
		l_a, l_b := mathstuff.GetMRB(l.Locations)
		ll_a, ll_b := mathstuff.GetMRB(ll.Locations)

		newRoot.Locations = append(newRoot.Locations, &domain.Location[T]{
			ChildPointer: l,
			LimitA:       l_a,
			LimitB:       l_b,
		}, &domain.Location[T]{
			ChildPointer: ll,
			LimitA:       ll_a,
			LimitB:       ll_b,
		})

		l.Parent = newRoot
		ll.Parent = newRoot
		rtree.Root = newRoot
		return
	}

	// Not root, then, could we added it without conflict?
	parent := l.Parent
	ll_a, ll_b := mathstuff.GetMRB(ll.Locations)
	parent.Locations = append(parent.Locations, &domain.Location[T]{
		ChildPointer: ll,
		LimitA:       ll_a,
		LimitB:       ll_b,
	})

	// Update the MBR of the parent node
	parent.UpdateMBR()

	if len(parent.Locations) > rtree.MaxValues {
		newParent, siblingParent := rtree.SplitNodeQuadraticCost(parent)
		rtree.AdjustTree(newParent, siblingParent)
	}
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

		l_1A, l_1B = mathstuff.GetMRB(l_1)
		l_2A, l_2B = mathstuff.GetMRB(l_2)

		l_1_coords = &domain.Location[T]{
			LimitA: l_1A,
			LimitB: l_1B,
		}

		l_2_coords = &domain.Location[T]{
			LimitA: l_2A,
			LimitB: l_2B,
		}
	)

	for i, v := range origin {
		l_1_area := mathstuff.CalculateArea(mathstuff.MinSquare(l_1_coords, v))
		l_2_area := mathstuff.CalculateArea(mathstuff.MinSquare(l_2_coords, v))

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

func (rtree *RTree[T]) SplitNodeQuadraticCost(l *domain.Node[T]) (*domain.Node[T], *domain.Node[T]) {

	l_1 := make([]*domain.Location[T], 0)
	l_2 := make([]*domain.Location[T], 0)

	totalEntries := l.Locations

	// Selecting seeds
	a, aI, b, bI := rtree.PickSeeds(totalEntries)
	l_1 = append(l_1, a)
	l_2 = append(l_2, b)

	// Removing seeds
	totalEntries = utils.DeleteElements(totalEntries, aI, bI)
	for len(totalEntries) > 0 {
		totalEntries, l_1, l_2 = rtree.PickNext(totalEntries, l_1, l_2)
	}

	// Balancing, the last three elements are the ones with bigger areas.
	if len(l_1) > rtree.MinValues {
		l_2 = append(l_2, l_1[rtree.MinValues:]...)
		l_1 = l_1[0:rtree.MinValues]
	}

	if len(l_2) > rtree.MinValues {
		l_1 = append(l_1, l_2[rtree.MinValues:]...)
		l_2 = l_2[0:rtree.MinValues]
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

	l.Locations = append(l.Locations, &domain.Location[T]{
		Value:        value,
		ChildPointer: nil,
		LimitA:       latLon,
		LimitB:       latLon,
	})

	if l.Parent != nil {
		l.Parent.UpdateMBR()
	}

	if len(l.Locations) > rtree.MaxValues {
		log.Printf("Okey.. we insert more than the possible for %v \n", value)
		l_1, l_2 := rtree.SplitNodeQuadraticCost(l)

		l.Locations = l_1.Locations
		l.Parent = l_1.Parent
		l.UpdateMBR()
		rtree.AdjustTree(l, l_2)
	}
}

func (rtree *RTree[T]) Search(point *domain.LatLon, node *domain.Node[T]) []*domain.Location[T] {
	if node.IsLeaf() {
		return node.Locations
	}

	// Plan A - Lit, preciso
	for _, val := range node.Locations {
		if mathstuff.IsPointCovered(val, point) {
			return rtree.Search(point, val.ChildPointer)
		}
	}

	// Plan B In range - Area

	var (
		minAreaIndex = 0
		minArea      = math.Inf(0)
	)

	for i, val := range node.Locations {
		areaBefore := mathstuff.CalculateAreaV2(val, val)
		areaAfter := mathstuff.CalculateAreaV2(val, &domain.Location[T]{
			LimitA: point,
			LimitB: point,
		})

		delta := math.Abs(areaBefore - areaAfter)

		if delta < minArea {
			minArea = delta
			minAreaIndex = i
		}
	}

	return rtree.Search(point, node.Locations[minAreaIndex].ChildPointer)

	// return nil
}
