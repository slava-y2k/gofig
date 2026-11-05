// Command basic demonstrates loading, modifying, and reading configuration.
package main

import (
	"fmt"
	"log"

	"github.com/slava-y2k/gonfig"
	"github.com/slava-y2k/gonfig/provider/keyvalue"
)

func main() {
	src := map[string]any{
		"app": map[string]any{"name": "gonfig", "debug": true},
		"db":  map[string]any{"host": "localhost", "port": 5432},
	}
	var g *gonfig.Gonfig
	var err error
	g, err = gonfig.NewWithProvider(keyvalue.New(src).WithDeepCopy(true))
	if err != nil {
		log.Fatal(err)
	}

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
