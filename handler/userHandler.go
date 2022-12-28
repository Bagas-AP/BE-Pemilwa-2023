package handler

import (
	"TestVote/middleware"
	"TestVote/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Login(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	r.POST("/login", func(c *gin.Context) {
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

		var save model.Users

		// setting tahun
		var temp string
		if user.Account.NIM[0:2] == "19" {
			temp = "2019"
		} else if user.Account.NIM[0:2] == "20" {
			temp = "2020"
		} else if user.Account.NIM[0:2] == "21" {
			temp = "2021"
		} else if user.Account.NIM[0:2] == "22" {
			temp = "2022"
		}

		//var input model.Users
		save = model.Users{
			NIM:       user.Account.NIM,
			Nama:      user.Account.Nama,
			Prodi:     user.Account.ProgramStudi,
			Tahun:     temp,
			Foto:      fmt.Sprintf("https://siakad.ub.ac.id/dirfoto/foto/foto_20%s/%s.jpg", user.Account.NIM[0:2], user.Account.NIM),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&save); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong with user creation",
				"error":   err.Error.Error(),
			})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"id":  save.ID,
			"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
		})
		godotenv.Load("../.env")
		strToken, err := token.SignedString([]byte(os.Getenv("TOKEN_G")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Welcome, here's your token. don't lose it ;)",
			"data": gin.H{
				"token": strToken,
			},
		})
	})

	r.GET("/profile", middleware.Authorization(), func(c *gin.Context) {
		id, isIdExists := c.Get("id")

		if !isIdExists {
			c.JSON(http.StatusForbidden, gin.H {
				"success": !isIdExists,
				"message": "user belum registrasi",
				"statusCode": http.StatusForbidden,
				"error": nil,
			})
			return
		}
		
		var user model.Users

		err := db.Where("id = ?", id).Take(&user)
		if err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"success": false,
				"message": "error ketika melakukan query data user",
				"statusCode": http.StatusInternalServerError,
				"error": err.Error.Error(),
			})
			return
		}

		if err.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H {
				"success": false,
				"message": "user tidak ditemukan",
				"statusCode": http.StatusNotFound,
				"error": nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "data user ditemukan",
			"statusCode": http.StatusOK,
			"data": user,
		})
	})
}
