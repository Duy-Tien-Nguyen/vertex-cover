package vc

import( "vertex_cover/graph"
	"time")

type DynamicProgrammingSolver struct{}

func DynamicProgramming(g *graph.Graph) (int, []int){
	dp := make([][2]int, g.N)
	visited := make([]bool, g.N)
	parent := make([]int, g.N)

	var DFS func(u int)
	DFS = func(u int) {
		visited[u] = true
		dp[u][0] = 0
		dp[u][1] = 1

		for _, v := range g.Adj[u] {
			if !visited[v] {
				parent[v] = u
				DFS(v)
				dp[u][0] += dp[v][1]
				dp[u][1] += min(dp[v][0], dp[v][1])
			}
		}
	}
	DFS(0)
	cover := []int{}
	var findCover func(u int, isInCover bool)
	findCover = func(u int, isInCover bool) {
		if isInCover {
			cover = append(cover, u)
		}
		for _, v := range g.Adj[u] {
			if parent[v] == u {
				if isInCover {
					findCover(v, false)
				} else {
					if dp[v][0] < dp[v][1] {
						findCover(v, true)
					} else {
						findCover(v, false)
					}
				}
			}
		}
	}
	findCover(0, dp[0][0] < dp[0][1])
	return min(dp[0][0], dp[0][1]), cover
}

func (s *DynamicProgrammingSolver) Name() string {
	return "Dynamic Programming"
}

func (s *DynamicProgrammingSolver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	coverSize, cover := DynamicProgramming(g)
	duration := time.Since(start)
	return cover[:coverSize], duration
}