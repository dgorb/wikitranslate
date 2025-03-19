package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgorb/wikitranslate/utils"
	"github.com/gocolly/colly"
)

type Translator struct {
	Promt          string   `json:"prompt"`
	AvailableLangs []string `json:"availableLanguages"`
}

func Translate(c *colly.Collector, inputLang, outputLang, input string) (string, error) {
	url := makeUrl(inputLang, input)
	var translation string
	availableLanguages := map[string]string{}

	c.OnHTML(`.interlanguage-link-target`, func(e *colly.HTMLElement) {
		langCode := e.Attr("lang")
		availableLanguages[langCode] = utils.LanguageByCode(langCode)
		if langCode == outputLang {
			translation = strings.Split(e.Attr("title"), " – ")[0]
		}
	})
	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	if translation == "" {
		return "", errors.New("Translation not found")
	}

	return translation, nil
}

func makeUrl(lang, text string) string {
	return fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, text)
}
