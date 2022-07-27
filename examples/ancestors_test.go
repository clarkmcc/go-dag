package main

import (
	"fmt"
	"github.com/clarkmcc/go-dag"
)

func ExampleAncestors() {
	g := dag.AcyclicGraph{}

	// Initialize two vertices and an edge between them
	g.Add(0)
	g.Add(1)
	g.Connect(dag.BasicEdge(0, 1))

	ancestors, _ := g.Ancestors(0)
	fmt.Println(ancestors.List())
	//Output:[1]
}
