# Tokenizer Utils

Very simple set of utilities to tokenize a string using hugging face bindings.  The underlying library for the bindings are [daulet/tokenizers](https://github.com/daulet/tokenizers/).

## Tokenizer CLI

```bash
tokenizer -h
Usage of tokenizer:
  -add_special_tokens
        Add special tokens
  -model string
        The path to the model
```

The CLI is a simple command line interface that tokenizes a string using hugging face bindings.  It reads from STDIN.

## Tokenizer Lambda

This is a Lambda function that tokenizes a string using hugging face bindings. It is meant to be fronted by an AWS Application Load Balancer.

## Tokenizer HTTP Server

This is a simple HTTP server that tokenizes a string using hugging face bindings. It is a standalone server that can be run locally or in a container.
