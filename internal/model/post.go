package model

import "time"

type Vote struct {
	UserID string `json:"user" bson:"user"`
	Vote   int    `json:"vote" bson:"vote"`
}

type TextPostInput struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Text     string `json:"text"`
}

type URLPostInput struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type Post struct {
	ID               string    `json:"id" bson:"id"`
	Score            int       `json:"score" bson:"score"`
	Views            int       `json:"views" bson:"views"`
	Type             string    `json:"type" bson:"type"`
	Title            string    `json:"title" bson:"title"`
	Author           Author    `json:"author" bson:"author"`
	Category         string    `json:"category" bson:"category"`
	Text             string    `json:"text,omitempty" bson:"text"`
	URL              string    `json:"url,omitempty" bson:"url"`
	Votes            []Vote    `json:"votes" bson:"votes"`
	Comments         []Comment `json:"comments" bson:"comments"`
	Created          string    `json:"created" bson:"created"`
	UpvotePercentage int       `json:"upvotePercentage" bson:"upvotePercentage"`
}

func NewTextPost(postID string, input TextPostInput, author Author) Post {
	return Post{
		Score:            0,
		Views:            0,
		Type:             input.Type,
		Title:            input.Title,
		Author:           author,
		Category:         input.Category,
		Text:             input.Text,
		Votes:            make([]Vote, 0),
		Comments:         make([]Comment, 0),
		Created:          time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		UpvotePercentage: 0,
		ID:               postID,
	}
}

func NewURLPost(postID string, input URLPostInput, author Author) Post {
	return Post{
		Score:            0,
		Views:            0,
		Type:             input.Type,
		Title:            input.Title,
		Author:           author,
		Category:         input.Category,
		URL:              input.URL,
		Votes:            make([]Vote, 0),
		Comments:         make([]Comment, 0),
		Created:          time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		UpvotePercentage: 0,
		ID:               postID,
	}
}

func (p *Post) Upvote(userID string) *Post {
	found := false
	for i, vote := range p.Votes {
		if vote.UserID == userID {
			p.Votes[i].Vote = 1
			p.Score += 2
			found = true
			break
		}
	}
	if !found {
		p.Votes = append(p.Votes, Vote{UserID: userID, Vote: 1})
		p.Score += 1
	}
	return p
}

func (p *Post) Downvote(userID string) *Post {
	found := false
	for i, vote := range p.Votes {
		if vote.UserID == userID {
			p.Votes[i].Vote = -1
			p.Score -= 2
			found = true
			break
		}
	}
	if !found {
		p.Votes = append(p.Votes, Vote{UserID: userID, Vote: -1})
		p.Score -= 1
	}
	return p
}

func (p *Post) Unvote(userID string) *Post {
	for i, vote := range p.Votes {
		if vote.UserID == userID {
			p.Votes = append(p.Votes[:i], p.Votes[i+1:]...)
			p.Score -= vote.Vote
			break
		}
	}
	return p
}

func (p *Post) RecalculatePercentage() *Post {
	upvotes := 0
	for _, vote := range p.Votes {
		if vote.Vote == 1 {
			upvotes++
		}
	}
	if len(p.Votes) == 0 {
		p.UpvotePercentage = 0
	} else {
		p.UpvotePercentage = 100 * upvotes / len(p.Votes)
	}
	return p
}
