package vc

import (
    "time"
    "vertex_cover/graph"
)

type Solver interface {
    Name() string
    Solve(g *graph.Graph) (cover []int, duration time.Duration)
}
