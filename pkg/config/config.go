package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQUrl       string `env:"RABBITMQ_URL"            envDefault:"localhost:9092"`
	DSN               string `env:"DSN"                     envDefault:"data/main.db"`
	CreateVoteReqType string `env:"CREATE_REQ_TYPE"         envDefault:"request.vote.create"`
	CreatedVoteType   string `env:"CREATED_VOTE_EVENT_TYPE" envDefault:"vote.created"`
	DeleteVoteReqType string `env:"DELETE_REQ_TYPE"         envDefault:"request.vote.delete"`
	DeletedVoteType   string `env:"DELETED_VOTE_EVENT_TYPE" envDefault:"vote.deleted"`
	RequestTopicName  string `env:"REQUEST_TOPIC_NAME"      envDefault:"requests"`
	OutputTopicName   string `env:"OUTPUT_TOPIC_NAME"       envDefault:"events"`
	GroupID           string `env:"GROUP_ID"                envDefault:"vote"`
	RedisUrl          string `env:"REDIS_URL"               envDefault:"localhost:3056"`
}

func Load() (*Config, error) {
	cfg := Config{}

	_ = godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
