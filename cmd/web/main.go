package main
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
=======
>>>>>>> 3b57f93b198708c114101c90251b028fa1eb643b
