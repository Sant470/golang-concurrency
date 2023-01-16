package main

import "fmt"

func naturals(out chan<- int) {
	for i := 0; i < 100; i++ {
		out <- i
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for el := range in {
		out <- el * el
	}
	close(out)
}

func printer(in <-chan int) {
	for e := range in {
		fmt.Println(e)
	}
}

func main() {
	n := make(chan int)
	s := make(chan int)
	go naturals(n)
	go squarer(s, n)
	printer(s)
}