package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/user/gin-microservice-boilerplate/pkg/db"
)

// Registration returns the ProviderRegistration for MongoDB.
func Registration() db.ProviderRegistration {
	return db.ProviderRegistration{
		Open:  openMongo,
		Close: closeMongo,
	}
}

func openMongo(ctx context.Context, cfg map[string]string) (interface{}, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg["uri"]))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client.Database(cfg["database"]), nil
}

func closeMongo(ctx context.Context, conn interface{}) error {
	database := conn.(*mongo.Database)
	return database.Client().Disconnect(ctx)
}
