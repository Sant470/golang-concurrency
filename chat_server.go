package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadCaster() {
	users := make(map[client]bool)
	select {
	case msg := <-messages:
		for user := range users {
			user<- msg
		}
	case user := <-entering:
		users[user] = true
	case user := <- leaving:
		delete(users, user)
		close(user)
	}
}

func handle(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch<- "you are : "+ who
	messages <- who + "has arrived"
	entering <- ch
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		messages <- who + ":" + scanner.Text()
	}
	leaving <- ch  
	messages <- who + ":"+ "has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for message := range ch {
		fmt.Println(conn, message)
	}
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadCaster()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handle(conn)
	}
}
