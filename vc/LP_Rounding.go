package vc

import ("vertex_cover/graph"
	 "gonum.org/v1/gonum/optimize/convex/lp"
	 "gonum.org/v1/gonum/mat"
	 "log"
	 "math/rand"
	"time")

type LP_RoundingSolver struct{}
func (s *LP_RoundingSolver) Name() string {
	return "LP Rounding"
}
func (s *LP_RoundingSolver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	cover := LPRounding(g)
	duration := time.Since(start)
	result := make([]int, 0, len(cover))
	for v := range cover {
		result = append(result, v)
	}
	return result, duration
}


func LPRounding(G *graph.Graph) []int {
	edges:= G.Edges()
	n:= G.N

	c:= make([]float64, n)
	for i := 0; i < n; i++ {
		c[i] = 1
	}

	A:=mat.NewDense((len(edges)), n, nil)
	b:=make([]float64, len(edges))
	for i, edge := range edges {
		u := edge.U
		v := edge.V
		A.Set(i, u, 1)
		A.Set(i, v, 1)
		b[i] = 1
	}

	_, x, err:= lp.Simplex(c, A, b, 1e-6, nil)
	if err != nil {
		log.Fatalf("LPRounding LP solver failed: %v", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cover []int
	for v:=0;v<G.N;v++{
		randomValue:=r.Float64()
		if randomValue<=x[v]{
			cover= append(cover, v)
		}
	}
	return cover
}