package main

import (
	"TestVote/database"
	"TestVote/handler"
	"TestVote/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var cORS = func() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	// Connect Database
	db := database.Open()
	if db != nil {
		println("Nice, DB Connected")
	}

	// Gin Framework
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	r.Use(cORS())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "alive",
		})
	})

	// api untuk add senat baru
	r.POST("/post-senat", func(c *gin.Context) {
		var inputSenat model.CalonSenat

		err := c.BindJSON(&inputSenat)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message":    "format input tidak valid",
				"error":      err.Error(),
				"success":    false,
				"statusCode": http.StatusBadRequest,
			})
			return
		}

		newSenat := model.CalonSenat{
			NamaSenat: inputSenat.NamaSenat,
			Foto:      inputSenat.Foto,
		}

		if err := db.Create(&newSenat); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":    "failed when creating a new data user",
				"success":    false,
				"error":      err.Error.Error(),
				"statusCode": http.StatusInternalServerError,
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":    "successfully add new data",
			"error":      nil,
			"data":       newSenat.NamaSenat,
			"statusCode": http.StatusCreated,
		})
	})

	// api untuk get detail user by id
	r.GET("/get-detail/:id", func(c *gin.Context) {
		id, _ := c.Params.Get("id")

		var getUser model.Users

		if err := db.Preload("CalonSenat").Preload("CalonKepala").Where("id = ?", id).Take(&getUser); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":    "failed get detail data",
				"success":    false,
				"error":      err.Error.Error(),
				"statusCode": http.StatusInternalServerError,
			})
			return
		}

		// type DetailUser struct {
		// 	NIM string
		// 	Nama string
		// 	Foto string
		// 	IsVote bool
		// 	NamaSenat string
		// 	NamaKepala string
		// }

		// getDetail := DetailUser {
		// 	NIM: getUser.NIM,
		// 	Nama: getUser.Nama,
		// 	Foto: getUser.Foto,
		// 	IsVote: getUser.IsVote,
		// 	NamaSenat: getUser.CalonSenat.Nama,
		// 	NamaKepala: getUser.CalonKepala.Nama,
		// }

		c.JSON(http.StatusOK, gin.H{
			"message":    "successfully add new data",
			"error":      nil,
			"data":       getUser,
			"statusCode": http.StatusOK,
		})
	})

	// pengelompokan api serta fungsinya
	r.Group("/api")
	handler.Login(db, r)
	handler.Vote(db, r)

	if err := r.Run(":8081"); err != nil {
		log.Fatal(err.Error())
		return
	}
}
