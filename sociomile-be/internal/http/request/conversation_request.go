package request

type ReplyRequest struct {
	Message string `json:"message" validate:"required"`
}
