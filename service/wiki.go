package service

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/dgorb/wikitranslate/utils"
	"github.com/gocolly/colly"
)

func initColly() *colly.Collector {
	return colly.NewCollector(colly.CacheDir("./cache"))
}

type Wiki struct {
	collector *colly.Collector
}

func NewWiki() *Wiki {
	c := initColly()
	return &Wiki{
		collector: c,
	}
}

func (w *Wiki) GetSummary(lang, input string) (string, error) {
	url := makeUrl(lang, input)

	var summary string
	var foundFirstParagraph bool
	w.collector.OnHTML(".mw-parser-output p", func(e *colly.HTMLElement) {
		if !foundFirstParagraph {
			summary = strings.TrimSpace(e.Text)
			if summary != "" {
				re := regexp.MustCompile(`\[.*?\]`)
				summary = re.ReplaceAllString(summary, "")
				foundFirstParagraph = true
			}
		}
	})
	err := w.collector.Visit(url)
	if err != nil {
		log.Printf("Error: %s", err)
		return "", err
	}
	return summary, nil
}

func (w *Wiki) GetTranslation(inputLang, outputLang, input string) (string, error) {
	url := makeUrl(inputLang, input)
	availableLanguages := map[string]string{}

	var translation string
	w.collector.OnHTML(`.interlanguage-link-target`, func(e *colly.HTMLElement) {
		langCode := e.Attr("lang")
		availableLanguages[langCode] = utils.LanguageByCode(langCode)
		if langCode == outputLang {
			translation = strings.Split(e.Attr("title"), " – ")[0]
		}
	})
	err := w.collector.Visit(url)
	if err != nil {
		return "", err
	}
	if translation == "" {
		return "", errors.New("Translation not found")
	}

	return translation, nil
}

func makeUrl(lang, text string) string {
	encodedText := url.QueryEscape(text)
	return fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, encodedText)
}
