package main

import (
	"context"
	"log"

	"github.com/DimTur/lp_api_gateway/cmd/serve"
	_ "github.com/DimTur/lp_api_gateway/docs"
)

// @title           Learning Platform API
// @version         0.1.0
// @description     The project is only in its initial stages.

// @contact.name   API Support

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

// @security ApiKeyAuth
func main() {
	ctx := context.Background()

	cmd := serve.NewServeCmd()

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("smth went wrong: %s", err)
	}
}
