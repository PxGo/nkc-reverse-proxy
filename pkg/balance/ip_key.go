package balance

import (
	farmhash "github.com/leemcloughlin/gofarmhash"
)

type IPKeyBalancer struct {
	nodes []string
}

func NewIPKeyBalancer() *IPKeyBalancer {
	b := new(IPKeyBalancer)
	b.nodes = []string{}
	return b
}

func (b *IPKeyBalancer) Add(names ...string) {
	b.nodes = append(b.nodes, names...)
}

// must pass in a IP string, like 127.0.0.1
func (b *IPKeyBalancer) Get(key string) string {
	count := len(b.nodes)
	if count == 1 {
		return b.nodes[0]
	}
	index := int(farmhash.Hash32([]byte(key))) % count
	return b.nodes[index]
}
