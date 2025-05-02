package vc

import ("vertex_cover/graph"
	"github.com/lukpank/go-glpk/glpk"
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
    edges := G.Edges()
    n := G.N
    m := len(edges)

    prob := glpk.New()
    defer prob.Delete()

    // 2) Thêm n biến x1..xn, 0 ≤ xi ≤ 1, objective: min Σ xi
    prob.AddCols(int(n))
    for i := int(1); i <= int(n); i++ {
        prob.SetColBnds(i, glpk.DB, 0.0, 1.0) // DB = double-bounded
        prob.SetObjCoef(i, 1.0)              // min Σ xi
    }
    prob.SetObjDir(glpk.MIN)

    // 3) Thêm m ràng buộc: x_u + x_v ≥ 1
    prob.AddRows(int(m))
    for i, e := range edges {
        rid := int(i + 1)
        // thiết lập bnds: LO = lower-bound; x_u + x_v ≥ 1
        prob.SetRowBnds(rid, glpk.LO, 1.0, 0.0)
        // gán hệ số tại row rid: xi và xj
        idx := []int32{int32(e.U + 1), int32(e.V + 1)}
        val := []float64{1.0, 1.0}
        prob.SetMatRow(rid, idx, val)
    }

    // 4) Giải LP bằng Simplex
    if err := prob.Simplex(nil); err != nil {
        log.Fatalf("GLPK Simplex failed: %v", err)
    }
    // 5) Randomized Rounding
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
    cover := make([]int, 0, n)
    for i := int(1); i <= int(n); i++ {
        xi := prob.ColPrim(i) // nghiệm xi ∈ [0,1]
        if rnd.Float64() <= xi {
            cover = append(cover, int(i-1)) // chuyển về 0-based
        }
    }
    return cover
}