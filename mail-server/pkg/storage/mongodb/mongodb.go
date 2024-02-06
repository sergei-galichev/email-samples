package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type StorageMongoDB struct {
	dbName string
	client *mongo.Client
}

func NewStorage() *StorageMongoDB {

	uri := os.Getenv("URI")
	dbName := os.Getenv("DB_NAME")
	count := 0
	var client *mongo.Client
	var err error

	log.Println("....... Setting up Connection to MongoDB .......")

	for {
		client, err = connect(uri)
		if err != nil {
			log.Println(err)
			log.Println("MongoDB is not connected")
			count++
		} else {
			log.Println("MongoDB is connected")
			break
		}

		if count >= 5 {
			log.Fatal("Cannot connect to MongoDB")
		}

		log.Println("Wait:.... Mail App Database Retrying to connect ....")
		time.Sleep(5 * time.Second)
		continue
	}

	return &StorageMongoDB{
		dbName: dbName,
		client: client,
	}
}

func connect(uri string) (*mongo.Client, error) {
	dbCtx, dbCtxCancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer dbCtxCancel()

	log.Println(uri)
	client, err := mongo.Connect(dbCtx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.New("Error to connect to MongoDB")
	}

	err = client.Ping(dbCtx, nil)
	if err != nil {
		return nil, errors.New("Error to ping MongoDB")
	}

	return client, nil
}

func (s *StorageMongoDB) Client() *mongo.Client {
	return s.client
}

func (s *StorageMongoDB) Collection(name string) *mongo.Collection {
	return s.client.Database(s.dbName).Collection(name)
}

func (s *StorageMongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.client.Disconnect(ctx)
	if err != nil {
		log.Fatal("Could not disconnect from MongoDB: ", err)
	}
}
