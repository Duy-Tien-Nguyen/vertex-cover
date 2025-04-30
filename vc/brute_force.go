// package main
package vc

import (
    // "flag"                  // Đã bị comment, cần dùng để đọc -i
    // "fmt"
    "time"
    "vertex_cover/graph"    // Đảm bảo đúng module path
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
func BruteForceVertexCover(g *graph.Graph) *BitMask { // Fix: `Grapg` → `Graph`
    n := g.N
    var bestCover *BitMask
    minSize := n+1

    limit := 1 << n // 2^n
    for i:=0; i<limit ; i++ {
        bm:= NewBitMask(n)
        for j:=0; j<n ;j++{
            if (i>>j)&1==1{
                bm.Set(j)
            }
        }
        if isCover(bm, g) {
            if bm.Count() < minSize {
                minSize = bm.Count()
                bestCover = bm.Clone()
            }   
        }     
    }
    if bestCover == nil {
        bestCover = NewBitMask(n)
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


// func main() {
//     path := flag.String("i", "", "Input edge list")
//     flag.Parse()
//     if *path == "" {
//         fmt.Println("Usage: go run brute_force.go -i input.txt")
//         return
//     }

//     g, err := graph.ParseEdgeList(*path)
//     if err != nil {
//         panic(err)
//     }

//     start := time.Now()
//     cover := BruteForceVertexCover(g)
//     dt := time.Since(start)

//     fmt.Printf("BruteForce: |C| = %d, time = %s\n", cover.Count(), dt)
//     fmt.Printf("Cover vertices: %v\n", cover)
// }
