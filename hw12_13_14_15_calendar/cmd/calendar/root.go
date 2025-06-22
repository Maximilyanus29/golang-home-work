/*
Copyright © 2025 Maxim Ryabtsev <Max-r2010@mail.ru>
*/
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(ErrCodeExecuteFiled)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Fprintln(os.Stderr, "flag config is required")
		os.Exit(ErrCodeNoConfig)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "read config error")
		os.Exit(ErrCodeNoConfig)
	}
}

func action(cmd *cobra.Command, args []string) {
	config := config.NewConfig()
	logg := logger.New(config.Logger.Level, os.Stderr)

	calendar := app.New(logg, getStorage(config))

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
		os.Exit(ErrCodeExecuteFiled) //nolint:gocritic
	}
}

func getStorage(config config.Config) app.Storage {
	switch config.Storage.Type {
	case "memory":
		return memorystorage.New()
	case "sql":
		return sqlstorage.New()
	}
	return nil
}
