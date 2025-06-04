package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/spf13/cobra"
)

var (
	Verbose           bool
	CfgFile           string
	ErrConfigIsNotSet = errors.New("config is not set")

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
		Run: runRootCMD,
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
)

func main() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().StringVar(&CfgFile, "config", "../../configs/config.toml", "Path to configuration file default: ../../config/config.toml")
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}

func runRootCMD(cmd *cobra.Command, args []string) {
	config := config.NewConfig(CfgFile)
	logg := logger.New(config.Logger.Level)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

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
