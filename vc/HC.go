package vc

import (
	"math/rand"
	"time"
	"vertex_cover/graph"
)

type HC_Solver struct{}
func (s *HC_Solver) Name() string {
	return "Hill Climbing"
}
func (s *HC_Solver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	cover := HillClimbVC(g, 1000, 100, g.N/2)
	duration := time.Since(start)
	result := make([]int, 0, len(cover))
	for v := range cover {
		result = append(result, v)
	}
	return result, duration
}

func InitialSolution(g *graph.Graph) *BitMask {
	cover := NewBitMask(g.N)
	for _, edge := range g.Edges() {
		u, v := edge.U, edge.V
		if !cover.Get(u) && !cover.Get(v) {
			cover.Set(u)
		}
	}
	return cover
}

func IsValidFast(mask *BitMask, g *graph.Graph) bool {
	for _, edge := range g.Edges() {
		u, v := edge.U, edge.V
		if !mask.Get(u) && !mask.Get(v) {
			return false
		}
	}
	return true
}

func GenerateMove(current *BitMask, g *graph.Graph, moveType string) *BitMask {
	n := g.N
	Cp := current.Clone()
	switch moveType {
	case "add":
		for {
			u := rand.Intn(n)
			if !current.Get(u) {
				Cp.Set(u)
				break
			}
		}
	case "remove":
		for {
			u := rand.Intn(n)
			if current.Get(u) {
				Cp.Unset(u)
				break
			}
		}
	case "swap":
		var u, v int
		for {
			v = rand.Intn(n)
			if current.Get(v) {
				break
			}
		}
		for {
			u = rand.Intn(n)
			if !current.Get(u) {
				break
			}
		}
		Cp.Unset(v)
		Cp.Set(u)
	}
	return Cp
}

func CountUncoveredEdges(cover *BitMask, g *graph.Graph) int {
	count := 0
	for _, edge := range g.Edges() {
		u, v := edge.U, edge.V
		if !cover.Get(u) && !cover.Get(v) {
			count++
		}
	}
	return count
}

func RandomChoiceWeighted(choices []string, probs []float64) string {
	rand.Seed(time.Now().UnixNano())
	p := rand.Float64()
	sum := 0.0
	for i, prob := range probs {
		sum += prob
		if p < sum {
			return choices[i]
		}
	}
	return choices[len(choices)-1]
}

func HillClimbVC(g *graph.Graph, maxIter, maxMovePerIter, targetSize int) []int {
	currentCover := InitialSolution(g)
	bestCover := currentCover.Clone()

	for iter := 0; iter < maxIter; iter++ {
		N := []*BitMask{}
		moveCount := 0

		uncovered := CountUncoveredEdges(currentCover, g)

		for moveCount < maxMovePerIter {
			var moveType string
			if uncovered > 0 {
				moveType = RandomChoiceWeighted([]string{"add", "swap"}, []float64{0.7, 0.3})
			} else if currentCover.Count() > targetSize {
				moveType = RandomChoiceWeighted([]string{"remove", "swap"}, []float64{0.7, 0.3})
			} else {
				moveType = RandomChoiceWeighted([]string{"add", "remove", "swap"}, []float64{0.33, 0.33, 0.34})
			}

			Cp := GenerateMove(currentCover, g, moveType)
			if IsValidFast(Cp, g) {
				N = append(N, Cp)
			}
			moveCount++
		}

		if len(N) > 0 {
			bestNeighbor := N[0]
			for _, c := range N[1:] {
				if c.Count() < bestNeighbor.Count() {
					bestNeighbor = c
				}
			}
			if bestNeighbor.Count() < bestCover.Count() {
				currentCover = bestNeighbor
				bestCover = bestNeighbor
			} else {
				break
			}
		} else {
			break
		}
	}
	res := make([]int, 0, bestCover.Count())
	for v := 0; v < g.N; v++ {
		if bestCover.Get(v) {
			res = append(res, v)
		}
	}

	return res
}

