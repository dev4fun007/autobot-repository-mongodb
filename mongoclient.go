package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ClientMongoDb struct {
	*mongo.Client
}

func NewClientMongoDb() *ClientMongoDb {
	return &ClientMongoDb{}
}

func (receiver *ClientMongoDb) ConnectClient(ctx context.Context, dbUri string) error {
	clientOptions := options.Client().ApplyURI(dbUri)
	ctxTimeout, cancel := context.WithTimeout(ctx, DbConnectionTimeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctxTimeout, clientOptions)
	if err != nil {
		return err
	}
	receiver.Client = client
	return nil
}

func (receiver *ClientMongoDb) DisconnectClient(ctx context.Context) error {
	return receiver.Disconnect(ctx)
}
