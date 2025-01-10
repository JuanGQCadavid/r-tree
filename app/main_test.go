package main

import (
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
