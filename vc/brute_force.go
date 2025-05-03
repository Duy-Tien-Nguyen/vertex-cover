package vc

import (
    "time"
    "math/big"
    "vertex_cover/graph"    
)

type BruteForceSolver struct{}

func isCover(mask *BitMask, g *graph.Graph) bool {
    for u:=0 ;u<g.N;u++{
        for _, v := range g.Adj[u]{
            if u<v&& !mask.Get(u)&& !mask.Get(v){
                return false
            }
        }
    }
    return true
}

// BruteForceVertexCover thử tất cả subsets của {0..n-1}, trả về cover nhỏ nhất
func BruteForceVertexCover(g *graph.Graph) *BitMask {
    n := g.N
    var bestCover *BitMask
    minSize := n+1

    limit :=new(big.Int).Lsh(big.NewInt(1), uint(n)) // 2^n
    i := new(big.Int)
    one := big.NewInt(1)

    for i.Cmp(limit) < 0 { // i < limit (so sánh big.Int)
		bm := NewBitMask(n)
		for j := 0; j < n; j++ {
			// Check if the j-th bit of i is set
			temp := new(big.Int).Set(i) // Create a copy of i
			temp.Rsh(temp, uint(j))    // Right shift by j
			bit := new(big.Int).And(temp, one)       // AND with 1

			if bit.Cmp(one) == 0 { // if (i >> j) & 1 == 1
				bm.Set(j)
			}
		}
		if isCover(bm, g) {
			if bm.Count() < minSize {
				minSize = bm.Count()
				bestCover = bm.Clone()
			}
		}
		i.Add(i, one) // i++ (cộng big.Int)
	}
	return bestCover
}



func (s *BruteForceSolver) Name() string {
    return "Brute Force"
}

func (s *BruteForceSolver) Solve(g *graph.Graph) ([]int, time.Duration) {
    start := time.Now()
    bm := BruteForceVertexCover(g)
    duration := time.Since(start)

    var result []int
    for i := 0; i < g.N; i++ {
        if bm.Get(i) {
            result = append(result, i)
        }
    }
    return result, duration
}
