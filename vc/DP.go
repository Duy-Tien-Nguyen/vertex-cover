package vc

import (
	// "fmt"
	"math"
)

// Hàm DFS
func DFS(u int, tree [][]int, visited []bool, dp [][]int) {
	visited[u] = true
	dp[u][0] = 0 // Không chọn đỉnh u
	dp[u][1] = 1 // Chọn đỉnh u

	for _, v := range tree[u] {
		if !visited[v] {
			DFS(v, tree, visited, dp)
			dp[u][0] += dp[v][1]
			dp[u][1] += int(math.Min(float64(dp[v][0]), float64(dp[v][1])))
		}
	}
}

// Hàm xây dựng danh sách kề từ danh sách cạnh
func BuildAdjList(n int, edges [][2]int) [][]int {
	tree := make([][]int, n)
	for _, edge := range edges {
		u := edge[0]
		v := edge[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	return tree
}

// Hàm giải chính
func SolveVertexCoverFromEdges(n int, edges [][2]int, root int) int {
	tree := BuildAdjList(n, edges)
	visited := make([]bool, n)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, 2)
	}

	DFS(root, tree, visited, dp)
	return int(math.Min(float64(dp[root][0]), float64(dp[root][1])))
}

// Ví dụ sử dụng
// func main() {
// 	n := 5
// 	edges := [][2]int{
// 		{0, 1},
// 		{0, 2},
// 		{1, 3},
// 		{1, 4},
// 	}

// 	root := 0
// 	result := SolveVertexCoverFromEdges(n, edges, root)
// 	fmt.Println("Kích thước vertex cover nhỏ nhất là:", result)
// }
