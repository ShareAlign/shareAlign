package main

import (
	"encoding/json"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
)

type GraphQL struct {
	Schema *graphql.Schema
}

type GraphQLReq struct {
	Query         string                 `json: query`
	OperationName string                 `json: operationName`
	Variables     map[string]interface{} `json: variables`
}

func (h GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Header().Set("Access-Control-Allow-Origin", "https://sharealign.com")

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	w.Header().Set("Access-Control-Allow-Credentials", "true")

	var req GraphQLReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
	}

	ctx := r.Context()
	res := h.Schema.Exec(ctx, req.Query, req.OperationName, req.Variables)

	json, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}

	w.Write(json)
}
