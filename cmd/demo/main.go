package main

import (
	"TheWozard/standardinator/pkg/config"
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

	extractor, err := config.GetExtractor("json", []byte("{}"), file)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		result, err := extractor.Next()
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
