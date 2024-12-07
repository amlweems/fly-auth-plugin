package main

import (
	"context"
	"encoding/json"
	"os"
)

func main() {
	aud := os.Getenv("GOOGLE_EXTERNAL_ACCOUNT_AUDIENCE")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(NewProvider().Token(context.Background(), aud))
}
