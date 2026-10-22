package bootstrap

import (
	"app-services-go/cmd/api/http"
	"app-services-go/cmd/api/http/rest"
	"app-services-go/cmd/api/http/view"
	"app-services-go/configs"
	"context"

	_ "github.com/go-sql-driver/mysql"
)

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
