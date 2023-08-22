package repository

import (
	"github.com/gamaput/go-redeem/model"
	"gorm.io/gorm"
)

type prizeRepository struct {
	DB *gorm.DB
}

type PrizeRepository interface {
	GetRandomPrize() (model.Prize, error)
	CreatePrize(prize *model.Prize) error
	GetAllPrizes() ([]model.Prize, error)
	UpdatePrize(prize *model.Prize) error
	DeletePrize(model.Prize) (model.Prize, error)
	GetPrizeByID(uint) (model.Prize, error)
}

func NewPrizeRepository(db *gorm.DB) PrizeRepository {
	return &prizeRepository{
		DB: db,
	}
}

func (pr *prizeRepository) GetRandomPrize() (model.Prize, error) {
	var prize model.Prize
	if err := pr.DB.Model(&model.Prize{}).Scopes(prizeIsAvailable).Order("RAND()").First(&prize).Error; err != nil {
		// Jika tidak ada hadiah yang tersedia, prize.ID akan diatur menjadi 0
		if err == gorm.ErrRecordNotFound {
			prize.ID = 0
			return prize, nil
		}
		return model.Prize{}, err
	}
	return prize, nil
}

func (pr *prizeRepository) CreatePrize(prize *model.Prize) error {
	if err := pr.DB.Create(&prize).Error; err != nil {
		return err
	}
	return nil
}

func (pr *prizeRepository) GetAllPrizes() ([]model.Prize, error) {
	var prizes []model.Prize
	result := pr.DB.Find(&prizes)
	if result.Error != nil {
		return nil, result.Error
	}
	return prizes, nil
}

func prizeIsAvailable(db *gorm.DB) *gorm.DB {
	return db.Where("quantity > 0")
}

func (r *prizeRepository) UpdatePrize(prize *model.Prize) error {
	if err := r.DB.Model(&model.Prize{}).Where("id =?", prize.ID).Updates(map[string]interface{}{
		"name":     prize.Name,
		"quantity": prize.Quantity, // Menggunakan nilai kuantitas yang baru
	}).Error; err != nil {
		return err
	}
	return nil
}

// DeletePrize deletes a prize
func (r *prizeRepository) DeletePrize(prize model.Prize) (model.Prize, error) {
	if err := r.DB.Delete(&prize, prize.ID).Error; err != nil {
		return prize, err
	}
	return prize, r.DB.Delete(&prize).Error
}

func (r *prizeRepository) GetPrizeByID(id uint) (prize model.Prize, err error) {
	return prize, r.DB.First(&prize, id).Error
}
