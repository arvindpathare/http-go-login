package api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type response interface {
	GetResponse() string
}

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ","))
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	out := []string{}
	for word, occurance := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occurance))
	}
	return fmt.Sprintf("%s", strings.Join(out, ","))
}

func (a API) DoGetRequest(requestURL string) (response, error) {

	response, err := a.Client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("http Get error: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}
	if response.StatusCode != 200 {
		fmt.Printf("Invalid status code %d: %s\n", response.StatusCode, body)
		os.Exit(1)
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Unmarshal error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("words unmarshal error: %s", err),
			}
		}
		return words, nil
	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("occurrence unmarshal error: %s", err),
			}
		}
		return occurrence, nil
	}
	return nil, nil
}
