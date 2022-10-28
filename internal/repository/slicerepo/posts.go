package slicerepo

import (
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"sync"
)

type postsRepo struct {
	mutex sync.RWMutex
	posts []model.Post
}

func NewPostsRepo() *postsRepo {
	return &postsRepo{
		posts: make([]model.Post, 0),
	}
}

func (r *postsRepo) GetAllPosts() ([]model.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.posts, nil
}

func (r *postsRepo) GetPostsByCategory(category string) ([]model.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	posts := make([]model.Post, 0)
	for _, existedPost := range r.posts {
		if existedPost.Category == category {
			posts = append(posts, existedPost)
		}
	}
	return posts, nil
}

func (r *postsRepo) GetPostsByAuthor(username string) ([]model.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	posts := make([]model.Post, 0)
	for _, existedPost := range r.posts {
		if existedPost.Category == username {
			posts = append(posts, existedPost)
		}
	}

	return posts, nil
}

func (r *postsRepo) AddPost(newPost model.Post) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.posts = append(r.posts, newPost)

	return nil
}

func (r *postsRepo) GetPostByID(postID string) (model.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, post := range r.posts {
		if post.ID == postID {
			return post, nil
		}
	}

	return model.Post{}, customerr.PostNotFoundByID{PostID: postID}
}

func (r *postsRepo) GetPostByIDAndUpdateViews(postID string) (model.Post, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, post := range r.posts {
		if post.ID == postID {
			r.posts[i].Views += 1
			return r.posts[i], nil
		}
	}

	return model.Post{}, customerr.PostNotFoundByID{PostID: postID}
}

func (r *postsRepo) DeletePost(postID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for idx, existedPost := range r.posts {
		if existedPost.ID == postID {
			r.posts = append(r.posts[:idx], r.posts[idx+1:]...)
			return nil
		}
	}

	return customerr.PostNotFoundByID{PostID: postID}
}

func (r *postsRepo) AddComment(postID string, comment model.Comment) (model.Post, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for idx, existedPost := range r.posts {
		if existedPost.ID == postID {
			r.posts[idx].Comments = append(existedPost.Comments, comment)
			return r.posts[idx], nil
		}
	}

	return model.Post{}, customerr.PostNotFoundByID{PostID: postID}
}

func (r *postsRepo) GetCommentByID(postID, commentID string) (model.Comment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, existedPost := range r.posts {
		if existedPost.ID == postID {
			for _, existedComment := range existedPost.Comments {
				if existedComment.ID == commentID {
					return existedComment, nil
				}
			}
			return model.Comment{}, customerr.CommentNotFoundByID{CommentID: commentID, PostID: postID}
		}
	}

	return model.Comment{}, customerr.PostNotFoundByID{PostID: postID}
}

func (r *postsRepo) DeleteComment(postID, commentID string) (model.Post, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for idx, existedPost := range r.posts {
		if existedPost.ID == postID {
			for i, comment := range existedPost.Comments {
				if comment.ID == commentID {
					r.posts[idx].Comments = append(existedPost.Comments[:i], existedPost.Comments[i+1:]...)
					return r.posts[idx], nil
				}
			}
			return model.Post{}, customerr.CommentNotFoundByID{CommentID: commentID, PostID: postID}
		}
	}

	return model.Post{}, customerr.PostNotFoundByID{PostID: postID}
}

func (r *postsRepo) UpdateVotes(postID string, score int, upvotePercentage int, votes []model.Vote) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, post := range r.posts {
		if post.ID == postID {
			r.posts[i].Score = score
			r.posts[i].UpvotePercentage = upvotePercentage
			r.posts[i].Votes = votes
			return nil
		}
	}

	return customerr.PostNotFoundByID{PostID: postID}
}
