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

}
