package main

import (
	"fmt"
	"log"
	"sort"
	"github.com/shawnfeng/consistent"

	"github.com/shawnfeng/sutil/stime"
	"github.com/shawnfeng/sutil/slog"


	"github.com/serialx/hashring"
)

func initHash0() (*consistent.Consistent, []string) {
	c := consistent.New()
	var vs []string
	for i := 0; i < 20; i++ {
		for j := 0; j < 100; j++ {
			vs = append(vs, fmt.Sprintf("SERV%d-%d", i, j))
			c.Add(fmt.Sprintf("SERV%d-%d", i, j))

		}
	}


	for i := 0; i < 10; i++ {
		for j := 0; j < 1; j++ {
			vs = append(vs, fmt.Sprintf("SERV%d-%d", i, j))
			c.Add(fmt.Sprintf("SERV%d-%d", i, j))
		}
	}

	return c, vs
}


func initHash1() *consistent.Consistent {
	c := consistent.New()
	var vs []string
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			vs = append(vs, fmt.Sprintf("SERV%d-%d", i, j))
		}
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 1; j++ {
			vs = append(vs, fmt.Sprintf("SERV%d-%d", i, j))
		}
	}

	c.Set(vs)
	return c
}


func initHash2() (*consistent.Consistent, []string) {

	var vs []string
	for i := 0; i < 20; i++ {
		for j := 0; j < 100; j++ {
			vs = append(vs, fmt.Sprintf("SERV%d-%d", i, j))
		}
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 1; j++ {
			vs = append(vs, fmt.Sprintf("SERV%d-%d", i, j))
		}
	}

	c := consistent.NewWithElts(vs)
	return c, vs
}



func test0() {
	st := stime.NewTimeStat()
	c, elts := initHash0()
	tm := st.Duration()
	slog.Infoln("init hash0 tm", tm, tm)
	st.Reset()
	c2, elts2 := initHash2()
	tm = st.Duration()
	slog.Infoln("init hash1 tm", tm, tm)

	cmem := c.Members()
	cmem2 := c2.Members()


	sort.Strings(cmem)
	sort.Strings(cmem2)

	sort.Strings(elts)
	sort.Strings(elts2)

	fmt.Printf("elts mem %d %d\n", len(elts), len(elts2))
	if len(elts) != len(elts2) {
		slog.Errorf("not equal")
		return
	}

	for idx, m := range elts {
		if m != elts2[idx] {
			slog.Errorf("not equal elts")
			return
		}
	}

	//slog.Infof("yes equal elts %s", elts)


	fmt.Printf("mem %d %d\n", len(cmem), len(cmem2))
	if len(cmem) != len(cmem2) {
		slog.Errorf("not equal")
		return
	}

	for idx, m := range cmem {
		if m != cmem2[idx] {
			slog.Errorf("not equal memeber")
			return

		}
	}

	//slog.Infof("yes equal %s", cmem)

	count := hashit(c, cmem)
	count2 := hashit(c2, cmem2)

	if len(count) != len(count2) {
		slog.Errorf("hash not equal")
		return
	}

	for k, v := range count {
		v2, ok := count2[k]
		if !ok {
			slog.Errorln("hash not equal", k)
			return
		}

		if v != v2 {
			slog.Errorln("hash not equal value", k, v, v2)
			return
		}

		slog.Infoln("equal ok", k, v, v2)
	}

	slog.Infof("YES YES FUCK FUCK")

}

func hashit(c *consistent.Consistent, cmem []string) map[string]int {

	count := make(map[string]int)
	for _, e := range cmem {
		count[e] = 0
	}


	loopcn := 2000000
	st := stime.NewTimeStat()

	for i := 0; i < loopcn; i++ {
		server, err := c.Get(fmt.Sprintf("%d", i))
		if err != nil {
			log.Fatal(err)
		}
		count[server]++
	}
	tm := st.Duration()
	fmt.Printf("loop:%d tm:%s dur:%d avg:%d\n", loopcn, tm, tm, int(tm)/loopcn)

	/*
	for k, v := range count {
		fmt.Println(k, v)
	}
    */

	return count


}

func test1() {
	count := make(map[string]int)
	weights := make(map[string]int)
	for i := 0; i < 100000; i++ {
		s := fmt.Sprintf("SERV%d", i)
		weights[s] = 100
		count[s] = 0
	}

	ring := hashring.NewWithWeights(weights)

	loopcn := 1000000
	st := stime.NewTimeStat()

	for i := 0; i < loopcn; i++ {
		server, _ := ring.GetNode(fmt.Sprintf("%d", i))
		count[server]++
	}
	tm := st.Duration()
	fmt.Printf("loop:%d tm:%s dur:%d avg:%d\n", loopcn, tm, tm, int(tm)/loopcn)


	//for k, v := range count {
	//	fmt.Println(k, v)
	//}


}

func main() {
	test0()

}
