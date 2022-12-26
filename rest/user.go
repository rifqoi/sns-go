package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	snsgo "github.com/rifqoi/sns-go"
	"github.com/rifqoi/sns-go/rest/responses"
)

type UserController struct {
	userSvc snsgo.UserService
}

func NewUserController(userSvc snsgo.UserService) *UserController {
	return &UserController{
		userSvc: userSvc,
	}
}

func (u *UserController) Register(r *chi.Mux) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/signup", u.RegisterUser)
		r.Get("/", u.ListAllUsers)
		r.Get("/email/{email}", u.GetUserByEmail)
		r.Get("/uuid/{uuid}", u.GetUserByID)
	})
}

type CreateUserRequest struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Age         int    `json:"age"`
	Gender      Gender `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

func (u *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := render.Decode(r, &req)
	if err != nil {
		responses.ErrorResponse(w, r, "invalid request", snsgo.WrapErrorf(err, snsgo.ErrorCodeInvalidArg, "json decoder"))
		return
	}

	err = u.userSvc.RegisterUser(r.Context(), snsgo.CreateRequest{
		Username:    req.Username,
		Email:       req.Email,
		Name:        req.Name,
		Password:    req.Password,
		Gender:      req.Gender.Convert(),
		Age:         req.Age,
		DateOfBirth: req.DateOfBirth,
	})

	if err != nil {
		responses.ErrorResponse(w, r, "failed to add user", err)
		return
	}

	responses.SuccessResponse(w, r, http.StatusCreated, "user created succesfully", nil)
}

type GetUserResponse struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Username    string    `json:"username,omitempty"`
	Name        string    `json:"name,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	DateOfBirth string    `json:"date_of_birth,omitempty"`
	Age         int       `json:"age,omitempty"`
	Email       string    `json:"email,omitempty"`
}

func (u *UserController) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.userSvc.GetUsers(r.Context())
	if err != nil {
		responses.ErrorResponse(w, r, "cannot get users", err)
		return
	}

	var listUsers []GetUserResponse
	for _, user := range users {
		listUsers = append(listUsers, GetUserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			Gender:      user.Gender.String(),
			DateOfBirth: user.DateOfBirth.String(),
			Age:         user.Age,
			Email:       user.Email,
		})
	}

	responses.SuccessResponse(w, r, http.StatusOK, "get users succesfully", listUsers)
}

func (u *UserController) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if len(email) < 1 {
		responses.ErrorResponse(w, r, "invalid request", snsgo.NewErrorf(snsgo.ErrorCodeInvalidArg, "invalid request: need email"))
		return
	}

	user, err := u.userSvc.GetUserByEmail(r.Context(), email)
	if err != nil {
		responses.ErrorResponse(w, r, "cannot get user", err)
		return
	}

	resp := GetUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Gender:      user.Gender.String(),
		DateOfBirth: user.DateOfBirth.String(),
		Age:         user.Age,
		Email:       user.Email,
	}

	responses.SuccessResponse(w, r, http.StatusOK, "user get succesfully", resp)
}

func (u *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "uuid")

	uuid, err := uuid.Parse(id)
	fmt.Println(uuid)
	if err != nil {
		responses.ErrorResponse(w, r, "invalid user id", snsgo.NewErrorf(snsgo.ErrorCodeInvalidArg, "invalid user id: %v", id))
		return
	}

	user, err := u.userSvc.GetUserByID(r.Context(), uuid)
	if err != nil {
		responses.ErrorResponse(w, r, "cannot get user", err)
		return
	}

	resp := GetUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Gender:      user.Gender.String(),
		DateOfBirth: user.DateOfBirth.String(),
		Age:         user.Age,
		Email:       user.Email,
	}

	responses.SuccessResponse(w, r, http.StatusOK, "user get succesfully", resp)
}
