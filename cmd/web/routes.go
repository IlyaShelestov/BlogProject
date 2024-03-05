package main

import (
	"net/http"

	"blog_project/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/signup", dynamic.ThenFunc(app.userSignUp))
	router.Handler(http.MethodPost, "/signup", dynamic.ThenFunc(app.userSignUpPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/blog", protected.ThenFunc(app.blocksView))
	router.Handler(http.MethodPost, "/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/profile", protected.ThenFunc(app.userProfile))
	router.Handler(http.MethodPost, "/profile", protected.ThenFunc(app.userProfilePost))

	admin := protected.Append(app.requireAdmin)

	router.Handler(http.MethodGet, "/admin", admin.ThenFunc(app.admin))
	router.Handler(http.MethodPost, "/admin", admin.ThenFunc(app.adminPost))

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standart.Then(router)
}
