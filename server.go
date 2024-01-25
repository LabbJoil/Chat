package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	service "github.com/LabbJoil/Chat/Services"
	"github.com/LabbJoil/Chat/graph"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	newServiceChat := service.ChatInteraction{}
	newServiceChat.ReceivedMessages = make(map[string][]string)
	if err := newServiceChat.ConnectDB(); err != nil {
		return
	}

	host := viper.GetString("qraphql.host")
	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("qraphql.port")
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{ServiceChat: newServiceChat}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://%s:%s/ for GraphQL playground", host, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initConfig() error {
	viper.AddConfigPath("Config")
	viper.SetConfigName("mainConfiguration")
	return viper.ReadInConfig()
}
