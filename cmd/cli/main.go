package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/daulet/tokenizers"
	"github.com/jmoney/tokenizer-server/internal/tokenize"
)

var (
	elog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	model = flag.String("model", os.Getenv("MODEL"), "The path to the model")
)

func main() {
	flag.Parse()
	tk, err := tokenizers.FromFile(fmt.Sprintf("%s/tokenizer.json", *model))
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
