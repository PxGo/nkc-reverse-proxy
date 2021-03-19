package balance

import (
	"math/rand"
	"time"
)

type RandomBalancer struct {
	node []string
}

func NewRandomBalancer() *RandomBalancer {
	r := new(RandomBalancer)
	r.node = []string{}
	return r
}

func (r *RandomBalancer) Add(names ...string) {
	r.node = append(r.node, names...)
}

// must pass a number string
func (r *RandomBalancer) Get() string {
	rs := rand.NewSource(time.Now().Unix())
	generator := rand.New(rs)
	index := generator.Intn(len(r.node))
	return r.node[index]
}
