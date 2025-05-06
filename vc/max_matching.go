package vc

import (
	"time"
	"vertex_cover/graph"
)

type MaxMatchingSolver struct{}
func (s *MaxMatchingSolver) Name() string {
	return "Max Matching"
}
func (s *MaxMatchingSolver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	cover := maxMatchingVertexCover(g)
	duration := time.Since(start)
	result := make([]int, 0, len(cover))
	for v := range cover {
		result = append(result, v)
	}
	return result, duration
}

func maxMatchingVertexCover(G *graph.Graph) map[int]struct{} {
	// Tạo đồ thị bipartite từ đồ thị G
	C:=make(map[int]struct{})
	visitedEdges := make(map[[2]int]bool)

	for u:=0; u<G.N;u++{
		for _,v:=range G.Adj[u]{
			if u<v{
				visitedEdges[[2]int{u,v}]=false
			}else{
				visitedEdges[[2]int{v,u}]=false
			}
		}
	}

	for {
		found := false
		for u:=0;u<G.N;u++{
			for _, v := range G.Adj[u] {
				if u < v && !visitedEdges[[2]int{u, v}] {
					_, foundU := C[u]
					_, foundV := C[v]
					if !foundU && !foundV {
						// Thêm cả hai đỉnh vào C
						C[u] = struct{}{}
						C[v] = struct{}{}
						// Đánh dấu tất cả các cạnh kề với u hoặc v là “đã xét”
						for _, neighbor := range G.Adj[u] {
							if u < neighbor {
								visitedEdges[[2]int{u, neighbor}] = true
							} else {
								visitedEdges[[2]int{neighbor, u}] = true
							}
						}
						for _, neighbor := range G.Adj[v] {
							if v < neighbor {
								visitedEdges[[2]int{v, neighbor}] = true
							} else {
								visitedEdges[[2]int{neighbor, v}] = true
							}
						}
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}
		allVisited := true
		for _, visited := range visitedEdges {
			if !visited {
				allVisited = false
				break
			}
		}
		if allVisited {
			break
		}
	}

	return C
}