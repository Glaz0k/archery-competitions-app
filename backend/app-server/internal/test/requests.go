package test

import (
	"app-server/internal/server/router"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendRequest(method, host string, body io.Reader, headers map[string]string) (*http.Response, error) {
	url := "http://" + host
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return resp, nil
}

func AddCupRequest(url, title, address, season string) (string, error) {
	requestBody := map[string]string{
		"title":   title,
		"address": address,
		"season":  season,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("could not marshal request body: %w", err)
	}

	resp, err := SendRequest("POST", url+router.CreateCup, bytes.NewReader(jsonBody), map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected status OK, got %v", resp.Status)
	}

	var responseBody bytes.Buffer
	if _, err := responseBody.ReadFrom(resp.Body); err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	return responseBody.String(), nil
}

func AddCompetitionRequest(cupId int, url, stage, startDate, endDate string, isEnded bool) (string, error) {
	requestBody := map[string]interface{}{
		"cup_id":     cupId,
		"stage":      stage,
		"start_date": startDate,
		"end_date":   endDate,
		"is_ended":   isEnded,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("could not marshal request body: %w", err)
	}

	resp, err := SendRequest("POST", url+router.CreateCompetition, bytes.NewReader(jsonBody), map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return "", fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Читаем тело ответа для диагностики
		return "", fmt.Errorf("expected status OK, got %v: %s", resp.Status, string(body))
	}

	// Читаем тело ответа
	var responseBody bytes.Buffer
	if _, err := responseBody.ReadFrom(resp.Body); err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	return responseBody.String(), nil
}
