package balance

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
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
	for _, tag := range names {
		sp := strings.Split(tag, " ")
		if len(sp) == 1 {
			r.node = append(r.node, sp[0])
		} else {
			weight, err := strconv.ParseInt(sp[1], 10, 64)
			if err != nil {
				fmt.Printf("加权随机选择目标服务器时出错\n%s\n%s", err.Error(), debug.Stack())
				os.Exit(1)
			}
			for i := 0; i < int(weight); i++ {
				r.node = append(r.node, sp[0])
			}
		}
	}
}

func (r *RandomBalancer) Get() string {
	rs := rand.NewSource(time.Now().Unix())
	generator := rand.New(rs)
	index := generator.Intn(len(r.node))
	return r.node[index]
}
