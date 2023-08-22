package controller

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gamaput/go-redeem/model"
	"github.com/gamaput/go-redeem/repository"
	"github.com/gin-gonic/gin"
)

// PrizeController adalah controller untuk mengelola hadiah

type PrizeController interface {
	GenerateRandomPrize(c *gin.Context)
	GetRandomPrize(c *gin.Context)
	CreatePrize(c *gin.Context)
	GetAllPrizes(c *gin.Context)
	DeletePrize(c *gin.Context)
	UpdatePrize(c *gin.Context)
	GetPrizeByID(c *gin.Context)
}

type prizeController struct {
	Repo repository.PrizeRepository
}

func NewPrizeController(repo repository.PrizeRepository) PrizeController {
	return prizeController{
		Repo: repo,
	}
}
func (pc prizeController) GenerateRandomPrize(c *gin.Context) {
	var request struct {
		Quantity int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prizes, err := pc.Repo.GetAllPrizes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Memastikan bahwa jumlah hadiah yang di-generate tidak melebihi total hadiah yang tersedia
	if request.Quantity > len(prizes) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough prizes available"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(prizes), func(i, j int) { prizes[i], prizes[j] = prizes[j], prizes[i] })

	randomPrizes := prizes[:request.Quantity]

	c.JSON(http.StatusOK, randomPrizes)
}

// GetRandomPrize mengambil hadiah secara acak dan mengembalikannya dalam respons JSON
func (pc prizeController) GetRandomPrize(c *gin.Context) {
	prize, err := pc.Repo.GetRandomPrize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prize)
}

func (pc prizeController) CreatePrize(c *gin.Context) {
	var input model.Prize
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi input
	if err := validateCreatePrizeInput(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prize := model.Prize{
		Name:     input.Name,
		Quantity: input.Quantity,
	}

	if err := pc.Repo.CreatePrize(&prize); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, prize)
}

func validateCreatePrizeInput(input model.Prize) error {
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	return nil
}

func (pr prizeController) GetAllPrizes(c *gin.Context) {

	products, err := pr.Repo.GetAllPrizes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdatePrize updates a prize
func (c prizeController) UpdatePrize(ctx *gin.Context) {
	var prize model.Prize
	id := ctx.Param("prize")

	if err := ctx.ShouldBindJSON(&prize); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	prizeID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prize ID"})
		return
	}

	existingPrize, err := c.Repo.GetPrizeByID(uint(prizeID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Prize not found"})
		return
	}

	// Check if the requested quantity is negative
	if prize.Quantity < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	// Update the existing prize object
	existingPrize.Name = prize.Name
	existingPrize.Quantity = prize.Quantity

	if err := c.Repo.UpdatePrize(&existingPrize); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prize"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Prize updated successfully", "prize": existingPrize})
}

// DeletePrize deletes a prize
func (c prizeController) DeletePrize(ctx *gin.Context) {
	var prize model.Prize
	id := ctx.Param("prize")
	intID, _ := strconv.Atoi(id)
	prize.ID = uint(intID)
	prize, err := c.Repo.DeletePrize(prize)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prize ID"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Prize deleted successfully"})
}

func (pc prizeController) GetPrizeByID(c *gin.Context) {
	id := c.Param("prize")
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prize ID"})
		return
	}

	prize, err := pc.Repo.GetPrizeByID(uint(intID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize not found"})
		return
	}

	c.JSON(http.StatusOK, prize)
}
