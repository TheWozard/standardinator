package main

import (
	"TheWozard/standardinator/pkg/config"
	"TheWozard/standardinator/pkg/core"
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

	tokenizer, err := config.NewTokenizer("json", file)
	if err != nil {
		fmt.Println(err)
		return
	}

	iter := core.NewIterator(tokenizer)

	for {
		result, err := iter.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		fmt.Println(result)
	}
}
