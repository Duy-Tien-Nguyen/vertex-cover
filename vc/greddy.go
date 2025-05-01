package vc
import ("vertex_cover/graph"
"time"
)

type GreedySolver struct{}
func (s *GreedySolver) Name() string {
	return "Greedy"
}

func GreedyVertexCover(g *graph.Graph) *BitMask {
    cover := NewBitMask(g.N)

    for len(g.Edges()) > 0 {
        maxDegree := -1
        maxVertex := -1
        for v := 0; v < g.N; v++ {
            if len(g.Adj[v]) > maxDegree {
                maxDegree = len(g.Adj[v])
                maxVertex = v
            }
        }

        if maxVertex == -1 {
            break 
        }
        // Thêm đỉnh vào cover
        cover.Set(maxVertex)
        g.RemoveVertex(maxVertex)
    }

    return cover
}

func (s *GreedySolver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	cover := GreedyVertexCover(g)
	duration := time.Since(start)

	result := make([]int, 0)
	for i := 0; i < g.N; i++ {
		if cover.Get(i) {
			result = append(result, i)
		}
	}

	return result, duration
}


