package repository

import (
	"context"
	"errors"
	"github.com/dw-parameter-service/internal/db"
	"github.com/dw-parameter-service/internal/db/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

type VoucherRepository struct {
	Model      *entity.Voucher
	Pagination *entity.PaginatedRequest
}

func NewVoucherRepository() VoucherRepository {
	return VoucherRepository{
		Model:      new(entity.Voucher),
		Pagination: new(entity.PaginatedRequest),
	}
}

func (r *VoucherRepository) Create() (interface{}, error) {

	result, err := db.Mongo.Collection.Voucher.InsertOne(context.TODO(), r.Model)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil

}

func (r *VoucherRepository) SoftDelete() error {

	result, err := db.Mongo.Collection.Voucher.UpdateOne(
		context.TODO(),
		bson.D{
			{"partnerId", r.Model.PartnerID},
			{"code", r.Model.Code},
		},
		bson.D{
			{"$set", bson.D{
				{"updatedAt", r.Model.UpdatedAt},
				{"deactivatedAt", r.Model.DeactivatedAt},
			}},
		})

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New(
			"deactivation failed, cannot find voucher with current partnerId and code")
	}

	return nil

}

func (r *VoucherRepository) FindAll() (interface{}, int64, int64, error) {

	filter := bson.D{}
	switch r.Pagination.Status {
	case entity.VoucherStatusActive:
		filter = bson.D{{"deactivatedAt", bson.D{{"$eq", primitive.Null{}}}}}
	case entity.VoucherStatusDeactivated:
		filter = bson.D{{"deactivatedAt", bson.D{{"$ne", primitive.Null{}}}}}
	default:
	}

	filter = append(filter, bson.D{{"partnerId", r.Pagination.PartnerID}}...)

	skipValue := (r.Pagination.Page - 1) * r.Pagination.Size

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	cursor, err := db.Mongo.Collection.Voucher.Find(
		ctx,
		filter,
		options.Find().SetSkip(skipValue).SetLimit(r.Pagination.Size),
	)

	if err != nil {
		return nil, 0, 0, err
	}

	totalDocs, _ := db.Mongo.Collection.Voucher.CountDocuments(ctx, filter)
	var vouchers []entity.Voucher
	if err = cursor.All(context.TODO(), &vouchers); err != nil {
		return nil, 0, 0, err
	}

	if len(vouchers) == 0 {
		return nil, 0, 0, errors.New("empty results or last pages has been reached")
	}

	totalPages := math.Ceil(float64(totalDocs) / float64(r.Pagination.Size))
	return &vouchers, totalDocs, int64(totalPages), nil

}

func (r *VoucherRepository) FindByID(id interface{}) (interface{}, error) {
	result := new(entity.Voucher)

	err := db.Mongo.Collection.Voucher.FindOne(
		context.TODO(),
		bson.D{{"_id", id}},
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *VoucherRepository) FindByVoucherID(id interface{}) (interface{}, error) {
	result := new(entity.Voucher)

	err := db.Mongo.Collection.Voucher.FindOne(
		context.TODO(),
		bson.D{{"voucherId", id}},
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *VoucherRepository) FindByCode() (interface{}, error) {
	result := new(entity.Voucher)

	err := db.Mongo.Collection.Voucher.FindOne(
		context.TODO(),
		bson.D{
			{"partnerId", r.Model.PartnerID},
			{"code", r.Model.Code},
		}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *VoucherRepository) Update() error {

	result, err := db.Mongo.Collection.Voucher.UpdateOne(
		context.TODO(),
		bson.D{
			{"partnerId", r.Model.PartnerID},
			{"code", r.Model.Code},
		},
		bson.D{
			{"$set", bson.D{
				{"description", r.Model.Description},
				{"amount", r.Model.Amount},
				{"price", r.Model.Price},
				{"updatedAt", r.Model.UpdatedAt},
			}},
		})

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New(
			"update voucher failed, cannot find voucher with current partnerId and code")
	}

	return nil

}
