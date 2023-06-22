package validator

import (
	"errors"
	"github.com/dw-parameter-service/internal/db/entity"
)

func ValidateRequest(payload interface{}) (interface{}, error) {
	var msg []string

	switch p := payload.(type) {
	case *entity.Voucher:
		if p.PartnerID == "" {
			msg = append(msg, "partnerId cannot be empty.")
		}

		if p.Code == "" {
			msg = append(msg, "code cannot be empty.")
		}

		if p.Amount == 0 {
			msg = append(msg, "amount must be greater than 0.")
		}

		if p.Description == "" {
			msg = append(msg, "description cannot be empty.")
		}

	default:
	}

	if len(msg) > 0 {
		return msg, errors.New("request validation status failed")
	}
	return msg, nil

}
