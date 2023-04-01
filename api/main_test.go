package api

import (
	"context"
	"os"
	"testing"
	"time"

	db "blogapi/db/sqlc"
	"blogapi/pkg/utils"
	"blogapi/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	ctx := context.Background()
	config := util.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(ctx, config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}