# Graph Connection

Library provides a method to find the shortest paths that interconnect multiple nodes in a graph

```go
vertexes, err := gc.ShortestPaths[T](
    nodesToJoin []T{}, 
    outgoingVertexes func([]T, []T) ([]Vertex[T], error), 
    maxIterations)
```
