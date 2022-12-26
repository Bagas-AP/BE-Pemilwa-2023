package handler

import (
	"TestVote/middleware"
	"TestVote/model"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func AdminKepala(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")

	// create calon kepala
	r.POST("/admin/kepala", middleware.Authorization(), func(c *gin.Context) {
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

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "get form err: " + err.Error(),
			})
			return
		}

		rand.Seed(time.Now().Unix())

		str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

		shuff := []rune(str)

		rand.Shuffle(len(shuff), func(i, j int) {
			shuff[i], shuff[j] = shuff[j], shuff[i]
		})
		file.Filename = string(shuff)

		if err := c.SaveUploadedFile(file, "./Images/"+file.Filename); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Success": false,
				"error":   "upload file err: " + err.Error(),
			})
			return
		}

		godotenv.Load("../.env")
		newKepala := model.CalonKepala{
			Nama: c.PostForm("nama"),
			Foto: os.Getenv("BASE_URL") + "/api/admin/kepala/" + file.Filename,
		}

		if err := db.Create(&newKepala); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "error when inserting a new kepala",
				"error":   err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "a new group has successfully created",
			"error":   nil,
			"nama":    newKepala.Nama,
		})

		// get all calon kepala
		r.GET("/admin/kepala", middleware.Authorization(), func(c *gin.Context) {
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
			var kepala []model.CalonKepala

			if result := db.Find(&kepala); result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when querying the database.",
					"error":   result.Error.Error(),
				})
				return
			}

			var ret []model.CalonKepala

			for _, value := range kepala {
				var temp model.CalonKepala
				temp.IDKepala = value.IDKepala
				temp.Nama = value.Nama
				temp.Foto = value.Foto
				ret = append(ret, temp)
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "query completed.",
				"data":    ret,
			})
		})

	})

	// get calon kepala by id
	r.GET("/admin/kepala/:id", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user model.Users

		if !user.ISAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "unauthorized access :(",
				"error":   nil,
			})
			return
		}

		if err := db.Where("id = ?", ID).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query completed.",
			"data":    user,
		})

	})

	// untuk memperbarui data kepala by id
	r.PATCH("/admin/kepala/:id", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user model.Users

		if !user.ISAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "unauthorized access :(",
				"error":   nil,
			})
			return
		}

		if err := db.Where("id = ?", ID).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
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

		file, err := c.FormFile("foto")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "get form err: " + err.Error(),
			})
			return
		}

		rand.Seed(time.Now().Unix())

		str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

		shuff := []rune(str)

		rand.Shuffle(len(shuff), func(i, j int) {
			shuff[i], shuff[j] = shuff[j], shuff[i]
		})
		file.Filename = string(shuff)

		if err := c.SaveUploadedFile(file, "./images/"+file.Filename); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Success": false,
				"error":   "upload file err: " + err.Error(),
			})
			return
		}

		godotenv.Load("../.env")
		parsedId, _ := strconv.ParseUint(id, 10, 32)
		newKepala := model.CalonKepala{
			IDKepala: uint(parsedId),
			Nama:     c.PostForm("nama"),
			Foto:     os.Getenv("BASE_URL") + "/api/admin/kepala/" + file.Filename,
		}

		result := db.Where("id = ?", id).Model(&newKepala).Select("*").Updates(newKepala)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if result = db.Where("id = ?", id).Take(&newKepala); result.Error != nil {
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
				"message": "group not found.",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Update successful.",
			"data":    newKepala,
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
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}

		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}

		kepala := model.CalonKepala{
			IDKepala: uint(parsedId),
		}

		if result := db.Delete(&kepala); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete successful.",
		})
	})

}
