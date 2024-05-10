package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/GetterSethya/library"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	JWTSecret     string
	RefreshSecret string
	GrpcClient    library.UserClient
}

func NewAuthService(jwtSecret string, grpcClient library.UserClient, refreshSecret string) *AuthService {

	return &AuthService{
		JWTSecret:     jwtSecret,
		GrpcClient:    grpcClient,
		RefreshSecret: refreshSecret,
	}
}

func (s *AuthService) RegisterRoutes(r *mux.Router) {

	//v1/auth/login
	r.HandleFunc("/login", library.CreateHandler(s.handleLoginAuth)).Methods(http.MethodPost)

	//v1/auth/register
	r.HandleFunc("/register", library.CreateHandler(s.handleRegisterAuth)).Methods(http.MethodPost)

	//v1/auth/refresh
	r.HandleFunc("/refresh", library.CreateHandler(s.handleRefreshAuth)).Methods(http.MethodPost)

}

func (s *AuthService) handleRegisterAuth(w http.ResponseWriter, r *http.Request) (int, error) {

	log.Println("hit handle register auth")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid user detail")
	}

	defer r.Body.Close()

	user := &RegisterUser{}

	if err := json.Unmarshal(body, user); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid user detail")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	uuid := uuid.NewString()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println("Error when hashing password:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	// panggil lewat grpc
	in := &library.CreateUserReq{
		Id:           uuid,
		Username:     user.Username,
		Name:         user.Name,
		HashPassword: string(hashPassword),
	}

	_, err = s.GrpcClient.CreateUser(r.Context(), in)
	if err != nil {
		log.Println("Error when calling s.GrpcClient:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("user created!", nil)

	library.WriteJson(w, http.StatusCreated, resp)

	return http.StatusCreated, nil
}

func (s *AuthService) handleLoginAuth(w http.ResponseWriter, r *http.Request) (int, error) {

	log.Println("hit handle register auth")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid user detail")
	}

	defer r.Body.Close()

	user := &LoginUser{}
	if err := json.Unmarshal(body, user); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid user creds")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	// panggil grpc getUserByUsername
	in := &library.GetUserByUsernameReq{
		Username: user.Username,
	}

	userDb, err := s.GrpcClient.GetUserPasswordByUsername(r.Context(), in)
	if err != nil {
		log.Println("Error when calling GetUserByUsername:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDb.HashPassword), []byte(user.Password)); err != nil {
		return http.StatusBadRequest, fmt.Errorf("Invalid username/password")
	}

	// generate jwt token, refresh token
	jwtToken, err := library.CreateJWT(userDb.Id, s.JWTSecret, time.Now().Add(6*time.Hour))
	refreshToken, err := library.CreateJWT(userDb.Id, s.RefreshSecret, time.Now().Add(24*time.Hour))

	library.WriteJson(w, http.StatusOK, map[string]interface{}{"accessToken": jwtToken, "refreshToken": refreshToken})

	return http.StatusOK, nil
}

func (s *AuthService) handleRefreshAuth(w http.ResponseWriter, r *http.Request) (int, error) {

	///////////////////////////
	//TODO handleRefreshAuth//
	/////////////////////////

	log.Println("hit handle refresh auth")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid user detail")
	}

	defer r.Body.Close()

	re := &Refresh{}
	if err := json.Unmarshal(body, re); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid refresh token")
	}

	token, err := library.ValidateJWT(re.RefreshToken, s.RefreshSecret)
	if err != nil {
		return http.StatusUnauthorized, fmt.Errorf("Invalid refresh token")
	}

	userId := token.Claims.(jwt.MapClaims)["sub"].(string)

	// generate refresh and access token
	jwtToken, err := library.CreateJWT(userId, s.JWTSecret, time.Now().Add(6*time.Hour))
	refreshToken, err := library.CreateJWT(userId, s.RefreshSecret, time.Now().Add(24*time.Hour))

	library.WriteJson(w, http.StatusOK, map[string]interface{}{"accessToken": jwtToken, "refreshToken": refreshToken})

	return http.StatusOK, nil
}
