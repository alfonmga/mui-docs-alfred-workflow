package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	searchRes, err := searchIdx.Search("Button")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, searchRes)
}

func TestUnmarshalSearchResHits(t *testing.T) {
	searchRes, err := searchIdx.Search("guidelines")
	if err != nil {
		t.Fatal(err)
	}
	var hits []SearchResultHit
	err = searchRes.UnmarshalHits(&hits)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, hits)
	assert.EqualValues(t, "https://mui.com/material-ui/react-button/#main-content", hits[0].URL)
}

func TestGenerateAlfredResultItemFromSearchResHit(t *testing.T) {
	searchRes, err := searchIdx.Search("Button")
	if err != nil {
		t.Fatal(err)
	}
	var hits []SearchResultHit
	err = searchRes.UnmarshalHits(&hits)
	if err != nil {
		t.Fatal(err)
	}
	alfredResultItemData := generateAlfredResultItemFromSearchResHit(hits[0])
	assert.EqualValues(t, "https://mui.com/material-ui/react-button/#main-content", alfredResultItemData.URL)
	assert.EqualValues(t, "Components > Button [Material UI]", alfredResultItemData.Title)
}
