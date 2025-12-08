package vdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Embeddings(input string, model string) ([][]float64, error) {
	data, err := json.Marshal(EmbedRequest{
		Model: model,
		Input: input,
	})
	if err != nil {
		return nil, err
	}

	sdata := bytes.NewReader(data)
	resp, err := http.Post("http://127.0.0.1:11434/api/embed", "application/json", sdata)
	if err != nil {
		return nil, err
	}

	var bb bytes.Buffer
	io.Copy(&bb, resp.Body)
	resp.Body.Close()

	var m *EmbedResponse
	err = json.Unmarshal(bb.Bytes(), &m)
	if err != nil {
		return nil, err
	}

	if len(m.Embeddings) == 0 {
		return nil, fmt.Errorf("odd embed response: %v", bb.String())
	}

	return m.Embeddings, nil
}

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
}
