package vc

import ("vertex_cover/graph"
"time")

type PrimalDualSolver struct{}
func (s *PrimalDualSolver) Name() string {
	return "Primal-Dual"
}
func (s *PrimalDualSolver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	cover := primalDualVertexCover(g)
	duration := time.Since(start)
	result := make([]int, 0, len(cover))
	for v := range cover {
		result = append(result, v)
	}
	return result, duration
}

func primalDualVertexCover(G *graph.Graph) map[int]struct{} {
    x := make(map[int]int)   // Biến primal, 0 = chưa chọn v
    y := make(map[[2]int]int) // Biến dual
    C := make(map[int]struct{}) // Tập vertex cover

    // Khởi tạo giá trị ban đầu cho x và y
    for v := 0; v < G.N; v++ {
        x[v] = 0 // chưa chọn v
    }
    for _, e := range G.Edges() {
        y[[2]int{e.U, e.V}] = 0 // giá trị dual ban đầu là 0
    }
    // Lặp cho đến khi mọi cạnh được cover
    for {
        covered := true
        for _, e := range G.Edges() {
            u, v := e.U, e.V
            if x[u] == 0 && x[v] == 0 {
                // Cạnh (u,v) chưa được cover
                y[[2]int{u, v}]++
                x[u] = 1
                x[v] = 1
                C[u] = struct{}{}
                C[v] = struct{}{}
                covered = false
            }
        }
        if covered {
            break
        }
    }
    return C
}
