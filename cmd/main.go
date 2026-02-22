package main

import (
	"github.com/prasannakumar414/somq/http"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	server := http.NewServer(8090, *logger)
	err := server.Serve()
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
