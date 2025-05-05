package vc

import (
	"math/rand"
	"time"
	"vertex_cover/graph"
)

type GASolver struct {
	PopSize     int     
	Generations int     
	TournamentK int     
	Pcrossover  float64 
	Pmutation   float64 
}

func (g *GASolver) Name() string { return "GeneticAlgorithm" }

func (g *GASolver) Solve(graph *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	n := graph.N
	// Default parameters if zero
	if g.PopSize == 0 {
		g.PopSize = max(50, n/5)
	}
	if g.Generations == 0 {
		g.Generations = max(100, n/5)
	}
	if g.TournamentK == 0 {
		g.TournamentK = max(3, n/20)
	}
	if g.Pcrossover <= 0 {
		g.Pcrossover = 0.8
	}
	if g.Pmutation <= 0 {
		g.Pmutation = 1.0 / float64(n)
	}

	rand.Seed(time.Now().UnixNano())

	// 1. Initialize population
	pop := generateInitialPopulation(g.PopSize, n)
	fitness := make([]float64, g.PopSize)
	for i := range pop {
		fitness[i] = evaluateFitness(pop[i], graph)
	}
	// best solution so far
	best := pop[argMax(fitness)].Clone()

	// 2. Evolution loop
	for gen := 0; gen < g.Generations; gen++ {
		// 2.1–2.6: create offspring
		newPop := make([]*BitMask, 0, g.PopSize)
		for len(newPop) < g.PopSize {
			// selection
			p1 := tournamentSelect(pop, fitness, g.TournamentK)
			p2 := tournamentSelect(pop, fitness, g.TournamentK)
			// crossover
			c1, c2 := crossover(p1, p2, g.Pcrossover)
			// mutation
			mutate(c1, g.Pmutation)
			mutate(c2, g.Pmutation)
			// repair
			repair(c1, graph)
			repair(c2, graph)
			// local search: greedy remove
			greedyRemoveLS(c1, graph)
			greedyRemoveLS(c2, graph)
			// local search: memetic bit-flip
			memeticBitFlipLS(c1, graph)
			memeticBitFlipLS(c2, graph)
			newPop = append(newPop, c1, c2)
		}
		// 2.7 Evaluate new population
		newFit := make([]float64, len(newPop))
		bestNewIdx := 0
		bestNewFit := 0.0
		for i, ind := range newPop {
			newFit[i] = evaluateFitness(ind, graph)
			if newFit[i] > bestNewFit {
				bestNewFit = newFit[i]
				bestNewIdx = i
			}
		}
		// 2.8 Update best
		if bestNewFit > evaluateFitness(best, graph) {
			best = newPop[bestNewIdx].Clone()
		}
		// 2.9 Elitism: replace worst with best
		worstIdx := argMin(newFit)
		newPop[worstIdx] = best.Clone()
		pop = newPop
		fitness = newFit
	}

	// 3. Return best as slice of vertices
	res := make([]int, 0, best.Count())
	for v := 0; v < n; v++ {
		if best.Get(v) {
			res = append(res, v)
		}
	}
	return res, time.Since(start)
}

// generateInitialPopulation creates random bitmasks then repairs them
func generateInitialPopulation(size, n int) []*BitMask {
	pop := make([]*BitMask, size)
	for i := 0; i < size; i++ {
		bm := NewBitMask(n)
		// chọn ngẫu nhiên subset
		order := rand.Perm(n)
		for j := 0; j < n; j++ {
			if rand.Float64() < 0.5 {
				bm.Set(order[j])
			}
		}
		pop[i] = bm
	}
	return pop
}

// evaluateFitness returns fitness = 1/|cover| (maximizing)
func evaluateFitness(ind *BitMask, graph *graph.Graph) float64 {
	edges := graph.Edges()
	for _, e := range edges {
		if !ind.Get(e.U) && !ind.Get(e.V) {
			return 0.0
		}
	}
	c := ind.Count()
	return 1.0 / float64(c) // càng ít đỉnh càng tốt
}

// tournamentSelect picks best of k random individuals
func tournamentSelect(pop []*BitMask, fit []float64, k int) *BitMask {
	best := rand.Intn(len(pop))
	for i := 1; i < k; i++ {
		r := rand.Intn(len(pop))
		if fit[r] > fit[best] {
			best = r
		}
	}
	return pop[best].Clone()
}

// crossover one-point
func crossover(p1, p2 *BitMask, pc float64) (*BitMask, *BitMask) {
	n := p1.n
	if rand.Float64() > pc {
		return p1.Clone(), p2.Clone()
	}
	cut := rand.Intn(n-1) + 1
	c1 := NewBitMask(n)
	c2 := NewBitMask(n)
	for i := 0; i < n; i++ {
		if i < cut {
			if p1.Get(i) {
				c1.Set(i)
			}
			if p2.Get(i) {
				c2.Set(i)
			}
		} else {
			if p2.Get(i) {
				c1.Set(i)
			}
			if p1.Get(i) {
				c2.Set(i)
			}
		}
	}
	return c1, c2
}

// mutate flips bits with probability pm
func mutate(ind *BitMask, pm float64) {
	n := ind.n
	for i := 0; i < n; i++ {
		if rand.Float64() < pm {
			ind.Toggle(i)
		}
	}
}

// repair ensures all edges are covered
func repair(ind *BitMask, graph *graph.Graph) {
	edges := graph.Edges()
	for _, e := range edges {
		if !ind.Get(e.U) && !ind.Get(e.V) {
			if rand.Float64() < 0.5 {
				ind.Set(e.U)
			} else {
				ind.Set(e.V)
			}
		}
	}
}

// greedyRemoveLS tries to remove each vertex safely
func greedyRemoveLS(ind *BitMask, graph *graph.Graph) {
	for v := 0; v < ind.n; v++ {
		if ind.Get(v) {
			ind.Unset(v)
			valid := true
			edges := graph.Edges()
			for _, e := range edges {
				if !ind.Get(e.U) && !ind.Get(e.V) {
					valid = false
					break
				}
			}
			if !valid {
				ind.Set(v)
			}
		}
	}
}

// memeticBitFlipLS flips bits greedily until no improvement
func memeticBitFlipLS(ind *BitMask, graph *graph.Graph) {
	improved := true
	for improved {
		improved = false
		prevCount := ind.Count()
		for v := 0; v < ind.n; v++ {
			ind.Toggle(v)
			valid := true
			edges := graph.Edges()
			for _, e := range edges {
				if !ind.Get(e.U) && !ind.Get(e.V) {
					valid = false
					break
				}
			}
			if valid && ind.Count() < prevCount {
				improved = true
				break
			}
			ind.Toggle(v) // revert if no improvement
		}
	}
}

// argMax returns index of max element
func argMax(arr []float64) int {
	best := 0
	for i := 1; i < len(arr); i++ {
		if arr[i] > arr[best] {
			best = i
		}
	}
	return best
}

// argMin returns index of min element
func argMin(arr []float64) int {
	best := 0
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[best] {
			best = i
		}
	}
	return best
}
