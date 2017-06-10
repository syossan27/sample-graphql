package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"

	"github.com/graphql-go/graphql"
	"io/ioutil"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var data map[string]user

var userType = graphql.NewObject(
	graphql.ObjectConfig {
		Name: "User",
		Fields: graphql.Fields {
			"id": &graphql.Field {
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig {
		Name: "Query",
		Fields: graphql.Fields {
			"user": &graphql.Field {
				Type: userType,
				Args: graphql.FieldConfigArgument {
					"id": &graphql.ArgumentConfig {
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := p.Args["id"].(string)
					if isOK {
						return data[idQuery], nil
					}
					return nil, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig {
		Query: queryType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params {
		Schema: schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}

	return result
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Query
	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(r.Body)
	query := bufBody.String()

	result := executeQuery(query, schema)
	json.NewEncoder(w).Encode(result)
}

func main() {
	_ = importJSONDataFromFile("data.json", &data)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func importJSONDataFromFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}

	return
}
