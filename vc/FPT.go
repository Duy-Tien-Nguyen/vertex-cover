package vc

import (
	"time"
	"vertex_cover/graph"
)

type FPTSolver struct{}

func Kernelization(g *graph.Graph, k int) (*graph.Graph, map[int]struct{}, int, map[int]int) {
	coverForced := make(map[int]struct{})
	origToNew := make(map[int]int)
	remaining := make([]int, 0)

	for i := 0; i < g.N; i++ {
		remaining = append(remaining, i)
	}

	changed := true
	for changed {
		changed = false
		for v := 0; v < g.N; v++ {
			if len(g.Adj[v]) == 0 {
				g.RemoveVertex(v)
				changed = true
			}
		}
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
		for v := 0; v < g.N; v++ {
			if len(g.Adj[v]) > k {
				coverForced[v] = struct{}{}
				g.RemoveVertex(v)
				k--
				changed = true
			}
		}
	}

	// Tạo lại map từ đỉnh ban đầu -> vị trí trong graph sau kernel hóa
	newID := 0
	for v := 0; v < len(g.Adj); v++ {
		if len(g.Adj[v]) >0 {
			origToNew[v] = newID
			newID++
		}
	}

	return g, coverForced, k, origToNew
}


func graphToMasks(g *graph.Graph) []*BitMask {
	n := g.N
	masks := make([]*BitMask, n)
	for u := 0; u < n; u++ {
		masks[u] = NewBitMask(n)
		for _, v := range g.Adj[u] {
			masks[u].Set(v)
		}
	}
	return masks
}

var bbBestSize int
var bbBestCover *BitMask

// Xóa đỉnh khỏi bitmask (để loại khỏi consideration)
func removeVertex(masks []*BitMask, v int) {
	for u := range masks {
		if masks[u].Get(v) {
			masks[u].Unset(v)
		}
	}
	masks[v] = NewBitMask(len(masks))
}

// Kiểm tra còn cạnh nào không
func hasEdge(masks []*BitMask) bool {
	for _, m := range masks {
		if m.Count() > 0 {
			return true
		}
	}
	return false
}

// Chọn một cạnh tùy ý (u,v) để phân nhánh
func pickEdge(masks []*BitMask) (int, int) {
	n := len(masks)
	for u := 0; u < n; u++ {
		block := masks[u]
		if block.Count() > 0 {
			for v := 0; v < block.n; v++ {
				if block.Get(v) {
					return u, v
				}
			}
		}
	}
	return -1, -1
}

// Greedy Matching để ước lượng lower bound
func matchingLB(masks []*BitMask) int {
	tmp := make([]*BitMask, len(masks))
	for i, m := range masks {
		tmp[i] = m.Clone()
	}
	count := 0
	for {
		u, v := pickEdge(tmp)
		if u < 0 {
			break
		}
		count++
		removeVertex(tmp, u)
		removeVertex(tmp, v)
	}
	return count
}

// Hàm đệ quy branch-and-bound
func dfsMask(masks []*BitMask, cover *BitMask, k int) {
	currSize := cover.Count()
	if currSize >= bbBestSize || currSize > k {
		return
	}
	lb := matchingLB(masks)
	if currSize+lb >= bbBestSize {
		return
	}
	if !hasEdge(masks) {
		bbBestSize = currSize
		bbBestCover = cover.Clone()
		return
	}
	u, v := pickEdge(masks)
	// Nhánh chọn u
	masks1 := make([]*BitMask, len(masks))
	for i := range masks { masks1[i] = masks[i].Clone() }
	cover1 := cover.Clone()
	cover1.Set(u)
	removeVertex(masks1, u)
	dfsMask(masks1, cover1, k)
	// Nhánh chọn v
	masks2 := make([]*BitMask, len(masks))
	for i := range masks { masks2[i] = masks[i].Clone() }
	cover2 := cover.Clone()
	cover2.Set(v)
	removeVertex(masks2, v)
	dfsMask(masks2, cover2, k)
}

// Chạy B&B trên bitmask, với đỉnh bị ép buộc và k ban đầu
func solveBBMask(masks []*BitMask, coverForced *BitMask, k int) (*BitMask, int) {
	bbBestSize = coverForced.Count()
	bbBestCover = coverForced.Clone()
	dfsMask(masks, coverForced.Clone(), k)
	return bbBestCover, bbBestSize
}

// Hàm chính để tìm Vertex Cover nhỏ nhất bằng FPT
func FindFPTMinimumVertexCover(g *graph.Graph) []int {
	n := g.N
	low, high := 0, n
	var globalCover *BitMask
	globalSize := n + 1
	// tìm nhị phân trên kích thước k
	for low <= high {
		mid := (low + high) / 2
		gker, coverForcedMap, kRemain, origToNew := Kernelization(g.Clone(), mid)
		masks:= graphToMasks(gker)
		forcedMask := NewBitMask(n)
		for v := range coverForcedMap {
            if newV, ok := origToNew[v]; ok {
                forcedMask.Set(newV)
            }
        }
		bestMask, size := solveBBMask(masks, forcedMask, kRemain)
		if size <= mid {
			globalCover = bestMask.Clone()
			globalSize = size
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	res := make([]int, 0, globalSize)
	for v := 0; v < n; v++ {
		if globalCover.Get(v) {
			res = append(res, v)
		}
	}
	return res
}

func (s *FPTSolver) Name() string { return "FPT" }

func (s *FPTSolver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	cover := FindFPTMinimumVertexCover(g)
	duration := time.Since(start)
	return cover, duration
}


