package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store *SqliteStorage
}

const defaultProfile = "http://localhost/v1/image/thumbnail/1714794135-a06a41d8-6351-4dbb-9141-a7e2ace86a35.jpg"

func NewUserService(store *SqliteStorage) *UserService {

	return &UserService{
		Store: store,
	}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {

	//v1/user/id/{id}
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleGetUserById))).Methods(http.MethodGet)
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleDeleteUserById))).Methods(http.MethodDelete)

	//v1/user/username/{username}
	r.HandleFunc("/username/{username}", library.CreateHandler(library.JWTMiddleware(s.handleGetUserByUsername))).Methods(http.MethodGet)

	//v1/user/update
	r.HandleFunc("/update", library.CreateHandler(library.JWTMiddleware(s.handleUpdateUserName))).Methods(http.MethodPost)

	//v1/user/update_password
	r.HandleFunc("/update_password", library.CreateHandler(library.JWTMiddleware(s.handleUpdateUserPassword))).Methods(http.MethodPost)

}

func (s *UserService) handleDeleteUserById(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle delete user by id")
	//////////////////////////////
	//TODO handleDeleteUserById//
	////////////////////////////

	return http.StatusOK, nil
}

func (s *UserService) handleGetUserById(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle get user by id")

	vars := mux.Vars(r)
	id := vars["id"]

	user := &ReturnUser{}

	err := s.Store.GetUserById(id, user)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("User did not exists/not found")
	}

	resp := library.NewResp("Success", map[string]interface{}{"user": user})
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *UserService) handleGetUserByUsername(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle get user by id")

	vars := mux.Vars(r)
	username := vars["username"]

	user := &ReturnUser{}

	err := s.Store.GetUserByUsername(username, user)
	if err != nil {
		log.Println("Error when getting user by username", err)
		return http.StatusNotFound, fmt.Errorf("User did not exists/not found")
	}

	resp := library.NewResp("Success", map[string]interface{}{"user": user})
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *UserService) handleUpdateUserPassword(w http.ResponseWriter, r *http.Request) (int, error) {

	log.Println("hit handle update user password")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
	}

	defer r.Body.Close()

	user := &User{}

	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("Error when umarshaling json", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid user detail")
	}

	userIdJWT := library.GetUserIdFromJWT(r)

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////
	newPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashPassword), 12)
	if err != nil {
		log.Println("Error when hashing password:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	if err := s.Store.UpdateUserPasswordById(string(newPassword), userIdJWT); err != nil {
		log.Println("Error when updating username:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("User updated!", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *UserService) handleUpdateUserName(w http.ResponseWriter, r *http.Request) (int, error) {

	log.Println("hit handle update user name")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
	}

	defer r.Body.Close()

	user := &User{}

	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("Error when umarshaling json", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid user detail")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	userIdJWT := library.GetUserIdFromJWT(r)

	if err := s.Store.UpdateUserNameById(user.Name, userIdJWT); err != nil {
		log.Println("Error when updating username:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("User updated!", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}
