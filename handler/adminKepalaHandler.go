package handler

import (
	"TestVote/middleware"
	"TestVote/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func AdminKepala(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")

	// create calon kepala
	r.POST("/admin/loginKepala", middleware.Authorization(), func(c *gin.Context) {
		user := NewUser()
		var input Mahasiswa
		err := c.ShouldBindJSON(&input)
		if err != nil {
			log.Println(err.Error())
			return
		}

		x := input.Nim[5:8]

		if x != "020" && x != "030" && x != "021" && x != "031" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "NIM anda tidak valid!",
			})
			return
		}

		err = user.Login(input.Nim, input.Password)

		if err != nil {
			log.Println(err.Error())
			return
		}

		err = user.GetData()
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = user.Logout()
		if err != nil {
			log.Println(err.Error())
			return
		}

		save := model.CalonKepala{
			Nama: user.Account.Nama,
			Foto: fmt.Sprintf("https://siakad.ub.ac.id/dirfoto/foto/foto_20%s/%s.jpg", user.Account.NIM[0:2], user.Account.NIM),
		}

		if err := db.Create(&save); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong with user creation",
				"error":   err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Selamat Datang Calon Kepala Baru!",
			"data":    save,
		})
		return
	})

	// get calon kepala by id
	r.GET("/admin/kepala/:id", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user model.Users
		if err := db.Where("id = ?", ID).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}

		if !user.ISAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "unauthorized access :(",
				"error":   nil,
			})
			return
		}

		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"Success": false,
				"message": "id is not available",
			})
			return
		}

		var kepala model.CalonKepala

		if result := db.Where("id_kepala = ?", id).Take(&kepala); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query completed.",
			"data":    kepala,
		})

	})

	// untuk memperbarui data kepala by id
	r.PATCH("/admin/kepala/:id", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user model.Users
		if err := db.Where("id = ?", ID).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}

		if !user.ISAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "unauthorized access :(",
				"error":   nil,
			})
			return
		}

		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"Success": false,
				"message": "id is not available",
			})
			return
		}

		var kepala model.CalonKepala

		if result := db.Where("id_kepala = ?", id).Take(&kepala); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		var input model.CalonKepala
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid input",
				"error":   err.Error(),
			})
			return
		}

		//if input.CalonKepalaID == nil || input.CalonSenatID == nil || *input.CalonKepalaID > 2 || *input.CalonSenatID > 2 {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"success": false,
		//		"message": "input is invalid",
		//		"error":   "calon kepala dan calon senat tidak boleh kosong",
		//	})
		//	return
		//}

		update := model.CalonKepala{
			IDKepala: kepala.IDKepala,
			Nama:     input.Nama,
			Foto:     input.Foto,
		}

		if err := db.Select("*").Model(&kepala).Updates(update).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "update completed.",
			"data":    kepala,
		})

	})

	// untuk menghapus data kepala by id
	r.DELETE("/admin/kepala/:id", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user model.Users
		if err := db.Where("id = ?", ID).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}

		if !user.ISAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "unauthorized access :(",
				"error":   nil,
			})
			return
		}

		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"Success": false,
				"message": "id is not available",
			})
			return
		}

		var kepala model.CalonKepala

		if result := db.Where("id_kepala = ?", id).Take(&kepala); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if err := db.Delete(&kepala).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting the database.",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "delete completed.",
			"data":    kepala,
		})
	})

}
