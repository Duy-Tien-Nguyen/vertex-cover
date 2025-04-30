package vc

import (
	// "fmt"
)

// Hàm primal-dual vertex cover
func PrimalDualVertexCover(n int, edges []Edge) []int {
	x := make([]int, n)                     // Biến primal: x[v] = 1 nếu v ∈ cover
	y := make(map[Edge]int)                // Biến dual: y[e] = số lần cạnh e được "kích hoạt"
	C := make(map[int]bool)                // Tập kết quả: vertex cover

	for _, e := range edges {
		u, v := e.u, e.v
		// Nếu cả 2 đỉnh chưa nằm trong cover (x[u] = 0 và x[v] = 0)
		if x[u] == 0 && x[v] == 0 {
			y[e]++          // Tăng biến dual y[e]
			x[u] = 1        // Đưa u vào cover
			x[v] = 1        // Đưa v vào cover
			C[u] = true
			C[v] = true
		}
	}

	// Chuyển tập C từ map về slice
	result := []int{}
	for v := range C {
		result = append(result, v)
	}

	return result
}

// // Hàm chính để chạy thử
// func main() {
// 	// Số đỉnh
// 	n := 5
// 	// Danh sách cạnh (vô hướng)
// 	edges := []Edge{
// 		{0, 1},
// 		{0, 2},
// 		{1, 3},
// 		{3, 4},
// 	}

// 	cover := PrimalDualVertexCover(n, edges)

// 	fmt.Println("Vertex Cover tìm được:")
// 	for _, v := range cover {
// 		fmt.Printf("%d ", v)
// 	}
// 	fmt.Println()
// }
