package db

import (
	db "blogapi/db/sqlc"
	"blogapi/util"
	"database/sql"
	"fmt"
)

func ConnectToStore(config util.Config) db.Store { //*db.Queries {

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		fmt.Printf("error db: %v\n", err)
	}

	return db.NewStore(conn)
	//return db.New(conn)
}
