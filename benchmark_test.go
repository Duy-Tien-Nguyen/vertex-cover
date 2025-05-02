package vc_test

import (
    "testing"
	_"net/http/pprof"
    "vertex_cover/graph"
    "vertex_cover/vc"
	"log"
	"net/http"
)

var graphName = "data/johnson8-2-4.txt"
// var graphName = "data/C250-9.txt"
// var graphName = "data/tree_graph_10000.txt"
//var graphName = "data/tree_graph.txt"

var g, err = graph.ParseEdgeList(graphName)

// func BenchmarkFPTSolver(b *testing.B) {
//     if err != nil {
//         b.Fatal(err)
//     }

//     solver := &vc.FPTSolver{}

//     for i := 0; i < b.N; i++ {
//         gClone := g.Clone() // mỗi lần clone để tránh reuse
//         solver.Solve(gClone)
//     }
// }

// func BenchmarkGreedySolver(b *testing.B) {
//     if err != nil {
//         b.Fatal(err)
//     }

//     solver := &vc.GreedySolver{}

//     for i := 0; i < b.N; i++ {
//         solver.Solve(g.Clone())
//     }
// }

func BenchmarkSolver(b *testing.B){
	if err != nil {
		b.Fatal(err)
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

			graphClone := g.Clone()

			// Benchmark each solver
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = solver.Solve(graphClone) // Call solver's solve function
			}
		})
	}
}