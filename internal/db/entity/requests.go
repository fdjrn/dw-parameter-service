package entity

type PaginatedRequest struct {
	PartnerID string `json:"partnerId,omitempty"`

	Code string `json:"code,omitempty"`

	// Voucher status: active/deactivated
	Status string `json:"status,omitempty"`

	Page int64 `json:"page,omitempty"`

	Size int64 `json:"size,omitempty"`
}
