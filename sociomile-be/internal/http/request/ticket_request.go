package request

type EscalateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type UpdateTicketStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=open in_progress resolved closed"`
}
