package main

import (
	"blog_project/internal/models"
	"log/slog"
)

<<<<<<< HEAD

import (
	"blog_project/internal/models"
	"html/template"
	"log/slog"
)

type application struct {
	debug          bool
	logger         *slog.Logger
	snippets       models.SnippetModelInterface // Use our new interface type.
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}
