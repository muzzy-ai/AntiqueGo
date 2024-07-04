package app

import (
	"flag"
	// "fmt"
	"log"
	// "net/http"
	"os"

	// "AntiqueGo/database/seeders"
	"AntiqueGo/app/controllers"

	// "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	// "github.com/urfave/cli"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
)



func getEnv(key,fallback string) string{
	if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}


func Run(){
	appConfig:= controllers.AppConfig{}
	dbConfig := controllers.DBConfig{}
	server := controllers.Server{}

	
	err:=godotenv.Load()
	if err != nil{
		log.Fatal("Error loading.env file")
	}



	appConfig.AppName = getEnv("APP_NAME","AntiqueGo")
	appConfig.AppEnv = getEnv("APP_ENV","development")
	appConfig.AppPort = getEnv("APP_PORT","8080")
	appConfig.AppURL = getEnv("APP_URL","http://localhost:8080")

	
	dbConfig.DBName = getEnv("DB_NAME","AntiqueGo")
	dbConfig.DBUser = getEnv("DB_USER","postgres")
	dbConfig.DBPassword = getEnv("DB_PASS","admin")
	dbConfig.DBHost = getEnv("DB_HOST","localhost")
	dbConfig.DBPort = getEnv("DB_PORT","5432")

	sessionKey := getEnv("SESSION_KEY","")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY is not set")
	}
	controllers.SetSessionStore(sessionKey)

	flag.Parse()
	arg := flag.Arg(0)
	if len(arg) > 0 {
        server.InitCommands(appConfig, dbConfig)
    }else {
		server.Initialize(appConfig,dbConfig)
		server.Run(":"+appConfig.AppPort)
	}


}