package handlers

import (
	"errors"
	"fmt"
	"github.com/dw-parameter-service/internal/db/entity"
	"github.com/dw-parameter-service/internal/db/repository"
	"github.com/dw-parameter-service/internal/handlers/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type VoucherHandler struct {
	Repository repository.VoucherRepository
}

func NewVoucherHandler() VoucherHandler {
	return VoucherHandler{Repository: repository.NewVoucherRepository()}
}

func (v *VoucherHandler) validateRequest(payload *entity.Voucher) error {
	if payload.PartnerID == "" {
		return errors.New("partnerId cannot be empty")
	}

	if payload.Code == "" {
		return errors.New("code cannot be empty")
	}

	return nil
}

func (v *VoucherHandler) Create(c *fiber.Ctx) error {
	// new account struct
	payload := new(entity.Voucher)

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate
	errMsg, err := validator.ValidateRequest(payload)
	if err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    errMsg,
		})
	}

	v.Repository.Model = payload

	_, err = v.Repository.FindByCode()
	// code exists
	if err == nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "voucher code already exists",
			Data: fiber.Map{
				"code":      payload.Code,
				"partnerId": payload.PartnerID,
			},
		})
	}

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	tStamp := time.Now().UnixMilli()
	v.Repository.Model.CreatedAt = tStamp
	v.Repository.Model.UpdatedAt = tStamp

	lastId, err := v.Repository.Create()
	if err != nil {
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	newVoucher, err := v.Repository.FindByID(lastId)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(entity.Responses{
		Success: true,
		Message: "voucher successfully created",
		Data:    newVoucher,
	})
}

func (v *VoucherHandler) Deactivate(c *fiber.Ctx) error {
	// new account struct
	payload := new(entity.Voucher)

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err := v.validateRequest(payload)
	if err != nil {
		var msg []string
		msg = append(msg, err.Error())
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "request validation status failed",
			Data:    msg,
		})
	}

	tStamp := time.Now().UnixMilli()
	v.Repository.Model = payload
	v.Repository.Model.UpdatedAt = tStamp
	v.Repository.Model.DeactivatedAt = tStamp

	err = v.Repository.SoftDelete()
	if err != nil {
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "voucher successfully deactivated",
		Data:    nil,
	})
}

func (v *VoucherHandler) Update(c *fiber.Ctx) error {
	// new account struct
	payload := new(entity.Voucher)

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err := v.validateRequest(payload)
	if err != nil {
		var msg []string
		msg = append(msg, err.Error())
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "request validation status failed",
			Data:    msg,
		})
	}

	tStamp := time.Now().UnixMilli()
	v.Repository.Model = payload
	v.Repository.Model.UpdatedAt = tStamp
	v.Repository.Model.DeactivatedAt = tStamp

	err = v.Repository.Update()
	if err != nil {
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	data, err := v.Repository.FindByCode()
	if err != nil {
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: "cannot fetch updated voucher detail",
			Data:    nil,
		})
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "voucher successfully updated",
		Data:    data,
	})
}

func (v *VoucherHandler) GetAllVoucher(c *fiber.Ctx) error {
	var payload = new(entity.PaginatedRequest)

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	msgResponse := "vouchers successfully fetched"
	payload.Status = strings.ToLower(payload.Status)

	validStatus := map[string]interface{}{
		entity.VoucherStatusAll:         0,
		entity.VoucherStatusActive:      1,
		entity.VoucherStatusDeactivated: 2,
	}

	if payload.Status != "" {
		if _, ok := validStatus[payload.Status]; !ok {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "invalid status value. its only accept all, active or deactivated",
				Data:    nil,
			})
		}
		msgResponse = fmt.Sprintf("%s vouchers successfully fetched", payload.Status)
	}

	// set default value
	if payload.Page == 0 {
		payload.Page = 1
	}

	if payload.Size == 0 {
		payload.Size = 10
	}

	v.Repository.Pagination = payload
	vouchers, total, pages, err := v.Repository.FindAll()
	if err != nil {
		return c.Status(500).JSON(entity.ResponsePaginated{
			Success: false,
			Message: err.Error(),
			Data:    entity.ResponsePaginatedData{},
		})
	}

	return c.Status(200).JSON(entity.ResponsePaginated{
		Success: true,
		Message: msgResponse,
		Data: entity.ResponsePaginatedData{
			Result:      vouchers,
			Total:       total,
			PerPage:     payload.Size,
			CurrentPage: payload.Page,
			LastPage:    pages,
		},
	})

}

func (v *VoucherHandler) GetVoucherByID(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	result, err := v.Repository.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "voucher detail not founds",
				Data:    nil,
			})
		}
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "voucher detail successfully fetched",
		Data:    result,
	})
}

func (v *VoucherHandler) GetVoucherByCode(c *fiber.Ctx) error {
	// new account struct
	payload := new(entity.Voucher)

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err := v.validateRequest(payload)
	if err != nil {
		var msg []string
		msg = append(msg, err.Error())
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "request validation status failed",
			Data:    msg,
		})
	}

	v.Repository.Model = payload
	result, err := v.Repository.FindByCode()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "voucher detail not found",
				Data:    nil,
			})
		}

		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "voucher detail successfully fetched",
		Data:    result,
	})
}
