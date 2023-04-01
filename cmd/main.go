package main

import (
	"context"
	"database/sql"
	"fmt"

	"blogapi/api"
	db "blogapi/db/sqlc"
	"blogapi/internal/services"
	"blogapi/util"

	_ "github.com/lib/pq"
)

func main() {

	ctx := context.Background()

	config, err := util.LoadConfig(".")
	if err != nil {
		fmt.Println(err)
		//log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		fmt.Println(err)
		//log.Fatal().Err(err).Msg("cannot connect to db")
	}

	store := db.NewStore(conn) //.ConnectToStore(config)

	services.InitUser(ctx, store)

	runGinServer(ctx, config, store)

	// s, err := api.NewServer(ctx, config, store)
	// if err != nil {
	// 	fmt.Println(err)
	// 	//log.Fatal().Err(err).Msg("cannot load config")
	// }
	// s.Start(config.HTTPServerAddress)

}

func runGinServer(ctx context.Context, config util.Config, store db.Store) {
	server, err := api.NewServer(ctx, config, store)
	if err != nil {
		//log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		//log.Fatal().Err(err).Msg("cannot start server")
	}
}
