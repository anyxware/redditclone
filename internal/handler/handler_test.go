package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	mock "redditclone/internal/handler/mocks"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"redditclone/pkg/cookie"
	"redditclone/pkg/token"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func initHandler(ctrl *gomock.Controller, service *mock.MockappService) *Handler {
	cookieStorage := cookie.NewMapStorage()
	sessions := cookie.NewManager(cookieStorage)
	signer := token.NewSigner("love")
	return NewHandler(signer, sessions, service)
}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("POST", "/register", strings.NewReader("{\"username\":\"van\",\"password\":\"qqq\"}")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().RegisterUser(model.Credential{Username: "van", Password: "qqq"})
				handler.signUp(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				matched, _ := regexp.MatchString("^{\"token\": \".*\"}$", string(body))
				return matched
			},
		},
		{
			request: httptest.NewRequest("POST", "/register", strings.NewReader("invalid json")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				handler.signUp(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"bad request\"}\n")
				return reflect.DeepEqual(body, data)
			},
		},
		{
			request: httptest.NewRequest("POST", "/register", strings.NewReader("{\"username\":\"van\"}")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				handler.signUp(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"body\",\"param\":\"password\",\"value\":\"\",\"msg\":\"field is required\"}]}\n")
				return reflect.DeepEqual(body, data)
			},
		},
		{
			request: httptest.NewRequest("POST", "/register", strings.NewReader("{\"username\":\"van\",\"password\":\"qqq\"}")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().
					RegisterUser(model.Credential{Username: "van", Password: "qqq"}).
					Return(model.User{}, customerr.UserAlreadyExists{Username: "van"})
				handler.signUp(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"body\",\"param\":\"username\",\"value\":\"van\",\"msg\":\"already exists\"}]}\n")
				return reflect.DeepEqual(body, data)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("POST", "/register", strings.NewReader("{\"username\":\"van\",\"password\":\"qqq\"}")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().LoginUser(model.Credential{Username: "van", Password: "qqq"})
				handler.signIn(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				matched, _ := regexp.MatchString("^{\"token\": \".*\"}$", string(body))
				return matched
			},
		},
		{
			request: httptest.NewRequest("POST", "/register", strings.NewReader("invalid json")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				handler.signIn(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"bad request\"}\n")
				return reflect.DeepEqual(body, data)
			},
		},
		{
			request: httptest.NewRequest("POST", "/register", strings.NewReader("{\"username\":\"van\"}")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				handler.signIn(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"body\",\"param\":\"password\",\"value\":\"\",\"msg\":\"field is required\"}]}\n")
				return reflect.DeepEqual(body, data)
			},
		},
		{
			request: httptest.NewRequest("POST", "/register", strings.NewReader("{\"username\":\"van\",\"password\":\"qqq\"}")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().
					LoginUser(model.Credential{Username: "van", Password: "qqq"}).
					Return(model.User{}, customerr.WrongCredential{Username: "van"})
				handler.signIn(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"wrong credential\"}\n")
				return reflect.DeepEqual(body, data)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestGetAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("GET", "/api/posts", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().GetAllPosts().Return([]model.Post{{ID: "1"}, {ID: "2"}, {ID: "3"}}, nil)
				handler.getAllPosts(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal([]model.Post{{ID: "1"}, {ID: "2"}, {ID: "3"}})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/posts", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().GetAllPosts().Return(nil, errors.New("internal error"))
				handler.getAllPosts(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"internal error\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestGetPostsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("GET", "/api/posts/category", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().GetPostsByCategory("funny").Return([]model.Post{{ID: "1"}, {ID: "2"}, {ID: "3"}}, nil)
				r = mux.SetURLVars(r, map[string]string{"category": "funny"})
				handler.getPostsByCategory(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal([]model.Post{{ID: "1"}, {ID: "2"}, {ID: "3"}})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/posts/category", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"category": "bad_category"})
				handler.getPostsByCategory(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte(
					"{\"errors\":[{\"location\":\"path\",\"param\":\"category\",\"value\":\"bad_category\",\"msg\":\"category must has specific type\"}]}\n",
				)
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/posts/category", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().GetPostsByCategory("funny").Return(nil, errors.New("internal error"))
				r = mux.SetURLVars(r, map[string]string{"category": "funny"})
				handler.getPostsByCategory(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"internal error\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestGetPostsByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("GET", "/api/user/username", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().GetPostsByAuthor("username").Return([]model.Post{{ID: "1"}, {ID: "2"}, {ID: "3"}}, nil)
				r = mux.SetURLVars(r, map[string]string{"username": "username"})
				handler.getPostsByUsername(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal([]model.Post{{ID: "1"}, {ID: "2"}, {ID: "3"}})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/user/username", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"username": ""})
				handler.getPostsByUsername(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"path\",\"param\":\"username\",\"value\":\"\",\"msg\":\"username must be a non-empty string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/user/username", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().GetPostsByAuthor("username").Return(nil, errors.New("internal error"))
				r = mux.SetURLVars(r, map[string]string{"username": "username"})
				handler.getPostsByUsername(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"internal error\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	textPostInput := model.TextPostInput{
		Type:     "text",
		Category: "funny",
		Title:    "bebe",
		Text:     "bubu",
	}
	textPostInputBody, _ := json.Marshal(textPostInput)

	urlPostInput := model.URLPostInput{
		Type:     "link",
		Category: "funny",
		Title:    "bebe",
		URL:      "bubu",
	}
	urlPostInputBody, _ := json.Marshal(urlPostInput)

	untypedPostInput := model.TextPostInput{
		Category: "funny",
		Title:    "bebe",
		Text:     "bubu",
	}
	untypedPostInputBody, _ := json.Marshal(untypedPostInput)

	textPostInputWithoutCategory := model.TextPostInput{
		Type:     "text",
		Category: "kek",
		Title:    "bebe",
		Text:     "bubu",
	}
	textPostInputWithoutCategoryBody, _ := json.Marshal(textPostInputWithoutCategory)

	invalidTextPostInput := model.TextPostInput{
		Category: "funny",
		Type:     "text",
		Title:    "bebe",
	}
	invalidTextPostInputBody, _ := json.Marshal(invalidTextPostInput)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("POST", "/api/posts", bytes.NewReader(textPostInputBody)),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.
					EXPECT().
					CreateTextPost(model.TextPostInput{
						Category: "funny",
						Type:     "text",
						Title:    "bebe",
						Text:     "bubu",
					}, model.User{ID: "1"}).Return(model.Post{ID: "1"}, nil)
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createPost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal(model.Post{ID: "1"})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("POST", "/api/posts", strings.NewReader("invalid json")),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createPost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"bad request\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("POST", "/api/posts", bytes.NewReader(urlPostInputBody)),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.
					EXPECT().
					CreateURLPost(model.URLPostInput{
						Category: "funny",
						Type:     "link",
						Title:    "bebe",
						URL:      "bubu",
					}, model.User{ID: "1"}).Return(model.Post{ID: "1"}, nil)
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createPost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal(model.Post{ID: "1"})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("POST", "/api/posts", bytes.NewReader(textPostInputWithoutCategoryBody)),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createPost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"body\",\"param\":\"category\",\"value\":\"kek\",\"msg\":\"category must has specific type\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("POST", "/api/posts", bytes.NewReader(untypedPostInputBody)),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createPost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"body\",\"param\":\"type\",\"value\":\"\",\"msg\":\"type must be a text or a link\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("POST", "/api/posts", bytes.NewReader(invalidTextPostInputBody)),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createPost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"body\",\"param\":\"text\",\"value\":\"\",\"msg\":\"text must be a non-empty string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("POST", "/api/posts", bytes.NewReader(textPostInputBody)),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.
					EXPECT().
					CreateTextPost(model.TextPostInput{
						Category: "funny",
						Type:     "text",
						Title:    "bebe",
						Text:     "bubu",
					}, model.User{ID: "1"}).Return(model.Post{}, errors.New("internal error"))
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createPost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"internal error\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestGetPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("GET", "/api/post/111111111111111111111111", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().GetPostByID("111111111111111111111111").Return(model.Post{ID: "111111111111111111111111"}, nil)
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				handler.getPost(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal(model.Post{ID: "111111111111111111111111"})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/post/1", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "1"})
				handler.getPost(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"path\",\"param\":\"post_id\",\"value\":\"1\",\"msg\":\"post_id must be a hexadecimal 24-symbols string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/post/111111111111111111111111", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().GetPostByID("111111111111111111111111").Return(model.Post{}, customerr.PostNotFoundByID{PostID: "111111111111111111111111"})
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				handler.getPost(w, r)
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"post not found\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestDeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("DELETE", "/api/post/111111111111111111111111", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().DeletePost("111111111111111111111111", model.User{ID: "1"}).Return(nil)
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.deletePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\": \"success\"}")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("DELETE", "/api/post/1", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "1"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.deletePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"path\",\"param\":\"post_id\",\"value\":\"1\",\"msg\":\"post_id must be a hexadecimal 24-symbols string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("DELETE", "/api/post/111111111111111111111111", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().DeletePost("111111111111111111111111", model.User{ID: "1"}).Return(customerr.PostNotFoundByID{PostID: "111111111111111111111111"})
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.deletePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"post not found\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestCreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest(
				"POST",
				"/api/post/111111111111111111111111",
				bytes.NewReader([]byte("{\"comment\": \"comment\"}")),
			),
			writer: httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().
					AddComment("111111111111111111111111", "comment", model.User{ID: "1"}).
					Return(model.Post{ID: "111111111111111111111111"}, nil)
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createComment(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal(model.Post{ID: "111111111111111111111111"})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest(
				"POST",
				"/api/post/1",
				bytes.NewReader([]byte("{\"comment\": \"comment\"}")),
			),
			writer: httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "1"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createComment(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"path\",\"param\":\"post_id\",\"value\":\"1\",\"msg\":\"post_id must be a hexadecimal 24-symbols string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest(
				"POST",
				"/api/post/111111111111111111111111",
				bytes.NewReader([]byte("invalid json")),
			),
			writer: httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createComment(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"bad request\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest(
				"POST",
				"/api/post/111111111111111111111111",
				bytes.NewReader([]byte("{\"comment\": \"\"}")),
			),
			writer: httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createComment(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"body\",\"param\":\"comment\",\"value\":\"\",\"msg\":\"comment must be a non-empty string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest(
				"POST",
				"/api/post/111111111111111111111111",
				bytes.NewReader([]byte("{\"comment\": \"comment\"}")),
			),
			writer: httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().
					AddComment("111111111111111111111111", "comment", model.User{ID: "1"}).
					Return(model.Post{}, customerr.PostNotFoundByID{PostID: "111111111111111111111111"})
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.createComment(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"post not found\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest(
				"DELETE",
				"/api/post/111111111111111111111111/111111111111111111111111",
				nil,
			),
			writer: httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().DeleteComment(
					"111111111111111111111111",
					"111111111111111111111111",
					model.User{ID: "1"},
				).Return(model.Post{ID: "111111111111111111111111"}, nil)
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111", "comment_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.deleteComment(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal(model.Post{ID: "111111111111111111111111"})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest(
				"DELETE",
				"/api/post/111111111111111111111111/1",
				nil,
			),
			writer: httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111", "comment_id": "1"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.deleteComment(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"path\",\"param\":\"comment_id\",\"value\":\"1\",\"msg\":\"comment_id must be a hexadecimal 24-symbols string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest(
				"DELETE",
				"/api/post/111111111111111111111111/111111111111111111111111",
				nil,
			),
			writer: httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().DeleteComment(
					"111111111111111111111111",
					"111111111111111111111111",
					model.User{ID: "1"},
				).Return(model.Post{}, customerr.CommentNotFoundByID{PostID: "111111111111111111111111", CommentID: "111111111111111111111111"})
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111", "comment_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.deleteComment(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"comment not found\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestUpvotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("GET", "/api/post/111111111111111111111111/upvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().UpvotePost(
					"111111111111111111111111",
					model.User{ID: "1"},
				).Return(model.Post{ID: "111111111111111111111111"}, nil)
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.upvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal(model.Post{ID: "111111111111111111111111"})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/post/1/upvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "1"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.upvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"path\",\"param\":\"post_id\",\"value\":\"1\",\"msg\":\"post_id must be a hexadecimal 24-symbols string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/post/111111111111111111111111/upvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().UpvotePost(
					"111111111111111111111111",
					model.User{ID: "1"},
				).Return(model.Post{}, customerr.PostNotFoundByID{PostID: "111111111111111111111111"})
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.upvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"post not found\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestDownvotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("GET", "/api/post/111111111111111111111111/downvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().DownvotePost(
					"111111111111111111111111",
					model.User{ID: "1"},
				).Return(model.Post{ID: "111111111111111111111111"}, nil)
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.downvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal(model.Post{ID: "111111111111111111111111"})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/post/1/downvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "1"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.downvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"path\",\"param\":\"post_id\",\"value\":\"1\",\"msg\":\"post_id must be a hexadecimal 24-symbols string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/post/111111111111111111111111/downvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().DownvotePost(
					"111111111111111111111111",
					model.User{ID: "1"},
				).Return(model.Post{}, customerr.PostNotFoundByID{PostID: "111111111111111111111111"})
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.downvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"post not found\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}

func TestUnvotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockappService(ctrl)
	handler := initHandler(ctrl, service)

	cases := []struct {
		request *http.Request
		writer  *httptest.ResponseRecorder
		run     func(w *httptest.ResponseRecorder, r *http.Request) *http.Response
		check   func(body []byte) bool
	}{
		{
			request: httptest.NewRequest("GET", "/api/post/111111111111111111111111/unvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().UnvotePost(
					"111111111111111111111111",
					model.User{ID: "1"},
				).Return(model.Post{ID: "111111111111111111111111"}, nil)
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.unvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data, _ := json.Marshal(model.Post{ID: "111111111111111111111111"})
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/post/1/unvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				r = mux.SetURLVars(r, map[string]string{"post_id": "1"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.unvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"errors\":[{\"location\":\"path\",\"param\":\"post_id\",\"value\":\"1\",\"msg\":\"post_id must be a hexadecimal 24-symbols string\"}]}\n")
				return reflect.DeepEqual(data, body)
			},
		},
		{
			request: httptest.NewRequest("GET", "/api/post/111111111111111111111111/unvote", nil),
			writer:  httptest.NewRecorder(),
			run: func(w *httptest.ResponseRecorder, r *http.Request) *http.Response {
				service.EXPECT().UnvotePost(
					"111111111111111111111111",
					model.User{ID: "1"},
				).Return(model.Post{}, customerr.PostNotFoundByID{PostID: "111111111111111111111111"})
				r = mux.SetURLVars(r, map[string]string{"post_id": "111111111111111111111111"})
				ctx := context.WithValue(r.Context(), "user", model.User{ID: "1"})
				handler.unvotePost(w, r.WithContext(ctx))
				return w.Result()
			},
			check: func(body []byte) bool {
				data := []byte("{\"message\":\"post not found\"}\n")
				return reflect.DeepEqual(data, body)
			},
		},
	}

	for i, item := range cases {
		resp := item.run(item.writer, item.request)
		body, _ := ioutil.ReadAll(resp.Body)
		if !item.check(body) {
			t.Errorf("[%d] unexpected body: %s", i, string(body))
		}
	}
}
