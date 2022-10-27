package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.DebugMode)
	os.Exit(m.Run())
}

// func newTestServer(t *testing.T, store db.Store) *Server {
// 	config := util.Config{
// 		TokenSymmetricKey:   util.RandomString(32),
// 		AccessTokenDuration: time.Minute,
// 	}

// 	server, err := NewServer(config, store)
// 	require.NoError(t, err)

// 	return server
// }
