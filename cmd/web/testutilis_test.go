package main

import (
	"blog_project/internal/models"
	"blog_project/internal/models/mocks"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"io"
	"log/slog"
	"testing"
	"time"
)

func newTestApplication(t *testing.T) *application {
	// Create an instance of the template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	// And a form decoder.
	formDecoder := form.NewDecoder()
	// And a session manager instance. Note that we use the same settings as
	// production, except that we *don't* set a Store for the session manager.
	// If no store is set, the SCS package will default to using a transient
	// in-memory store, which is ideal for testing purposes.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	return &application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		blocks:         &models.BlockModel{Collection: blocksCollection},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

}
