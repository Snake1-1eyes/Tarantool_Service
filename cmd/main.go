package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpDelivery "github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/crudoshlep/delivery/http"
	tarantoolRepo "github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/crudoshlep/repo/tarantool"
	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/crudoshlep/usecase"
	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	log := logger.GetLogger()
	defer log.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  "tarantool:3301",
		User:     "storage",
		Password: "passw0rd",
	}
	opts := tarantool.Opts{
		Timeout: time.Second * 5,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		log.Fatal("Ошибка подключения к Tarantool",
			zap.Error(err),
			zap.String("address", dialer.Address),
		)
	}
	defer conn.Close()
	log.Info("Успешное подключение к Tarantool", zap.String("address", dialer.Address))

	repo := tarantoolRepo.NewTarantoolRepo(conn, log)
	useCase := usecase.NewKVUseCase(repo)
	handler := httpDelivery.NewHandler(useCase, log)

	r := mux.NewRouter()
	r.HandleFunc("/kv", handler.Create).Methods("POST")
	r.HandleFunc("/kv/{id}", handler.Get).Methods("GET")
	r.HandleFunc("/kv/{id}", handler.Update).Methods("PUT")
	r.HandleFunc("/kv/{id}", handler.Delete).Methods("DELETE")
	log.Info("Маршруты зарегистрированы")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Сервер запущен", zap.String("port", ":8080"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Ошибка запуска сервера",
				zap.Error(err),
			)
		}
	}()

	<-done
	log.Info("Получен сигнал завершения")

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Ошибка при graceful shutdown",
			zap.Error(err),
			zap.Duration("timeout", 30*time.Second),
		)
	} else {
		log.Info("Сервер успешно остановлен")
	}
}
