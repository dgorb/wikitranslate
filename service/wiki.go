package service

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/dgorb/wikitranslate/utils"
	"github.com/gocolly/colly"
)

func GetSummary(c *colly.Collector, lang, input string) (string, error) {
	url := makeUrl(lang, input)

	var summary string
	var foundFirstParagraph bool
	c.OnHTML(".mw-parser-output p", func(e *colly.HTMLElement) {
		if !foundFirstParagraph {
			summary = strings.TrimSpace(e.Text)
			if summary != "" {
				re := regexp.MustCompile(`\[.*?\]`)
				summary = re.ReplaceAllString(summary, "")
				foundFirstParagraph = true
			}
		}
	})
	err := c.Visit(url)
	if err != nil {
		log.Printf("Error: %s", err)
		return "", err
	}
	return summary, nil
}

func GetTranslation(c *colly.Collector, inputLang, outputLang, input string) (string, error) {
	url := makeUrl(inputLang, input)
	availableLanguages := map[string]string{}

	var translation string
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
