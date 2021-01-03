package main

import (
	"fmt"

	"github.com/nivista/bptree/bptree"
)

var zeroes = make([]byte, 4096)

func main() {

	tree := bptree.New()
	for {
		fmt.Print(">> ")

		var command string
		_, err := fmt.Scan(&command)
		if err != nil {
			continue
		}

		var key string
		_, err = fmt.Scan(&key)
		if err != nil {
			continue
		}

		switch command {
		case "I":
			var value string
			_, err = fmt.Scan(&value)
			if err != nil {
				continue
			}

			err = tree.Insert([]byte(key), []byte(value))
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("OK")
			}
		case "C":
			contains, err := tree.Contains([]byte(key))
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println(contains)
			}

		case "D":
			err := tree.Delete([]byte(key))
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("OK")
			}
		case "G":
			val, err := tree.Get([]byte(key))
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println(string(val))
			}
		case "E":
			break
		default:
			continue
		}
	}
}
