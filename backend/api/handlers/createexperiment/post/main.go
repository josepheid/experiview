package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/josepheid/experiview/api/handlers/createexperiment"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	region := "eu-west-2"

	ctx := context.Background()
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	experimentsTableName := os.Getenv("EXPERIMENTS_TABLE_NAME")

	// If the environment variable is not set
	if experimentsTableName == "" {
		logger.Error("environment variable EXPERIMENTS_TABLE_NAME is not set", "error", err)
		os.Exit(1)
	}

	ddbc := dynamodb.NewFromConfig(config)

	h, err := createexperiment.NewHandler(logger, ddbc, experimentsTableName)
	if err != nil {
		logger.Error("could not create handler", slog.Any("error", err))
		os.Exit(1)
	}
	function := httpadapter.New(h).ProxyWithContext
	lambda.Start(function)
}
