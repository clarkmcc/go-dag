package main

import (
	"fmt"
	"github.com/clarkmcc/go-dag"
)

func ExampleWalk() {
	g := dag.AcyclicGraph{}

	// Initialize two vertices and an edge between them
	g.Add(0)
	g.Add(1)
	g.Connect(dag.BasicEdge(0, 1))

	// Create a set representing the vertices that we want to start walking from
	s := make(dag.Set)
	s.Add(0)

	_ = g.DepthFirstWalk(s, func(vertex dag.Vertex, i int) error {
		fmt.Println(vertex)
		return nil
	})
	//Output:
	//0
	//1
}
