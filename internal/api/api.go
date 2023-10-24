package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const ApiUrl string = "https://api.mentor.duden.de/api/grammarcheck"

type SpellAdvice struct {
	ErrorCode     string   `json:"errorCode"`
	ErrorMessage  string   `json:"errorMessage"`
	ShortMessage  string   `json:"shortMessage"`
	Length        int      `json:"length"`
	Offset        int      `json:"offset"`
	OriginalError string   `json:"originalError"`
	Proposals     []string `json:"proposals"`
	Synonyms      []string `json:"synonyms"`
}

type responseData struct {
	SpellAdvices []SpellAdvice
}

type response struct {
	Data responseData
}

type payload struct {
	Text string `json:"text"`
}

func FetchErrors(text string) ([]SpellAdvice, error) {
	payload := payload{
		Text: text,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encoding payload: %w", err)
	}

	resp, err := http.Post(ApiUrl, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("fetching api: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response stream: %w", err)
	}

	var respJson response
	err = json.Unmarshal(respData, &respJson)
	if err != nil {
		return nil, fmt.Errorf("parsing JSON response: %w", err)
	}

	return respJson.Data.SpellAdvices, nil
}
