package graphconnection_test

import (
	"testing"

	gc "github.com/artyomturkin/graph-connection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type graph struct {
	adjacent map[int][]int
}

func (g *graph) GetVertexes(nodes, seen []int) ([]gc.Vertex[int], error) {
	verts := []gc.Vertex[int]{}

	for _, n := range nodes {
		for _, e := range g.adjacent[n] {
			if !gc.Contains(seen, e) {
				verts = append(verts, gc.Vertex[int]{From: n, To: e})
			}
		}
	}

	return verts, nil
}

var g = graph{
	adjacent: map[int][]int{
		0:  {4},
		1:  {4},
		4:  {0, 1, 6},
		2:  {5},
		3:  {5},
		5:  {2, 3, 6},
		6:  {4, 5, 7},
		7:  {6, 8},
		8:  {7, 9, 10},
		9:  {8, 11, 12},
		10: {8, 13, 14},
		11: {9},
		12: {9},
		13: {10},
		14: {10},
	},
}

func TestSearchCloseBy(t *testing.T) {
	join := []int{0, 1}

	expected := []gc.Vertex[int]{
		{From: 1, To: 4},
		{From: 0, To: 4},
	}

	vertexes, err := gc.ShortestPaths(join, g.GetVertexes, 10)
	require.NoError(t, err)

	assert.ElementsMatch(t, expected, vertexes)
}

func TestSearchFar(t *testing.T) {
	join := []int{0, 14}

	expected := []gc.Vertex[int]{
		// Starting from 0
		{From: 0, To: 4},
		{From: 4, To: 6},
		{From: 6, To: 7},
		// Starting from 14
		{From: 14, To: 10},
		{From: 10, To: 8},
		{From: 8, To: 7},
	}

	vertexes, err := gc.ShortestPaths(join, g.GetVertexes, 10)
	require.NoError(t, err)

	assert.ElementsMatch(t, expected, vertexes)
}

func TestSearchThree(t *testing.T) {
	join := []int{0, 3, 8}

	expected := []gc.Vertex[int]{
		// Starting from 0
		{From: 0, To: 4},
		{From: 4, To: 6},
		// Starting from 3
		{From: 3, To: 5},
		{From: 5, To: 6},
		// Starting from 8
		{From: 8, To: 7},
		{From: 7, To: 6},
	}

	vertexes, err := gc.ShortestPaths(join, g.GetVertexes, 10)
	require.NoError(t, err)

	assert.ElementsMatch(t, expected, vertexes)
}
