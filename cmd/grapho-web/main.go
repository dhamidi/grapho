package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dhamidi/grapho"
)

func draftPost(app *grapho.Grapho, w http.ResponseWriter, req *http.Request) error {
	id, err := grapho.NewPostId(req.FormValue("post_id"))
	if err != nil {
		return err
	}

	command := &grapho.DraftPostCommand{
		PostId: id,
		Now:    time.Now(),
		Title:  req.FormValue("title"),
		Body:   req.FormValue("body"),
	}

	if err := app.DraftPost(command); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func listDrafts(app *grapho.Grapho, w http.ResponseWriter, req *http.Request) error {
	enc := json.NewEncoder(w)
	drafts, err := app.ListDrafts()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return enc.Encode(drafts)
}

func respondWithError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "%s\n", err.Error())
}

func main() {
	ErrMethodNotSupported := fmt.Errorf("Method not supported")
	app := grapho.NewGrapho(os.Getenv("GRAPHO_ENV"))
	http.HandleFunc("/drafts", func(w http.ResponseWriter, req *http.Request) {
		err := (error)(nil)

		switch req.Method {
		case "POST":
			err = draftPost(app, w, req)
		case "GET":
			err = listDrafts(app, w, req)
		default:
			err = ErrMethodNotSupported
		}

		if err != nil {
			respondWithError(w, err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
