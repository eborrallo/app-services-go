package bootstrap

import (
	"app-services-go/cmd/api/http/rest"
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

	ctx, srv := rest.New(context.Background(), cfg, c, q, r)
	return srv.Run(ctx)
}
