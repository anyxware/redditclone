package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"redditclone/internal/model"
	"redditclone/pkg/cookie"
	"redditclone/pkg/httpvalidator"
	"redditclone/pkg/token"
)

const (
	authorizationHeader = "Authorization"
)

type authService interface {
	RegisterUser(cred model.Credential) (model.User, error)
	LoginUser(cred model.Credential) (model.User, error)
}

type postsService interface {
	GetAllPosts() ([]model.Post, error)
	GetPostsByCategory(category string) ([]model.Post, error)
	GetPostsByAuthor(username string) ([]model.Post, error)
	CreateTextPost(input model.TextPostInput, usr model.User) (model.Post, error)
	CreateURLPost(input model.URLPostInput, usr model.User) (model.Post, error)
	GetPostByID(postID string) (model.Post, error)
	DeletePost(postID string, usr model.User) error
	AddComment(postID string, commentText string, usr model.User) (model.Post, error)
	DeleteComment(postID, commentID string, usr model.User) (model.Post, error)
	UpvotePost(postID string, usr model.User) (model.Post, error)
	DownvotePost(postID string, usr model.User) (model.Post, error)
	UnvotePost(postID string, usr model.User) (model.Post, error)
}

type usersService interface {
	GetUserByID(userID string) (model.User, error)
}

type appService interface {
	authService
	postsService
	usersService
}

type Handler struct {
	sessions  cookie.Manager
	signer    token.Signer
	validator httpvalidator.Validator
	service   appService
}

func NewHandler(signer token.Signer, sessions cookie.Manager, service appService) *Handler {
	validator := httpvalidator.NewValidator()
	handler := &Handler{signer: signer, sessions: sessions, validator: validator, service: service}
	handler.initValidator()
	return handler
}

func (h *Handler) CreateRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/hamster.jpg")
	})
	fs := http.FileServer(http.Dir("./static/html/"))
	router.Handle("/", fs)
	fs = http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.HandleFunc("/api/register", h.signUp).Methods("POST")
	router.HandleFunc("/api/login", h.signIn).Methods("POST")

	router.HandleFunc("/api/posts/", h.getAllPosts).Methods("GET")
	router.HandleFunc("/api/posts/{category}", h.getPostsByCategory).Methods("GET")
	router.HandleFunc("/api/post/{post_id}", h.getPost).Methods("GET")

	routerForAuthorized := router.PathPrefix("/api").Subrouter()
	routerForAuthorized.Use(h.authorizeMiddleware)
	routerForAuthorized.HandleFunc("/posts", h.createPost).Methods("POST")
	routerForAuthorized.HandleFunc("/post/{post_id}", h.deletePost).Methods("DELETE")
	routerForAuthorized.HandleFunc("/post/{post_id}", h.createComment).Methods("POST")
	routerForAuthorized.HandleFunc("/post/{post_id}/{comment_id}", h.deleteComment).Methods("DELETE")
	routerForAuthorized.HandleFunc("/post/{post_id}/upvote", h.upvotePost).Methods("GET")
	routerForAuthorized.HandleFunc("/post/{post_id}/downvote", h.downvotePost).Methods("GET")
	routerForAuthorized.HandleFunc("/post/{post_id}/unvote", h.unvotePost).Methods("GET")

	router.HandleFunc("/api/user/{username}", h.getPostsByUsername).Methods("GET")

	router.UseEncodedPath().NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/index.html")
	})

	router.Use(h.recoverPanicMiddleware)
	router.Use(h.accessLogMiddleware)

	return router
}
