package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/neonxp/geezer"
	"github.com/neonxp/geezer/services/mongodb"

	"github.com/neonxp/geezer/example/services/hello"
)

var (
	listen   string
	mongoURI string
	mongoDB  string
)

func main() {
	flag.StringVar(&listen, "listen", ":3000", "Host and port to listen (ex: '0.0.0.0:3000')")
	flag.StringVar(&mongoURI, "mongo", "mongodb://localhost:27017/", "MongoDB connection uri (ex: 'mongodb://user:pass@sample.host:27017/')")
	flag.StringVar(&mongoDB, "mongo_db", "geezer", "Database name")
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	// MongoDB connection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Disconnect(context.Background())
	db := client.Database(mongoDB)

	app := geezer.NewHttpKernel()

	hello.RegisterHooks(app)                              // Register hooks
	_ = app.Register(hello.ServiceName, &hello.Service{}) // Register service as external handler

	_ = app.Register("test", mongodb.New[Product](db.Collection("test"))) // Register mongodb crud service

	log.Printf("Started on %s\n", listen)
	srv := http.Server{Addr: listen, Handler: app}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
