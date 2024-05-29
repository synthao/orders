package sso

import (
	"context"
	ssov1 "github.com/synthao/orders/gen/go/sso"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ClientParams struct {
	fx.In

	Config *Config
	Logger *zap.Logger
}

type Client struct {
	api ssov1.ServiceClient
}

func newConfig() (*Config, error) {
	provider, err := config.NewYAML(config.File("config.yml"))
	if err != nil {
		return nil, err
	}

	var c Config

	if err = provider.Get("sso").Populate(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

var Module = fx.Options(
	fx.Provide(
		newConfig,
		newGRPCClient,
		newSSOClient,
	),
	fx.Invoke(
		func(lc fx.Lifecycle, cnf *Config, grpcConn *grpc.ClientConn) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return nil
				},
				OnStop: func(ctx context.Context) error {
					if err := grpcConn.Close(); err != nil {
						return err
					}
					return nil
				},
			})
		},
	),
)
