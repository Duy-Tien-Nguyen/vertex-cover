package main

// import "fmt"

func isCover(set []int, edges []Edge) bool {
	covered := make(map[Edge]bool)
	for _, e := range edges {
		covered[e] = false
	}
	for _, e := range edges {
		for _, v := range set {
			if e.u == v || e.v == v {
				covered[e] = true
			}
		}
	}
	for _, c := range covered {
		if !c {
			return false
		}
	}
	return true
}

func BruteForceVertexCover(n int, edges []Edge) []int {
	best := []int{}
	for mask := 0; mask < (1 << n); mask++ {
		var set []int
		for i := 0; i < n; i++ {
			if (mask>>i)&1 == 1 {
				set = append(set, i)
			}
		}
		if isCover(set, edges) {
			if len(best) == 0 || len(set) < len(best) {
				best = set
			}
		}
	}
	return best
}

// func main() {
// 	edges := []Edge{{0, 1}, {0, 2}, {1, 3}, {2, 3}}
// 	n := 4
// 	fmt.Println("Minimum Vertex Cover:", BruteForceVertexCover(n, edges))
// }
