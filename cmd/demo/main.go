package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("./test.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	decoder := json.NewDecoder(file)

	token, err := decoder.Token()
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	for token != nil {
		fmt.Println(token)
		token, err = decoder.Token()
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		switch t := token.(type) {
		case json.Delim:
			if t == '{' || t == '[' {
				fmt.Println("OPEN")
			} else {
				fmt.Println("CLOSE")
			}
		default:
			fmt.Println(t)
		}
	}
}
