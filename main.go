package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kadblue/placeclone/middleware"
	"github.com/kadblue/placeclone/placeclone"
	"log"
)

var adapter *gorillamux.GorillaMuxAdapter

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	//client, err := cache.NewClient("/cachedb")
	//if err != nil {
	//	fmt.Println("cache client could not be created")
	//}

	//client.ClearAll()

	if err != nil {
		fmt.Println(err.Error())
	}

	DbCli := dynamodb.NewFromConfig(cfg)
	sessionStore := sessions.NewCookieStore([]byte("aksjdfjjlasdfjlkjlasdf"))
	userpoolId := "ap-south-1_DTkRR7wmN"
	middlewareServer := middleware.NewAuthMiddlewareServer(sessionStore, nil, DbCli, userpoolId)

	mainRouter := mux.NewRouter()

	mainRouter.Use(middleware.CorsMiddleware)

	placecloneServerOptions := &placeclone.Options{
		DbCli: DbCli,
		Store: sessionStore,
		//CacheCli:       client,
		AuthMiddleware: middlewareServer,
	}

	//authServerOptions := &auth.Options{
	//	DbCli: DbCli,
	//	Store: sessionStore,
	//}

	placeclone.AddSubrouter(placecloneServerOptions, mainRouter)
	//auth.AddSubrouter(authServerOptions, mainRouter)

	//http.Handle("/", mainRouter)
	//
	//http.ListenAndServe(":8000", nil)

	adapter = gorillamux.New(mainRouter)

}

func LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	fmt.Println("req:", req.Path)
	resp, err := adapter.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&req))

	return *resp.Version1(), err

}

func main() {

	lambda.Start(LambdaHandler)
}
