package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/GetterSethya/library"
	"github.com/GetterSethya/userProto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReplyService struct {
	Store                 *SqliteStorage
	UserServiceGrpcClient userProto.UserClient
}

func NewReplyService(store *SqliteStorage, userGrpcClient userProto.UserClient) *ReplyService {
	return &ReplyService{
		Store:                 store,
		UserServiceGrpcClient: userGrpcClient,
	}
}

func (s *ReplyService) RegisterRoutes(r *mux.Router) {
	// create reply -> http://localhost/v1/reply/{postId}/create
	r.HandleFunc("/{postId}/create", library.CreateHandler(library.JWTMiddleware(s.handleCreateReply))).Methods(http.MethodPost, http.MethodOptions)

	// create child reply -> http://localhost/v1/reply/{replyId}/create_child
	r.HandleFunc("/{replyId}/create_child", library.CreateHandler(library.JWTMiddleware(s.handleCreateChildReply))).Methods(http.MethodPost, http.MethodOptions)

	// update reply -> http://localhost/v1/reply/{replyId}/update
	r.HandleFunc("/{replyId}/update", library.CreateHandler(library.JWTMiddleware(s.handleUpdateReply))).Methods(http.MethodPost, http.MethodOptions)

	// delete reply -> http://localhost/v1/reply/{replyId}/delete
	r.HandleFunc("/{replyId}/update", library.CreateHandler(library.JWTMiddleware(s.handleDeleteReply))).Methods(http.MethodDelete, http.MethodOptions)

	// get replies by postId -> http://localhost/v1/reply/post/{postId}
	r.HandleFunc("/post/{postId}", library.CreateHandler(library.JWTMiddleware(s.handleGetReplyByPostId))).Methods(http.MethodGet, http.MethodOptions)

	// get replies by parentId -> http://localhost/v1/reply/child/{parentId}
	r.HandleFunc("/child/{parentId}", library.CreateHandler(library.JWTMiddleware(s.handleGetReplyByParentId))).Methods(http.MethodGet, http.MethodOptions)

	// get replies by replyId -> http://localhost/v1/reply/{replyId}
	r.HandleFunc("/{replyId}", library.CreateHandler(library.JWTMiddleware(s.handleGetReplyByReplyId))).Methods(http.MethodGet, http.MethodOptions)
}

func (s *ReplyService) handleGetReplyByPostId(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	postId := vars["postId"]
	urlQuery := r.URL.Query()
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

	replies := &[]Reply{}

	if err := s.Store.GetReplyByPostId(postId, int64(intCursor), int32(intLimit), replies); err != nil {
		return http.StatusNotFound, fmt.Errorf("reply did not exists")
	}

	var metaCursor int64

	if len(*replies) == 0 {
		metaCursor = 0
	} else {
		metaCursor = (*replies)[len(*replies)-1].CreatedAt
	}

	meta := struct {
		Cursor int64 `json:"cursor"`
	}{
		Cursor: metaCursor,
	}

	resp := library.NewResp("success", map[string]interface{}{
		"reply": replies,
		"meta":  meta,
	})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *ReplyService) handleGetReplyByParentId(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	parentId := vars["parentId"]
	urlQuery := r.URL.Query()
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

	replies := &[]Reply{}

	if err := s.Store.GetReplyByParentId(parentId, int64(intCursor), int32(intLimit), replies); err != nil {
		return http.StatusNotFound, fmt.Errorf("reply did not exists")
	}

	var metaCursor int64

	if len(*replies) == 0 {
		metaCursor = 0
	} else {
		metaCursor = (*replies)[len(*replies)-1].CreatedAt
	}

	meta := struct {
		Cursor int64 `json:"cursor"`
	}{
		Cursor: metaCursor,
	}

	resp := library.NewResp("success", map[string]interface{}{
		"reply": replies,
		"meta":  meta,
	})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *ReplyService) handleGetReplyByReplyId(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	replyId := vars["replyId"]
	reply := &Reply{}

	if err := s.Store.GetReplyById(replyId, reply); err != nil {
		return http.StatusNotFound, fmt.Errorf("reply did not exists")
	}

	resp := library.NewResp("success", map[string]interface{}{
		"reply": reply,
	})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *ReplyService) handleDeleteReply(w http.ResponseWriter, r *http.Request) (int, error) {
	userId := library.GetUserIdFromJWT(r)
	vars := mux.Vars(r)
	replyId := vars["replyId"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid reply detail")
	}
	defer r.Body.Close()

	reqReply := &Reply{}

	if err := json.Unmarshal(body, reqReply); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid reply detail")
	}

	replyDb := &Reply{}
	if err := s.Store.GetReplyById(replyId, replyDb); err != nil {
		return http.StatusNotFound, fmt.Errorf("reply doesn't exists")
	}

	if replyDb.IdUser != userId {
		return http.StatusForbidden, fmt.Errorf("forbidden")
	}

	if err := s.Store.DeleteReply(replyId); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("reply deleted", nil)
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *ReplyService) handleUpdateReply(w http.ResponseWriter, r *http.Request) (int, error) {
	userId := library.GetUserIdFromJWT(r)
	body, err := io.ReadAll(r.Body)
	vars := mux.Vars(r)
	replyId := vars["replyId"]

	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid reply detail")
	}
	defer r.Body.Close()

	reqReply := &Reply{}

	if err := json.Unmarshal(body, reqReply); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid reply detail")
	}

	replyDb := &Reply{}
	if err := s.Store.GetReplyById(replyId, replyDb); err != nil {
		return http.StatusNotFound, fmt.Errorf("reply doesnt exists")
	}

	if replyDb.IdUser != userId {
		return http.StatusForbidden, fmt.Errorf("forbidden")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	if err := s.Store.UpdateReply(replyId, reqReply.Body); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("reply updated", nil)
	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *ReplyService) handleCreateChildReply(w http.ResponseWriter, r *http.Request) (int, error) {
	userId := library.GetUserIdFromJWT(r)
	uuid := uuid.NewString()
	vars := mux.Vars(r)
	replyId := vars["replyId"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid reply detail")
	}
	defer r.Body.Close()

	reply := &Reply{}

	if err := json.Unmarshal(body, reply); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid reply detail")
	}

	in := &userProto.GetUserByIdReq{
		Id: userId,
	}

	userDb, err := s.UserServiceGrpcClient.GetUserById(r.Context(), in)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	if err := s.Store.CreateReply(uuid, reply.Body, userId, userDb.Username, userDb.Name, userDb.Profile, reply.IdPost, replyId); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("child reply created", nil)
	library.WriteJson(w, http.StatusCreated, resp)

	return http.StatusCreated, nil
}

func (s *ReplyService) handleCreateReply(w http.ResponseWriter, r *http.Request) (int, error) {
	userId := library.GetUserIdFromJWT(r)
	uuid := uuid.NewString()
	vars := mux.Vars(r)
	postId := vars["postId"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid reply detail")
	}
	defer r.Body.Close()

	reply := &Reply{}

	if err := json.Unmarshal(body, reply); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid reply detail")
	}

	in := &userProto.GetUserByIdReq{
		Id: userId,
	}

	userDb, err := s.UserServiceGrpcClient.GetUserById(r.Context(), in)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	if err := s.Store.CreateReply(uuid, reply.Body, userId, userDb.Username, userDb.Name, userDb.Profile, postId, nil); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("reply created", nil)
	library.WriteJson(w, http.StatusCreated, resp)

	return http.StatusCreated, nil
}
