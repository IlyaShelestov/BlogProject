package main

import (
	"blog_project/internal/models"
	"blog_project/internal/validator"
	"errors"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "about.tmpl")
}

func (app *application) blocksView(w http.ResponseWriter, r *http.Request) {
	blocks, err := app.blocks.GetAll()
	if err != nil {
		app.serverError(w, r, err)
	}

	data := app.newTemplateData(r)
	data.Blocks = blocks

	app.render(w, r, http.StatusOK, "blog.tmpl", data)
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	// Render the sign-up form template
	app.render(w, r, http.StatusOK, "signin.tmpl")
}

func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form userSignupForm
	form.UserName = r.PostForm.Get("username")
	form.Password = r.PostForm.Get("password")

	form.CheckField(validator.NotBlank(form.UserName), "username", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.UserName, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUsername) {
			form.AddFieldError("username", "Username address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "login.tmpl")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "login.tmpl")
}

type userSignupForm struct {
	UserName            string `form:"name"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
type userLoginForm struct {
	UserName            string `form:"name"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
