package entity

// Voucher
// adalah struct yang digunakan untuk manajemen voucher
type Voucher struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	VoucherID     string `json:"voucherId" bson:"voucherId"`
	PartnerID     string `json:"partnerId" bson:"partnerId"`
	Code          string `json:"code" bson:"code"`
	Description   string `json:"description" bson:"description"`
	Amount        int64  `json:"amount" bson:"amount"`
	Price         int64  `json:"price" bson:"price"`
	CreatedAt     int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt     int64  `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeactivatedAt int64  `json:"deactivatedAt,omitempty" bson:"deactivatedAt,omitempty"`
	Status        string `json:"status,omitempty" bson:"-"`
}

const (
	VoucherStatusActive      = "active"
	VoucherStatusDeactivated = "deactivated"
	VoucherStatusAll         = "all"
)
