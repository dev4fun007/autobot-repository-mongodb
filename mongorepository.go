package mongodb

import (
	"context"
	"github.com/dev4fun007/autobot-common"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	MongoRepositoryTag = "MongoRepository"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(client *ClientMongoDb, dbName, repoCollection string) MongoRepository {
	collection := client.Database(dbName).Collection(repoCollection)
	return MongoRepository{
		collection: collection,
	}
}

func (repo MongoRepository) Save(ctx context.Context, value interface{}) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, DbRequestTimeout*time.Second)
	defer cancel()
	_, err := repo.collection.InsertOne(ctxTimeout, value)
	if err != nil {
		log.Error().Str(common.LogComponent, MongoRepositoryTag).Err(err).Msg("error saving object")
		return err
	}
	return nil
}

func (repo MongoRepository) SaveAll(ctx context.Context, value []interface{}) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := repo.collection.InsertMany(ctxTimeout, value)
	if err != nil {
		log.Error().Str(common.LogComponent, MongoRepositoryTag).Err(err).Msg("error saving batch")
		return err
	}
	log.Debug().Str(common.LogComponent, MongoRepositoryTag).Int("batch-size", len(value)).Msg("batch saved")
	return nil
}

//{"base_config.name": name, "base_config.strategy_type": strategyType}
func (repo MongoRepository) Update(ctx context.Context, filter interface{}, value interface{}) error {
	filter = filter.(bson.M)
	ctxTimeout, cancel := context.WithTimeout(ctx, DbRequestTimeout*time.Second)
	defer cancel()
	res := repo.collection.FindOneAndReplace(ctxTimeout, filter, value)
	if res.Err() != nil {
		log.Error().Str(common.LogComponent, MongoRepositoryTag).Err(res.Err()).Msg("error updating record")
		return res.Err()
	}
	return nil
}

//{"base_config.name": name, "base_config.strategy_type": strategyType}
func (repo MongoRepository) Delete(ctx context.Context, filter interface{}) error {
	filter = filter.(bson.M)
	ctxTimeout, cancel := context.WithTimeout(ctx, DbRequestTimeout*time.Second)
	defer cancel()
	_, err := repo.collection.DeleteOne(ctxTimeout, filter)
	if err != nil {
		return err
	}
	return nil
}

//{"base_config.name": name, "base_config.strategy_type": strategyType}
func (repo MongoRepository) Get(ctx context.Context, filter interface{}) (interface{}, error) {
	filter = filter.(bson.M)
	ctxTimeout, cancel := context.WithTimeout(ctx, DbRequestTimeout*time.Second)
	defer cancel()
	response := repo.collection.FindOne(ctxTimeout, filter)
	if response.Err() != nil {
		return nil, response.Err()
	}
	// bson.Raw is sent
	return response.DecodeBytes()
}

//{"base_config.strategy_type": strategyType}
func (repo MongoRepository) GetAllByFilter(ctx context.Context, filter interface{}) []interface{} {
	filter = filter.(bson.M)
	ctxTimeout, cancel := context.WithTimeout(ctx, DbRequestTimeout*time.Second)
	defer cancel()
	cur, err := repo.collection.Find(ctxTimeout, filter)
	if err != nil {
		return nil
	}
	list := make([]interface{}, 0)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		list = append(list, cur.Current)
	}
	if err = cur.Err(); err != nil {
		return nil
	}
	return list
}
