package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dgorb/wikitranslate/utils"
	"github.com/gocolly/colly"
)

type Result struct {
	Translation string `json:"translation"`
	Summary     string `json:"summary"`
}

func Translate(c *colly.Collector, inputLang, outputLang, input string) (*Result, error) {
	var result Result
	url := makeUrl(inputLang, input)
	availableLanguages := map[string]string{}

	c.OnHTML(`.interlanguage-link-target`, func(e *colly.HTMLElement) {
		langCode := e.Attr("lang")
		availableLanguages[langCode] = utils.LanguageByCode(langCode)
		if langCode == outputLang {
			result.Translation = strings.Split(e.Attr("title"), " – ")[0]
		}
	})

	var foundFirstParagraph bool
	c.OnHTML(".mw-parser-output p", func(e *colly.HTMLElement) {
		if !foundFirstParagraph {
			summary := strings.TrimSpace(e.Text)
			if summary != "" {
				re := regexp.MustCompile(`\[.*?\]`)
				result.Summary = re.ReplaceAllString(summary, "")
				foundFirstParagraph = true
			}
		}
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	if result.Translation == "" {
		return nil, errors.New("Translation not found")
	}

	return &result, nil
}

func makeUrl(lang, text string) string {
	return fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, text)
}
