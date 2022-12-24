package main

import (
	"TestVote/database"
	"TestVote/middleware"
	"TestVote/model"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

var (
	loginUrl  = "https://siam.ub.ac.id/index.php/"  //POST
	siamUrl   = "https://siam.ub.ac.id/"            //GET
	logoutUrl = "https://siam.ub.ac.id/logout.php/" //GET

	Version = "0.1.0"

	ErrorNotLoggedIn = errors.New("please login first")
	ErrorLoggedIn    = errors.New("already logged in")
)

type User struct {
	c       *colly.Collector
	Account struct {
		NIM          string
		Nama         string
		Fakultas     string
		ProgramStudi string
	}
	LoginStatus bool
}

// constructor
func NewUser() User {
	return User{c: colly.NewCollector(), LoginStatus: false}
}

func (s *User) Login(us string, ps string) error {
	if s.LoginStatus {
		return ErrorLoggedIn
	}

	var errLogin error
	var doc *goquery.Document

	s.c.OnResponse(func(r *colly.Response) {
		doc, errLogin = goquery.NewDocumentFromReader(strings.NewReader(string(r.Body)))
		if errLogin != nil {
			errLogin = errors.New("couldn't read response body")
			return
		}
		temp := errors.New(strings.TrimSpace(doc.Find("small.error-code").Text()))
		if temp != nil {
			if len(temp.Error()) != 0 {
				errLogin = temp
				return
			}
		}
	})
	err := s.c.Post(loginUrl, map[string]string{
		"username": us,
		"password": ps,
		"login":    "Masuk",
	})

	if err != nil {
		if err.Error() != "Found" {
			return err
		}
	}
	if errLogin != nil {
		if len(errLogin.Error()) != 0 {
			return errLogin
		}
	}
	s.LoginStatus = true
	return nil
}

func (s *User) GetData() error {
	//scraping data mahasiswas
	result := make([]string, 8)
	s.c.OnHTML("div[class=\"bio-info\"]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			each := strings.TrimSpace(h.Text)
			if each != "PDDIKTI KEMDIKBUDDetail" {
				result[i] = h.Text
			}
		})
	})
	err := s.c.Visit(siamUrl)
	if err != nil {
		return err
	}

	s.Account.NIM = result[0]
	s.Account.Nama = result[1]
	// result2 = Jenjang/Fakultas--/--
	jenjangFakultas := strings.Split(result[2][16:], "/")
	s.Account.Fakultas = jenjangFakultas[1]
	s.Account.ProgramStudi = result[4][13:]
	return nil
}

// make sure to defer this method after login, so the phpsessionid won't be misused
func (s *User) Logout() error {
	if !s.LoginStatus {
		return ErrorNotLoggedIn
	}
	s.c.Visit(logoutUrl)
	return nil
}

type mahasiswa struct {
	Nim      string `json:"nim"`
	Password string `json:"password"`
}

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

	//Database
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
	r.POST("/auth", func(c *gin.Context) {
		user := NewUser()
		var input mahasiswa
		err := c.ShouldBindJSON(&input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		// if input.Nim == "215150200111006" {
		// 	fmt.Println("NIM anda tidak valid")
		// 	return
		// }
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
		// type Data struct {
		// 	NIM          string
		// 	Nama         string
		// 	Fakultas     string
		// 	ProgramStudi string
		// 	FotoProfile  string
		// }

		//var data = Data{
		//	NIM:          user.Account.NIM,
		//	Nama:         user.Account.Nama,
		//	Fakultas:     user.Account.Fakultas,
		//	ProgramStudi: user.Account.ProgramStudi,
		//	FotoProfile:  fmt.Sprintf("https://siakad.ub.ac.id/dirfoto/foto/foto_20%s/%s.jpg", user.Account.NIM[0:2], user.Account.NIM),
		//}

		//c.JSON(200, gin.H{
		//	"success": true,
		//	"data":    data,
		//})

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

		//c.JSON(http.StatusOK, gin.H{
		//	"success": true,
		//	"message": "User created successfully",
		//	"data":    save,
		//})
	})

	r.POST("/vote", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		type vote struct {
			CalonKepalaID int `gorm:"default:null" json:"calonKepala"`
			CalonSenatID  int `gorm:"default:null" json:"calonSenat"`
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

		if input.CalonKepalaID == 0 || input.CalonSenatID == 0 || input.CalonKepalaID > 2 || input.CalonSenatID > 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "input is invalid",
				"error":   "calon kepala dan calon senat tidak boleh kosong",
			})
			return
		}

		mahasiswa := model.Users{
			CalonKepalaID: input.CalonKepalaID,
			CalonSenatID:  input.CalonSenatID,
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
			"Calon Kepala": mahasiswa.CalonKepalaID,
			"Calon Senat":  mahasiswa.CalonSenatID,
		})
	})

	if err := r.Run(":8081"); err != nil {
		log.Fatal(err.Error())
		return
	}
}
