package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	r, err := search("A91C85C")
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, notFoundErr) {
			return // ignore not found error
		}
		os.Exit(1) // alert on other errors
	}
	fmt.Println(r)
	os.Exit(1) // alert on success
}

func search(eircode string) (*response, error) {
	parsedURL, err := url.Parse("https://service.siro.ie/search-eircode")
	if err != nil {
		return nil, err
	}
	q := make(url.Values)
	q.Set("query", eircode)
	parsedURL.RawQuery = q.Encode()
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r *response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if len(r.Suggestions) == 0 {
		return nil, notFoundErr
	}
	return r, nil
}

var notFoundErr = errors.New("no suggestions found")

type response struct {
	Query       string       `json:"query"`
	Suggestions []suggestion `json:"suggestions"`
}

type suggestion struct {
	Value string `json:"value"`
}

func (r *response) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("query: %s\n", r.Query))
	sb.WriteString("suggestions:\n")
	for _, s := range r.Suggestions {
		sb.WriteString(fmt.Sprintf("  - %s\n", s.Value))
	}
	return sb.String()
}
