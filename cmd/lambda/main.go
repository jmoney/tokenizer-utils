package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/daulet/tokenizers"
	"github.com/google/uuid"
	"github.com/jmoney/tokenizer-server/internal/tokenize"
)

var (
	elog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	wlog = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	ilog = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)

	model = os.Getenv("MODEL")
)

func main() {
	tk, err := tokenizers.FromFile(fmt.Sprintf("%s/tokenizer.json", model))
	if err != nil {
		elog.Fatal(err)
	}

	lambda.StartWithOptions(handleRequest, lambda.WithEnableSIGTERM(func() {
		ilog.Println("SIGTERM received, shutting down")
		tk.Close()
	}), lambda.WithContext(context.WithValue(context.Background(), tokenize.ContextKey("tokenizer"), tk)))
}

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {

	requestID, err := uuid.Parse(request.Headers["X-Request-ID"])
	if err != nil {
		id := uuid.New()
		wlog.Printf("Error parsing request ID \"%s\".  Using \"%s\"\n", request.Headers["X-Request-ID"], id.String())
		requestID = id
	} else {
		ilog.Printf("Request ID: %s\n", requestID.String())
	}

	tokenizerRequest := tokenize.TokenizerRequest{}
	err = json.Unmarshal([]byte(request.Body), &tokenizerRequest)
	if err != nil {
		wlog.Printf("Error unmarshalling request body: %s\n", err)
		errorResponse := tokenize.ErrorResponse{
			ID:      requestID.String(),
			Message: "Error unmarshalling request body",
			Object:  "error",
			Type:    "invalid_request",
			Code:    400,
		}
		resp, _ := json.Marshal(errorResponse)
		return events.ALBTargetGroupResponse{
			Body:              string(resp),
			StatusCode:        400,
			StatusDescription: "400 Bad Request",
			IsBase64Encoded:   false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	tokenizerResponse := tokenize.Tokenize(ctx, &tokenizerRequest)

	resp, _ := json.Marshal(*tokenizerResponse)

	return events.ALBTargetGroupResponse{
		StatusCode: 200,
		Body:       string(resp),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
