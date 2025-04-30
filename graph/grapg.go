package graph

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type Graph struct {
    N   int
    Adj [][]int
}

func ParseDIMACS(path string) (*Graph, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var n, m int
    sc := bufio.NewScanner(f)
    for sc.Scan() {
        line := sc.Text()
        if strings.HasPrefix(line, "p ") {
            fmt.Sscanf(line, "p edge %d %d", &n, &m)
            break
        }
    }

    g := &Graph{N: n, Adj: make([][]int, n)}
    for sc.Scan() {
        line := sc.Text()
        if !strings.HasPrefix(line, "e ") {
            continue
        }
        var u, v int
        fmt.Sscanf(line, "e %d %d", &u, &v)
        u--; v-- // chuyển về 0-based index
        g.Adj[u] = append(g.Adj[u], v)
        g.Adj[v] = append(g.Adj[v], u)
    }
    return g, sc.Err()
}
