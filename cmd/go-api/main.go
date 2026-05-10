package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drona-gyawali/learn-go/internal/config"
	"github.com/drona-gyawali/learn-go/internal/http/handlers/student"
)



func main() {
	cfg := config.MustLoad()

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students/create/", student.New("welcome to server"))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM,  syscall.SIGINT)

	slog.Info("Server started successfully")
	go func () {
		err:=http.ListenAndServe(cfg.Address, router)
		if err != nil {
			log.Fatal("Error occured while listening to the server")
		}
	}()

	<- done

	slog.Error("server is closing")

	ctx , cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	server_:=http.Server{
		Addr: cfg.Address,
		Handler: router,
	}


	err:=server_.Shutdown(ctx)
	if err != nil {
		slog.Error("error occured while shutting down %s", slog.String("error", err.Error()))
	}
	slog.Info("server closed")

}