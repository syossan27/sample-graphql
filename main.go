package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"bytes"

	"github.com/graphql-go/graphql"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Schema
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(r.Body)
	query := bufBody.String()

	params := graphql.Params{Schema: schema, RequestString: query}
	d := graphql.Do(params)
	if len(d.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %v", d.Errors)
	}
	rJSON, _ := json.Marshal(d)
	fmt.Printf("%s \n", rJSON)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
