package bootstrap

import (
	"app-services-go/cmd/api/http/rest"
	"app-services-go/cmd/api/http/view"
	"app-services-go/cmd/api/http"
	"app-services-go/configs"
	"context"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	engine *gin.Engine
}
func Run() error {
	c, q, r, e := Container()
	var cfg, err = configs.GetServerConfig()
	if err != nil {
		return err
	}

	if e != nil {
		return e
	}

	ctx, srv := http.New(context.Background(), cfg)

	rest.RegisterRoutes(srv.Engine, c, q, r)
	view.RegisterRoutes(srv.Engine)

	return srv.Run(ctx)
}
