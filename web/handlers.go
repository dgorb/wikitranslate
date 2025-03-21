package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgorb/wikitranslate/service"
	"github.com/dgorb/wikitranslate/utils"
	"github.com/gocolly/colly"
)

func GetLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.Languages())
}

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	inputLang := r.URL.Query().Get("inputLang")
	outputLang := r.URL.Query().Get("outputLang")
	input := r.URL.Query().Get("input")

	if inputLang == "" || outputLang == "" || input == "" {
		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

	if !utils.ValidLanguage(inputLang) {
		http.Error(w, fmt.Sprintf("%s is an invalid language", inputLang), http.StatusBadRequest)
		return
	}

	if !utils.ValidLanguage(outputLang) {
		http.Error(w, fmt.Sprintf("%s is an invalid language", outputLang), http.StatusBadRequest)
		return
	}

	c := colly.NewCollector()
	translation, err := service.Translate(c, inputLang, outputLang, input)
	if err != nil {
		log.Printf("Error translating: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(translation)
}
