package handler

import (
	"TestVote/middleware"
	"TestVote/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminUser(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")

	// get all user
	r.GET("/admin/mahasiswa", middleware.Authorization(), func(c *gin.Context) {
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

		var users []model.Users
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

	// get mahasiswa by name or nim
	r.POST("/admin/mahasiswa", middleware.Authorization(), func(c *gin.Context) {
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

		name, _ := c.GetQuery("name")
		nim, _ := c.GetQuery("nim")

		q := c.Request.URL.Query()

		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 1:
			pageSize = 10
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize

		var queryResults []model.Users

		if res := db.Where("nama LIKE ?", "%"+name+"%").Where("nim LIKE ?", "%"+nim+"%").Offset(offset).Limit(pageSize).Find(&queryResults); res.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "user is not found.",
				"error":   res.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data": gin.H{
				"query": gin.H{
					"name": name,
					"nim":  nim,
				},
				"result": queryResults,
			},
		})

	})

	// get mahasiswa by id
	r.GET("/admin/mahasiswa/:id", middleware.Authorization(), func(c *gin.Context) {
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

		var mahasiswa model.Users

		if result := db.Where("id = ?", id).Take(&mahasiswa); result.Error != nil {
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
			"data":    mahasiswa,
		})

	})

	// update mahasiswa by id
	r.PATCH("/admin/mahasiswa/:id", middleware.Authorization(), func(c *gin.Context) {
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

		var mahasiswa model.Users

		if result := db.Where("id = ?", id).Take(&mahasiswa); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		var input model.Users
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

		var temp bool

		if input.CalonKepalaID == nil && input.CalonSenatID == nil {
			temp = false
		} else if input.CalonKepalaID != nil || input.CalonSenatID != nil {
			temp = false
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Dilarang iseng mengubah data!!!",
			})
		} else {
			temp = true
		}

		update := model.Users{
			ID:            mahasiswa.ID,
			Nama:          mahasiswa.Nama,
			NIM:           mahasiswa.NIM,
			Foto:          mahasiswa.Foto,
			IsVote:        temp,
			ISAdmin:       mahasiswa.ISAdmin,
			Prodi:         mahasiswa.Prodi,
			Tahun:         mahasiswa.Tahun,
			UpdatedAt:     time.Now(),
			CalonKepalaID: input.CalonKepalaID,
			CalonSenatID:  input.CalonSenatID,
		}

		if err := db.Select("*").Model(&mahasiswa).Updates(update).Error; err != nil {
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
			"data":    mahasiswa,
		})

	})

	// delete mahasiswa by id
	r.DELETE("/admin/mahasiswa/:id", middleware.Authorization(), func(c *gin.Context) {
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

		var mahasiswa model.Users

		if result := db.Where("id = ?", id).Take(&mahasiswa); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if err := db.Delete(&mahasiswa).Error; err != nil {
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
			"data":    mahasiswa,
		})

	})
}
