package main

import (
	"TheWozard/standardinator/pkg/config"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {

	configFile := flag.String("config", "./config.json", "The path to configuration file")
	dataFile := flag.String("config", "./test.json", "The path to configuration file")

	flag.Parse()

	file, err := os.Open(*dataFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	conf, err := config.NewConfigFromFile(*configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	decoder, err := conf.GetDecoder()
	if err != nil {
		fmt.Println(err)
		return
	}

	extractor := decoder.New(file)
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
