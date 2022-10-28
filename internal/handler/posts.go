package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
)

func writePosts(w http.ResponseWriter, posts []model.Post) error {
	resp, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	if _, err = w.Write(resp); err != nil {
		return err
	}

	return nil
}

func writePost(w http.ResponseWriter, p model.Post) error {
	resp, err := json.Marshal(p)
	if err != nil {
		return err
	}

	if _, err = w.Write(resp); err != nil {
		return err
	}

	return nil
}

func (h *Handler) getAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAllPosts()
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePosts(w, posts); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getPostsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	if errs := h.validator.ValidatePathValue("category", category); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	posts, err := h.service.GetPostsByCategory(category)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePosts(w, posts); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getPostsByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if errs := h.validator.ValidatePathValue("username", username); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	posts, err := h.service.GetPostsByAuthor(username)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePosts(w, posts); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(model.User)

	var input map[string]string
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.handleError(w, customerr.RequestNotParsed{Message: err.Error()})
		return
	}

	if errs := h.validator.ValidateBody("PostInput", input); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	var validationSchema string
	if input["type"] == "link" {
		validationSchema = "URLPostInput"
	} else if input["type"] == "text" {
		validationSchema = "TextPostInput"
	}

	if errs := h.validator.ValidateBody(validationSchema, input); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	var newPost model.Post
	if input["type"] == "link" {
		newPost, err = h.service.CreateURLPost(model.URLPostInput{
			Category: input["category"],
			Type:     input["type"],
			Title:    input["title"],
			URL:      input["url"],
		}, usr)
	} else if input["type"] == "text" {
		newPost, err = h.service.CreateTextPost(model.TextPostInput{
			Category: input["category"],
			Type:     input["type"],
			Title:    input["title"],
			Text:     input["text"],
		}, usr)
	}

	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePost(w, newPost); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	if errs := h.validator.ValidatePathValue("post_id", postID); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	existedPost, err := h.service.GetPostByID(postID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePost(w, existedPost); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(model.User)

	vars := mux.Vars(r)
	postID := vars["post_id"]

	if errs := h.validator.ValidatePathValue("post_id", postID); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	if err := h.service.DeletePost(postID, usr); err != nil {
		h.handleError(w, err)
		return
	}

	resp := []byte(fmt.Sprintf("{\"message\": \"success\"}"))
	if _, err := w.Write(resp); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(model.User)

	vars := mux.Vars(r)
	postID := vars["post_id"]

	if errs := h.validator.ValidatePathValue("post_id", postID); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	var input map[string]string
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.handleError(w, customerr.RequestNotParsed{Message: err.Error()})
		return
	}

	if errs := h.validator.ValidateBody("Comment", input); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	comment := input["comment"]

	existedPost, err := h.service.AddComment(postID, comment, usr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePost(w, existedPost); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(model.User)

	vars := mux.Vars(r)
	postID := vars["post_id"]
	commentID := vars["comment_id"]

	postIDValidationErrs := h.validator.ValidatePathValue("post_id", postID)
	commentIDValidationErrs := h.validator.ValidatePathValue("comment_id", commentID)
	if errs := append(postIDValidationErrs, commentIDValidationErrs...); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	existedPost, err := h.service.DeleteComment(postID, commentID, usr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePost(w, existedPost); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) upvotePost(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(model.User)

	vars := mux.Vars(r)
	postID := vars["post_id"]

	if errs := h.validator.ValidatePathValue("post_id", postID); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	existedPost, err := h.service.UpvotePost(postID, usr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePost(w, existedPost); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) downvotePost(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(model.User)

	vars := mux.Vars(r)
	postID := vars["post_id"]

	if errs := h.validator.ValidatePathValue("post_id", postID); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	existedPost, err := h.service.DownvotePost(postID, usr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePost(w, existedPost); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) unvotePost(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(model.User)

	vars := mux.Vars(r)
	postID := vars["post_id"]

	if errs := h.validator.ValidatePathValue("post_id", postID); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	existedPost, err := h.service.UnvotePost(postID, usr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = writePost(w, existedPost); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
