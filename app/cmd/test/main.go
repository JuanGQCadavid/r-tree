package main

import (
	"fmt"
	"math"
)

// Rectangle represents the minimum bounding rectangle (MBR).
type Rectangle struct {
	MinX, MinY, MaxX, MaxY float64
}

// Intersect checks if two rectangles intersect.
func (r *Rectangle) Intersect(other *Rectangle) bool {
	return r.MinX <= other.MaxX && r.MaxX >= other.MinX &&
		r.MinY <= other.MaxY && r.MaxY >= other.MinY
}

// Expand adjusts the rectangle to include another rectangle.
func (r *Rectangle) Expand(other *Rectangle) {
	r.MinX = math.Min(r.MinX, other.MinX)
	r.MinY = math.Min(r.MinY, other.MinY)
	r.MaxX = math.Max(r.MaxX, other.MaxX)
	r.MaxY = math.Max(r.MaxY, other.MaxY)
}

// Area calculates the area of the rectangle.
func (r *Rectangle) Area() float64 {
	return (r.MaxX - r.MinX) * (r.MaxY - r.MinY)
}

// Node represents a node in the R-tree.
type Node struct {
	MBR      Rectangle
	Children []*Node
	IsLeaf   bool
}

// RTree represents the R-tree structure.
type RTree struct {
	Root        *Node
	MaxChildren int
}

// NewRTree initializes a new R-tree.
func NewRTree(maxChildren int) *RTree {
	return &RTree{
		Root: &Node{
			MBR:      Rectangle{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
			IsLeaf:   true,
			Children: []*Node{},
		},
		MaxChildren: maxChildren,
	}
}

// Insert inserts a rectangle into the R-tree.
func (tree *RTree) Insert(rect Rectangle) {
	tree.insert(tree.Root, rect)

	// If the root node overflows, split it.
	if len(tree.Root.Children) > tree.MaxChildren {
		tree.splitRoot()
	}
}

func (tree *RTree) insert(node *Node, rect Rectangle) {
	if node.IsLeaf {
		// Add the rectangle as a new child node.
		node.Children = append(node.Children, &Node{
			MBR:    rect,
			IsLeaf: true,
		})
		node.MBR.Expand(&rect)
	} else {
		// Choose the best child to insert into (least area enlargement).
		bestChild := node.Children[0]
		minEnlargement := enlargement(&bestChild.MBR, &rect)

		for _, child := range node.Children[1:] {
			enlargement := enlargement(&child.MBR, &rect)
			if enlargement < minEnlargement {
				bestChild = child
				minEnlargement = enlargement
			}
		}

		// Recursively insert into the best child.
		tree.insert(bestChild, rect)
		node.MBR.Expand(&rect)
	}
}

func (tree *RTree) splitRoot() {
	oldRoot := tree.Root
	tree.Root = &Node{
		MBR:      Rectangle{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
		IsLeaf:   false,
		Children: []*Node{oldRoot},
	}
	tree.splitNode(tree.Root, oldRoot)
}

func (tree *RTree) splitNode(parent, node *Node) {
	// Simplified linear split (not optimized for balance).
	mid := len(node.Children) / 2
	left := &Node{
		MBR:      Rectangle{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
		Children: node.Children[:mid],
		IsLeaf:   node.IsLeaf,
	}
	right := &Node{
		MBR:      Rectangle{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
		Children: node.Children[mid:],
		IsLeaf:   node.IsLeaf,
	}

	for _, child := range left.Children {
		left.MBR.Expand(&child.MBR)
	}
	for _, child := range right.Children {
		right.MBR.Expand(&child.MBR)
	}

	parent.Children = append(parent.Children, left, right)
	parent.MBR.Expand(&left.MBR)
	parent.MBR.Expand(&right.MBR)
}

func enlargement(mbr, rect *Rectangle) float64 {
	expanded := *mbr
	expanded.Expand(rect)
	return expanded.Area() - mbr.Area()
}

func main() {
	rtree := NewRTree(4)

	rtree.Insert(Rectangle{MinX: 1, MinY: 1, MaxX: 2, MaxY: 2})
	rtree.Insert(Rectangle{MinX: 2, MinY: 2, MaxX: 3, MaxY: 3})
	rtree.Insert(Rectangle{MinX: 3, MinY: 3, MaxX: 4, MaxY: 4})
	rtree.Insert(Rectangle{MinX: 4, MinY: 4, MaxX: 5, MaxY: 5})
	rtree.Insert(Rectangle{MinX: 5, MinY: 5, MaxX: 6, MaxY: 6})

	fmt.Println("R-tree constructed.")
}
