package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/GetterSethya/imageProto"
	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store           *SqliteStorage
	RabbitMQ        *library.RabbitMq
	ImageGrpcClient imageProto.UserClient
}

const defaultProfile = "1714794135-a06a41d8-6351-4dbb-9141-a7e2ace86a35.jpg"

func NewUserService(store *SqliteStorage, producer *library.RabbitMq, imageGrpcClient imageProto.UserClient) *UserService {
	return &UserService{
		Store:           store,
		RabbitMQ:        producer,
		ImageGrpcClient: imageGrpcClient,
	}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	// v1/user/{id}
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleGetUserById))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleDeleteUserById))).Methods(http.MethodDelete, http.MethodOptions)

	// v1/user/username/{username}
	r.HandleFunc("/username/{username}", library.CreateHandler(library.JWTMiddleware(s.handleGetUserByUsername))).Methods(http.MethodGet, http.MethodOptions)

	// v1/user/update -> update user data by jwt
	r.HandleFunc("/update", library.CreateHandler(library.JWTMiddleware(s.handleUpdateUserByJWT))).Methods(http.MethodPost, http.MethodOptions)

	// v1/user/update_password
	r.HandleFunc("/update_password", library.CreateHandler(library.JWTMiddleware(s.handleUpdateUserPassword))).Methods(http.MethodPost, http.MethodOptions)
}

func (s *UserService) handleDeleteUserById(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle delete user by id")

	vars := mux.Vars(r)
	urlIdUser := vars["id"]
	idUser := library.GetUserIdFromJWT(r)

	if urlIdUser != idUser {
		return http.StatusUnauthorized, fmt.Errorf("Unauthorized")
	}

	if err := s.Store.DeleteUserById(idUser); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("User didnot exists")
		}

		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("User deleted!", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *UserService) handleGetUserById(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle get user by id")

	vars := mux.Vars(r)
	id := vars["id"]

	user := &ReturnUser{}

	err := s.Store.GetUserById(id, user)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("user did not exists/not found")
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
		return http.StatusNotFound, fmt.Errorf("user did not exists/not found")
	}

	resp := library.NewResp("Success", map[string]interface{}{"user": user})
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *UserService) handleUpdateUserPassword(w http.ResponseWriter, r *http.Request) (int, error) {
	type changePassReq struct {
		CurrentPassword    string `json:"currentPassword"`
		NewPassword        string `json:"newPassword"`
		ConfirmNewPassword string `json:"confirmNewPassword"`
	}

	log.Println("hit handle update user password")

	userIdJWT := library.GetUserIdFromJWT(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
	}

	defer r.Body.Close()

	changePass := &changePassReq{}

	if err := json.Unmarshal(body, &changePass); err != nil {
		log.Println("Error when umarshaling json", err)
		return http.StatusBadRequest, fmt.Errorf("invalid user detail")
	}

	if changePass.NewPassword != changePass.ConfirmNewPassword {
		return http.StatusBadRequest, fmt.Errorf("new password didnot match with confirm new password")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////
	newPassword, err := bcrypt.GenerateFromPassword([]byte(changePass.NewPassword), 12)
	if err != nil {
		log.Println("Error when hashing password:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	if err := s.Store.UpdateUserPasswordById(string(newPassword), userIdJWT); err != nil {
		log.Println("Error when updating username:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("User password updated!", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *UserService) handleUpdateUserByJWT(w http.ResponseWriter, r *http.Request) (int, error) {
	// semua input user pake formData
	type UserDetailChangeEvent struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Profile string `json:"profile"`
	}

	userIdJWT := library.GetUserIdFromJWT(r)

	log.Println("hit handle update user name")

	// ambil image dari formdata
	err := r.ParseMultipartForm(2 * 1024 * 1024)
	if err != nil {
		log.Println("Error when parsing request formdata:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid image/image is too big")
	}

	reqName := r.FormValue("reqName")
	if reqName == "" {
		log.Println("Error when getting reqName, invalid/missing form data")
		return http.StatusBadRequest, fmt.Errorf("invalid/missing form data")
	}

	userData := &ReturnUser{}

	err = s.Store.GetUserById(userIdJWT, userData)
	if err != nil {
		log.Println("Error when getting user data", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	newUserData := &User{}
	newUserData.Name = reqName

	file, handler, err := r.FormFile("reqImage")
	if err != nil {
		newUserData.Profile = userData.Profile
	} else {
		defer file.Close()
		bytesFile, err := io.ReadAll(file)
		if err != nil {
			log.Println("Error when reading file from form data:", err)
			return http.StatusBadRequest, fmt.Errorf("invalid post detail")
		}

		createImageGrpcReq := &imageProto.CreateImageReq{
			ImageFile: bytesFile,
			FileName:  handler.Filename,
		}

		imageGrpcResp, err := s.ImageGrpcClient.CreateImage(r.Context(), createImageGrpcReq)
		if err != nil {
			log.Println("Error when calling createImage grpc:", err)
			return http.StatusInternalServerError, fmt.Errorf("something went wrong")
		}
		newUserData.Profile = imageGrpcResp.GetFilename()
	}

	log.Printf("newUserData: %+v", newUserData)

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	if err := s.Store.UpdateUserNameAndProfile(newUserData.Name, newUserData.Profile, userIdJWT); err != nil {
		log.Println("Error when updating username:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	ch, err := s.RabbitMQ.Conn.Channel()
	if err != nil {
		log.Println("Error when creating channel in user service handler")
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"userServiceExchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error when declaring exchange:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	event := UserDetailChangeEvent{
		Id:      userIdJWT,
		Name:    newUserData.Name,
		Profile: newUserData.Profile,
	}

	publishBody, err := json.Marshal(event)
	if err != nil {
		log.Println("Error when marshaling event:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	err = ch.PublishWithContext(
		r.Context(),
		"userServiceExchange",
		"user.detail.change",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         publishBody,
		},
	)
	if err != nil {
		log.Println("Error when publishing user.detail.change event")
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("User updated!", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}
