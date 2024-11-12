package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

func init() {
	log.SetFlags(0)
}

func main() {
	priv, pub, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatalln(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(map[string]string{
		"publicKey":  pub,
		"privateKey": priv,
	})
}
