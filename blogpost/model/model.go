package model

import (
	"log"

	validation "github.com/go-ozzo/ozzo-validation"
)

// BlogPost represents a blog post.
type BlogPost struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	Slug              string `json:"slug"`
	Description       string `json:"description"`
	Body              string `json:"body"`
	Revision          int64  `json:"revision"`
	CreationTimestamp int64  `json:"creationTimestamp"`
	UpdateTimestamp   int64  `json:"updateTimestamp"`
}

// Cursor represents the cursor of a blog post page.
type Cursor struct {
	ID                string `json:"id"`
	CreationTimestamp int64  `json:"creationTimestamp"`
}

// Page represents a page of blog posts.
type Page struct {
	Posts  []BlogPost `json:"posts"`
	Cursor Cursor     `json:"cursor"`
}

// Validate checks if a BlogPost instance is valid and returns an error. If it's valid, it returns nil.
func (post BlogPost) Validate() []string {
	err := validation.ValidateStruct(
		&post,
		validation.Field(
			&post.ID,
			validation.Required.Error("The id is required."),
			validation.Length(0, 250).Error("The id may have up to 250 characters.")),
		validation.Field(
			&post.Title,
			validation.Required.Error("The title is required."),
			validation.Length(0, 250).Error("The title may have up to 250 characters.")),
		validation.Field(
			&post.Slug,
			validation.Required.Error("The slug is required."),
			validation.Length(0, 250).Error("The slug may have up tp 250 characters.")),
		validation.Field(
			&post.Description,
			validation.Required.Error("The description is required."),
			validation.Length(0, 250).Error("The description may have up tp 250 characters.")),
		validation.Field(
			&post.Body,
			validation.Required.Error("The body is required.")),
		validation.Field(
			&post.Revision,
			validation.Required.Error("The revision is required.")))

	if err == nil {
		return []string{}
	}

	validationErrors, ok := err.(validation.Errors)
	if !ok {
		log.Fatal("An unexpected error occurred while validating a model: ", err)
		return []string{"An unexpected error occurred."}
	}

	result := make([]string, len(validationErrors))
	i := 0
	for _, err = range validationErrors {
		result[i] = err.Error()
		i++
	}

	return result
}
