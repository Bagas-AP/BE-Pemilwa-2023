package handler

import (
	"TestVote/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
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

		//var input model.Users
		save := model.Users{
			NIM:       user.Account.NIM,
			Nama:      user.Account.Nama,
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
		return
	})
}
