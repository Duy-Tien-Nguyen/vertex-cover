package vc

// import (
// 	"fmt"
// )

// type Edge struct {
// 	u, v int
// }

// Đếm bậc của các đỉnh dựa trên danh sách cạnh
func buildDegree(edges []Edge) map[int]int {
	degree := make(map[int]int)
	for _, e := range edges {
		degree[e.u]++
		degree[e.v]++
	}
	return degree
}

func GreedyVertexCover(edges []Edge) []int {
	// Sao chép cạnh để xử lý mà không thay đổi đầu vào
	E := make([]Edge, len(edges))
	copy(E, edges)

	cover := make(map[int]bool)
	for len(E) > 0 {
		// Tính bậc các đỉnh hiện tại
		degree := buildDegree(E)

		// Tìm đỉnh có bậc cao nhất
		var maxV int
		maxDegree := -1
		for v, d := range degree {
			if d > maxDegree {
				maxDegree = d
				maxV = v
			}
		}

		// Thêm đỉnh vào tập cover
		cover[maxV] = true

		// Xóa các cạnh kề với maxV
		newE := []Edge{}
		for _, e := range E {
			if e.u != maxV && e.v != maxV {
				newE = append(newE, e)
			}
		}
		E = newE
	}

	// Chuyển cover từ map sang slice
	result := []int{}
	for v := range cover {
		result = append(result, v)
	}
	return result
}

// // Ví dụ sử dụng
// func main() {
// 	edges := []Edge{
// 		{0, 1},
// 		{0, 2},
// 		{1, 3},
// 		{2, 3},
// 		{3, 4},
// 		{4, 5},
// 	}

// 	cover := GreedyVertexCover(edges)
// 	fmt.Println("Vertex Cover:", cover)
// }
