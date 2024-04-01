package main

import (
	"fmt"
	"strings"

	"github.com/dgorb/wikitranslate/utils"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	fmt.Println(translate(c, "en", "uk", "Guitar"))
}

type Translator struct {
	Promt          string   `json:"prompt"`
	AvailableLangs []string `json:"availableLanguages"`
}

func makeUrl(lang, text string) string {
	return fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, text)
}

func translate(c *colly.Collector, inputLang, outputLang, input string) string {
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

	c.Visit(url)

	fmt.Println(availableLanguages)

	return translation
}
