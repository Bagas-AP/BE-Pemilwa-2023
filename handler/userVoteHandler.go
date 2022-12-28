package handler

import (
	"TestVote/middleware"
	"TestVote/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Vote(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	// api untuk melakukan voting
	r.POST("/vote", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		type vote struct {
			CalonKepalaID *int `gorm:"default:null" json:"calonKepala"`
			CalonSenatID  *int `gorm:"default:null" json:"calonSenat"`
		}

		var cek model.Users

		if search := db.Where("id = ?", ID).Take(&cek); search.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"statusCode": http.StatusInternalServerError,
				"error":   search.Error.Error(),
			})
			return
		}

		if cek.IsVote {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Anda sudah memilih",
				"statusCode": http.StatusForbidden,
				"error": nil,
			})
			return
		}

		var input vote

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "input is invalid",
				"statusCode": http.StatusBadRequest,
				"error":   err.Error(),
			})
			return
		}

		if input.CalonKepalaID == nil || input.CalonSenatID == nil || *input.CalonKepalaID > 2 || *input.CalonSenatID > 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "input is invalid",
				"statusCode": http.StatusBadRequest,
				"error":   "calon kepala dan calon senat tidak boleh kosong",
			})
			return
		}

		mahasiswa := model.Users{
			IsVote:        true,
			CalonKepalaID: input.CalonKepalaID,
			CalonSenatID:  input.CalonSenatID,
			WaktuVote:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		result := db.Where("id = ?", ID).Model(&mahasiswa).Updates(mahasiswa)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "error when updating mahasiswa data to database.",
				"statusCode": http.StatusInternalServerError,
				"error":   result.Error.Error(),
			})
			return
		}

		if result = db.Where("id = ?", ID).Take(&mahasiswa); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "error when querying data mahasiswa from database.",
				"statusCode": http.StatusInternalServerError,
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "mahasiswa not found.",
				"statusCode": http.StatusNotFound,
				"error": nil,
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":      true,
			"message":      "berhasil melakukan voting",
			"data": gin.H {
				"NIM":          mahasiswa.NIM,
				"Nama":         mahasiswa.Nama,
				"voted":        mahasiswa.IsVote,
				"Calon Kepala": mahasiswa.CalonKepalaID,
				"Calon Senat":  mahasiswa.CalonSenatID,
			},
			"statusCode": http.StatusCreated,
			"error" : nil,
		})
	})
}
