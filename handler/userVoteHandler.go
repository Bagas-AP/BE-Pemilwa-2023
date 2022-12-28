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
	r.POST("/vote", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		type vote struct {
			CalonKepalaID *int `gorm:"default:null" json:"calonKepala"`
			CalonSenatID  *int `gorm:"default:null" json:"calonSenat"`
		}

		var input vote

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "input is invalid",
				"error":   err.Error(),
			})
			return
		}

		if input.CalonKepalaID == nil || input.CalonSenatID == nil || *input.CalonKepalaID > 2 || *input.CalonSenatID > 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "input is invalid",
				"error":   "calon kepala dan calon senat tidak boleh kosong",
			})
			return
		}

		cek := model.Users{}

		search := db.Where("id = ?", ID).Take(&cek)
		if search = db.Where("id = ?", ID).Take(&cek); search.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   search.Error.Error(),
			})
			return
		}

		if cek.IsVote == true {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Anda sudah memilih",
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
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if result = db.Where("id = ?", ID).Take(&mahasiswa); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "mahasiswa not found.",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success":      true,
			"message":      "successfully updated data.",
			"NIM":          mahasiswa.NIM,
			"Nama":         mahasiswa.Nama,
			"voted":        mahasiswa.IsVote,
			"Calon Kepala": mahasiswa.CalonKepalaID,
			"Calon Senat":  mahasiswa.CalonSenatID,
		})
	})
}
