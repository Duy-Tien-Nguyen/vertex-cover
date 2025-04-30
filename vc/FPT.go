package main

import "fmt"

type Edge struct {
    u, v int
}

// Hàm kiểm tra xem liệu đồ thị G có thể có vertex cover với kích thước <= k không và trả về cover
func VertexCover(edges []Edge, k int, covered map[int]bool) ([]int, bool) {
    if len(edges) == 0 {
        // Nếu không còn cạnh nào, trả về danh sách vertex cover rỗng
        result := []int{}
        for v := range covered {
            result = append(result, v)
        }
        return result, true // Tìm được vertex cover hợp lệ
    }

    if k < 0 {
        return nil, false // Nếu k < 0, không thể tìm được vertex cover hợp lệ
    }

    // Chọn một cạnh bất kỳ (u, v) từ danh sách cạnh
    e := edges[0]
    u, v := e.u, e.v

    // Nhánh 1: Chọn u vào cover
    newEdges := removeEdges(edges, u)
    newCovered := copyMap(covered)
    newCovered[u] = true
    cover, found := VertexCover(newEdges, k-1, newCovered)
    if found {
        return cover, true
    }

    // Nhánh 2: Chọn v vào cover
    newEdges = removeEdges(edges, v)
    newCovered = copyMap(covered)
    newCovered[v] = true
    cover, found = VertexCover(newEdges, k-1, newCovered)
    if found {
        return cover, true
    }

    return nil, false
}

// Hàm xóa các cạnh kề với một đỉnh u
func removeEdges(edges []Edge, u int) []Edge {
    var result []Edge
    for _, e := range edges {
        if e.u != u && e.v != u {
            result = append(result, e)
        }
    }
    return result
}

// Hàm sao chép bản đồ
func copyMap(original map[int]bool) map[int]bool {
    newMap := make(map[int]bool)
    for k, v := range original {
        newMap[k] = v
    }
    return newMap
}

func main() {
    edges := []Edge{
        {0, 1},
        {0, 2},
        {1, 3},
        {2, 3},
    }

    k := 2
    covered := make(map[int]bool)

    cover, result := VertexCover(edges, k, covered)
    if result {
        fmt.Println("Vertex cover:", cover)
    } else {
        fmt.Println("Không thể tìm vertex cover với kích thước <= k.")
    }
}
