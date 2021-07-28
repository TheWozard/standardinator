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

	extractors, err := config.NewExtractionConfigFromFile("./config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	extractor, err := extractors.GetExtractor(file)
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
