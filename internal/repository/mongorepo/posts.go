package mongorepo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
)

type postsRepo struct {
	posts *mongo.Collection
}

func NewPostsRepo(collection *mongo.Collection) *postsRepo {
	return &postsRepo{posts: collection}
}

func (r *postsRepo) GetAllPosts() ([]model.Post, error) {
	posts := make([]model.Post, 0)
	filter := bson.M{}
	cursor, err := r.posts.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var post model.Post
		err = cursor.Decode(&post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *postsRepo) GetPostsByCategory(category string) ([]model.Post, error) {
	posts := make([]model.Post, 0)
	filter := bson.M{"category": category}
	cursor, err := r.posts.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var post model.Post
		err = cursor.Decode(&post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *postsRepo) GetPostsByAuthor(username string) ([]model.Post, error) {
	posts := make([]model.Post, 0)
	filter := bson.M{"author.username": username}
	cursor, err := r.posts.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var post model.Post
		err = cursor.Decode(&post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *postsRepo) AddPost(post model.Post) error {
	_, err := r.posts.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

func (r *postsRepo) GetPostByID(postID string) (model.Post, error) {
	var post model.Post
	filter := bson.M{"id": postID}
	err := r.posts.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Post{}, customerr.PostNotFoundByID{PostID: postID}
		}
		return model.Post{}, err
	}
	return post, nil
}

func (r *postsRepo) GetPostByIDAndUpdateViews(postID string) (model.Post, error) {
	var post model.Post
	filter := bson.M{"id": postID}
	update := bson.M{"$inc": bson.M{"views": 1}}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.posts.FindOneAndUpdate(context.TODO(), filter, update, opt).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Post{}, customerr.PostNotFoundByID{PostID: postID}
		}
		return model.Post{}, err
	}
	return post, nil
}

func (r *postsRepo) DeletePost(postID string) error {
	filter := bson.M{"id": postID}
	_, err := r.posts.DeleteOne(context.TODO(), filter)
	return err
}

func (r *postsRepo) AddComment(postID string, comment model.Comment) (model.Post, error) {
	var post model.Post
	filter := bson.M{"id": postID}
	update := bson.M{"$push": bson.M{"comments": comment}}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.posts.FindOneAndUpdate(context.TODO(), filter, update, opt).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Post{}, customerr.PostNotFoundByID{PostID: postID}
		}
		return model.Post{}, err
	}

	return post, nil
}

func (r *postsRepo) GetCommentByID(postID, commentID string) (model.Comment, error) {
	var post model.Post
	filter := bson.M{"id": postID}
	err := r.posts.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Comment{}, customerr.PostNotFoundByID{PostID: postID}
		}
		return model.Comment{}, err
	}

	for _, c := range post.Comments {
		if c.ID == commentID {
			return c, nil
		}
	}

	return model.Comment{}, customerr.CommentNotFoundByID{PostID: postID, CommentID: commentID}
}

func (r *postsRepo) DeleteComment(postID, commentID string) (model.Post, error) {
	var post model.Post
	filter := bson.M{"id": postID}
	update := bson.M{"$pull": bson.M{"comments": bson.M{"id": commentID}}}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.posts.FindOneAndUpdate(context.TODO(), filter, update, opt).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Post{}, customerr.PostNotFoundByID{PostID: postID}
		}
		return model.Post{}, err
	}

	return post, nil
}

func (r *postsRepo) UpdateVotes(postID string, score int, upvotePercentage int, votes []model.Vote) error {
	filter := bson.M{"id": postID}
	update := bson.M{"$set": bson.M{"score": score, "votes": votes, "upvotePercentage": upvotePercentage}}
	_, err := r.posts.UpdateOne(context.TODO(), filter, update)
	return err
}
