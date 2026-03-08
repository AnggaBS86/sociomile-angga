package request

type ChannelWebhookRequest struct {
	TenantID           int64  `json:"tenant_id" validate:"required,gt=0"`
	CustomerExternalID string `json:"customer_external_id" validate:"required"`
	Message            string `json:"message" validate:"required"`
}
