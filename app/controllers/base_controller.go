package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"AntiqueGo/app/models"
	"AntiqueGo/database/seeders"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/cli"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB        *gorm.DB
	Router    *mux.Router
	AppConfig *AppConfig
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
	AppURL  string
}

type DBConfig struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
}

type PageLink struct {
	Page          int64
	Url           string
	IsCurrentPage bool
}

type PaginationLinks struct {
	CurrentPage string
	NextPage    string
	PrevPage    string
	TotalRows   int64
	TotalPage   int64
	Links       []PageLink
}

type PaginationParams struct {
	Path        string
	TotalRows   int64
	PerPage     int64
	CurrentPage int64
}

type Result struct {
	Code 	int 			`json:"code"`
	Data 	interface{} 	`json:"data"`
	Message string 			`json:"message"`

}

var store *sessions.CookieStore
var sessionShoppingCart = "shopping-cart-session"
var sessionFlash = "flash-session"
var sessionUser = "user-session"

func SetSessionStore(sessionKey string) {
	store = sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // misalnya, sesuaikan dengan kebutuhan Anda
		HttpOnly: true,
		Secure:   true, // set true jika Anda memakai HTTPS
		SameSite: http.SameSiteStrictMode,
	}
}

func (s *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("welcome  to " + appConfig.AppName)

	s.InitializeDb(dbConfig)
	s.Router = mux.NewRouter()
	s.InitializeRoutes()
	s.InitializeAppConfig(appConfig)
	// seeders.DBSeed(s.DB)

}

func (s *Server) InitializeDb(dbConfig DBConfig) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	s.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

}

func (s *Server) InitializeAppConfig(appConfig AppConfig) {
	s.AppConfig = &appConfig
}

func (s *Server) dbMigrate() {
	for _, model := range models.RegisterModel() {
		err := s.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("database migrated successfully")

}

func (s *Server) Run(addr string) {
	fmt.Printf("listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, s.Router))
}

func (s *Server) InitCommands(config AppConfig, dbConfig DBConfig) {
	s.InitializeDb(dbConfig)

	cmdApp := cli.NewApp()
	cmdApp.Commands = []cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				s.dbMigrate()
				return nil
			},
		},
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(s.DB)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPaginationLinks(config *AppConfig, params PaginationParams) (PaginationLinks, error) {
	var links []PageLink

	totalPage := int64(math.Ceil(float64(params.TotalRows) / float64(params.PerPage)))
	for i := 1; int64(i) <= totalPage; i++ {
		links = append(links, PageLink{
			Page:          int64(i),
			Url:           fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(i)),
			IsCurrentPage: int64(i) == params.CurrentPage,
		})

	}

	var nextPage int64
	var prevPage int64

	prevPage = 1
	nextPage = totalPage

	if params.CurrentPage > 2 {
		prevPage = params.CurrentPage - 1
	}

	if params.CurrentPage < totalPage {
		nextPage = params.CurrentPage + 1
	}
	return PaginationLinks{
		CurrentPage: fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(params.CurrentPage)),
		NextPage:    fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(nextPage)),
		PrevPage:    fmt.Sprintf("%s/%s/?page=%s", config.AppURL, params.Path, fmt.Sprint(prevPage)),
		TotalRows:   params.TotalRows,
		TotalPage:   totalPage,
		Links:       links,
	}, nil

}

func (s *Server) GetProvinces() ([]models.Province, error) {
	url := fmt.Sprintf("%sprovince?key=%s", os.Getenv("API_ONGKIR_BASE_URL"), os.Getenv("API_ONGKIR_KEY"))
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get provinces: status code %d", response.StatusCode)
	}

	var provinceResponse models.ProvinceResponse
	err = json.NewDecoder(response.Body).Decode(&provinceResponse)
	if err != nil {
		return nil, err
	}

	return provinceResponse.ProvinceData.Results, nil
}

func (s *Server) GetCitiesByProvinceID(provinceID string) ([]models.City, error) {
	url := fmt.Sprintf("%scity?key=%s&province=%s", os.Getenv("API_ONGKIR_BASE_URL"), os.Getenv("API_ONGKIR_KEY"), provinceID)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get cities: status code %d", response.StatusCode)
	}

	var cityResponse models.CityResponse
	err = json.NewDecoder(response.Body).Decode(&cityResponse)
	if err != nil {
		return nil, err
	}

	return cityResponse.CityData.Results, nil

}

func (s *Server) CalculateShippingFee(shippingParams models.ShippingFeeParams)([]models.ShippingFeeOption,error) {
	if shippingParams.Origin == "" || shippingParams.Destination == "" || shippingParams.Weight <= 0 || shippingParams.Courier == "" {
        return nil, errors.New("invalid params")
    }

    params := url.Values{}
    params.Add("key", os.Getenv("API_ONGKIR_KEY"))
    params.Add("origin", shippingParams.Origin)
    params.Add("destination", shippingParams.Destination)
    params.Add("weight", strconv.Itoa(shippingParams.Weight))
    params.Add("courier", shippingParams.Courier)

    response, err := http.PostForm(os.Getenv("API_ONGKIR_BASE_URL")+"cost", params)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get shipping fee: status code %d", response.StatusCode)
    }

    var ongkirResponse models.OngkirResponse
    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    jsonErr := json.Unmarshal(body, &ongkirResponse)
    if jsonErr != nil {
        return nil, jsonErr
    }

    var shippingFeeOptions []models.ShippingFeeOption
    for _, result := range ongkirResponse.OngkirData.Results {
        for _, cost := range result.Costs {
            shippingFeeOptions = append(shippingFeeOptions, models.ShippingFeeOption{
                Service: cost.Service,
                Fee:     cost.Cost[0].Value,
            })
        }
    }

    return shippingFeeOptions, nil
}

func SetFlash(w http.ResponseWriter, r *http.Request,name string,value string){
	session,err := store.Get(r,sessionFlash)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	session.AddFlash(value,name)
	session.Save(r,w)
}

func GetFlash(w http.ResponseWriter, r *http.Request, name string) []string{
	session,err := store.Get(r,sessionFlash)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return nil
    }
	fm := session.Flashes(name)
	if len(fm)==0 {
		return nil
	}
	session.Save(r,w)
	var flashes []string
	for _,fl := range fm {
        flashes = append(flashes, fl.(string))
    }
	return flashes

}

func IsLoggedIn(r *http.Request)bool{
	session,_:=store.Get(r,sessionUser)
	if session.Values["id"]==nil {
        return false
    }
	return true
}

func ComparePassword(password string, hashedPassword string)bool{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))==nil
}

func MakePassword(password string) (string, error) {
   hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost )

   return string(hashedPassword), err
}


func (s *Server) CurrentUser(w http.ResponseWriter, r *http.Request) *models.User{
	if !IsLoggedIn(r){
		return nil
	}

	userModel := models.User{}

	session,_:= store.Get(r,sessionUser)
	user,err := userModel.FindByID(s.DB, session.Values["id"].(string))
	if err != nil {
		session.Values["id"] = nil
		session.Save(r,w)
		return nil
	}
	return user
}

