package main

import (
	"context"

	"github.com/dirty-bro-tech/peers-touch-go/core/config"
	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/source/file"
)

/*
*
config:

	postgresSQLDsn: host=localhost user=user_hello password=passport_hello dbname=db_hello port=5432 sslmode=disable TimeZone=Asia/Shanghai
*/
type genConfig struct {
	PostgresSQL string `pconf:"postgresSQLDsn"`
}

type Value struct {
	Config genConfig `pconf:"config"`
}

var (
	value Value
)

func LoadConfig() {
	config.RegisterOptions(&value)
	err := config.NewConfig(
		config.WithSources(
			file.NewSource(
				option.WithRootCtx(context.Background()),
				file.WithPath("./config.yml"),
			),
		),
	).Init()
	if err != nil {
		panic(err)
	}
}
