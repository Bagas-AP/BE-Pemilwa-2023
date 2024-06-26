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

func AdminSenat(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")

	// get all user
	r.GET("/admin/senat", middleware.Authorization(), func(c *gin.Context) {
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

		var users []model.CalonSenat
		if res := db.Find(&users); res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   res.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query completed.",
			"users":   users,
		})
	})

	// create calon senat
	r.POST("/admin/loginSenat", middleware.Authorization(), func(c *gin.Context) {
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

		save := model.CalonSenat{
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
			"message": "Selamat Datang Calon Senat Baru!",
			"data":    save,
		})
		return
	})

	// get calon senat by id
	r.GET("/admin/senat/:id", middleware.Authorization(), func(c *gin.Context) {
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

		var senat model.CalonSenat

		if result := db.Where("id_senat = ?", id).Take(&senat); result.Error != nil {
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
			"data":    senat,
		})

	})

	// untuk memperbarui data senat by id
	r.PATCH("/admin/senat/:id", middleware.Authorization(), func(c *gin.Context) {
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

		var senat model.CalonSenat

		if result := db.Where("id_senat = ?", id).Take(&senat); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		var input model.CalonSenat
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

		update := model.CalonSenat{
			IDSenat: senat.IDSenat,
			Nama:    input.Nama,
			Foto:    input.Foto,
		}

		if err := db.Select("*").Model(&senat).Updates(update).Error; err != nil {
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
			"data":    senat,
		})

	})

	// untuk menghapus data senat by id
	r.DELETE("/admin/senat/:id", middleware.Authorization(), func(c *gin.Context) {
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

		var senat model.CalonSenat

		if result := db.Where("id_senat = ?", id).Take(&senat); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if err := db.Delete(&senat).Error; err != nil {
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
			"data":    senat,
		})
	})

}
