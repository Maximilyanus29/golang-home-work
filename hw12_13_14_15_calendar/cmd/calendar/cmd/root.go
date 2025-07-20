/*
Copyright © 2025 Maxim Ryabtsev <Max-r2010@mail.ru>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/config"
	filelogger "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/logger/file"
	internalhttp "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "calendar",
		Short: `Сервис "Календарь" представляет собой максимально упрощенный сервис для хранения календарных событий и отправки уведомлений.`,
		Long: `Сервис предполагает возможность:
* добавить/обновить событие;
* получить список событий на день/неделю/месяц;
* получить уведомление за N дней до события.

Сервис НЕ предполагает:
* авторизации;
* разграничения доступа;
* web-интерфейса.`,
		Run: action,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func action(cmd *cobra.Command, args []string) {
	config := config.NewConfig(&cfgFile)

	logg := getLogger(config.Logger)

	calendar := app.New(logg, getStorage(config.Storage))

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

func getStorage(storageConf config.StorageConf) app.Storage {
	switch storageConf.Type {
	case "memory":
		return memorystorage.New()
	case "sql":
		return sqlstorage.New()
	}

	app.Exit("storage type unknown")
	return nil
}

func getLogger(loggerConf config.LoggerConf) app.Logger {
	switch loggerConf.Type {
	case "file":
		return filelogger.New(loggerConf)
	}

	app.Exit("logger type unknown")
	return nil
}
