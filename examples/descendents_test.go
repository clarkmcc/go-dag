package main

import (
	"fmt"
	"github.com/clarkmcc/go-dag"
)

func ExampleDescendents() {
	g := dag.AcyclicGraph{}

	// Initialize two vertices and an edge between them
	g.Add(0)
	g.Add(1)
	g.Connect(dag.BasicEdge(0, 1))

	descendents, _ := g.Descendents(1)
	fmt.Println(descendents.List())
	//Output:[0]
}
