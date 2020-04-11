package article

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type ArticleParams struct {
	Url        *url.URL
	ArchivedAt time.Time
}

func Params(args []string) (ArticleParams, error) {
	if len(args) < 2 {
		return ArticleParams{}, errors.New("Not enough arguments")
	}

	rawUrl := args[1]
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return ArticleParams{}, errors.New(fmt.Sprintf("Failed to parse URL %v", err))
	}

	if uri.Scheme != "http" && uri.Scheme != "https" {
		return ArticleParams{}, errors.New(fmt.Sprintf("URL in unsupported format %v", rawUrl))
	}

	if len(args) > 2 {
		archivedAt, err := time.Parse("2006-01-02", args[2])
		if err != nil {
			if epoch, err := strconv.ParseInt(args[2], 10, 64); err != nil {
				return ArticleParams{}, err
			} else {
				archivedAt := time.Unix(epoch, 0)
				return ArticleParams{Url: uri, ArchivedAt: archivedAt}, nil
			}
		}
		return ArticleParams{Url: uri, ArchivedAt: archivedAt}, nil
	} else {
		return ArticleParams{Url: uri, ArchivedAt: time.Now()}, nil
	}
}
