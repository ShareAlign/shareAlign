package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"share-align/resolver"

	graphql "github.com/graph-gophers/graphql-go"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	file, err := os.Open("schema/schema.graphql")
	if err != nil {
		log.Fatalf("Failed to read schema: %s", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read schema: %s", err)
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	schema, err := graphql.ParseSchema(string(bytes), &resolver.QueryResolver{}, opts...)
	if err != nil {
		log.Fatal(err)
	}

	h := GraphQL{
		Schema: schema,
	}

	mux := http.NewServeMux()
	mux.Handle("/", h)

	go http.ListenAndServe(":80", http.HandlerFunc(redirect))

	s := &http.Server{
		Addr:              ":443",
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
	}

	log.Printf("Listening for requests on %s", s.Addr)

	if err = s.ListenAndServeTLS(
		"/etc/letsencrypt/live/sharealign.com/fullchain.pem",
		"/etc/letsencrypt/live/api.mazede.com/privkey.pem"); err != nil {
		log.Println("server.ListenAndServeTLS:", err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host

	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
