package main

import (
	"context"
	"database/sql"
	"os"

	"blogapi/api"
	db "blogapi/db/sqlc"
	"blogapi/internal/services"
	"blogapi/util"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	var log = logrus.New()

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
		//log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.GinMode == "debug" {
		log.Level = logrus.TraceLevel
		log.Warn("you are in development mode!")
	} else {
		log.Level = logrus.InfoLevel
		log.Warn("you are in production mode!")
	}

	log.Out = os.Stdout
	
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn) //.ConnectToStore(config)

	services.InitUser(ctx, store)

	runGinServer(log, ctx, config, store)

}

func runGinServer(log *logrus.Logger, ctx context.Context, config util.Config, store db.Store) {
	server, err := api.NewServer(log, ctx, config, store)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
