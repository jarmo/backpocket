package main

import (
	"fmt"
	"os"

	"github.com/jarmo/backpocket/article"
)

func main() {
	url, err := article.URL(os.Args)
	if err != nil {
		fmt.Println(err)
		fmt.Println("\nUSAGE: backpocket ARTICLE_URL")
		os.Exit(1)
	}

	os.MkdirAll(article.RootDir, os.ModePerm)
	fmt.Println(article.Create(url))
}
