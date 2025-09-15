package server

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/salahfarzin/logger"
	"github.com/salahfarzin/roja-shop/configs"
	"github.com/salahfarzin/roja-shop/middlewares"
	"github.com/salahfarzin/roja-shop/pkg/db"
	"github.com/salahfarzin/roja-shop/router"
	"go.uber.org/zap"
)

func Run() {
	cfg := configs.New()
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("unable to parse ennvironment variables: %v", err)
	}

	// Initialize logger
	curPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory:", err)
	}
	currentDate := time.Now().Format("2006-01-02")
	loggerFile := "app-" + currentDate + ".log"
	logger.Init(&zap.Config{
		OutputPaths: []string{filepath.Join(curPath, cfg.Log.Path, loggerFile)},
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
	})

	conn, err := db.InitSqlite(cfg.DB.URL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	app := fiber.New()

	app.Use(middlewares.Config(*cfg))

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()

		logger.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("agent", c.Get("User-Agent")),
			zap.String("body", string(c.Body())),
			zap.Int("status", c.Response().StatusCode()),
		)

		return err
	})

	router.New(app, conn)

	app.Listen(fmt.Sprintf("%s:%d", cfg.URL, cfg.Port))
}
