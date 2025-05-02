package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shoggothforever/beefine/config"
	"strings"
	"text/template"
)

type DSResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
		PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
		PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func Analysis(userContent string) string {
	tp := template.New("dpTemplate")
	var buf bytes.Buffer
	tp.Parse(config.Mcfg.Template)
	cfg := config.Mcfg
	ss := strings.Split(userContent, "\n")
	cfg.UserContent = strings.Join(ss, " ")
	err := tp.Execute(&buf, cfg)
	if err != nil {
		fmt.Println(err)
	}
	payload := strings.NewReader(buf.String())
	client := &http.Client{}
	url := "https://api.deepseek.com/chat/completions"
	method := "POST"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+cfg.Key)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	var rsp DSResponse
	json.Unmarshal(body, &rsp)
	if len(rsp.Choices) > 0 {
		fmt.Println("feedback content: ", rsp.Choices[0].Message.Content)
		return rsp.Choices[0].Message.Content
	}
	return ""
}
