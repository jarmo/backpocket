package main

import (
	"fmt"
	"os"

	"github.com/jarmo/pocketize/article"
)

func main() {
	url, err := article.URL(os.Args)
	if err != nil {
		fmt.Println(err)
		fmt.Println("\nUSAGE: pocketize ARTICLE_URL")
		os.Exit(1)
	}

	os.MkdirAll(article.RootDir, os.ModePerm)

	fmt.Println(article.Create(url))
}
