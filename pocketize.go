package main

import (
	"fmt"
	"os"
)

const articlesRootDir = "articles"

func main() {
	url, err := ArticleURL(os.Args)
	if err != nil {
		//fmt.Println(err)
		fmt.Println("\nUSAGE: pocketize ARTICLE_URL")
		os.Exit(1)
	}

	os.MkdirAll(articlesRootDir, os.ModePerm)

	fmt.Println(CreateArticle(url))
}
