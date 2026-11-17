package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

func init() {
	log.SetFlags(0)
}

func main() {
	subscriber := flag.String("subscriber", "", "VAPID subscriber contact URI (mailto:... or https://...)")
	flag.Parse()

	priv, pub, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatalln(err)
	}

	out := map[string]string{
		"publicKey":  pub,
		"privateKey": priv,
	}
	if *subscriber != "" {
		out["subscriber"] = *subscriber
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		log.Fatalln(err)
	}
}
