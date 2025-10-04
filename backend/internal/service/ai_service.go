package service

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type AIService interface {
    GenerateSummary(content string) (string, error)
}

type aiService struct {
    apiKey string
}

func NewAIService(apiKey string) AIService {
    return &aiService{apiKey: apiKey}
}

func (s *aiService) GenerateSummary(content string) (string, error) {
    if s.apiKey == "" {
        return "AI summary not available.", nil
    }
    
    url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", s.apiKey)
    
    payload := map[string]interface{}{
        "contents": []map[string]interface{}{
            {
                "parts": []map[string]interface{}{
                    {
                        "text": fmt.Sprintf("Write a short, clear, and concise summary highlighting the main message of the following text:\n\n%s", content),
                    },
                },
            },
        },
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %v", err)
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %v", err)
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %v", err)
    }
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return "", fmt.Errorf("failed to parse response: %v", err)
    }
    
    candidates, ok := result["candidates"].([]interface{})
    if !ok || len(candidates) == 0 {
        return "Unable to generate summary. Please try again.", nil
    }
    
    candidate := candidates[0].(map[string]interface{})
    contentData, ok := candidate["content"].(map[string]interface{})
    if !ok {
        return "Unable to extract summary content.", nil
    }
    
    parts, ok := contentData["parts"].([]interface{})
    if !ok || len(parts) == 0 {
        return "Unable to extract summary parts.", nil
    }
    
    part := parts[0].(map[string]interface{})
    summary, ok := part["text"].(string)
    if !ok {
        return "Unable to extract summary text.", nil
    }
    
    return summary, nil
}