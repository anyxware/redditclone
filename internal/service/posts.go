package service

import (
	"github.com/sirupsen/logrus"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"redditclone/pkg/hexid"
)

type postsRepo interface {
	GetAllPosts() ([]model.Post, error)
	GetPostsByCategory(category string) ([]model.Post, error)
	GetPostsByAuthor(username string) ([]model.Post, error)
	AddPost(newPost model.Post) error
	GetPostByID(postID string) (model.Post, error)
	GetPostByIDAndUpdateViews(postID string) (model.Post, error)
	DeletePost(postID string) error
	AddComment(postID string, comment model.Comment) (model.Post, error)
	GetCommentByID(postID, commentID string) (model.Comment, error)
	DeleteComment(postID, commentID string) (model.Post, error)
	UpdateVotes(postID string, score int, upvotePercentage int, votes []model.Vote) error
}

func (s *service) GetAllPosts() ([]model.Post, error) {
	return s.postsRepo.GetAllPosts()
}

func (s *service) CreateTextPost(input model.TextPostInput, usr model.User) (model.Post, error) {
	postID, err := hexid.Generate()
	if err != nil {
		return model.Post{}, err
	}

	post := model.NewTextPost(postID, input, model.Author{ID: usr.ID, Username: usr.Username})
	if err = s.postsRepo.AddPost(post); err != nil {
		return model.Post{}, err
	}

	logrus.Infoln("new text post created")

	return post, nil
}

func (s *service) CreateURLPost(input model.URLPostInput, usr model.User) (model.Post, error) {
	postID, err := hexid.Generate()
	if err != nil {
		return model.Post{}, err
	}

	post := model.NewURLPost(postID, input, model.Author{ID: usr.ID, Username: usr.Username})
	if err = s.postsRepo.AddPost(post); err != nil {
		return model.Post{}, err
	}

	logrus.Infoln("new url post created")

	return post, nil
}

func (s *service) GetPostsByCategory(category string) ([]model.Post, error) {
	return s.postsRepo.GetPostsByCategory(category)
}

func (s *service) GetPostsByAuthor(username string) ([]model.Post, error) {
	return s.postsRepo.GetPostsByAuthor(username)
}

func (s *service) GetPostByID(postID string) (model.Post, error) {
	return s.postsRepo.GetPostByIDAndUpdateViews(postID)
}

func (s *service) DeletePost(postID string, usr model.User) error {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	post, err := s.postsRepo.GetPostByID(postID)
	if err != nil {
		return err
	}

	if post.Author.ID != usr.ID {
		return customerr.NotOwner{Username: usr.Username}
	}

	if err = s.postsRepo.DeletePost(postID); err != nil {
		return err
	}

	logrus.Infoln("post deleted")

	return nil
}

func (s *service) AddComment(postID string, commentText string, usr model.User) (model.Post, error) {
	commentID, err := hexid.Generate()
	if err != nil {
		return model.Post{}, err
	}

	comment := model.NewComment(commentID, commentText, model.Author{ID: usr.ID, Username: usr.Username})
	post, err := s.postsRepo.AddComment(postID, comment)
	if err != nil {
		return model.Post{}, err
	}

	logrus.Infoln("comment added")

	return post, nil
}

func (s *service) DeleteComment(postID, commentID string, usr model.User) (model.Post, error) {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	comment, err := s.postsRepo.GetCommentByID(postID, commentID)
	if err != nil {
		return model.Post{}, err
	}

	if comment.Author.ID != usr.ID {
		return model.Post{}, customerr.NotOwner{Username: usr.Username}
	}

	post, err := s.postsRepo.DeleteComment(postID, commentID)
	if err != nil {
		return model.Post{}, err
	}

	logrus.Infoln("comment deleted")

	return post, nil
}

func (s *service) UpvotePost(postID string, usr model.User) (model.Post, error) {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	post, err := s.postsRepo.GetPostByID(postID)
	if err != nil {
		return model.Post{}, err
	}

	post.Upvote(usr.ID).RecalculatePercentage()

	if err = s.postsRepo.UpdateVotes(postID, post.Score, post.UpvotePercentage, post.Votes); err != nil {
		logrus.Errorln(err)
		return model.Post{}, err
	}

	logrus.Infoln("post upvoted")

	return post, nil
}

func (s *service) DownvotePost(postID string, usr model.User) (model.Post, error) {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	post, err := s.postsRepo.GetPostByID(postID)
	if err != nil {
		return model.Post{}, err
	}

	post.Downvote(usr.ID).RecalculatePercentage()

	if err = s.postsRepo.UpdateVotes(postID, post.Score, post.UpvotePercentage, post.Votes); err != nil {
		return model.Post{}, err
	}

	logrus.Infoln("post downvoted")

	return post, nil
}

func (s *service) UnvotePost(postID string, usr model.User) (model.Post, error) {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	post, err := s.postsRepo.GetPostByID(postID)
	if err != nil {
		return model.Post{}, err
	}

	post.Unvote(usr.ID).RecalculatePercentage()

	if err = s.postsRepo.UpdateVotes(postID, post.Score, post.UpvotePercentage, post.Votes); err != nil {
		return model.Post{}, err
	}

	logrus.Infoln("post unvoted")

	return post, nil
}
