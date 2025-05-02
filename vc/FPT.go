package vc

import (
	"time"
	"vertex_cover/graph"
)

type FPTSolver struct{}

func Kernelization(g *graph.Graph, k int) (*graph.Graph, map[int]struct{}, int, map[int]int) {
	coverForced := make(map[int]struct{})
	origToNew := make(map[int]int)
	removed := make([]bool, g.N)
	changed := true

	for changed {
		changed = false

		for v := 0; v < g.N; v++ {
			if !removed[v] && len(g.Adj[v]) == 0 {
				g.RemoveVertex(v)
				removed[v] = true
				changed = true
			}
		}

		for v := 0; v < g.N; v++ {
			if !removed[v] && len(g.Adj[v]) == 1 {
				u := g.Adj[v][0]
				coverForced[u] = struct{}{}
				if !removed[v] {
					g.RemoveVertex(v)
					removed[v] = true
				}
				if !removed[u] {
					g.RemoveVertex(u)
					removed[u] = true
				}
				k--
				if k < 0 {
					return g, coverForced, k, origToNew
				}
				changed = true
			}
		}

		for v := 0; v < g.N; v++ {
			if !removed[v] && len(g.Adj[v]) > k {
				coverForced[v] = struct{}{}
				g.RemoveVertex(v)
				removed[v] = true
				k--
				if k < 0 {
					return g, coverForced, k, origToNew
				}
				changed = true
			}
		}
	}

	newID := 0
	for v := 0; v < g.N; v++ {
		if !removed[v] && len(g.Adj[v]) > 0 {
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

func removeVertex(masks []*BitMask, v int) {
	for u := range masks {
		if masks[u].Get(v) {
			masks[u].Unset(v)
		}
	}
	masks[v] = NewBitMask(len(masks))
}

func hasEdge(masks []*BitMask) bool {
	for _, m := range masks {
		if m.Count() > 0 {
			return true
		}
	}
	return false
}

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
	if u < 0 || v < 0 {
		return
	}
	masks1 := make([]*BitMask, len(masks))
	for i := range masks {
		masks1[i] = masks[i].Clone()
	}
	cover1 := cover.Clone()
	cover1.Set(u)
	removeVertex(masks1, u)
	dfsMask(masks1, cover1, k-1)

	masks2 := make([]*BitMask, len(masks))
	for i := range masks {
		masks2[i] = masks[i].Clone()
	}
	cover2 := cover.Clone()
	cover2.Set(v)
	removeVertex(masks2, v)
	dfsMask(masks2, cover2, k-1)
}

func solveBBMask(g *graph.Graph, k int) *BitMask {
	masks := graphToMasks(g)
	bbBestSize = k + 1
	bbBestCover = nil
	dfsMask(masks, NewBitMask(len(masks)), k)
	return bbBestCover
}

func FindFPTMinimumVertexCover(g *graph.Graph) map[int]struct{} {
	lo, hi := 0, len(g.Edges())
	var result map[int]struct{}

	for lo <= hi {
		mid := (lo + hi) / 2
		gker, coverForcedMap, kRemain, origToNew := Kernelization(g.Clone(), mid)
		if gker == nil {
			lo = mid + 1
			continue
		}
		subset := solveBBMask(gker, kRemain)
		if subset == nil {
			lo = mid + 1
			continue
		}
		finalCover := make(map[int]struct{})
		for v := range coverForcedMap {
			finalCover[v] = struct{}{}
		}
		for i := 0; i < len(gker.Adj); i++ {
			if subset.Get(i) {
				for orig, newIdx := range origToNew {
					if newIdx == i {
						finalCover[orig] = struct{}{}
						break
					}
				}
			}
		}
		result = finalCover
		hi = mid - 1
	}
	return result
}

func (s *FPTSolver) Name() string { return "FPT" }

func (s *FPTSolver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	cover := FindFPTMinimumVertexCover(g)
	duration := time.Since(start)
	result := make([]int, 0, len(cover))
	for v := range cover {
		result = append(result, v)
	}
	return result, duration
}