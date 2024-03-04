package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//app.render(w, r, http.StatusOK, "home.tmpl", "Hello")
}

func (app *application) blocksView(w http.ResponseWriter, r *http.Request) {
	blocks, err := app.blocks.GetAll()
	if err != nil {
		app.serverError(w, r, err)
	}

	data := app.newTemplateData(r)
	data.Blocks = blocks

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
