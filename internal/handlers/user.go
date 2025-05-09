package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"gitgub.com/tilherme/quicknotes/internal/repositories"
	"gitgub.com/tilherme/quicknotes/utils"
)

type UserHandle struct {
	repo repositories.UserRepo
}

func NewUserHandler(repo repositories.UserRepo) *UserHandle {
	return &UserHandle{repo: repo}
}
func (uh *UserHandle) SigninForm(w http.ResponseWriter, r *http.Request) error {
	return render(w, http.StatusOK, "signin.html", nil)
}
func (uh *UserHandle) Signin(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	name := r.PostFormValue("name")
	data := newUserRequest(email, password, name)
	data.Email = email
	data.Password = password
	data.Name = name

	if strings.TrimSpace(data.Password) == "" {
		data.AddFieldErrors("password", "Senha é obrigatorio")
	}

	if !isEmailValid(data.Email) {
		data.AddFieldErrors("email", "Email é invalido")
	}
	if !data.Valid() {
		render(w, http.StatusUnprocessableEntity, "signin.html", data)
		return nil
	}

	user, err := uh.repo.FindByEmail(r.Context(), data.Email)
	if err != nil {
		data.AddFieldErrors("email", "Credenciais invalidas")
		return render(w, http.StatusUnprocessableEntity, "signin.html", data)
	}
	if !utils.ValidatePassword(data.Password, user.Password.String) {
		data.AddFieldErrors("email", "Credenciais invalidas")
		return render(w, http.StatusUnprocessableEntity, "signin.html", data)
	}
	if !user.Active.Bool {
		data.AddFieldErrors("email", "Email não confirmado")
		return render(w, http.StatusUnprocessableEntity, "signin.html", data)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
func (uh *UserHandle) SignupForm(w http.ResponseWriter, r *http.Request) error {
	return render(w, http.StatusOK, "signup.html", nil)
}

func (uh *UserHandle) Signup(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	name := r.PostFormValue("name")
	data := newUserRequest(email, password, name)
	data.Email = email
	data.Password = password
	data.Name = name

	if strings.TrimSpace(data.Password) == "" {
		data.AddFieldErrors("password", "Senha é obrigatorio")
	}
	if len(strings.TrimSpace(data.Password)) < 6 {
		data.AddFieldErrors("password", "Senha precisa ter 6 digitos ou mais")
	}
	if !isEmailValid(data.Email) || strings.TrimSpace(data.Email) == "" {
		data.AddFieldErrors("email", "Email é invalido")
	}
	if !data.Valid() {
		render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return nil
	}
	hash, err := utils.GenerateFromPassword(data.Password)
	if err != nil {
		return err
	}
	hashToken := utils.GenerateToken()
	user, token, err := uh.repo.Create(r.Context(), data.Email, hash, data.Name, hashToken)
	if err == repositories.ErrDuplicateEmail {
		data.AddFieldErrors("email", "Email, ja esta em uso")
		return render(w, http.StatusUnprocessableEntity, "signup.html", data)

	}
	if err != nil {
		return err
	}
	fmt.Println(token)
	fmt.Println(user.Email)
	return render(w, http.StatusOK, "signup-sucess.html", token)
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func (uh *UserHandle) Confirm(w http.ResponseWriter, r *http.Request) error {
	token := r.PathValue("token")
	err := uh.repo.ConfirmUserByToken(r.Context(), token)
	msg := "Cadastro realizado com sucesso"
	if err != nil {
		msg = "Token invalido ou cadastro já confirmado"
	}
	return render(w, http.StatusOK, "user-confirm.html", msg)
}
