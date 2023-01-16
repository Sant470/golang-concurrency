package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Func func(url string) (interface{}, error)

type entry struct {
	res result
	ready chan struct{}
}

type result struct {
	value interface{}
	err error 
}

type memo struct {
	f Func
	cache map[string]*entry
	mu sync.Mutex
}

func NewMemo(f Func) *memo {
	return &memo{f: f, cache: make(map[string]*entry)}
}

func (m *memo) Get(url string) (interface{}, error) {
	m.mu.Lock()
	e := m.cache[url]
	if e == nil {
		e = &entry{ready: make(chan struct{})}
		m.cache[url] = e 
		m.mu.Unlock()
		e.res.value, e.res.err = m.f(url)
		close(e.ready)
	} else {
		m.mu.Unlock()
		<-e.ready
	}
	return e.res.value, e.res.err
}

func httpGetBody(url string) (interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err 
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func main() {
	flag.Parse()
	m := NewMemo(httpGetBody)
	var wg sync.WaitGroup
	urls := flag.Args()
	for _, url:= range urls {
		go func(url string){
			res, _ := m.Get(url)
			fmt.Println("res: ", string(res.([]byte)[:10]))
			wg.Done()
		}(url)
	}
	wg.Wait()
}
