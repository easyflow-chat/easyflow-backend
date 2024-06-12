package main

import (
	"easflow-backend/config"
	"easflow-backend/database"
)

func main() {
	cfg := config.LoadDefaultConfig()
	dbInst, err := database.NewDatabaseInst(cfg.DatabaseURL, &cfg.GormConfig)
	if err != nil {
		panic(err)
	}
	dbInst.Acquire()
	dbInst.Migrate()

	defer dbInst.Release() // En	sure to release when done
}
