package vc

import (
	"math"
	"math/rand"
	"time"
	"vertex_cover/graph"
)

type ACO struct {}

func (s *ACO) Name() string {
	return "ACO"
}
func (s *ACO) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	mAnts:=max(10, g.N/2)
	maxIter:=min(2000, 20*g.N)
	tau0:=1.0/float64(g.N)
	rh0:=0.3
	alpha:=1.0
	beta:=2.0
	Q:=float64(g.N)
	cover := runACO(g, mAnts, maxIter, tau0, rh0, Q, alpha, beta)
	duration := time.Since(start)
	result := make([]int, 0, len(cover))
	for v := range cover {
		result = append(result, v)
	}
	return result, duration
}

func runACO(G *graph.Graph, mAnts, maxIter int, tau0, rho, Q, alpha, beta float64) []int {
	rand.Seed(time.Now().UnixNano())
	n := G.N
	edges := G.Edges()
	// 1. Initialize pheromone and heuristic
	tau := make([]float64, n)
	theta := make([]float64, n)
	for v := 0; v < n; v++ {
		tau[v] = tau0
		theta[v] = float64(len(G.Adj[v]))
	}
	// bestCover = all vertices initially
	bestMask := NewBitMask(n)
	for v := 0; v < n; v++ {
		bestMask.Set(v)
	}
	// 2. Main ACO loop
	for iter := 0; iter < maxIter; iter++ {
		// solutions and costs for this iteration
		solutions := make([]*BitMask, mAnts)
		costs := make([]int, mAnts)
		for k := 0; k < mAnts; k++ {
			// 2.1 build solution
			coverMask := NewBitMask(n)
			uncovered := append([]graph.Edge(nil), edges...)
			for len(uncovered) > 0 {
				// compute probabilities
				f := make([]float64, n)
				sumF := 0.0
				for v := 0; v < n; v++ {
					if !coverMask.Get(v) {
						f[v] = math.Pow(tau[v], alpha) * math.Pow(theta[v], beta)
						sumF += f[v]
					}
				}
				// 2.2 select v*
				vstar := rouletteSelect(f, sumF)
				coverMask.Set(vstar)
				// 2.3 remove incident edges
				nextUn := uncovered[:0]
				for _, e := range uncovered {
					if !coverMask.Get(e.U) && !coverMask.Get(e.V) {
						nextUn = append(nextUn, e)
					}
				}
				uncovered = nextUn
				// 2.4 update heuristic
				for v := 0; v < n; v++ {
					d := 0
					for _, e := range uncovered {
						if e.U == v || e.V == v {
							d++
						}
					}
					theta[v] = float64(d)
				}
			}
			// record solution
			solutions[k] = coverMask
			costs[k] = coverMask.Count()
			if costs[k] < bestMask.Count() {
				bestMask = coverMask.Clone()
			}
		}
		// 2.5 pheromone evaporation
		for v := 0; v < n; v++ {
			tau[v] *= (1 - rho)
		}
		// 2.6 pheromone deposition by ants
		for k := 0; k < mAnts; k++ {
			delta := Q / float64(costs[k])
			for v := 0; v < n; v++ {
				if solutions[k].Get(v) {
					tau[v] += delta
				}
			}
		}
		// 2.7 optional bestCover deposit
		bestDelta := Q / float64(bestMask.Count())
		for v := 0; v < n; v++ {
			if bestMask.Get(v) {
				tau[v] += bestDelta
			}
		}
	}
	// return bestCover as slice
	res := make([]int, 0, bestMask.Count())
	for v := 0; v < n; v++ {
		if bestMask.Get(v) {
			res = append(res, v)
		}
	}
	return res
}


// rouletteSelect chọn index dựa trên trọng số f và tổng sumF
func rouletteSelect(f []float64, sumF float64) int {
	r := rand.Float64() * sumF
	cum := 0.0
	for i, val := range f {
		cum += val
		if cum >= r {
			return i
		}
	}
	// fallback
	for i := len(f) - 1; i >= 0; i-- {
		if f[i] > 0 {
			return i
		}
	}
	return 0
}