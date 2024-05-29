package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/segmentio/kafka-go"
	"github.com/synthao/orders/internal/adapter/mysql/order/repository"
	"github.com/synthao/orders/internal/config"
	"github.com/synthao/orders/internal/database"
	sso2 "github.com/synthao/orders/internal/module/sso"
	port "github.com/synthao/orders/internal/port/http/order"
	"github.com/synthao/orders/internal/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net"
	"os"
)

func main() {
	fx.New(
		sso2.Module,
		fx.Provide(
			fiber.New,
			config.NewServerConfig,
			config.NewKafkaConfig,
			config.NewLoggerConfig,
			newLogger,
			newKafkaWrite,
			repository.NewRepository,
			service.NewService,
			port.NewHandler,
			database.NewConnection,
		),
		fx.Invoke(
			database.ApplyMigrations,
			newHTTPServer,
		),
	).Run()
}

func newLogger(cnf *config.Logger) (*zap.Logger, error) {
	atomicLogLevel, err := zap.ParseAtomicLevel(cnf.Level)
	if err != nil {
		return nil, err
	}

	atom := zap.NewAtomicLevelAt(atomicLogLevel.Level())
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		),
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
	), nil
}

func newHTTPServer(lc fx.Lifecycle, app *fiber.App, handler *port.Handler, cnf *config.Server, ssoClient *sso2.Client, db *sqlx.DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			app.Get("/ping", func(ctx *fiber.Ctx) error {
				return ctx.SendString("pong")
			})

			handler.InitRoutes(ssoClient)

			go app.Listen(net.JoinHostPort("", cnf.Port))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("==> stopping app ...")

			if err := db.Close(); err != nil {
				return fmt.Errorf("failed to close db: %w", err)
			}

			if err := app.Shutdown(); err != nil {
				return fmt.Errorf("failed to shutdown app: %w", err)
			}

			return nil
		},
	})
}

func newKafkaWrite(lc fx.Lifecycle, cnf *config.Kafka) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cnf.Address),
		Topic:                  cnf.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true, // May want to disable in production
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// Hook into lifecycle startup here if needed
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Make sure to close the writer when the service is stopped.
			return writer.Close()
		},
	})

	return writer
}
