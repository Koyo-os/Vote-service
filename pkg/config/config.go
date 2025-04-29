package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	KAFKA_URL         string `env:"KAFKA_URL" envDefault:"localhost:9092"`
	DSN               string `env:"DSN" envDefault:"data/main.db"`
	CreateVoteReqType string `env:"CREATE_REQ_TYPE" envDefault:"request.vote.create"`
	CreatedVoteType   string `env:"CREATED_VOTE_EVENT_TYPE" envDefault:"vote.created"`
	DeleteVoteReqType string `env:"DELETE_REQ_TYPE" envDefault:"request.vote.delete"`
	DeletedVoteType   string `env:"DELETED_VOTE_EVENT_TYPE" envDefault:"vote.deleted"`
	RequestTopicName  string `env:"REQUEST_TOPIC_NAME" envDefault:"requests"`
	TopicName         string `env:"TOPIC_NAME" envDefault:"events"`
	GroupID           string `env:"GROUP_ID" envDefault:"vote"`
	RedisSentinelUrl  string `env:"REDIS_SENTINEL_HOST" envDefault:"redis-sentinel"`
	RedisPassword     string `env:"REDIS_PASSWORD" envDefault:"redis"`
}

func Load() (*Config, error) {
	cfg := Config{}

	_ = godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	cfg.RedisSentinelUrl = cfg.RedisSentinelUrl + ":26379"

	return &cfg, nil
}
