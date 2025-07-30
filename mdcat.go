package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <markdown_file>", filepath.Base(os.Args[0]))
	}
	mdPath := os.Args[1]
	content, err := ioutil.ReadFile(mdPath)
	if err != nil {
		log.Fatalf("read markdown: %v", err)
	}

	payload, err := json.Marshal(map[string]string{"text": string(content)})
	if err != nil {
		log.Fatalf("marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.github.com/markdown", bytes.NewReader(payload))
	if err != nil {
		log.Fatalf("create request: %v", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("github api error: %s", strings.TrimSpace(string(body)))
	}

	templateContent, err := ioutil.ReadFile("template.html")
	if err != nil {
		log.Fatalf("read template: %v", err)
	}

	html := strings.Replace(string(templateContent), "$MD_TITLE", filepath.Base(mdPath), 1)
	html = strings.Replace(html, "$MD_HTML", string(body), 1)

	outPath := filepath.Join(filepath.Dir(mdPath), strings.TrimSuffix(filepath.Base(mdPath), filepath.Ext(mdPath))+".html")
	if err := ioutil.WriteFile(outPath, []byte(html), 0644); err != nil {
		log.Fatalf("write output: %v", err)
	}
}
