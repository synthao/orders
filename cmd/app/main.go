package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/segmentio/kafka-go"
	"github.com/synthao/orders/internal/adapter/mysql/order/repository"
	"github.com/synthao/orders/internal/config"
	"github.com/synthao/orders/internal/database"
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
		fx.Provide(
			fiber.New,
			config.NewServerConfig,
			config.NewKafkaConfig,
			config.NewLoggerConfig,
			newLogger,
			kafkaProducer,
			repository.NewRepository,
			service.NewService,
			port.NewHandler,
			database.NewConnection,
		),
		fx.Invoke(createHTTPServer),
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

func createHTTPServer(lc fx.Lifecycle, app *fiber.App, handler *port.Handler, cnf *config.Server, db *sqlx.DB) {
	database.ApplyMigrations(db)

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	handler.InitRoutes()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(net.JoinHostPort("", cnf.Port))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("==> stopping server ...")

			return app.Shutdown()
		},
	})
}

func kafkaProducer(lc fx.Lifecycle, cnf *config.Kafka) *kafka.Writer {
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
