package api

import (
	"context"

	db "blogapi/db/sqlc"
)

var ctx context.Context
var queries *db.Queries

func Init(c context.Context, q *db.Queries) {
	ctx = c
	queries = q
}
