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

	"github.com/aniket-mahakalkar/student_api/internal/config"
	"github.com/aniket-mahakalkar/student_api/internal/http/handlers/student"
)


func main(){

	//load config

	cfg := config.MustLoad()


	//database setup
	//setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students",student.New())


	//setup server

	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}
	
	slog.Info("server started %s",slog.String( "address" , cfg.HTTPServer.Addr))





	done := make(chan os.Signal , 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func ()  {

		err := server.ListenAndServe()
		if err != nil {

			log.Fatal("failed to start server")
	}
		
	}()

	<- done

	slog.Info("shutting down the server")

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()



	err := server.Shutdown(ctx)

	if err != nil {
		
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	

	slog.Info("server shutdown successfully")


	

}