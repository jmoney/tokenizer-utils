package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/daulet/tokenizers"
	"github.com/jmoney/tokenizer-server/internal/tokenize"
)

var (
	elog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	model = flag.String("model", os.Getenv("MODEL"), "The path to the model")

	huggingFaceToken = os.Getenv("HF_TOKEN")
)

func main() {
	flag.Parse()

	modelPath := *model
	tokenizerPath := fmt.Sprintf("%s/tokenizer.json", *model)
	if !filepath.IsAbs(tokenizerPath) {

		home := os.Getenv("HOME")
		modelPath = strings.ToLower(modelPath)
		modelPath = fmt.Sprintf("%s/.tokenizers/%s", home, modelPath)
		tokenizerPath = fmt.Sprintf("%s/tokenizer.json", modelPath)

		if _, err := os.Stat(tokenizerPath); errors.Is(err, os.ErrNotExist) {
			client := http.DefaultClient
			huggingFacePath := fmt.Sprintf("https://huggingface.co/%s/raw/main/tokenizer.json", *model)
			req, _ := http.NewRequest("GET", huggingFacePath, nil)
			if huggingFaceToken != "" {
				req.Header.Set("Authorization", fmt.Sprintf(" Bearer %s", huggingFaceToken))
			}
			resp, err := client.Do(req)
			if err != nil {
				elog.Fatalln(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				elog.Fatalln(err)
			}
			if resp.StatusCode != 200 {
				elog.Fatal(string(body))
			}

			home := os.Getenv("HOME")
			modelPath = strings.ToLower(modelPath)
			modelPath = fmt.Sprintf("%s/.tokenizers/%s", home, modelPath)
			tokenizerPath = fmt.Sprintf("%s/tokenizer.json", modelPath)
			err = os.MkdirAll(modelPath, 0755)
			if err != nil {
				elog.Fatalln(err)
			}
			err = os.WriteFile(tokenizerPath, body, 0644)
			if err != nil {
				elog.Fatalln(err)
			}
		}
	}

	tk, err := tokenizers.FromFile(tokenizerPath)
	if err != nil {
		elog.Fatal(err)
	}
	defer tk.Close()

	stdin, err := io.ReadAll(os.Stdin)

	if err != nil {
		elog.Fatal(err)
	}
	str := string(stdin)
	prompt := strings.TrimSuffix(str, "\n")
	tokenizerRequest := tokenize.TokenizerRequest{
		Prompt: prompt,
	}
	tokenizerResponse := tokenize.Tokenize(context.WithValue(context.Background(), tokenize.ContextKey("tokenizer"), tk), &tokenizerRequest)
	resp, _ := json.Marshal(*tokenizerResponse)
	fmt.Printf("%s\n", string(resp))
}
