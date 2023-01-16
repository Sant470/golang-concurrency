package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var verbose = flag.Bool("v", false, "show verbose progress message")
var sem = make(chan struct{}, 20)

func walkDir(dir string, n *sync.WaitGroup, filesizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			n.Add(1)
			walkDir(subdir,n, filesizes)
		} else {
			filesizes <- entry.Size()
		} 
	}
}

func dirents(dir string) []os.FileInfo {
	sem <- struct{}{}
	defer func(){ <-sem }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil 
	}
	return entries
}

func main() {
	flag.Parse()
	var tick <- chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	filesizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, filesizes)
	}
	go func(){
		n.Wait()
		close(filesizes)
	}()
	// Print the results
	var nfiles, nbytes int64
	loop:
		for {
			select {
			case size, ok := <-filesizes:
				if !ok {
					break loop
				}
				nfiles ++
				nbytes += size
			case <-tick:
				printDiskUsage(nfiles, nbytes)
			}
		}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}