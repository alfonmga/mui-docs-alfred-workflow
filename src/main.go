package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	aw "github.com/deanishe/awgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	// Algolia application ID
	AppID = "TZGZ85B9TB"
	// Algolia search API key
	SearchAPIKey = "8177dfb3e2be72b241ffb8c5abafa899"
	// Algolia index name
	IndexName = "material-ui"
)

var wf *aw.Workflow
var searchIdx *search.Index

type SearchResultHit struct {
	URL             string                 `json:"url"`
	HierarchyLevels map[string]interface{} `json:"hierarchy"`
	Product         string                 `json:"product"`
}
type AlfredSearchResultItem struct {
	Title string
	URL   string
}

func init() {
	searchIdx = getSearchIndex()
}

func main() {
	wf = aw.New()
	wf.Run(run)
}
func run() {
	res, err := searchIdx.Search(wf.Args()[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hits, err := unmarshalSearchResHits(res)
	if err != nil {
		panic(err)
	}
	for _, hit := range hits {
		alfredResultItemData := generateAlfredResultItemFromSearchResHit(hit)
		wf.NewItem(alfredResultItemData.Title).
			Arg(alfredResultItemData.URL).
			Subtitle(alfredResultItemData.URL).
			Arg(alfredResultItemData.URL).
			UID(alfredResultItemData.URL).
			Valid(true)
	}
	wf.SendFeedback()
}

func getSearchIndex() *search.Index {
	return search.NewClient(AppID, SearchAPIKey).InitIndex(IndexName)
}
func unmarshalSearchResHits(searchRes search.QueryRes) ([]SearchResultHit, error) {
	var hits []SearchResultHit
	err := searchRes.UnmarshalHits(&hits)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal Algolia search result hits: %s", err)
	}
	return hits, nil
}
func generateAlfredResultItemFromSearchResHit(hit SearchResultHit) *AlfredSearchResultItem {
	title := ""
	count := 0 // count of hierarchy levels
	for _, v := range hit.HierarchyLevels {
		if v == nil {
			continue
		}
		if count > 0 {
			title = fmt.Sprintf("%s > %s", title, v)
		} else {
			title = fmt.Sprintf("%s", v)
		}
		count++
	}
	title = fmt.Sprintf("%s [%s]", title, getProductNameFromSearchResHit(hit.Product))
	return &AlfredSearchResultItem{
		Title: title,
		URL:   hit.URL,
	}
}
func getProductNameFromSearchResHit(s string) string {
	// Uppercase first letter of each word
	productName := cases.Title(language.Und, cases.NoLower).String(s)
	// Remove dashes
	productName = strings.Replace(productName, "-", " ", -1)
	// Replace "Ui" with "UI"
	productName = strings.Replace(productName, "Ui", "UI", -1)
	return productName
}
