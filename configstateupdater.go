package mongodb

import (
	"context"
	"github.com/dev4fun007/autobot-common"
	"go.mongodb.org/mongo-driver/bson"
)

type ConfigStateUpdater struct {
	repository common.Repository
}

func NewConfigStateUpdater(repository common.Repository) ConfigStateUpdater {
	return ConfigStateUpdater{
		repository: repository,
	}
}

func (u ConfigStateUpdater) UpdateConfig(ctx context.Context, name string, strategyType common.StrategyType, value interface{}) error {
	filter := bson.M{"base_config.name": name, "base_config.strategy_type": strategyType}
	return u.repository.Update(ctx, filter, value)
}
