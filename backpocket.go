package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
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

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	os.MkdirAll(article.RootDir, os.ModePerm)
	fmt.Println(article.Create(params))
}
