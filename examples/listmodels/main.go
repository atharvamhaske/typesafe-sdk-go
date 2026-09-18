// Command listmodels prints the available TypeSafe models.
package main

import (
	"context"
	"fmt"
	"log"

	typesafe "github.com/atharvamhaske/typesafe-sdk-go"
)

func main() {
	client, err := typesafe.NewClient() // reads TYPESAFE_API_KEY
	if err != nil {
		log.Fatal(err)
	}

	models, err := client.ListModels(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range models {
		fmt.Printf("%s\t%s\t%s\n", m.Name, m.ReleaseDate, m.Description)
	}
}
