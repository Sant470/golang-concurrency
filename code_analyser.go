
// return number of lines of codes, comments(single line) and blank space ...
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	_blank     = "blank"
	_comment   = "comment"
	_code      = "code"
	_semaphore = 20
)

type block string

func (b block) blankline() bool {
	c := string(bytes.TrimSpace([]byte(b)))
	return len(c) == 0
}

func (b block) comment() bool {
	c := string(bytes.TrimSpace([]byte(b)))
	return strings.HasPrefix(c, "//") || strings.HasPrefix(c, "#")
}

func (b block) classifyblock(out chan <- string) {
	switch {
	case b.blankline():
		out <- _blank
	case b.comment():
		out <- _comment
	default:
		out <- _code
	}
}

func main() {
	flag.Parse()
	var count, codes, blanks, comments int 
	fileName := flag.Arg(0)
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("error opening file %s", err)
		panic(err)
	}
	defer f.Close()

	// increase the buffer capacity to use parallelism...
	class := make(chan string)
	sem := make(chan bool, _semaphore)
	reader := bufio.NewReader(f)
	for  {
		// identifying block will be a challenge ...
		line, err  := reader.ReadString('\n')
		b := block(line)
		count ++ 
		go func(b block) {
			sem <- true
			b.classifyblock(class)
			<- sem
		}(b)
		// Error will be mostly EOF ...
		if err != nil {
			break 
		}
	}
	// cf stands for classification 
	// why we are not using range to listen on channel class? 
	for i:=0; i< count; i++ {
		cf := <-class
		switch cf {
		case _blank:
			blanks ++
		case _comment:
			comments ++
		case _code:
			codes ++
		}
	}
	fmt.Println("blanks :", blanks)
	fmt.Println("codes: ", codes)
	fmt.Println("comments: ", comments)
	fmt.Println("total: ", count)
}