package tokenize

import (
	"context"

	"github.com/daulet/tokenizers"
)

type ContextKey string

type TokenizerRequest struct {
	Prompt string `json:"prompt"`
}

type TokenizerResponse struct {
	Tokens []uint32 `json:"tokens"`
	Count  int      `json:"count"`
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
	ids, _ := tk.Encode(request.Prompt, true)

	tokenizerResponse := TokenizerResponse{
		Tokens: ids,
		Count:  len(ids),
	}

	return &tokenizerResponse
}
