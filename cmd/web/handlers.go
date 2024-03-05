package main

import (
	"blog_project/internal/models"
	"blog_project/internal/validator"
	"errors"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "home.tmpl", data)
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

type userSignupForm struct {
	UserName            string `form:"username"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.UserName), "username", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	exists, err := app.users.ExistsByUsername(form.UserName)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if exists {
		form.AddFieldError("username", "Username is already in use")
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

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

type userLoginForm struct {
	UserName            string `form:"username"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.UserName), "username", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.UserName, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Username or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	user, err := app.users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.User = user
	data.Form = accountPasswordUpdateForm{}

	app.render(w, r, http.StatusOK, "profile.tmpl", data)
}

type accountPasswordUpdateForm struct {
	CurrentPassword         string `form:"currentPassword"`
	NewPassword             string `form:"newPassword"`
	NewPasswordConfirmation string `form:"newPasswordConfirmation"`
	validator.Validator     `form:"-"`
}

func (app *application) userProfilePost(w http.ResponseWriter, r *http.Request) {
	var form accountPasswordUpdateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")
	if !form.Valid() {
		userID := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
		user, err := app.users.Get(userID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				app.serverError(w, r, err)
			}
			return
		}
		data := app.newTemplateData(r)
		data.User = user
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "profile.tmpl", data)
		return
	}
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.users.PasswordUpdate(userID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			userID := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
			user, err := app.users.Get(userID)
			if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				} else {
					app.serverError(w, r, err)
				}
				return
			}
			form.AddFieldError("currentPassword", "Current password is incorrect")
			data := app.newTemplateData(r)
			data.User = user
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "profile.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated!")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (app *application) admin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "admin.tmpl", data)
}

type blockForm struct {
	Action              string `form:"_method"`
	Title               string `form:"title"`
	Description         string `form:"description"`
	ImageURL1           string `form:"image_1"`
	ImageURL2           string `form:"image_2"`
	ImageURL3           string `form:"image_3"`
	ID                  int    `form:"id"`
	validator.Validator `form:"-"`
}

func (app *application) adminPost(w http.ResponseWriter, r *http.Request) {
	var form blockForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	switch form.Action {
	case "CREATE":
		err := app.blocks.Insert(form.Title, form.Description, form.ImageURL1, form.ImageURL2, form.ImageURL3)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.sessionManager.Put(r.Context(), "flash", "Block successfully created!")
	case "UPDATE":
		err := app.blocks.Update(form.ID, form.Title, form.Description, form.ImageURL1, form.ImageURL2, form.ImageURL3)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.sessionManager.Put(r.Context(), "flash", "Block successfully updated!")
	case "DELETE":
		err := app.blocks.Delete(form.ID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.sessionManager.Put(r.Context(), "flash", "Block successfully deleted!")
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
