package main

import (
	"fmt"
	"os"

	"github.com/jarmo/backpocket/article"
)

func main() {
	params, err := article.Params(os.Args)
	if err != nil {
		fmt.Println(err)
		fmt.Println("\nUSAGE: backpocket ARTICLE_URL [YYYY-MM-DD|SECONDS_FROM_EPOCH]")
		os.Exit(1)
	}

	os.MkdirAll(article.RootDir, os.ModePerm)
	fmt.Println(article.Create(params))
}
