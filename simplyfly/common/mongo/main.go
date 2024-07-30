package mongo

import (
	"context"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"os/signal"
	"simplifly/common/constants"
	"sync"
	"syscall"
)

var (
	client *mongo.Client
	once   sync.Once
)

func Instance() *mongo.Client {
	once.Do(func() {
		connectionString := constants.MongoConnectionString
		uri := connectionString
		var err error
		client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		//log.Print("mongo connection successful with uri - ", connectionString)

		go func() {
			select {
			case sig := <-c:
				log.Debugf("Got %s signal. Disconnecting mongodb...\n", sig)
				if err = client.Disconnect(context.TODO()); err != nil {
					log.Errorj(log.JSON{"err": err, "message": "Could not disconnect mongo"})
				}
				os.Exit(1)
			}
		}()
		// Ping the primary
		if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
			panic(err)
		}
		//log.Print("mongo ping successful with uri - ", connectionString)

		log.Debugf("Successfully connected and pinged mongodb")
	})
	return client
}
