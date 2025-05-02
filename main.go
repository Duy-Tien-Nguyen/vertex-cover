package main
// package vc

import (
    "vertex_cover/graph"
    "os"
    "vertex_cover/vc"
    "fmt"
    "sync"
)

func main() {
    // graphName := "data/johnson8-2-4.txt"
    // graphName := "data/C250-9.txt"
    // graphName := "data/C500-9.txt"
    graphName := "data/C1000-9.txt"
    // graphName := "data/C2000-9.txt"
    // graphName := "data/C4000-5.txt"
    // graphName := "data/tree_graph_10000.txt"
    // graphName := "data/tree_graph.txt"
    g, err := graph.ParseEdgeList(graphName)
    if err != nil {
        panic(err)
    }

    // in ra một số thống kê
    fmt.Println("Số đỉnh N =", g.N)
    edgeCount := 0
    for u, nbrs := range g.Adj {
        edgeCount += len(nbrs)
        fmt.Printf("đỉnh %d has %d neighbors\n", u, len(nbrs))
    }
    // vì ta lưu 2 chiều, chia đôi
    fmt.Println("Tổng số cạnh (undirected) =", edgeCount/2)

    solvers := []vc.Solver{
        // &vc.BruteForceSolver{},
        // &vc.DynamicProgrammingSolver{},
        // &vc.FPTSolver{},
        &vc.GreedySolver{},
        &vc.PrimalDualSolver{},
        &vc.MaxMatchingSolver{},
        &vc.LP_RoundingSolver{},
        &vc.MPA_Solver{},
        // &vc.ACO{},
        // &vc.GASolver{},
        &vc.HC_Solver{},
    }  
    var wg sync.WaitGroup  
    for _, solver := range solvers {
        wg.Add(1)
        go func(solver vc.Solver) {
            defer wg.Done()
            cover, duration := solver.Solve(g.Clone())
            fmt.Printf("%s: |C| = %d, time = %s\n", solver.Name(), len(cover), duration)
            filename := fmt.Sprintf("Benchmark/%s.txt", solver.Name())
            f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
            if err != nil {
                panic(err)
            }
            defer f.Close()
            if _, err := fmt.Fprintf(f,
                "%s: Graph: %s, |C| = %d, time = %s\n\n",
                solver.Name(), graphName, len(cover), duration,
            ); err != nil {
                panic(err)
            }
        }(solver)
    }
    wg.Wait()
}

