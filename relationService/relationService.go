package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/GetterSethya/library"
	"github.com/GetterSethya/userProto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RelationService struct {
	Store                 *SqliteStorage
	UserServiceGrpcClient userProto.UserClient
}

func NewRelationService(store *SqliteStorage, userGrpcClient userProto.UserClient) *RelationService {
	return &RelationService{
		Store:                 store,
		UserServiceGrpcClient: userGrpcClient,
	}
}

func (s *RelationService) RegisterRoutes(r *mux.Router) {
	// follow -> http://localhost/v1/relation/{userId}/follow
	r.HandleFunc("/{userId}/follow", library.CreateHandler(library.JWTMiddleware(s.handleFollow))).Methods(http.MethodPost, http.MethodOptions)

	// unfollow -> http://localhost/v1/relation/{userId}/unfollow
	r.HandleFunc("/{userId}/unfollow", library.CreateHandler(library.JWTMiddleware(s.handleUnfollow))).Methods(http.MethodPost, http.MethodOptions)

	// list follower -> http://localhost/v1/relation/{userId}/follower
	r.HandleFunc("/{userId}/follower", library.CreateHandler(library.JWTMiddleware(s.handleListFollower))).Methods(http.MethodGet, http.MethodOptions)

	// list following -> http://localhost/v1/relation/{userId}/following
	r.HandleFunc("/{userId}/following", library.CreateHandler(library.JWTMiddleware(s.handleListFollowing))).Methods(http.MethodGet, http.MethodOptions)

	// check is following -> http://localhost/v1/relation/{userId}/is-following
	r.HandleFunc("/{userId}/is-following", library.CreateHandler(library.JWTMiddleware(s.handleIsFollowing))).Methods(http.MethodGet, http.MethodOptions)

	// check is follower -> http://localhost/v1/relation/{userId}/is-follower
	r.HandleFunc("/{userId}/is-follower", library.CreateHandler(library.JWTMiddleware(s.handleIsFollower))).Methods(http.MethodGet, http.MethodOptions)
}

func (s *RelationService) handleFollow(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	targetId := vars["userId"]
	followerId := library.GetUserIdFromJWT(r)

	if targetId == followerId {
		return http.StatusBadRequest, fmt.Errorf("cannot follow your own account")
	}

	targetIn := &userProto.GetUserByIdReq{
		Id: targetId,
	}
	followerIn := &userProto.GetUserByIdReq{
		Id: followerId,
	}

	targetUserData, err := s.UserServiceGrpcClient.GetUserById(r.Context(), targetIn)
	if err != nil {
		log.Println("error when fetching target user data:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	followerUserData, err := s.UserServiceGrpcClient.GetUserById(r.Context(), followerIn)
	if err != nil {
		log.Println("error when fetching follower user data:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	unixEpoch := time.Now().Unix()

	if err := s.Store.CreateUser(
		targetUserData.Id,
		targetUserData.Username,
		targetUserData.Name,
		targetUserData.Profile,
		unixEpoch,
		unixEpoch,
	); err != nil {
		log.Println("error when creating target user data:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	if err := s.Store.CreateUser(
		followerUserData.Id,
		followerUserData.Username,
		followerUserData.Name,
		followerUserData.Profile,
		unixEpoch,
		unixEpoch,
	); err != nil {
		log.Println("error when creating follower user data:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	uidRelation := uuid.NewString()

	if err := s.Store.CreateRelation(uidRelation, targetUserData.Id, followerUserData.Id); err != nil {
		log.Println("failed when creating relation:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("success", nil)
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *RelationService) handleUnfollow(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	targetId := vars["userId"]
	followerId := library.GetUserIdFromJWT(r)

	if targetId == followerId {
		return http.StatusBadRequest, fmt.Errorf("cannot unfollow your own account")
	}

	if err := s.Store.DeleteRelation(targetId, followerId); err != nil {
		log.Println("failed when creating relation:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("success", nil)
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *RelationService) handleIsFollowing(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	targetId := vars["userId"]
	followerId := library.GetUserIdFromJWT(r)

	if targetId == followerId {
		return http.StatusBadRequest, fmt.Errorf("bad request")
	}

	isFollowing, err := s.Store.IsFollowing(targetId, followerId)
	if err != nil {
		log.Println("failed when creating relation:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("success", map[string]interface{}{"isFollowing": isFollowing})
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *RelationService) handleIsFollower(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	followerId := vars["userId"]
	targetId := library.GetUserIdFromJWT(r)

	if targetId == followerId {
		return http.StatusBadRequest, fmt.Errorf("bad request")
	}

	isFollowing, err := s.Store.IsFollowing(targetId, followerId)
	if err != nil {
		log.Println("failed when creating relation:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("success", map[string]interface{}{"isFollower": isFollowing})
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *RelationService) handleListFollower(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	urlQuery := r.URL.Query()
	userId := vars["userId"]
	limit := urlQuery.Get("limit")
	cursor := urlQuery.Get("cursor")

	if limit == "" {
		limit = "10"
	}

	if cursor == "" {
		cursor = "0"
	}

	intCursor, err := strconv.Atoi(cursor)
	if err != nil {
		intCursor = 0
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		intLimit = 10
	}

	followers := &[]User{}

	if err := s.Store.GetFollowers(int64(intCursor), userId, intLimit, followers); err != nil {
		log.Println("cannot get followers", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	return http.StatusOK, nil
}

func (s *RelationService) handleListFollowing(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	urlQuery := r.URL.Query()
	userId := vars["userId"]
	limit := urlQuery.Get("limit")
	cursor := urlQuery.Get("cursor")

	if limit == "" {
		limit = "10"
	}

	if cursor == "" {
		cursor = "0"
	}

	intCursor, err := strconv.Atoi(cursor)
	if err != nil {
		intCursor = 0
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		intLimit = 10
	}

	followers := &[]User{}

	if err := s.Store.GetFollowing(int64(intCursor), userId, intLimit, followers); err != nil {
		log.Println("cannot get following", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	return http.StatusOK, nil
}
