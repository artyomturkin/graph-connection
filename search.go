package graphconnection

import (
	"golang.org/x/exp/maps"
)

type Vertex[T comparable] struct {
	From T
	To   T
	Meta any
}

type GetVertexesFunc[T comparable] func([]T, []T) ([]Vertex[T], error)

func Contains[T comparable](slice []T, value T) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}

	return false
}

func walkHome[T comparable](pathsHome map[T][]Vertex[T], from T, walked []T) []T {
	if vs, ok := pathsHome[from]; ok {
		parents := []T{}

		for _, v := range vs {
			if !Contains(walked, v.From) {
				parents = append(parents, walkHome(pathsHome, v.From, append(walked, v.To))...)
			}
		}
		return parents
	}

	return []T{from}
}

func gatherPath[T comparable](pathsHome map[T][]Vertex[T], from T, walked []T, startNodes []T) []Vertex[T] {
	verts := []Vertex[T]{}

	if _, ok := pathsHome[from]; !ok {
		return verts
	}

	for _, v := range pathsHome[from] {
		if Contains(startNodes, v.From) {
			verts = append(verts, v)
			continue
		}

		if Contains(walked, v.To) {
			continue
		}

		parents := gatherPath(pathsHome, v.From, append(walked, v.To), startNodes)
		if len(parents) == 0 {
			continue
		}

		if Contains(startNodes, parents[len(parents)-1].From) {
			verts = append(verts, v)
			verts = append(verts, parents...)
		}
	}

	return verts
}

func ShortestPaths[T comparable](nodesToJoin []T, adjacentVertexes GetVertexesFunc[T], maxIterations int) ([]Vertex[T], error) {
	lookup := nodesToJoin
	seen := []T{}

	pathHome := map[T][]Vertex[T]{}
	joined := []Vertex[T]{}
	joinsChecked := 0
	destinations := map[T][]Vertex[T]{}

	for _, n := range nodesToJoin {
		pathHome[n] = []Vertex[T]{}
	}

	for i := 0; i < maxIterations; i++ {
		vertexes, err := adjacentVertexes(lookup, seen)
		if err != nil {
			return nil, err
		}

		seen = append(seen, lookup...)
		lookup = []T{}
		outs := map[T][]Vertex[T]{}

		for _, vertex := range vertexes {
			if _, ok := outs[vertex.To]; !ok {
				outs[vertex.To] = []Vertex[T]{}
			}
			outs[vertex.To] = append(outs[vertex.To], vertex)

			if _, ok := destinations[vertex.To]; !ok {
				destinations[vertex.To] = []Vertex[T]{}
			}
			destinations[vertex.To] = append(destinations[vertex.To], vertex)

			if _, exists := pathHome[vertex.To]; !exists {
				pathHome[vertex.To] = []Vertex[T]{}
			}

			if !Contains(nodesToJoin, vertex.To) {
				pathHome[vertex.To] = append(pathHome[vertex.To], vertex)
			}

			if Contains(seen, vertex.To) {
				joined = append(joined, vertex)
			}

			if len(outs[vertex.To]) > 1 {
				joined = append(joined, outs[vertex.To]...)
			}
		}

		found := map[T]bool{}

		if joinsChecked != len(joined) {
			for _, join := range joined[joinsChecked:] {
				roots := walkHome(pathHome, join.From, []T{})
				for _, root := range roots {
					found[root] = true
				}
			}

			if len(found) == len(nodesToJoin) {
				break
			}

			joinsChecked = len(joined)
		}

		lookup = append(lookup, maps.Keys(outs)...)
	}

	if len(joined) == 0 {
		vs := []Vertex[T]{}
		for _, n := range nodesToJoin {
			vs = append(vs, Vertex[T]{From: n})
		}
		return vs, nil
	}

	result := map[Vertex[T]]bool{}
	for _, join := range joined {
		dests := destinations[join.To]
		deduped := []Vertex[T]{}
		for _, v := range dests {
			if !Contains(deduped, v) {
				deduped = append(deduped, v)
			}
		}

		if len(deduped) == 1 {
			continue
		}

		result[join] = true
		for _, p := range gatherPath(pathHome, join.From, nil, nodesToJoin) {
			result[p] = true
		}
	}

	return maps.Keys(result), nil
}
