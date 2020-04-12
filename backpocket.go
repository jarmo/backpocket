package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/jarmo/backpocket/article"
	"github.com/jarmo/backpocket/config"
)

const VERSION = "1.0.0"

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}

	params, err := article.Params(os.Args)
	if err != nil {
		fmt.Println(err)
		printUsage()
		os.Exit(1)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	storageDir := config.Read().StorageDir
	os.MkdirAll(storageDir, os.ModePerm)
	fmt.Println(article.Create(storageDir, params))
}

func printUsage() {
		fmt.Println(fmt.Sprintf(`backpocket %s

USAGE:
  backpocket ARTICLE_URL [YYYY-MM-DD|SECONDS_FROM_EPOCH]
`, VERSION))
}
