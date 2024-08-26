package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"unicode/utf8"
)

// Ошибка и исправление из ответа сервиса
type Error struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

// CheckSpelling отправляет текст на проверку орфографии и возвращает исправленный текст
func (ll *LogicLayer) CheckSpelling(text string) (string, error) {
	urlSpeller := "https://speller.yandex.net/services/spellservice.json/checkText"

	formData := url.Values{}
	formData.Set("text", text)

	resp, err := http.PostForm(urlSpeller, formData)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var errors []Error
	if err := json.Unmarshal(body, &errors); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	correctedText := text
	corr := 0
	for _, err := range errors {
		correctedText, corr = replaceText(correctedText, err, corr)
	}
	return correctedText, nil
}

// replaceText заменяет ошибочные слова на исправленные
func replaceText(text string, err Error, corr int) (string, int) {
	if len(err.S) == 0 {
		return text, corr
	}
	correctedWord := err.S[0]
	runedText := []rune(text)

	startPos := err.Pos + corr
	endPos := startPos + err.Len

	corr = corr + (utf8.RuneCountInString(correctedWord) - err.Len)

	return string(runedText[:startPos]) + correctedWord + string(runedText[endPos:]), corr
}
