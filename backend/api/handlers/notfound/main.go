package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	handler := http.NotFoundHandler()
	function := httpadapter.New(handler).ProxyWithContext
	lambda.Start(function)
}
