package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed template.html
var templateFS embed.FS

const (
	version         = "1.0.0"
	defaultTemplate = "template.html"
)

type Config struct {
	inputFile    string
	outputFile   string
	templateFile string
	showHelp     bool
	showVersion  bool
}

func main() {
	config := parseFlags()

	if config.showHelp {
		printHelp()
		return
	}

	if config.showVersion {
		fmt.Printf("MDcat version %s\n", version)
		return
	}

	if config.inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: input file is required\n\n")
		printUsage()
		os.Exit(1)
	}

	if err := convertMarkdown(config); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.inputFile, "input", "", "Input markdown file (required)")
	flag.StringVar(&config.inputFile, "i", "", "Input markdown file (short form)")
	flag.StringVar(&config.outputFile, "output", "", "Output HTML file (default: same name as input with .html extension)")
	flag.StringVar(&config.outputFile, "o", "", "Output HTML file (short form)")
	flag.StringVar(&config.templateFile, "template", defaultTemplate, "HTML template file")
	flag.StringVar(&config.templateFile, "t", defaultTemplate, "HTML template file (short form)")
	flag.BoolVar(&config.showHelp, "help", false, "Show help message")
	flag.BoolVar(&config.showHelp, "h", false, "Show help message (short form)")
	flag.BoolVar(&config.showVersion, "version", false, "Show version information")
	flag.BoolVar(&config.showVersion, "v", false, "Show version information (short form)")

	// Custom usage function
	flag.Usage = printUsage

	flag.Parse()

	// Handle positional argument for backward compatibility
	if config.inputFile == "" && flag.NArg() > 0 {
		config.inputFile = flag.Arg(0)
	}

	return config
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [markdown_file]\n\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Convert Markdown files to HTML using GitHub's API.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s README.md\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s -i README.md -o output.html\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s --input README.md --output docs/readme.html\n", filepath.Base(os.Args[0]))
}

func printHelp() {
	fmt.Printf("MDcat - Markdown to HTML converter\n")
	fmt.Printf("Version %s\n\n", version)
	printUsage()
	fmt.Printf("\nDescription:\n")
	fmt.Printf("  MDcat converts Markdown files to HTML using GitHub's Markdown API,\n")
	fmt.Printf("  which provides GitHub-flavored Markdown rendering with syntax highlighting\n")
	fmt.Printf("  and other GitHub-specific features.\n\n")
	fmt.Printf("  The output HTML file will be styled using GitHub's Primer CSS framework\n")
	fmt.Printf("  and includes a dark/light theme toggle.\n\n")
}

func convertMarkdown(config Config) error {
	// Validate input file exists
	if _, err := os.Stat(config.inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file '%s' does not exist", config.inputFile)
	}

	// Read markdown content
	content, err := ioutil.ReadFile(config.inputFile)
	if err != nil {
		return fmt.Errorf("read markdown file: %v", err)
	}

	// Convert markdown to HTML using GitHub API
	htmlContent, err := convertMarkdownToHTML(string(content))
	if err != nil {
		return fmt.Errorf("convert markdown: %v", err)
	}

	// Determine output file path
	outputPath := config.outputFile
	if outputPath == "" {
		dir := filepath.Dir(config.inputFile)
		base := strings.TrimSuffix(filepath.Base(config.inputFile), filepath.Ext(config.inputFile))
		outputPath = filepath.Join(dir, base+".html")
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %v", err)
	}

	// Read template
	var templateContent []byte

	if config.templateFile == defaultTemplate {
		// Use embedded template
		var err error
		templateContent, err = templateFS.ReadFile(defaultTemplate)
		if err != nil {
			return fmt.Errorf("read embedded template: %v", err)
		}
	} else {
		// Use custom template file
		var err error
		templateContent, err = ioutil.ReadFile(config.templateFile)
		if err != nil {
			return fmt.Errorf("read template file '%s': %v", config.templateFile, err)
		}
	} // Replace template variables
	html := strings.Replace(string(templateContent), "$MD_TITLE", filepath.Base(config.inputFile), 1)
	html = strings.Replace(html, "$MD_HTML", htmlContent, 1)

	// Write output file
	if err := ioutil.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("write output file: %v", err)
	}

	fmt.Printf("Successfully converted '%s' to '%s'\n", config.inputFile, outputPath)
	return nil
}

func convertMarkdownToHTML(markdown string) (string, error) {
	payload, err := json.Marshal(map[string]string{"text": markdown})
	if err != nil {
		return "", fmt.Errorf("marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.github.com/markdown", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %v", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return string(body), nil
}
