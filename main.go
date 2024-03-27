package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	inputLang := "en"
	outputLang := "ru"

	input := "Guitar"

	c.OnHTML(`.interlanguage-link-target`, func(e *colly.HTMLElement) {
		if e.Attr("lang") == outputLang {
			fmt.Printf("%s -> %s\n", inputLang, outputLang)
			fmt.Printf("%s -> %s\n", input, strings.Split(e.Attr("title"), " – ")[0])
			fmt.Println(e.Attr("href"))

			// If we find translation, visit the link and pull
			// up the first paragraph
			// fmt.Println(e.ChildText("p"))
		}

	})

	c.Visit(fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", inputLang, input))
}
