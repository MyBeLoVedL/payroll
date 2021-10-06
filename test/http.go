package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://gobyexample.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Printf("status %v\n", resp.Status)
	s := bufio.NewScanner(resp.Body)
	for i := 0; i < 10; i++ {
		s.Scan()
		fmt.Println(s.Text())
	}
}
