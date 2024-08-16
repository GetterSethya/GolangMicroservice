package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/GetterSethya/imageProto"
	"github.com/GetterSethya/library"
	"github.com/GetterSethya/userProto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PostService struct {
	Store                  *SqliteStorage
	UserServiceGrpcClient  userProto.UserClient
	ImageServiceGrpcClient imageProto.UserClient
}

func NewUserService(store *SqliteStorage, userGrpcClient userProto.UserClient, imageGrpcClient imageProto.UserClient) *PostService {
	return &PostService{
		Store:                  store,
		UserServiceGrpcClient:  userGrpcClient,
		ImageServiceGrpcClient: imageGrpcClient,
	}
}

func (s *PostService) RegisterRoutes(r *mux.Router) {
	// v1/post/ --> bikin post
	r.HandleFunc("/", library.CreateHandler(library.JWTMiddleware(s.handleCreatePost))).Methods(http.MethodPost, http.MethodOptions)

	// v1/post/ --> delete post (soft delete, deletedAt nya diisi unixepoch) !Penting nanti di cek dulu apakah idUser dari jwt sama dengan idUser yang ada didalam post
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleDeletePost))).Methods(http.MethodDelete, http.MethodOptions)

	// v1/post/{id} --> update post by id (cuma update isi post) !ini di cek juga idUser nya
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleUpdatePost))).Methods(http.MethodPost, http.MethodOptions)

	// v1/post --> list post
	r.HandleFunc("/", library.CreateHandler(library.JWTMiddleware(s.handleListPost))).Methods(http.MethodGet, http.MethodOptions)

	// v1/post/user/{idUser} --> list post by user
	r.HandleFunc("/user/{idUser}", library.CreateHandler(library.JWTMiddleware(s.handleListPostByUser))).Methods(http.MethodGet, http.MethodOptions)

	// v1/post/{id} --> get post by id
	r.HandleFunc("/{id}", library.CreateHandler(library.JWTMiddleware(s.handleGetPostById))).Methods(http.MethodGet, http.MethodOptions)

	// v1/{id}/like --> create like
	r.HandleFunc("/{id}/like", library.CreateHandler(library.JWTMiddleware(s.handleCreateLike))).Methods(http.MethodGet, http.MethodOptions)

	// v1/{id}/cancel-like --> delete like
	r.HandleFunc("/{id}/cancel-like", library.CreateHandler(library.JWTMiddleware(s.handleDeleteLike))).Methods(http.MethodGet, http.MethodOptions)

	// v1/post/like/{userId} --> list liked by user
	r.HandleFunc("/post/like/{userId}", library.CreateHandler(library.JWTMiddleware(s.handleListPostLikedByUser))).Methods(http.MethodGet, http.MethodOptions)
}

func (s *PostService) handleListPostLikedByUser(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle list post liked by user")
	vars := mux.Vars(r)
	urlQuery := r.URL.Query()
	limit := urlQuery.Get("limit")
	cursor := urlQuery.Get("cursor")
	userId := vars["userId"]

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

	posts := &[]Post{}

	if err := s.Store.ListLikedByUser(int64(intCursor), userId, int32(intLimit), posts); err != nil {

		log.Println("Error when getting list liked posts:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	var metaCursor int64

	if len(*posts) == 0 {
		metaCursor = 0
	} else {
		metaCursor = (*posts)[len(*posts)-1].CreatedAt
	}

	meta := struct {
		Cursor int64 `json:"cursor"`
	}{
		Cursor: metaCursor,
	}

	resp := library.NewResp("success", map[string]interface{}{
		"posts": posts,
		"meta":  meta,
	})

	library.WriteJson(w, http.StatusOK, resp)

	if err := uuid.Validate(userId); err != nil {
		log.Println("Invalid user uuid url")
		return http.StatusBadRequest, fmt.Errorf("user didnot exists")
	}

	return http.StatusOK, nil
}

func (s *PostService) handleCreateLike(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle create like")
	vars := mux.Vars(r)
	postId := vars["id"]
	userId := library.GetUserIdFromJWT(r)

	if err := uuid.Validate(postId); err != nil {
		log.Println("Invalid post uuid url")
		return http.StatusBadRequest, fmt.Errorf("post didnot exists")
	}

	likeId := uuid.NewString()

	if err := s.Store.CreateLike(likeId, userId, postId); err != nil {
		log.Println("Error when creating like", err)
		return http.StatusBadRequest, fmt.Errorf("error when creating like")
	}

	resp := library.NewResp("success", nil)

	library.WriteJson(w, http.StatusCreated, resp)

	return http.StatusCreated, nil
}

func (s *PostService) handleDeleteLike(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle delete like")
	vars := mux.Vars(r)
	postId := vars["id"]
	userId := library.GetUserIdFromJWT(r)

	if err := uuid.Validate(postId); err != nil {
		log.Println("Invalid post uuid url")
		return http.StatusBadRequest, fmt.Errorf("post didnot exists")
	}

	if err := s.Store.DeleteLike(postId, userId); err != nil {
		log.Println("Error when deleting like", err)
		return http.StatusBadRequest, fmt.Errorf("error when deleting like")
	}

	resp := library.NewResp("success", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleGetPostById(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle get post by id")
	vars := mux.Vars(r)
	postId := vars["id"]
	currentUser := library.GetUserIdFromJWT(r)

	if err := uuid.Validate(postId); err != nil {
		log.Println("Invalid post uuid url")
		return http.StatusBadRequest, fmt.Errorf("Post didnot exists")
	}

	post := &Post{}

	if err := s.Store.GetPostById(currentUser, postId, post); err != nil {
		log.Println("getPostById err:", err)
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("Post didnot exists")
		}

		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("Success", map[string]interface{}{"post": post})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleListPostByUser(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle list post by user")

	urlQuery := r.URL.Query()
	limit := urlQuery.Get("limit")
	cursor := urlQuery.Get("cursor")
	vars := mux.Vars(r)
	profileId := vars["idUser"]

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

	posts := &[]Post{}

	if err := s.Store.ListPostByUser(int64(intCursor), profileId, int32(intLimit), posts); err != nil {

		log.Println("Error when getting listPost:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	var metaCursor int64

	if len(*posts) == 0 {
		metaCursor = 0
	} else {
		metaCursor = (*posts)[len(*posts)-1].CreatedAt
	}

	meta := struct {
		Cursor int64 `json:"cursor"`
	}{
		Cursor: metaCursor,
	}

	resp := library.NewResp("success", map[string]interface{}{
		"posts": posts,
		"meta":  meta,
	})

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleListPost(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle list post")

	urlQuery := r.URL.Query()
	cursor := urlQuery.Get("cursor")
	limit := urlQuery.Get("limit")
	currentUser := library.GetUserIdFromJWT(r)

	if cursor == "" {
		cursor = "0"
	}
	if limit == "" {
		limit = "10"
	}

	intCursor, err := strconv.Atoi(cursor)
	if err != nil {
		intCursor = 0
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		intLimit = 10
	}

	posts := &[]Post{}

	if err := s.Store.ListPost(currentUser, int64(intCursor), int32(intLimit), posts); err != nil {
		log.Println("Error when getting listPost:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	var metaCursor int64

	if len(*posts) == 0 {
		metaCursor = 0
	} else {
		metaCursor = (*posts)[len(*posts)-1].CreatedAt
	}

	meta := struct {
		Cursor int64 `json:"cursor"`
	}{
		Cursor: metaCursor,
	}

	resp := library.NewResp("success", map[string]interface{}{
		"posts": posts,
		"meta":  meta,
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
		return http.StatusBadRequest, fmt.Errorf("invalid post detail")
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, post); err != nil {
		log.Println("Error when unmarshaling body:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid post detail")
	}

	fetchPost := &Post{}
	if err := s.Store.GetPostById(userId, postId, fetchPost); err != nil {
		log.Println("getPostById err:", err)
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("Post didnot exists")
		}

		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	if fetchPost.IdUser != userId {
		log.Println("Post userid did not match userid from jwt")
		return http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	if err := s.Store.UpdatePostBody(postId, post.Body, userId); err != nil {
		log.Println("Error when updating post body:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
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

	err := s.Store.GetPostById(userId, postId, post)
	if err != nil {
		log.Println("getPostById err:", err)
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("post didnot exists")
		}

		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	if userId != post.IdUser {
		log.Println("userid from jwt didnot match iduser post")
		return http.StatusForbidden, fmt.Errorf("Forbidden")
	}

	if err := s.Store.DeletePostById(postId, userId); err != nil {
		log.Println("Error when deleting post by id:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	resp := library.NewResp("Post deleted successfully", nil)

	library.WriteJson(w, http.StatusOK, resp)

	return http.StatusOK, nil
}

func (s *PostService) handleCreatePost(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Println("hit handle create post")

	idUser := library.GetUserIdFromJWT(r)
	postImage := ""
	userIn := &userProto.GetUserByIdReq{
		Id: idUser,
	}
	uuid := uuid.NewString()

	// ambil image dari formdata
	err := r.ParseMultipartForm(2 * 1024 * 1024)
	if err != nil {
		log.Println("Error when parsing request formdata:", err)
		return http.StatusBadRequest, fmt.Errorf("invalid formdata, or missing the required field")
	}

	reqBody := r.FormValue("reqBody")
	if reqBody == "" {
		log.Println("Error when getting reqBody, invalid/missing form data")
		return http.StatusBadRequest, fmt.Errorf("invalid/missing form data")
	}

	file, handler, err := r.FormFile("reqImage")
	if err != nil {
		log.Println("Error when creating file handler:", err)
	} else {
		defer file.Close()
		bytesFile, err := io.ReadAll(file)
		if err != nil {
			log.Println("Error when reading file from form data:", err)
			return http.StatusBadRequest, fmt.Errorf("invalid post detail")
		}

		imageIn := &imageProto.CreateImageReq{
			ImageFile: bytesFile,
			FileName:  handler.Filename,
		}

		imageGrpcResp, err := s.ImageServiceGrpcClient.CreateImage(r.Context(), imageIn)
		if err != nil {
			log.Println("Error when dialing image grpc client with CreateImage method:", err)
			return http.StatusInternalServerError, fmt.Errorf("something went wrong")
		}
		postImage = imageGrpcResp.GetFilename()
	}

	userGrpcResp, err := s.UserServiceGrpcClient.GetUserById(r.Context(), userIn)
	if err != nil {
		log.Println("Error when dialing grpc client with getUserById method:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	post := &Post{
		Id:       uuid,
		IdUser:   idUser,
		Username: userGrpcResp.GetUsername(),
		Name:     userGrpcResp.GetName(),
		Profile:  userGrpcResp.GetProfile(),
		Image:    postImage,
		Body:     reqBody,
	}

	/////////////////////////////
	//TODO validasi input user//
	///////////////////////////

	if err := s.Store.CreatePost(post.Id, post.Image, post.Body, post.IdUser, post.Username, post.Name, post.Profile); err != nil {
		log.Println("Error when creating post:", err)
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")

	}

	resp := library.NewResp("post created!", nil)
	library.WriteJson(w, http.StatusCreated, resp)

	return http.StatusCreated, nil
}
