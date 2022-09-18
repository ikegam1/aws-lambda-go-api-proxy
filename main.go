package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambdaALB

//var echoLambda *echoadapter.EchoLambda

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

var e = createMux()

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("echo cold start")

	e.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello world!")
	})
	e.POST("/t", func(c echo.Context) error {
		log.Println("routing t")
		return c.JSON(http.StatusOK, "hello world!")
	})

	echoLambda = echoadapter.NewALB(e)
	//echoLambda = echoadapter.New(e)
}

func handler(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	//func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	log.Printf("req: %#v", req)
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("reqBody: %#v", string(reqBody))
	}))
	e.HTTPErrorHandler = JSONErrorHandler

	return e
}

func JSONErrorHandler(err error, c echo.Context) {
	msg := err.Error()
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	c.JSON(code, msg)
}
