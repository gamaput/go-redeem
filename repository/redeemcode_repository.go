package repository

import (
	"github.com/gamaput/go-redeem/model"
	"gorm.io/gorm"
)

type redeemCodeRepository struct {
	DB *gorm.DB
}

// RedeemCodeRepository represents the redeem code repository contract
type RedeemCodeRepository interface {
	SaveRedeemCode(redeemCode *model.RedeemCode) error
	GetRedeemCodeByCode(code string) (*model.RedeemCode, error)
	UpdateRedeemCode(redeemCode *model.RedeemCode) error
	RedeemCode(redeemCode *model.RedeemCode) error
	CreateRedeemCode(redeemCode *model.RedeemCode) error
	GetAllRedeems() (redeems []model.RedeemCode, err error)
}

// NewRedeemCodeRepository returns a new instance of RedeemCodeRepository
func NewRedeemCodeRepository(db *gorm.DB) RedeemCodeRepository {
	return &redeemCodeRepository{
		DB: db,
	}
}

func (r *redeemCodeRepository) SaveRedeemCode(redeemCode *model.RedeemCode) error {
	if err := r.DB.Create(redeemCode).Error; err != nil {
		return err
	}
	return nil
}

func (r *redeemCodeRepository) CreateRedeemCode(redeemCode *model.RedeemCode) error {
	result := r.DB.Create(redeemCode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *redeemCodeRepository) GetRedeemCodeByCode(code string) (*model.RedeemCode, error) {
	var redeemCode model.RedeemCode
	if err := r.DB.Where("code = ?", code).First(&redeemCode).Error; err != nil {
		return nil, err
	}
	return &redeemCode, nil
}

func (r *redeemCodeRepository) UpdateRedeemCode(redeemCode *model.RedeemCode) error {
	result := r.DB.Model(&model.RedeemCode{}).Where("id = ?", redeemCode.ID).Updates(map[string]interface{}{
		"is_redeemed": redeemCode.IsRedeemed,
		"name":        redeemCode.Name,
		"no_ktp":      redeemCode.NoKTP,
		"city":        redeemCode.City,
		"address":     redeemCode.Address,
		"phone_no":    redeemCode.PhoneNo,
		"prize_id":    redeemCode.PrizeID,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *redeemCodeRepository) RedeemCode(redeemCode *model.RedeemCode) error {
	if err := r.DB.Model(&model.RedeemCode{}).Where("id = ?", redeemCode.ID).Update("is_redeemed", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *redeemCodeRepository) GetAllRedeems() (redeems []model.RedeemCode, err error) {
	return redeems, r.DB.Find(&redeems).Error
}
