package controller

import (
	"net/http"

	"github.com/gamaput/go-redeem/model"
	"github.com/gamaput/go-redeem/repository"
	"github.com/gamaput/go-redeem/utils"
	"github.com/gin-gonic/gin"
)

type RedeemCodeController struct {
	RedeemCodeRepo repository.RedeemCodeRepository
	PrizeRepo      repository.PrizeRepository
}

func NewRedeemCodeController(redeemCodeRepo repository.RedeemCodeRepository, prizeRepo repository.PrizeRepository) *RedeemCodeController {
	return &RedeemCodeController{
		RedeemCodeRepo: redeemCodeRepo,
		PrizeRepo:      prizeRepo,
	}
}

func (c *RedeemCodeController) GenerateCode(ctx *gin.Context) {
	code := utils.GenerateUniqueCode()

	redeemCode := &model.RedeemCode{
		Code:       code,
		IsRedeemed: false,
		Name:       "",
		NoKTP:      "",
		City:       "",
		Address:    "",
		PhoneNo:    "",
	}

	err := c.RedeemCodeRepo.SaveRedeemCode(redeemCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save code"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

func (c RedeemCodeController) RedeemCode(ctx *gin.Context) {
	var redeemCode model.RedeemCode

	if err := ctx.ShouldBindJSON(&redeemCode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Memeriksa apakah semua data terisi
	if redeemCode.Name == "" || redeemCode.NoKTP == "" || redeemCode.City == "" || redeemCode.Address == "" || redeemCode.PhoneNo == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	existingRedeemCode, err := c.RedeemCodeRepo.GetRedeemCodeByCode(redeemCode.Code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redeem code"})
		return
	}

	// Validasi apakah redeem code sudah digunakan sebelumnya
	if existingRedeemCode.IsRedeemed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Redeem code has already been redeemed"})
		return
	}

	// Lakukan tindakan yang sesuai dengan kode redeem yang valid dan informasi pengguna
	existingRedeemCode.Name = redeemCode.Name
	existingRedeemCode.NoKTP = redeemCode.NoKTP
	existingRedeemCode.City = redeemCode.City
	existingRedeemCode.Address = redeemCode.Address
	existingRedeemCode.PhoneNo = redeemCode.PhoneNo

	// Mendapatkan hadiah secara acak dari database
	randomPrize, err := c.PrizeRepo.GetRandomPrize()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get random prize"})
		return
	}

	// Jika stok hadiah tersedia, kurangi stok dan update prize
	if randomPrize.Quantity > 0 {
		randomPrize.Quantity--
		err = c.PrizeRepo.UpdatePrize(&randomPrize)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prize quantity"})
			return
		}
	} else {
		// Jika stok hadiah tidak tersedia, set ID hadiah sebagai 0
		randomPrize.ID = 0
	}

	// Set ID hadiah yang diberikan kepada pengguna ke dalam kolom RedeemCode
	existingRedeemCode.PrizeID = randomPrize.ID

	// Menandai redeem code sebagai sudah diredeem
	existingRedeemCode.IsRedeemed = true

	// Simpan perubahan ke dalam database
	err = c.RedeemCodeRepo.UpdateRedeemCode(existingRedeemCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update redeem code"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Redeem code successfully validated and marked as redeemed",
		"prize":   randomPrize,
	})

}

func (c *RedeemCodeController) GetAllRedeems(ctx *gin.Context) {

	redeems, err := c.RedeemCodeRepo.GetAllRedeems()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, redeems)
}
