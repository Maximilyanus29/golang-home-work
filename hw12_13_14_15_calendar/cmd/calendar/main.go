package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	initConfig()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	logg := logger.New(config.Logger.Level)

	storage := memorystorage.New()

	calendar, err := app.New(logg, storage)
	if err != nil {
		logg.Error("failed to start app: " + err.Error())
		os.Exit(1)
	}

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func initConfig() {
	if configFile != "" {
		// Используем файл конфигурации из флага.
		viper.SetConfigFile(configFile)

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Использование файла конфигурации:", viper.ConfigFileUsed())
		} else {
			panic(123)

			fmt.Println("Не удалось получить конфигуарцию")

		}
	} else {
		panic(123)
		fmt.Println("Необходимо передать кронфигурацию")
	}

	viper.AutomaticEnv()
}
