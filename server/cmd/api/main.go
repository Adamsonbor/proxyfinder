package main

import (
	"fmt"
	"net/http"
	chirouter "proxyfinder/internal/api/chi-router"
	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"
	gormstorage "proxyfinder/internal/storage/gorm-storage"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// INIT config
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	// INIT logger
	log := logger.New(cfg.Env)
	log.Info("Initializing with env: " + cfg.Env)
	
	// INIT gorm sqlite
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	// INIT storage 
	storage := gormstorage.New(db)
	
	// status := domain.Status{Name: "available"}
	// if err := storage.Create(&status); err != nil {
	// 	panic(err)
	// }

	// var statuses []domain.Status
	// if err := storage.GetAll(&statuses); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(statuses)
	//
	// var countries []domain.Country
	// if err := storage.GetAll(&countries); err != nil {
	// 	panic(err)
	// }
	// for _, v := range countries {
	// 	fmt.Println(v.Code, v.Name)
	// }
	router := chirouter.New(log, storage)
	http.ListenAndServe(":8080", router.Router)
}
