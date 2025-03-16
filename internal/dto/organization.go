package dto

type PrerequisiteRequest struct {
	Title string `json:"title" example:"Bachelor's degree in Computer Science" validate:"required,min=1"`
	Link  string `json:"link" example:"https://example.com" validate:"required,min=1"`
}

type PrerequisiteResponses struct {
	Value uint   `json:"value" example:"1"`
	Title string `json:"title" example:"Bachelor's degree in Computer Science"`
	Link  string `json:"link" example:"https://example.com"`
}
