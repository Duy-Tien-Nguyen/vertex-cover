package vc

import (
	"math"
	"math/rand"
	"sync"
	"time"
	"vertex_cover/graph"
)

type MPA_Solver struct {}
func (s *MPA_Solver) Name() string {
	return "MPA"
}
func (s *MPA_Solver) Solve(g *graph.Graph) ([]int, time.Duration) {
	start := time.Now()
	cover := MPA(g)
	duration := time.Since(start)
	result := make([]int, 0, len(cover))
	for v := range cover {
		result = append(result, v)
	}
	return result, duration
}

func MPA(G *graph.Graph) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	n:= G.N
	p:= 1.0/(math.Log(float64(n))*math.Log(float64(n))) //p=O(1/(log(n)*log(n)))
	C:=make(map[int]bool)

	E:=G.Edges()
	remaining:= make(map[graph.Edge]bool)
	for _, e := range E {
		remaining[e]=true
	}

	for len(remaining)>0{
		var mu sync.Mutex
		var wg sync.WaitGroup
		newCover:=make(map[int]bool)

		S:=[]graph.Edge{}
		for e:=range remaining{
			if r.Float64()<=p{
				S=append(S,e)
			}
		}
		// Sử lí song song
		wg.Add(len(S))
		for _, e := range S {
			go func (e graph.Edge){
				defer wg.Done()
				var chosen int
				if r.Intn(2)==0{
					chosen=e.U
				}else{
					chosen=e.V
				}
				mu.Lock()
				newCover[chosen]=true
				mu.Unlock()
			}(e)
		}
		wg.Wait()
		// Cập nhật E': loại bỏ cạnh mà cả hai đầu mút nằm trong cover
		for v:=range newCover{
			C[v]=true
		}
		newRemaining:=make(map[graph.Edge]bool)
		for e:=range remaining{
			if !C[e.U] && !C[e.V]{
				newRemaining[e]=true
			}
		}
		remaining=newRemaining
	}
	// Chuyển đổi map sang slice
	var cover []int
	for v := range C {
		cover = append(cover, v)
	}
	return cover
}	