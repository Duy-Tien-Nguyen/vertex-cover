package vc

import (
    "time"
    "vertex_cover/graph"
)

type FPTSolver struct{}

func Kernelization(g *graph.Graph, k int) (*graph.Graph, map[int]struct{}, int) {
    coverForced := make(map[int]struct{})
    changed := true
    for changed {
        changed = false
        // 1. Đỉnh cô lập
        for v := 0; v < g.N; v++ {
            if len(g.Adj[v]) == 0 {
                g.RemoveVertex(v)
                changed = true
            }
        }
        // 2. Đỉnh có bậc 1
        for v := 0; v < g.N; v++ {
            if len(g.Adj[v]) == 1 {
                u := g.Adj[v][0]
                coverForced[v] = struct{}{}
                g.RemoveVertex(v)
                g.RemoveVertex(u)
                k--
                changed = true
            }
        }
        // 3. Đỉnh có bậc lớn hơn k
        for v := 0; v < g.N; v++ {
            if len(g.Adj[v]) > k {
                coverForced[v] = struct{}{}
                g.RemoveVertex(v)
                k--
                changed = true
            }
        }
    }
    return g, coverForced, k
}

func getMaxDegree(g *graph.Graph) int {
    maxDegree := -1
    maxVertex := -1
    for v := 0; v < g.N; v++ {
        if len(g.Adj[v]) > maxDegree {
            maxDegree = len(g.Adj[v])
            maxVertex = v
        }
    }
    return maxVertex
}

func BranchAndBound(g *graph.Graph, coverCurr map[int]struct{}, kCurr int, bestCover map[int]struct{}, bestSize *int) {
    if len(g.Edges()) == 0 {
        if len(coverCurr) < *bestSize {
            for k := range bestCover {
                delete(bestCover, k)
            }
            for v := range coverCurr {
                bestCover[v] = struct{}{}
            }
            *bestSize = len(coverCurr)
        }
        return
    }

    if kCurr <= 0 || len(coverCurr)+g.N <= *bestSize {
        return
    }

    v := getMaxDegree(g)

    // Nhánh include v
    G1 := g.Clone()
    G1.RemoveVertex(v)
    coverIncludeV := make(map[int]struct{}, len(coverCurr)+1)
    for k := range coverCurr {
        coverIncludeV[k] = struct{}{}
    }
    coverIncludeV[v] = struct{}{}
    BranchAndBound(G1, coverIncludeV, kCurr-1, bestCover, bestSize)

    // Nhánh exclude v
    Nv := g.Adj[v]
    G2 := g.Clone()
    for _, neighbor := range Nv {
        G2.RemoveVertex(neighbor)
    }
    coverExcludeV := make(map[int]struct{}, len(coverCurr)+len(Nv))
    for k := range coverCurr {
        coverExcludeV[k] = struct{}{}
    }
    for _, neighbor := range Nv {
        coverExcludeV[neighbor] = struct{}{}
    }
    BranchAndBound(G2, coverExcludeV, kCurr-len(Nv), bestCover, bestSize)
}

func FindFPTMinimumVertexCover(g *graph.Graph) ([]int) {
    n := g.N  // số lượng đỉnh trong đồ thị
    low, high := 0, n // Tìm kiếm nhị phân trong khoảng [0, n]

    // Bước 2: Tìm kiếm nhị phân trên giá trị k
    bestCover :=make(map[int]struct{})
    var bestSize int
    var coverForced map[int]struct{}
    for low <= high {
        mid := (low + high) / 2
        // Bước 3: Thử nghiệm với k = mid
        gker, coverForced, k := Kernelization(g, mid)
        BranchAndBound(gker, coverForced, k, bestCover, &bestSize)

        // Cập nhật kết quả nếu cover nhỏ hơn hoặc bằng k
        if len(bestCover) <= mid {
            bestSize = len(bestCover)
            high = mid - 1  // Tìm kiếm phần dưới
        } else {
            low = mid + 1  // Tìm kiếm phần trên
        }
    }

    finalCover := make(map[int]struct{}, len(bestCover)+len(coverForced))
    for v := range bestCover {
        finalCover[v] = struct{}{}
    }
    for v := range coverForced {
        finalCover[v] = struct{}{}
    }

    cover := make([]int, 0, len(finalCover))
    for v := range finalCover {
        cover = append(cover, v)
    }

    return cover
}

func (s *FPTSolver) Name() string {
    return "FPT"
}

func (s *FPTSolver) Solve(g *graph.Graph) ([]int, time.Duration) {
    start := time.Now()
    cover := FindFPTMinimumVertexCover(g)
    duration := time.Since(start)

    return cover, duration
}
