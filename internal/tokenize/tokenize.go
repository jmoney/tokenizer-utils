package tokenize

import (
	"context"

	"github.com/jmoney/tokenizers"
)

type ContextKey string

type TokenizerRequest struct {
	Prompt string `json:"prompt"`
}

type TokenizerResponse struct {
	TokenIds []uint32 `json:"token_ids"`
	Tokens   []string `json:"tokens"`
	Stats    Stats    `json:"stats"`
}

type Stats struct {
	Count int `json:"count"`
}

type ErrorResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Object  string `json:"object"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}

func Tokenize(ctx context.Context, request *TokenizerRequest) *TokenizerResponse {
	tk := ctx.Value(ContextKey("tokenizer")).(*tokenizers.Tokenizer)
	ids, tokens := tk.Encode(request.Prompt, true)

	tokenizerResponse := TokenizerResponse{
		TokenIds: ids,
		Tokens:   tokens,
		Stats: Stats{
			Count: len(tokens),
		},
	}

	return &tokenizerResponse
}
