package main

import (
	"fmt"
	"sync"
)


func main() {
	var a, b int 
	var wg sync.WaitGroup
	wg.Add(2)
	go func(){
		defer wg.Done()
		a = 1
		fmt.Println("b: ", b)
	}()

	go func(){
		defer wg.Done()
		b = 1
		fmt.Println("a: ", a)
	}()
	wg.Wait()
}