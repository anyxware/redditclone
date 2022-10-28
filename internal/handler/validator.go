package handler

import (
	"redditclone/pkg/hexid"
	"redditclone/pkg/httpvalidator"
)

var (
	categories = [...]string{"music", "funny", "videos", "programming", "news", "fashion"}
)

func (h *Handler) initValidator() {
	postInputTmpl := httpvalidator.RequestBody{
		Fields: httpvalidator.Fields{
			"category": httpvalidator.BodyField{
				Required: true,
				Rules: []httpvalidator.Rule{
					{
						Description: "category must has specific type",
						Validate: func(category string) bool {
							for _, existedCategory := range categories {
								if existedCategory == category {
									return true
								}
							}
							return false
						},
					},
				},
			},
			"title": httpvalidator.BodyField{
				Required: true,
				Rules: []httpvalidator.Rule{
					{
						Description: "title must be a non-empty string",
						Validate: func(title string) bool {
							return len(title) > 0
						},
					},
				},
			},
			"type": httpvalidator.BodyField{
				Required: true,
				Rules: []httpvalidator.Rule{
					{
						Description: "type must be a text or a link",
						Validate: func(postType string) bool {
							return postType == "text" || postType == "link"
						},
					},
				},
			},
		},
	}

	textPostInputTmpl := httpvalidator.RequestBody{
		Fields: httpvalidator.Fields{
			"text": httpvalidator.BodyField{
				Required: true,
				Rules: []httpvalidator.Rule{
					{
						Description: "text must be a non-empty string",
						Validate: func(text string) bool {
							return len(text) > 0
						},
					},
				},
			},
		},
	}

	urlPostInputTmpl := httpvalidator.RequestBody{
		Fields: httpvalidator.Fields{
			"url": httpvalidator.BodyField{
				Required: true,
				Rules: []httpvalidator.Rule{
					{
						Description: "url must be a non-empty string",
						Validate: func(url string) bool {
							return len(url) > 0
						},
					},
				},
			},
		},
	}

	credentialTmpl := httpvalidator.RequestBody{
		Fields: httpvalidator.Fields{
			"username": httpvalidator.BodyField{
				Required: true,
				Rules: []httpvalidator.Rule{
					{
						Description: "username must be a non-empty string",
						Validate: func(username string) bool {
							return len(username) > 0
						},
					},
				},
			},
			"password": httpvalidator.BodyField{
				Required: true,
				Rules: []httpvalidator.Rule{
					{
						Description: "password must be a non-empty string",
						Validate: func(password string) bool {
							return len(password) > 0
						},
					},
				},
			},
		},
	}

	commentTmpl := httpvalidator.RequestBody{
		Fields: httpvalidator.Fields{
			"comment": httpvalidator.BodyField{
				Required: true,
				Rules: []httpvalidator.Rule{
					{
						Description: "comment must be a non-empty string",
						Validate: func(comment string) bool {
							return len(comment) > 0
						},
					},
				},
			},
		},
	}

	h.validator.AddBodyTemplate("PostInput", postInputTmpl)
	h.validator.AddBodyTemplate("TextPostInput", textPostInputTmpl)
	h.validator.AddBodyTemplate("URLPostInput", urlPostInputTmpl)
	h.validator.AddBodyTemplate("Credential", credentialTmpl)
	h.validator.AddBodyTemplate("Comment", commentTmpl)

	userIDValueRules := []httpvalidator.Rule{
		{
			Description: "user_id must be a hexadecimal 24-symbols string",
			Validate: func(id string) bool {
				return hexid.Validate(id)
			},
		},
	}

	postIDValueRules := []httpvalidator.Rule{
		{
			Description: "post_id must be a hexadecimal 24-symbols string",
			Validate: func(id string) bool {
				return hexid.Validate(id)
			},
		},
	}

	commentIDValueRules := []httpvalidator.Rule{
		{
			Description: "comment_id must be a hexadecimal 24-symbols string",
			Validate: func(id string) bool {
				return hexid.Validate(id)
			},
		},
	}

	categoryRules := []httpvalidator.Rule{
		{
			Description: "category must has specific type",
			Validate: func(category string) bool {
				for _, existedCategory := range categories {
					if existedCategory == category {
						return true
					}
				}
				return false
			},
		},
	}

	usernameRules := []httpvalidator.Rule{
		{
			Description: "username must be a non-empty string",
			Validate: func(username string) bool {
				return len(username) > 0
			},
		},
	}

	h.validator.AddPathValueTemplate("user_id", userIDValueRules)
	h.validator.AddPathValueTemplate("post_id", postIDValueRules)
	h.validator.AddPathValueTemplate("comment_id", commentIDValueRules)
	h.validator.AddPathValueTemplate("category", categoryRules)
	h.validator.AddPathValueTemplate("username", usernameRules)
}
