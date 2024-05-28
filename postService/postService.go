package main

import (
	"database/sql"
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

type PostService struct {
	Store                 *SqliteStorage
	UserServiceGrpcClient userProto.UserClient
}

func NewUserService(store *SqliteStorage, grpcClient userProto.UserClient) *PostService {

	return &PostService{
		Store:                 store,
		UserServiceGrpcClient: grpcClient,
	}
}

func (s *PostService) RegisterRoutes(r *mux.Router) {

	//v1/post/create --> bikin post
	r.HandleFunc("/", library.CreateHandler(library.JWTMiddleware(s.handleCreatePost))).Methods(http.MethodPost, http.MethodOptions)

	//v1/post/ --> delete post (soft delete, deletedAt nya diisi unixepoch) !Penting nanti di cek dulu apakah idUser dari jwt sama dengan idUser yang ada didalam post
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleDeletePost))).Methods(http.MethodDelete, http.MethodOptions)

	//v1/post/{id} --> update post by id (cuma update isi post) !ini di cek juga idUser nya
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleUpdatePost))).Methods(http.MethodPost, http.MethodOptions)

	//v1/post --> list post
	r.HandleFunc("/", library.CreateHandler(library.JWTMiddleware(s.handleListPost))).Methods(http.MethodGet, http.MethodOptions)

	//v1/post/user/{idUser} --> list post by user
	r.HandleFunc("/user/{idUser}", library.CreateHandler(library.JWTMiddleware(s.handleListPostByUser))).Methods(http.MethodGet, http.MethodOptions)

	//v1/post/{id} --> get post by id
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleGetPostById))).Methods(http.MethodGet, http.MethodOptions)

}

func (s *PostService) handleGetPostById(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle get post by id")
	vars := mux.Vars(r)
	postId := vars["id"]

	if err := uuid.Validate(postId); err != nil {
		log.Println("Invalid post uuid url")
		return http.StatusBadRequest, fmt.Errorf("Post didnot exists")
	}

	post := &Post{}

	if err := s.Store.GetPostById(postId, post); err != nil {
		log.Println("getPostById err:", err)
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("Post didnot exists")
		}

		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("Success", map[string]interface{}{"post": post})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleListPostByUser(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle list post by user")

	urlQuery := r.URL.Query()
	limit := urlQuery.Get("limit")
	offset := urlQuery.Get("offset")
	vars := mux.Vars(r)
	profileId := vars["idUser"]

	if limit == "" {
		limit = "10"
	}

	if offset == "" {
		offset = "0"
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		intOffset = 0
	}

	posts := &[]Post{}

	if err := s.Store.ListPostByUser(int64(intOffset), profileId, posts); err != nil {

		log.Println("Error when getting listPost:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("success", map[string]interface{}{
		"posts": posts,
	})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleListPost(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle list post")

	urlQuery := r.URL.Query()
	offset := urlQuery.Get("offset")

	if offset == "" {
		offset = "0"
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		intOffset = 0
	}

	posts := &[]Post{}

	if err := s.Store.ListPost(int64(intOffset), posts); err != nil {

		log.Println("Error when getting listPost:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("success", map[string]interface{}{
		"posts": posts,
	})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleUpdatePost(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle update post")

	vars := mux.Vars(r)
	postId := vars["id"]
	userId := library.GetUserIdFromJWT(r)
	post := &Post{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid post detail")
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, post); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid post detail")
	}

	fetchPost := &Post{}
	if err := s.Store.GetPostById(postId, fetchPost); err != nil {
		log.Println("getPostById err:", err)
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("Post didnot exists")
		}

		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	if fetchPost.IdUser != userId {
		log.Println("Post userid did not match userid from jwt")
		return http.StatusUnauthorized, fmt.Errorf("Unauthorized")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	if err := s.Store.UpdatePostBody(postId, post.Body, userId); err != nil {
		log.Println("Error when updating post body:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("Post updated successfully!", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleDeletePost(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle delete post")

	vars := mux.Vars(r)
	postId := vars["id"]
	userId := library.GetUserIdFromJWT(r)
	post := &Post{}

	err := s.Store.GetPostById(postId, post)
	if err != nil {
		log.Println("getPostById err:", err)
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("Post didnot exists")
		}

		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	if userId != post.IdUser {
		log.Println("userid from jwt didnot match iduser post")
		return http.StatusForbidden, fmt.Errorf("Forbidden")
	}

	if err := s.Store.DeletePostById(postId, userId); err != nil {
		log.Println("Error when deleting post by id:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("Post deleted successfully", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleCreatePost(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle create post")

	idUser := library.GetUserIdFromJWT(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid post detail")
	}

	defer r.Body.Close()

	in := &userProto.GetUserByIdReq{
		Id: idUser,
	}

	grpcResp, err := s.UserServiceGrpcClient.GetUserById(r.Context(), in)
	if err != nil {
		log.Println("Error when dialing grpc client with getUserById method:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	uuid := uuid.NewString()

	post := &Post{
		Id:       uuid,
		IdUser:   idUser,
		Username: grpcResp.GetUsername(),
		Name:     grpcResp.GetName(),
	}

	if err := json.Unmarshal(body, post); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("Invalid post detail")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	if err := s.Store.CreatePost(post.Id, post.Image, post.Body, post.IdUser, post.Username, post.Name); err != nil {
		log.Println("Error when creating post:", err)
		return http.StatusInternalServerError, fmt.Errorf("Something went wrong")
	}

	resp := library.NewResp("post created!", nil)

	library.WriteJson(w, http.StatusCreated, resp)

	return http.StatusCreated, nil

}
