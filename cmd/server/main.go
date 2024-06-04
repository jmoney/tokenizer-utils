package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jmoney/tokenizer-server/internal/tokenize"
	"github.com/jmoney/tokenizers"
)

var (
	elog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	wlog = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	ilog = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)

	model = flag.String("model", os.Getenv("MODEL"), "The path to the model")
	port  = flag.String("port", os.Getenv("PORT"), "The port to start the http server on")
)

func main() {
	flag.Parse()
	tk, err := tokenizers.FromFile(fmt.Sprintf("%s/tokenizer.json", *model))
	if err != nil {
		elog.Fatal(err)
	}
	defer tk.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequset)
	log.Printf("Start server on port: %s\n", *port)
	contextedMux := addTokenizerToContext(tk, mux)
	http.ListenAndServe(fmt.Sprintf(":%s", *port), contextedMux)
}

func addTokenizerToContext(tk *tokenizers.Tokenizer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), tokenize.ContextKey("tokenizer"), tk)))
	})
}

func handleRequset(w http.ResponseWriter, req *http.Request) {
	requestID := req.Header.Get("X-Request-ID")
	ilog.Printf("Request ID: %s\n", requestID)

	if req.Method != "POST" {
		wlog.Printf("requestID=%s, Method %s not allowed\n", requestID, req.Method)
		errorResponse := tokenize.ErrorResponse{
			ID:      requestID,
			Message: "Http Method Not Allowd",
			Object:  "error",
			Type:    "invalid_request",
			Code:    405,
		}

		resp, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		w.WriteHeader(405)
	}

	tokenizeRequest := tokenize.TokenizerRequest{}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		elog.Printf("requestID=%s, %s\n", requestID, err)
		errorResponse := tokenize.ErrorResponse{
			ID:      requestID,
			Message: "Internal Server Error",
			Object:  "error",
			Type:    "internal_error",
			Code:    500,
		}

		resp, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		w.WriteHeader(500)
	}
	defer req.Body.Close()
	err = json.Unmarshal(body, &tokenizeRequest)
	if err != nil {
		elog.Printf("requestID=%s, %s\n", requestID, err)
		errorResponse := tokenize.ErrorResponse{
			ID:      requestID,
			Message: "Internal Server Error",
			Object:  "error",
			Type:    "internal_error",
			Code:    500,
		}

		resp, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		w.WriteHeader(500)
	}

	tokenizeResposne := tokenize.Tokenize(req.Context(), &tokenizeRequest)
	resp, _ := json.Marshal(tokenizeResposne)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	w.WriteHeader(200)
}
