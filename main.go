package main

import (
    "vertex_cover/graph"
    "vertex_cover/vc"
    "fmt"
)

func main() {
    g, _ := graph.ParseDIMACS("input.dimacs")
    solver := vc.NewGreedySolver()
    cover, dt := solver.Solve(g)
    fmt.Println("Size:", len(cover), "Time:", dt)
}
