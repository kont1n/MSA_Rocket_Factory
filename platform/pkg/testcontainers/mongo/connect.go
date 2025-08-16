package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connectMongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.Errorf("failed to connect to mongo: %v", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Errorf("failed to ping mongo: %v", err)
	}

	return client, nil
}
