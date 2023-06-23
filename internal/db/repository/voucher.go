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

func (a *VoucherRepository) Create() (interface{}, error) {

	result, err := db.Mongo.Collection.Voucher.InsertOne(context.TODO(), a.Model)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil

}

func (a *VoucherRepository) SoftDelete() error {

	result, err := db.Mongo.Collection.Voucher.UpdateOne(
		context.TODO(),
		bson.D{
			{"partnerId", a.Model.PartnerID},
			{"code", a.Model.Code},
		},
		bson.D{
			{"$set", bson.D{
				{"updatedAt", a.Model.UpdatedAt},
				{"deactivatedAt", a.Model.DeactivatedAt},
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

// TODO: 3. fetch all voucher

/*
FindAll
return:

	data interfaces{},
	totalDocument int64,
	totalPages int,
	err error
*/
func (a *VoucherRepository) FindAll() (interface{}, int64, int64, error) {

	filter := bson.D{}
	switch a.Pagination.Status {
	case entity.VoucherStatusActive:
		filter = bson.D{{"deactivatedAt", bson.D{{"$eq", primitive.Null{}}}}}
	case entity.VoucherStatusDeactivated:
		filter = bson.D{{"deactivatedAt", bson.D{{"$ne", primitive.Null{}}}}}
	default:
	}

	filter = append(filter, bson.D{{"partnerId", a.Pagination.PartnerID}}...)

	skipValue := (a.Pagination.Page - 1) * a.Pagination.Size

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	cursor, err := db.Mongo.Collection.Voucher.Find(
		ctx,
		filter,
		options.Find().SetSkip(skipValue).SetLimit(a.Pagination.Size),
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

	totalPages := math.Ceil(float64(totalDocs) / float64(a.Pagination.Size))
	return &vouchers, totalDocs, int64(totalPages), nil

}

func (a *VoucherRepository) FindByID(id interface{}) (interface{}, error) {
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

func (a *VoucherRepository) FindByVoucherID(id interface{}) (interface{}, error) {
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

func (a *VoucherRepository) FindByCode() (interface{}, error) {
	result := new(entity.Voucher)

	err := db.Mongo.Collection.Voucher.FindOne(
		context.TODO(),
		bson.D{
			{"partnerId", a.Model.PartnerID},
			{"code", a.Model.Code},
		}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *VoucherRepository) Update() error {

	result, err := db.Mongo.Collection.Voucher.UpdateOne(
		context.TODO(),
		bson.D{
			{"partnerId", a.Model.PartnerID},
			{"code", a.Model.Code},
		},
		bson.D{
			{"$set", bson.D{
				{"description", a.Model.Description},
				{"amount", a.Model.Amount},
				{"price", a.Model.Price},
				{"updatedAt", a.Model.UpdatedAt},
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

//func (a *VoucherRepository) getDefaultFilter(account *entity.AccountBalance) bson.D {
//	filter := bson.D{
//		{"merchantId", account.MerchantID},
//		{"partnerId", account.PartnerID},
//	}
//
//	if account.TerminalID != "" {
//		filter = append(filter, bson.D{{"terminalId", account.TerminalID}}...)
//	}
//
//	if account.Type > 0 {
//		filter = append(filter, bson.D{{"type", account.Type}}...)
//	}
//
//	return filter
//}

//func (a *VoucherRepository) findByFilter(filter interface{}) (interface{}, error) {
//	account := new(entity.AccountBalance)
//
//	// filter condition
//	err := db.Mongo.Collection.Account.FindOne(context.TODO(), filter).Decode(&account)
//	if err != nil {
//		return nil, err
//	}
//
//	return account, nil
//}

/*
FindAll
function args:

	queryParams: string

return:

	code: int,
	accounts: interface{},
	length: int,
	err: error
*/
//func (a *VoucherRepository) FindAll(queryParams string) (int, interface{}, int, error) {
//	filter := bson.D{}
//
//	if queryParams != "" {
//		filter = bson.D{{"active", false}}
//		if queryParams == "true" {
//			filter = bson.D{{"active", true}}
//		}
//	}
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
//	defer cancel()
//
//	cursor, err := db.Mongo.Collection.Account.Find(ctx, filter,
//		options.Find().SetProjection(bson.D{{"secretKey", 0}, {"lastBalance", 0}}),
//	)
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, 0, err
//	}
//
//	var accounts []entity.AccountBalance
//	if err = cursor.All(context.TODO(), &accounts); err != nil {
//		return fiber.StatusInternalServerError, nil, 0, err
//	}
//
//	return fiber.StatusOK, &accounts, len(accounts), nil
//}

/*
FindAllPaginated
function args:

	*request.PaginatedAccountRequest

return:

	code int,
	data interfaces{},
	totalDocument int64,
	totalPages int,
	err error
*/
//func (a *VoucherRepository) FindAllPaginated(request *entity.PaginatedAccountRequest) (int, interface{}, int64, int64, error) {
//	var filter interface{}
//	switch request.Status {
//	case AccountStatusActive:
//		filter = bson.D{{"active", true}}
//	case AccountStatusDeactivated:
//		filter = bson.D{{"active", false}}
//	default:
//		filter = bson.D{}
//	}
//
//	skipValue := (request.Page - 1) * request.Size
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
//	defer cancel()
//
//	cursor, err := db.Mongo.Collection.Account.Find(
//		ctx,
//		filter,
//		options.Find().
//			SetProjection(bson.D{{"secretKey", 0}, {"lastBalance", 0}}).
//			SetSkip(skipValue).
//			SetLimit(request.Size),
//	)
//
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, 0, 0, err
//	}
//
//	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
//	var accounts []entity.AccountBalance
//	if err = cursor.All(context.TODO(), &accounts); err != nil {
//		return fiber.StatusInternalServerError, nil, 0, 0, err
//	}
//
//	if len(accounts) == 0 {
//		return fiber.StatusInternalServerError, nil, 0, 0, errors.New("empty results or last pages has been reached")
//	}
//
//	totalPages := math.Ceil(float64(totalDocs) / float64(request.Size))
//	return fiber.StatusOK, &accounts, totalDocs, int64(totalPages), nil
//}

// FindByID : id args accept interface{} or primitive.ObjectID
// make sure to convert it first
//func (a *VoucherRepository) FindByID(id interface{}, active bool) (int, interface{}, error) {
//	// filter condition
//	filter := bson.D{{"_id", id}}
//	if active {
//		filter = append(filter, bson.D{{"active", true}}...)
//	}
//
//	account, err := a.findByFilter(filter)
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}

//func (a *VoucherRepository) FindOne(account *entity.AccountBalance) (interface{}, error) {
//
//	result, err := a.findByFilter(GetDefaultAccountFilter(account))
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return nil, err
//		}
//		return nil, err
//	}
//
//	return result, nil
//}

//func (a *VoucherRepository) FindByUniqueID(id string, active bool) (int, interface{}, error) {
//
//	filter := bson.D{{"uniqueId", id}}
//	if active {
//		filter = bson.D{{"uniqueId", id}, {"active", true}}
//	}
//
//	account, err := a.findByFilter(filter)
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			//return fiber.StatusNotFound, nil, errors.New("account not found or it has been unregistered")
//			return fiber.StatusNotFound, nil, err
//		}
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}

//func (a *VoucherRepository) FindByActiveStatus(id string, status bool) (int, interface{}, error) {
//	//id, _ := primitive.ObjectIDFromHex(accountId)
//
//	account, err := a.findByFilter(bson.D{
//		{Key: "uniqueId", Value: id},
//		{Key: "active", Value: status},
//	})
//
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}
//
//func (a *VoucherRepository) DeactivateAccount(u *entity.UnregisterAccount) (int, error) {
//
//	account := new(entity.AccountBalance)
//	account.PartnerID = u.PartnerID
//	account.MerchantID = u.MerchantID
//	account.TerminalID = u.TerminalID
//
//	// update field
//	update := bson.D{
//		{"$set", bson.D{
//			{"active", false},
//			{"updatedAt", time.Now().UnixMilli()},
//		}},
//	}
//
//	result, err := db.Mongo.Collection.Account.UpdateOne(context.TODO(), GetDefaultAccountFilter(account), update)
//	if err != nil {
//		return fiber.StatusInternalServerError, err
//	}
//
//	if result.ModifiedCount == 0 {
//		return fiber.StatusBadRequest, errors.New(
//			"update failed, cannot find account with current uniqueId")
//	}
//
//	return fiber.StatusOK, nil
//}
//
//func (a *VoucherRepository) DeactivateMerchant(u *entity.UnregisterAccount) (int, error) {
//	//id, _ := primitive.ObjectIDFromHex(u.VoucherID)
//
//	// filter condition
//	filter := bson.D{{"merchantId", u.MerchantID}, {"partnerId", u.PartnerID}}
//
//	// update field
//	update := bson.D{
//		{"$set", bson.D{
//			{"active", false},
//			{"updatedAt", time.Now().UnixMilli()},
//		}},
//	}
//
//	result, err := db.Mongo.Collection.Account.UpdateOne(context.TODO(), filter, update)
//	if err != nil {
//		return fiber.StatusInternalServerError, err
//	}
//
//	if result.ModifiedCount == 0 {
//		return fiber.StatusBadRequest, errors.New(
//			"update failed, cannot find account with current uniqueId")
//	}
//
//	return fiber.StatusOK, nil
//}
//
//func (a *VoucherRepository) InsertDeactivatedAccount(account *entity.UnregisterAccount) (int, interface{}, error) {
//
//	result, err := db.Mongo.Collection.UnregisterAccount.InsertOne(context.TODO(), account)
//
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusCreated, result.InsertedID, nil
//
//}
//
//func (a *VoucherRepository) RemoveDeactivatedAccount(acc *entity.UnregisterAccount) (int, error) {
//	filter := bson.D{{"uniqueId", acc.VoucherID}}
//	result, err := db.Mongo.Collection.UnregisterAccount.DeleteOne(context.TODO(), filter)
//
//	if err != nil {
//		return fiber.StatusInternalServerError, err
//	}
//
//	if result.DeletedCount > 0 {
//		return fiber.StatusNoContent, nil
//	}
//
//	return fiber.StatusInternalServerError, errors.New("remove deactivated account failed, no document found")
//}

// ----------------- MERCHANTS ----------------

//func (a *VoucherRepository) FindAllMerchant(request *entity.PaginatedAccountRequest) (int, interface{}, int64, int64, error) {
//	//var filter interface{}
//
//	filter := bson.D{}
//	switch request.Status {
//	case AccountStatusActive:
//		filter = append(filter, bson.D{{"active", true}}...)
//	case AccountStatusDeactivated:
//		filter = append(filter, bson.D{{"active", false}}...)
//	default:
//	}
//
//	filter = append(filter, bson.D{{"type", 2}}...)
//
//	skipValue := (request.Page - 1) * request.Size
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
//	defer cancel()
//
//	cursor, err := db.Mongo.Collection.Account.Find(
//		ctx,
//		filter,
//		options.Find().
//			SetProjection(bson.D{{
//				"secretKey", 0},
//			//{"lastBalance", 0},
//			}).
//			SetSkip(skipValue).
//			SetLimit(request.Size),
//	)
//
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, 0, 0, err
//	}
//
//	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
//	var accounts []entity.AccountBalance
//	if err = cursor.All(context.TODO(), &accounts); err != nil {
//		return fiber.StatusInternalServerError, nil, 0, 0, err
//	}
//
//	if len(accounts) == 0 {
//		return fiber.StatusInternalServerError, nil, 0, 0, errors.New("empty results or last pages has been reached")
//	}
//
//	totalPages := math.Ceil(float64(totalDocs) / float64(request.Size))
//	return fiber.StatusOK, &accounts, totalDocs, int64(totalPages), nil
//}
//
//func (a *VoucherRepository) FindByMerchantID(ab *entity.AccountBalance) (int, interface{}, error) {
//
//	filter := bson.D{
//		{"merchantId", ab.MerchantID},
//		{"partnerId", ab.PartnerID},
//		{"terminalId", ab.TerminalID},
//		{"type", AccountTypeMerchant},
//	}
//
//	account, err := a.findByFilter(filter)
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return fiber.StatusNotFound, nil, err
//		}
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}
//
//func (a *VoucherRepository) FindByMerchantStatus(ab *entity.AccountBalance, active bool) (int, interface{}, error) {
//
//	filter := bson.D{
//		{"merchantId", ab.MerchantID},
//		{"partnerId", ab.PartnerID},
//	}
//
//	if active {
//		filter = append(filter, bson.D{{"active", true}}...)
//	}
//
//	account, err := a.findByFilter(filter)
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return fiber.StatusNotFound, nil, err
//		}
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}
//
///*
//CountMerchantMember
//out params:
//
//	totalDoc int
//	documents interface{}
//	err error
//*/
//func (a *VoucherRepository) CountMerchantMember(account entity.AccountBalance) (int64, error) {
//	filter := bson.D{
//		{"partnerId", account.PartnerID},
//		{"merchantId", account.MerchantID},
//		{"active", account.Active},
//	}
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
//	defer cancel()
//
//	totalDocs, err := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
//	if err != nil {
//		return 0, err
//	}
//
//	return totalDocs, nil
//}
//
//func (a *VoucherRepository) FindMemberByMerchant(account entity.AccountBalance) (int64, interface{}, error) {
//	filter := bson.D{
//		{"partnerId", account.PartnerID},
//		{"merchantId", account.MerchantID},
//		{"active", account.Active},
//	}
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
//	defer cancel()
//
//	cursor, err := db.Mongo.Collection.Account.Find(ctx, filter)
//
//	if err != nil {
//		return 0, nil, err
//	}
//
//	var accounts []entity.AccountBalance
//	if err = cursor.All(context.TODO(), &accounts); err != nil {
//		return 0, nil, err
//	}
//
//	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
//
//	if len(accounts) == 0 {
//		return 0, nil, errors.New("empty results or last pages has been reached")
//	}
//
//	return totalDocs, &accounts, nil
//}
//
//func (a *VoucherRepository) FindMemberByMerchantPaginated(request *entity.PaginatedAccountRequest) (int, interface{}, int64, int64, error) {
//	filter := bson.D{}
//	switch request.Status {
//	case AccountStatusActive:
//		filter = bson.D{{"active", true}}
//	case AccountStatusDeactivated:
//		filter = bson.D{{"active", false}}
//	default:
//	}
//
//	filter = append(filter, bson.D{
//		{"partnerId", request.PartnerID},
//		{"merchantId", request.MerchantID},
//		{"type", request.Type},
//	}...)
//
//	skipValue := (request.Page - 1) * request.Size
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
//	defer cancel()
//
//	cursor, err := db.Mongo.Collection.Account.Find(
//		ctx,
//		filter,
//		options.Find().
//			SetProjection(bson.D{{"secretKey", 0}, {"lastBalance", 0}}).
//			SetSkip(skipValue).
//			SetLimit(request.Size),
//	)
//
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, 0, 0, err
//	}
//
//	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
//	var accounts []entity.AccountBalance
//	if err = cursor.All(context.TODO(), &accounts); err != nil {
//		return fiber.StatusInternalServerError, nil, 0, 0, err
//	}
//
//	if len(accounts) == 0 {
//		return fiber.StatusInternalServerError, nil, 0, 0, errors.New("empty results or last pages has been reached")
//	}
//
//	totalPages := math.Ceil(float64(totalDocs) / float64(request.Size))
//	return fiber.StatusOK, &accounts, totalDocs, int64(totalPages), nil
//}
