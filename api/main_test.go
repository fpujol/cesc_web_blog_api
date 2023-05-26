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
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	log := logrus.New()
	ctx := context.Background()
	config := util.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(log, ctx, config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}