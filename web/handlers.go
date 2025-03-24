package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgorb/wikitranslate/service"
	"github.com/dgorb/wikitranslate/utils"
)

type TranslationResponse struct {
	Translation string `json:"translation"`
	Summary     string `json:"summary"`
}

func GetLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.Languages())
}

func TranslatePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	inputLang := r.URL.Query().Get("inputLang")
	outputLang := r.URL.Query().Get("outputLang")
	input := strings.TrimSpace(r.URL.Query().Get("input"))

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

	wiki := service.NewWiki()
	translation, err := wiki.GetTranslation(inputLang, outputLang, input)
	if err != nil {
		log.Printf("Error translating: %s", err)
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	summary, err := wiki.GetSummary(outputLang, translation)
	if err != nil {
		log.Printf("Error getting summary: %s", err)
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TranslationResponse{
		Translation: translation,
		Summary:     summary,
	})
}
