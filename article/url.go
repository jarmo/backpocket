package article

import (
	"errors"
	"fmt"
	"net/url"
)

func URL(args []string) (*url.URL, error) {
	if len(args) < 2 {
		return nil, errors.New("Not enough arguments")
	}

	rawUrl := args[1]
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse URL %v", err))
	}

	if uri.Scheme == "" {
		return nil, errors.New(fmt.Sprintf("URL in unsupported format %v", rawUrl))
	}

	return uri, nil
}
