// Command basic demonstrates loading key-value pairs from a key-value
// provider into gofig and reading them back.
package main

import (
	"fmt"
	"log"

	"github.com/slava-y2k/gofig"
	"github.com/slava-y2k/gofig/provider/keyvalue"
)

func main() {
	src := map[string]any{
		"app": map[string]any{"name": "gofig", "debug": true},
		"db":  map[string]any{"host": "localhost", "port": 5432},
	}
	g := gofig.New()
	// Load from a key-value provider. WithDeepCopy makes Read return a copy
	// on every call, so the served map is never mutated.
	if err := g.Load(keyvalue.New(src).WithDeepCopy(true)); err != nil {
		log.Fatal(err)
	}

	// MergeMap overrides, Set and Delete keys fluently.
	g.MergeMap(map[string]any{"db": map[string]any{"port": 5433}}).
		Set([]string{"app", "mode"}, "dev").
		Delete("app", "debug")

	host, _ := g.Get("db", "host")
	port, _ := g.Get("db", "port")
	mode, _ := g.Get("app", "mode")

	fmt.Println("host:", host)
	fmt.Println("port:", port)
	fmt.Println("mode:", mode)
	fmt.Println("raw: ", g.Raw())
}
