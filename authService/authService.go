package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/GetterSethya/library"
	"github.com/GetterSethya/userProto"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	JWTSecret             string
	RefreshSecret         string
	UserServiceGrpcClient userProto.UserClient
}

type tokenChan struct {
	token     string
	err       error
	tokenType string
}

func NewAuthService(jwtSecret string, grpcClient userProto.UserClient, refreshSecret string) *AuthService {
	return &AuthService{
		JWTSecret:             jwtSecret,
		UserServiceGrpcClient: grpcClient,
		RefreshSecret:         refreshSecret,
	}
}

func (s *AuthService) RegisterRoutes(r *mux.Router) {
	// v1/auth/login
	r.HandleFunc("/login", library.CreateHandler(s.handleLoginAuth)).Methods(http.MethodPost, http.MethodOptions)

	// v1/auth/register
	r.HandleFunc("/register", library.CreateHandler(s.handleRegisterAuth)).Methods(http.MethodPost, http.MethodOptions)

	// v1/auth/refresh
	r.HandleFunc("/refresh", library.CreateHandler(s.handleRefreshAuth)).Methods(http.MethodPost, http.MethodOptions)
}

func (s *AuthService) handleRegisterAuth(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle register auth")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid user detail")
	}

	defer r.Body.Close()

	user := &RegisterUser{}

	if err := json.Unmarshal(body, user); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid user detail")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	uuid := uuid.NewString()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println("Error when hashing password:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	// panggil lewat grpc
	in := &userProto.CreateUserReq{
		Id:           uuid,
		Username:     user.Username,
		Name:         user.Name,
		HashPassword: string(hashPassword),
	}

	_, err = s.UserServiceGrpcClient.CreateUser(r.Context(), in)
	if err != nil {
		log.Println("Error when calling s.GrpcClient:", err)
		return http.StatusInternalServerError, fmt.Errorf("username is invalid/already used")
	}

	resp := library.NewResp("user created!", nil)

	library.WriteJson(w, http.StatusCreated, resp)

	return http.StatusCreated, nil
}

func (s *AuthService) handleLoginAuth(w http.ResponseWriter, r *http.Request) (int, error) {
	ch := make(chan tokenChan)
	defer close(ch)

	log.Println("hit handle login auth")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid user detail")
	}

	defer r.Body.Close()

	user := &LoginUser{}
	if err := json.Unmarshal(body, user); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid user creds")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	// panggil grpc getUserByUsername
	in := &userProto.GetUserByUsernameReq{
		Username: user.Username,
	}

	userDb, err := s.UserServiceGrpcClient.GetUserPasswordByUsername(r.Context(), in)
	if err != nil {
		log.Println("Error when calling GetUserByUsername:", err)
		return http.StatusBadRequest, fmt.Errorf("username/password wrong")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDb.HashPassword), []byte(user.Password)); err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid username/password")
	}

	// generate jwt token, refresh token
	// expAccess := time.Now().Add(6 * time.Second)
	expAccess := time.Now().Add(6 * time.Hour)
	expRefresh := time.Now().Add(24 * time.Hour)
	go func() {
		jwtToken, err := library.CreateJWT(userDb.Id, s.JWTSecret, expAccess)
		ch <- tokenChan{token: jwtToken, err: err, tokenType: "access"}
	}()

	go func() {
		refreshToken, err := library.CreateJWT(userDb.Id, s.RefreshSecret, expRefresh)
		ch <- tokenChan{token: refreshToken, err: err, tokenType: "refresh"}
	}()

	var jwtToken string
	var refreshToken string
	var errs []error

	for i := 0; i < 2; i++ {
		res := <-ch
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		if res.tokenType == "access" {
			jwtToken = res.token
		} else {
			refreshToken = res.token
		}

	}

	if len(errs) > 0 {
		log.Println("Error when generating jwt access and refresh token:")
		for _, err := range errs {
			log.Println(err)
		}
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	access := &http.Cookie{
		Name:     "accessToken",
		Value:    jwtToken,
		Path:     "/",
		Expires:  expAccess,
		HttpOnly: true,
	}
	refresh := &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  expRefresh,
		HttpOnly: true,
	}

	resp := library.NewResp("User authenticated", map[string]interface{}{"accessToken": jwtToken, "refreshToken": refreshToken})

	http.SetCookie(w, access)
	http.SetCookie(w, refresh)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *AuthService) handleRefreshAuth(w http.ResponseWriter, r *http.Request) (int, error) {
	ch := make(chan tokenChan)

	///////////////////////////
	//TODO handleRefreshAuth//
	/////////////////////////

	log.Println("hit handle refresh auth")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid user detail")
	}

	defer r.Body.Close()

	re := &Refresh{}
	if err := json.Unmarshal(body, re); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid refresh token")
	}

	token, err := library.ValidateJWT(re.RefreshToken, s.RefreshSecret)
	if err != nil {
		return http.StatusUnauthorized, fmt.Errorf("invalid refresh token")
	}

	userId := token.Claims.(jwt.MapClaims)["sub"].(string)

	// generate jwt token, refresh token
	go func() {
		// jwtToken, err := library.CreateJWT(userId, s.JWTSecret, time.Now().Add(6*time.Second))
		jwtToken, err := library.CreateJWT(userId, s.JWTSecret, time.Now().Add(6*time.Hour))
		ch <- tokenChan{token: jwtToken, err: err, tokenType: "access"}
	}()

	go func() {
		refreshToken, err := library.CreateJWT(userId, s.RefreshSecret, time.Now().Add(24*time.Hour))
		ch <- tokenChan{token: refreshToken, err: err, tokenType: "refresh"}
	}()

	var jwtToken string
	var refreshToken string
	var errs []error

	for i := 0; i < 2; i++ {
		res := <-ch
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		if res.tokenType == "access" {
			jwtToken = res.token
		} else {
			refreshToken = res.token
		}

	}

	if len(errs) > 0 {
		log.Println("Error when generating jwt access and refresh token:")
		for _, err := range errs {
			log.Println(err)
		}
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("token refreshed", map[string]interface{}{"accessToken": jwtToken, "refreshToken": refreshToken})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}
