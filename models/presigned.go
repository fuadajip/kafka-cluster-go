package models

// CreateUserUploadPresignedRequest will return object from given data to generate presigned upload
type CreateUserUploadPresignedRequest struct {
	Filename     *string  `json:"filename" validate:"required"`
	DocumentType *string  `json:"document_type" validate:"required"` // [AVATAR, IC, IC_SELFIE, POLICY, etc]
	UserType     *string  `json:"user_type" validate:"required"`     // [AGENT, CORPORATE]
	Size         *float64 `json:"size" validate:"-"`
}

// CreateUserUploadPresignedResponse will return object from given data to generate presigned upload
type CreateUserUploadPresignedResponse struct {
	URL          *string  `json:"url"`
	Key          *string  `json:"key"`
	Filename     *string  `json:"filename"`
	MimeType     *string  `json:"mimetype"`
	Size         *float64 `json:"size"`
	DocumentType *string  `json:"document_type"`
	UserType     *string  `json:"user_type"`
}

type CreateUserUploadPresignedViewRequest struct {
	Key *string `json:"key" validate:"required"`
}

type CreateUserUploadPresignedViewResponse struct {
	URL *string `json:"url"`
	Key *string `json:"key"`
}
