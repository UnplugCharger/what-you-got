package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

// ProcessReceiptImage performs OCR on the given image using Tesseracts and returns the extracted text.
func ProcessReceiptImage(imagePath string) (string, error) {
	// Create a temporary output file name
	outputFile := "output.txt"

	// Run Tesseracts to perform OCR on the image
	_, err := exec.Command("tesseract", imagePath, outputFile).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running Tesseracts: %v", err)
	}

	// Read the extracted text from the output file
	text, err := os.ReadFile(outputFile + ".txt") // Tesseract adds .txt to the given output filename
	if err != nil {
		return "", fmt.Errorf("error reading output file: %v", err)
	}

	return string(text), nil
}

const openAIEndpoint = "https://api.openai.com/v1/engines/davinci/completions"
const apiKey = "fsshjsdhkyekwbd"

// RequestPayload defines the request structure for OpenAI's API
type RequestPayload struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

// ExtractReceiptInfo sends the raw OCR text to ChatGPT to generate a structured representation of the receipt.
func ExtractReceiptInfo(ocrText string) (map[string]interface{}, error) {
	prompt := "Extract structured receipt data from the following text:\n" + ocrText
	maxTokens := 500

	payload := RequestPayload{
		Prompt:    prompt,
		MaxTokens: maxTokens,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", openAIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(body)) // returns the API error message
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Extract the generated JSON representation from the response.
	if choices, exists := response["choices"].([]interface{}); exists && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if text, found := choice["text"].(string); found {
				var structuredReceipt map[string]interface{}
				err = json.Unmarshal([]byte(text), &structuredReceipt)
				if err != nil {
					return nil, err
				}
				return structuredReceipt, nil
			}
		}
	}

	return nil, errors.New("failed to extract structured receipt data")
}
