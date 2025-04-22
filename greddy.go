package main

// import "fmt"

type Edge struct {
	u, v int
}

func GreedyVertexCover(edges []Edge) []int {
	covered := make(map[int]bool)     // các đỉnh đã được chọn
	usedEdges := make(map[Edge]bool)  // các cạnh đã được phủ

	for _, e := range edges {
		if !usedEdges[e] && !covered[e.u] && !covered[e.v] {
			// Thêm cả hai đỉnh nếu chưa được phủ
			covered[e.u] = true
			covered[e.v] = true

			// Đánh dấu các cạnh có liên quan đến u hoặc v là đã phủ
			for _, other := range edges {
				if other.u == e.u || other.v == e.u || other.u == e.v || other.v == e.v {
					usedEdges[other] = true
				}
			}
		}
	}

	// Tạo kết quả
	result := []int{}
	for v := range covered {
		result = append(result, v)
	}
	return result
}

// func main() {
// 	edges := []Edge{
// 		{0, 1},
// 		{0, 2},
// 		{1, 3},
// 		{2, 3},
// 	}

	// cover := GreedyVertexCover(edges)
	// fmt.Println("Vertex cover:", cover)
// }
