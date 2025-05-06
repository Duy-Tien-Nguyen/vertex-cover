package vc_test

import (
    "testing"
    "vertex_cover/graph"
    "vertex_cover/vc"
)

// var graphName = "data/johnson8-2-4.txt"
// var graphName = "data/C250-9.txt"
// var graphName = "data/C500-9.txt"
// var graphName = "data/C1000-9.txt"
// var graphName = "data/tree_graph_10000.txt"
// var graphName = "data/tree_graph_5000.txt"
// var graphName = "data/tree_graph_28.txt"
var graphName = "data/tree_graph_250.txt"
// var graphName = "data/tree_graph_500.txt"
// var graphName = "data/tree_graph_1000.txt"

var g, err = graph.ParseEdgeList(graphName)

func BenchmarkSolver(b *testing.B){
	if err != nil {
		b.Fatal(err)
	}

	solvers := []vc.Solver{
		// &vc.BruteForceSolver{},
		// &vc.DynamicProgrammingSolver{},
		// &vc.FPTSolver{},
		// &vc.GreedySolver{},
		// &vc.PrimalDualSolver{},
		// &vc.MaxMatchingSolver{},
		// &vc.LP_RoundingSolver{},
		// &vc.MPA_Solver{},
		// &vc.ACO{},
		// &vc.GASolver{},
		&vc.HC_Solver{},
	}

	for _, solver := range solvers {

		solverName := solver.Name()

		b.Run(solverName, func(b *testing.B) {
			// Benchmark each solver
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				graphClone := g.Clone()
				_, _ = solver.Solve(graphClone) // Call solver's solve function
			}
		})
	}
}