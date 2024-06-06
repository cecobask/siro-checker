package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	contentTypeForm   = "application/x-www-form-urlencoded"
	envKeyEircode     = "EIRCODE"
	urlAddressLookup  = "https://siro.ie/address-lookup-result"
	urlSearchEircode  = "https://service.siro.ie/search-eircode"
	urlValueCounty    = "data[county]"
	urlValueEircode   = "data[eircode]"
	urlValuePremiseID = "data[premiseId]"
	urlValueQuery     = "query"
	urlValueTown      = "data[town]"
	urlValueValue     = "data[value]"
)

func main() {
	eircode, ok := os.LookupEnv(envKeyEircode)
	if !ok {
		panic(fmt.Sprintf("environment variable %q is not set", envKeyEircode))
	}
	ser, err := searchEircode(eircode)
	if err != nil {
		var noSuggestionsErr noSuggestionsError
		if errors.As(err, &noSuggestionsErr) {
			fmt.Println(noSuggestionsErr)
			return
		}
		panic(err)
	}
	providers, err := addressLookup(ser)
	if err != nil {
		var notAvailableErr notAvailableError
		if errors.As(err, &notAvailableErr) {
			fmt.Println(notAvailableErr)
			return
		}
		panic(err)
	}
	fmt.Println("SIRO is available via the following internet service providers:", strings.Join(providers, ", "))
	os.Exit(1) // alert on success
}

func searchEircode(eircode string) (*searchEircodeResponse, error) {
	parsedURL, err := url.Parse(urlSearchEircode)
	if err != nil {
		return nil, err
	}
	query := make(url.Values)
	query.Set(urlValueQuery, eircode)
	parsedURL.RawQuery = query.Encode()
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var ser *searchEircodeResponse
	if err = json.NewDecoder(resp.Body).Decode(&ser); err != nil {
		return nil, err
	}
	if len(ser.Suggestions) == 0 {
		return nil, noSuggestionsError{
			eircode: eircode,
		}
	}
	return ser, nil
}

func addressLookup(ser *searchEircodeResponse) ([]string, error) {
	suggestion := ser.Suggestions[0]
	form := make(url.Values)
	form.Set(urlValueValue, suggestion.Value)
	form.Set(urlValuePremiseID, suggestion.Data.PremiseID)
	form.Set(urlValueCounty, suggestion.Data.County)
	form.Set(urlValueTown, suggestion.Data.Town)
	form.Set(urlValueEircode, suggestion.Data.Eircode)
	resp, err := http.Post(urlAddressLookup, contentTypeForm, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	selection := doc.Find("div .retailers_block > div")
	if selection.Length() == 0 {
		return nil, notAvailableError{
			eircode: suggestion.Data.Eircode,
		}
	}
	providerSet := make(map[string]struct{})
	selection.Each(func(_ int, selection *goquery.Selection) {
		provider, ok := selection.Attr("data-provider-name")
		if !ok {
			panic("could not find internet service provider attribute")
		}
		providerSet[provider] = struct{}{}
	})
	providers := make([]string, 0, len(providerSet))
	for provider := range providerSet {
		providers = append(providers, provider)
	}
	slices.Sort(providers)
	return providers, nil
}

type noSuggestionsError struct {
	eircode string
}

func (e noSuggestionsError) Error() string {
	return fmt.Sprintf("no suggestions found for eircode %s", e.eircode)
}

type notAvailableError struct {
	eircode string
}

func (e notAvailableError) Error() string {
	return fmt.Sprintf("SIRO is not yet available at %s", e.eircode)
}

type searchEircodeResponse struct {
	Query       string `json:"query"`
	Suggestions []struct {
		Value string `json:"value"`
		Data  struct {
			PremiseID string `json:"premiseId"`
			County    string `json:"county"`
			Town      string `json:"town"`
			Eircode   string `json:"eircode"`
		} `json:"data"`
	} `json:"suggestions"`
}
