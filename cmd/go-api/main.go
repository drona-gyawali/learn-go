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
	"github.com/drona-gyawali/learn-go/internal/storage/sqllite"
)



func main() {
	cfg := config.MustLoad()

	storage , err := sqllite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Service Configured")
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students/create/", student.New("welcome to server", storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students/", student.GetStudentList(storage))
	router.HandleFunc("PATCH /api/students/", student.UpdateStudentView(storage))

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


	_err := server_.Shutdown(ctx)
	if _err != nil {
		slog.Error("error occured while shutting down %s", slog.String("error", _err.Error()))
	}
	slog.Info("server closed")

}