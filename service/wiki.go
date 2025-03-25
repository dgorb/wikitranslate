package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type wikipediaAPIResponse struct {
	Query struct {
		Pages map[string]struct {
			Extract   string `json:"extract"`
			Langlinks []struct {
				Lang  string `json:"lang"`
				Title string `json:"*"`
			} `json:"langlinks"`
		} `json:"pages"`
	} `json:"query"`
}

type WikiClient struct {
	c       *http.Client
	baseUrl string
}

func NewWikiClient() *WikiClient {
	baseUrl := "https://%s.wikipedia.org/w/api.php"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &WikiClient{
		baseUrl: baseUrl,
		c:       client,
	}
}

func (w *WikiClient) GetTranslation(ctx context.Context, inputLang, outputLang, input string) (string, error) {
	// https://www.mediawiki.org/wiki/API:Langlinks

	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("titles", input)
	params.Add("prop", "langlinks")
	params.Add("lllang", outputLang)
	params.Add("redirects", "1")

	result, err := w.makeRequest(ctx, inputLang, params)
	if err != nil {
		return "", err
	}

	for _, page := range result.Query.Pages {
		if len(page.Langlinks) > 0 {
			return page.Langlinks[0].Title, nil
		}
		log.Printf("Could not find translation for %s ([%s] -> [%s])", input, inputLang, outputLang)
		break
	}

	return "", errors.New("Translation not found")
}

func (w *WikiClient) GetSummary(ctx context.Context, lang, input string) (string, error) {
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("titles", input)
	params.Add("prop", "extracts")
	params.Add("exintro", "true")
	params.Add("explaintext", "true")

	result, err := w.makeRequest(ctx, lang, params)
	if err != nil {
		return "", err
	}

	var summary string
	for _, page := range result.Query.Pages {
		summary = strings.TrimSpace(page.Extract)
		break // Only process the first page
	}

	if summary == "" {
		return "", errors.New("Summary not found")
	}

	summaryParagraphs := strings.Split(summary, "\n")
	if len(summaryParagraphs) > 0 {
		return summaryParagraphs[0], nil
	}

	return summary, nil
}

func (w *WikiClient) makeRequest(ctx context.Context, lang string, params url.Values) (*wikipediaAPIResponse, error) {
	url := fmt.Sprintf(w.baseUrl, lang) + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := w.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non-OK status: %d", resp.StatusCode)
	}

	result := &wikipediaAPIResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}
