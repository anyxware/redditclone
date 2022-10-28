package mongorepo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"reflect"
	"testing"
)

// go test -v -coverprofile=cover.out && go tool cover -html=cover.out -o coverage.html; rm cover.out

func marshalPost(post model.Post) bson.D {
	var bsonData []byte
	bsonData, _ = bson.Marshal(post)

	var bsonD bson.D
	_ = bson.Unmarshal(bsonData, &bsonD)

	return bsonD
}

func marshalPosts(posts []model.Post) []bson.D {
	docs := make([]bson.D, 0)

	for _, post := range posts {
		bsonData, _ := bson.Marshal(post)
		var bsonD bson.D
		_ = bson.Unmarshal(bsonData, &bsonD)
		docs = append(docs, bsonD)
	}

	return docs
}

func compareErrorsMsg(err1 error, err2 error) bool {
	return fmt.Sprint(err1) == fmt.Sprint(err2)
}

func TestGetAllPosts(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedPosts []model.Post
		expectedErr   error
		run           func([]model.Post) ([]model.Post, error)
	}{
		{
			expectedPosts: []model.Post{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			expectedErr:   nil,
			run: func(expectedPosts []model.Post) ([]model.Post, error) {
				var posts []model.Post
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					docs := marshalPosts(expectedPosts)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch, docs...),
						mtest.CreateCursorResponse(0, "redditclone.posts", mtest.NextBatch),
					)
					posts, err = repo.GetAllPosts()
				})
				return posts, err
			},
		},
		{
			expectedPosts: make([]model.Post, 0),
			expectedErr:   mongo.CommandError{Message: "command failed"},
			run: func(expectedPosts []model.Post) ([]model.Post, error) {
				var posts []model.Post
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					posts, err = repo.GetAllPosts()
				})
				return posts, err
			},
		},
	}

	for i, item := range cases {
		posts, err := item.run(item.expectedPosts)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		for j, expectedPost := range item.expectedPosts {
			if !reflect.DeepEqual(expectedPost, posts[j]) {
				t.Errorf("[%d:%d] expected post: %+v, got: %+v", i, j, expectedPost, posts[j])
			}
		}
	}
}

func TestPostsByCategory(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedPosts []model.Post
		expectedErr   error
		run           func([]model.Post) ([]model.Post, error)
	}{
		{
			expectedPosts: []model.Post{
				{ID: "1", Category: "funny"},
				{ID: "2", Category: "funny"},
				{ID: "3", Category: "funny"},
			},
			expectedErr: nil,
			run: func(expectedPosts []model.Post) ([]model.Post, error) {
				var posts []model.Post
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					docs := marshalPosts(expectedPosts)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch, docs...),
						mtest.CreateCursorResponse(0, "redditclone.posts", mtest.NextBatch),
					)
					posts, err = repo.GetPostsByCategory("funny")
				})
				return posts, err
			},
		},
		{
			expectedPosts: make([]model.Post, 0),
			expectedErr:   mongo.CommandError{Message: "command failed"},
			run: func(expectedPosts []model.Post) ([]model.Post, error) {
				var posts []model.Post
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					posts, err = repo.GetPostsByCategory("funny")
				})
				return posts, err
			},
		},
	}

	for i, item := range cases {
		posts, err := item.run(item.expectedPosts)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		for j, expectedPost := range item.expectedPosts {
			if !reflect.DeepEqual(expectedPost, posts[j]) {
				t.Errorf("[%d:%d] expected post: %+v, got: %+v", i, j, expectedPost, posts[j])
			}
		}
	}
}

func TestPostsByAuthor(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedPosts []model.Post
		expectedErr   error
		run           func([]model.Post) ([]model.Post, error)
	}{
		{
			expectedPosts: []model.Post{
				{ID: "1", Author: model.Author{ID: "1", Username: "ivan"}},
				{ID: "2", Author: model.Author{ID: "1", Username: "ivan"}},
				{ID: "3", Author: model.Author{ID: "1", Username: "ivan"}},
			},
			expectedErr: nil,
			run: func(expectedPosts []model.Post) ([]model.Post, error) {
				var posts []model.Post
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					docs := marshalPosts(expectedPosts)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch, docs...),
						mtest.CreateCursorResponse(0, "redditclone.posts", mtest.NextBatch),
					)
					posts, err = repo.GetPostsByAuthor("ivan")
				})
				return posts, err
			},
		},
		{
			expectedPosts: make([]model.Post, 0),
			expectedErr:   mongo.CommandError{Message: "command failed"},
			run: func(expectedPosts []model.Post) ([]model.Post, error) {
				var posts []model.Post
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					posts, err = repo.GetPostsByAuthor("ivan")
				})
				return posts, err
			},
		},
	}

	for i, item := range cases {
		posts, err := item.run(item.expectedPosts)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		for j, expectedPost := range item.expectedPosts {
			if !reflect.DeepEqual(expectedPost, posts[j]) {
				t.Errorf("[%d:%d] expected post: %+v, got: %+v", i, j, expectedPost, posts[j])
			}
		}
	}
}

func TestAddPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedErr error
		run         func() error
	}{
		{
			expectedErr: nil,
			run: func() error {
				post := model.Post{}
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(mtest.CreateSuccessResponse())
					err = repo.AddPost(post)
				})
				return err
			},
		},
		{
			expectedErr: mongo.CommandError{Message: "command failed"},
			run: func() error {
				post := model.Post{}
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					err = repo.AddPost(post)
				})
				return err
			},
		},
	}

	for i, item := range cases {
		err := item.run()
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
	}
}

func TestGetPostByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedPost model.Post
		expectedErr  error
		run          func(model.Post) (model.Post, error)
	}{
		{
			expectedPost: model.Post{ID: "1"},
			expectedErr:  nil,
			run: func(expectedPost model.Post) (model.Post, error) {
				var post model.Post
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					doc := marshalPost(expectedPost)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch, doc),
					)
					post, err = repo.GetPostByID("1")
				})
				return post, err
			},
		},
		{
			expectedPost: model.Post{},
			expectedErr:  mongo.CommandError{Message: "command failed"},
			run: func(expectedPost model.Post) (model.Post, error) {
				var post model.Post
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					post, err = repo.GetPostByID("1")
				})
				return post, err
			},
		},
		{
			expectedPost: model.Post{},
			expectedErr:  customerr.PostNotFoundByID{PostID: "1"},
			run: func(expectedPost model.Post) (model.Post, error) {
				var post model.Post
				var err error
				mt.Run("post not found", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch),
						mtest.CreateCursorResponse(0, "redditclone.posts", mtest.NextBatch),
					)
					post, err = repo.GetPostByID("1")
				})
				return post, err
			},
		},
	}

	for i, item := range cases {
		post, err := item.run(item.expectedPost)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		if !reflect.DeepEqual(item.expectedPost, post) {
			t.Errorf("[%d] expected post: %+v, got: %+v", i, item.expectedPost, post)
		}
	}
}

func TestGetPostByIDAndUpdateViews(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedPost model.Post
		expectedErr  error
		run          func(post model.Post) (model.Post, error)
	}{
		{
			expectedPost: model.Post{ID: "1", Views: 1},
			expectedErr:  nil,
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					doc := marshalPost(post)
					mt.AddMockResponses(
						mtest.CreateSuccessResponse(bson.E{Key: "value", Value: doc}),
					)
					post, err = repo.GetPostByIDAndUpdateViews("1")
				})
				return post, err
			},
		},
		{
			expectedPost: model.Post{},
			expectedErr:  mongo.CommandError{Message: "command failed"},
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					post, err = repo.GetPostByIDAndUpdateViews("1")
				})
				return post, err
			},
		},
		{
			expectedPost: model.Post{},
			expectedErr:  customerr.PostNotFoundByID{PostID: "1"},
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("post not found", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch),
						mtest.CreateCursorResponse(0, "redditclone.posts", mtest.NextBatch),
					)
					post, err = repo.GetPostByIDAndUpdateViews("1")
				})
				return post, err
			},
		},
	}

	for i, item := range cases {
		post, err := item.run(item.expectedPost)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		if !reflect.DeepEqual(item.expectedPost, post) {
			t.Errorf("[%d] expected post: %+v, got: %+v", i, item.expectedPost, post)
		}
	}
}

func TestDeletePost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedErr error
		run         func() error
	}{
		{
			expectedErr: nil,
			run: func() error {
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(
						mtest.CreateSuccessResponse(),
					)
					err = repo.DeletePost("1")
				})
				return err
			},
		},
		{
			expectedErr: mongo.CommandError{Message: "command failed"},
			run: func() error {
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					err = repo.DeletePost("1")
				})
				return err
			},
		},
	}

	for i, item := range cases {
		err := item.run()
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
	}
}

func TestAddComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedPost model.Post
		expectedErr  error
		run          func(post model.Post) (model.Post, error)
	}{
		{
			expectedPost: model.Post{ID: "1", Comments: []model.Comment{{ID: "1"}}},
			expectedErr:  nil,
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					doc := marshalPost(post)
					mt.AddMockResponses(
						mtest.CreateSuccessResponse(bson.E{Key: "value", Value: doc}),
					)
					post, err = repo.AddComment("1", model.Comment{ID: "1"})
				})
				return post, err
			},
		},
		{
			expectedPost: model.Post{},
			expectedErr:  mongo.CommandError{Message: "command failed"},
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					post, err = repo.AddComment("1", model.Comment{ID: "1"})
				})
				return post, err
			},
		},
		{
			expectedPost: model.Post{},
			expectedErr:  customerr.PostNotFoundByID{PostID: "1"},
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("post not found", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch),
						mtest.CreateCursorResponse(0, "redditclone.posts", mtest.NextBatch),
					)
					post, err = repo.AddComment("1", model.Comment{ID: "1"})
				})
				return post, err
			},
		},
	}

	for i, item := range cases {
		post, err := item.run(item.expectedPost)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		if !reflect.DeepEqual(item.expectedPost, post) {
			t.Errorf("[%d] expected post: %+v, got: %+v", i, item.expectedPost, post)
		}
	}
}

func TestGetCommentByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedComment model.Comment
		expectedErr     error
		run             func(comment model.Comment) (model.Comment, error)
	}{
		{
			expectedComment: model.Comment{ID: "1"},
			expectedErr:     nil,
			run: func(comment model.Comment) (model.Comment, error) {
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					doc := marshalPost(model.Post{ID: "1", Comments: []model.Comment{{ID: "1"}}})
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch, doc),
					)
					comment, err = repo.GetCommentByID("1", "1")
				})
				return comment, err
			},
		},
		{
			expectedComment: model.Comment{},
			expectedErr:     mongo.CommandError{Message: "command failed"},
			run: func(comment model.Comment) (model.Comment, error) {
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					comment, err = repo.GetCommentByID("1", "1")
				})
				return comment, err
			},
		},
		{
			expectedComment: model.Comment{},
			expectedErr:     customerr.PostNotFoundByID{PostID: "1"},
			run: func(comment model.Comment) (model.Comment, error) {
				var err error
				mt.Run("post not found", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch),
						mtest.CreateCursorResponse(0, "redditclone.posts", mtest.NextBatch),
					)
					comment, err = repo.GetCommentByID("1", "1")
				})
				return comment, err
			},
		},
		{
			expectedComment: model.Comment{},
			expectedErr:     customerr.CommentNotFoundByID{PostID: "1", CommentID: "1"},
			run: func(comment model.Comment) (model.Comment, error) {
				var err error
				mt.Run("comment not found", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					doc := marshalPost(model.Post{ID: "1"})
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch, doc),
					)
					comment, err = repo.GetCommentByID("1", "1")
				})
				return comment, err
			},
		},
	}

	for i, item := range cases {
		comment, err := item.run(item.expectedComment)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		if !reflect.DeepEqual(item.expectedComment, comment) {
			t.Errorf("[%d] expected post: %+v, got: %+v", i, item.expectedComment, comment)
		}
	}
}

func TestDeleteComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedPost model.Post
		expectedErr  error
		run          func(post model.Post) (model.Post, error)
	}{
		{
			expectedPost: model.Post{ID: "1"},
			expectedErr:  nil,
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					doc := marshalPost(post)
					mt.AddMockResponses(
						mtest.CreateSuccessResponse(bson.E{Key: "value", Value: doc}),
					)
					post, err = repo.DeleteComment("1", "1")
				})
				return post, err
			},
		},
		{
			expectedPost: model.Post{},
			expectedErr:  mongo.CommandError{Message: "command failed"},
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					post, err = repo.DeleteComment("1", "1")
				})
				return post, err
			},
		},
		{
			expectedPost: model.Post{},
			expectedErr:  customerr.PostNotFoundByID{PostID: "1"},
			run: func(post model.Post) (model.Post, error) {
				var err error
				mt.Run("post not found", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1, "redditclone.posts", mtest.FirstBatch),
						mtest.CreateCursorResponse(0, "redditclone.posts", mtest.NextBatch),
					)
					post, err = repo.DeleteComment("1", "1")
				})
				return post, err
			},
		},
	}

	for i, item := range cases {
		post, err := item.run(item.expectedPost)
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
		if !reflect.DeepEqual(item.expectedPost, post) {
			t.Errorf("[%d] expected post: %+v, got: %+v", i, item.expectedPost, post)
		}
	}
}

func TestUpdateVotes(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	cases := []struct {
		expectedErr error
		run         func() error
	}{
		{
			expectedErr: nil,
			run: func() error {
				var err error
				mt.Run("success", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(
						mtest.CreateSuccessResponse(),
					)
					err = repo.UpdateVotes("1", 1, 100, []model.Vote{})
				})
				return err
			},
		},
		{
			expectedErr: mongo.CommandError{Message: "command failed"},
			run: func() error {
				var err error
				mt.Run("command failed", func(mt *mtest.T) {
					repo := NewPostsRepo(mt.Coll)
					mt.AddMockResponses(bson.D{{"ok", 0}})
					err = repo.UpdateVotes("1", 1, 100, []model.Vote{})
				})
				return err
			},
		},
	}

	for i, item := range cases {
		err := item.run()
		if !compareErrorsMsg(item.expectedErr, err) {
			t.Errorf("[%d] expected error: %s, got: %s", i, item.expectedErr, err)
		}
	}
}
