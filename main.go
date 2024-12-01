package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/AlvinTendio/minder/config"
	viper_cfg "github.com/AlvinTendio/minder/config/viper"
	common_http "github.com/AlvinTendio/minder/delivery/http"
	"github.com/AlvinTendio/minder/mysql"
	_ "github.com/go-sql-driver/mysql"

	minder_delivery "github.com/AlvinTendio/minder/minder/delivery/http"
	minder_repo "github.com/AlvinTendio/minder/minder/repository"
	minder_usecase "github.com/AlvinTendio/minder/minder/usecase"
)

var config core_config.Config

func getConfig() (core_config.Config, error) {
	config := viper_cfg.NewConfig("/opt/secret/minder")
	return config, nil
}

func main() {
	var err error
	config, err = getConfig()
	if err != nil {
		log.Println(err)
		return
	}

	dbServerName := config.GetString(`database.server.name`)
	dbHost := config.GetString(`database.host`)
	dbPort := config.GetString(`database.port`)
	dbUser := config.GetString(`database.user`)
	dbPass := config.GetString(`database.pass`)
	dbName := config.GetString(`database.name`)
	maxOpen := config.GetInt(`database.max.open`)
	maxIdle := config.GetInt(`database.max.idle`)
	maxLifetime := config.GetInt(`database.max.lifetime`)

	dbConfig := &mysql.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPass,
		Name:     dbName,
	}

	dbConn, err := mysql.DB(dbConfig,
		mysql.WithMysql(dbServerName, true, "Asia/Jakarta"),
		mysql.WithConnection(int(maxOpen), int(maxIdle),
			time.Duration(maxLifetime)*time.Minute, 0))
	if err != nil {
		log.Println(err)
		return
	}

	err = dbConn.Ping()
	if err != nil {
		log.Println("Error ping DB")
		panic(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Println(err)
		}
	}()

	addr := config.GetString("server.address")
	ctx := context.Background()

	// minder
	minderRepo := minder_repo.NewMinderRepositoryImpl(dbConn)
	minderUsecase := minder_usecase.NewMinderUsecaseImpl(minderRepo)
	minder_delivery.NewMinderHandler(minderUsecase)

	go func() {
		if err := serveHTTP(ctx, addr, config); err != nil {
			log.Println(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	log.Println("All server stopped!")
}

func serveHTTP(ctx context.Context, addr string,
	config core_config.Config) error {
	restHandler := common_http.NewRestHandlerService(config)
	log.Println(ctx, "http server started. Listening on port: "+addr)

	if err := http.ListenAndServe(addr,
		common_http.NewHandler(http.HandlerFunc(restHandler.Serve),
			common_http.WithCORS())); err != nil {
		log.Println("err ->", err)
		return err
	}

	return nil
}
