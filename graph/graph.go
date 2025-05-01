// graph/graph.go
package graph

import (
    "bufio"
    "fmt"
    // "io"
    "os"
    "strings"
    "strconv"
)

// Edge biểu diễn một cạnh vô hướng 0-based
type Edge struct { U, V int }

// Graph lưu adjacency list
type Graph struct {
    N   int
    Adj [][]int
}

func ParseEdgeList(path string) (*Graph, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    maxNode := 0
    edges := [][2]int{}

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        parts := strings.Fields(line)
        if len(parts) != 2 {
            return nil, fmt.Errorf("invalid edge line: %q", line)
        }

        u, err1 := strconv.Atoi(parts[0])
        v, err2 := strconv.Atoi(parts[1])
        if err1 != nil || err2 != nil {
            return nil, fmt.Errorf("invalid numbers in line: %q", line)
        }

        if u > maxNode {
            maxNode = u
        }
        if v > maxNode {
            maxNode = v
        }

        // chuyển về 0-based
        edges = append(edges, [2]int{u - 1, v - 1})
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }

    g := &Graph{
        N:   maxNode,
        Adj: make([][]int, maxNode),
    }

    for _, e := range edges {
        u, v := e[0], e[1]
        g.Adj[u] = append(g.Adj[u], v)
        g.Adj[v] = append(g.Adj[v], u) // undirected
    }

    return g, nil
}

// Edges trả về danh sách các Edge {u<v} (loại trùng lặp)
func (g *Graph) Edges() []Edge {
    seen := make(map[[2]int]bool)
    var out []Edge
    for u, nbrs := range g.Adj {
        for _, v := range nbrs {
            if u < v && !seen[[2]int{u, v}] {
                out = append(out, Edge{U: u, V: v})
                seen[[2]int{u, v}] = true
            }
        }
    }
    return out
}

func (g *Graph) GetMaxDegree() int {
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

func (g *Graph) Clone() *Graph{
    newGraph :=&Graph{
        N: g.N,
        Adj: make([][]int, g.N),
    }
    
    for i:=0;i<g.N;i++{
        newGraph.Adj[i]=append([]int(nil), g.Adj[i]...)
    }
    return newGraph
}


func (g *Graph) RemoveVertex(v int) {
    for u := 0; u < len(g.Adj); u++ {
        if u == v {
            continue
        }
        adj :=g.Adj[u]
        for i:=0;i<len(adj);i++{
            if adj[i]==v{
                adj = append(adj[:i], adj[i+1:]...)
                break
            }
        }
        g.Adj[u] = adj
    }
    g.Adj[v] = nil
}
