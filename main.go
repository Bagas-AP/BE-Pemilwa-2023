package main

import (
	"TestVote/database"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

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
	router := gin.Default()
	router.Use(cORS())
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "alive",
		})
	})
	router.POST("/auth", func(ctx *gin.Context) {
		user := NewUser()
		var input mahasiswa
		err := ctx.ShouldBindJSON(&input)
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
		type Data struct {
			NIM          string
			Nama         string
			Fakultas     string
			ProgramStudi string
			FotoProfile  string
		}
		var data = Data{
			NIM:          user.Account.NIM,
			Nama:         user.Account.Nama,
			Fakultas:     user.Account.Fakultas,
			ProgramStudi: user.Account.ProgramStudi,
			FotoProfile:  fmt.Sprintf("https://siakad.ub.ac.id/dirfoto/foto/foto_20%s/%s.jpg", user.Account.NIM[0:2], user.Account.NIM),
		}

		ctx.JSON(200, gin.H{
			"success": true,
			"data":    data,
		})
	})
	if err := router.Run(":5000"); err != nil {
		log.Fatal(err.Error())
		return
	}
}
