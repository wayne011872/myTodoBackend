package model

import (
	"context"
	"errors"

	"github.com/wayne011872/log"
	"github.com/wayne011872/microservice/cfg"
	"github.com/wayne011872/microservice/di"
	"github.com/wayne011872/pgx/conn"
)

type ModelDI interface {
	conn.PgxDI
	log.LoggerDI
}

type Config struct {
	Log log.Logger

	di      ModelDI
	pgxConn conn.PgxConn
}

func (c *Config) errorHandler(err error) {
	if err != nil {
		c.Log.ErrorPkg(err)
	}
}

func (c *Config) Close() error {
	var err error
	if c.pgxConn != nil {
		c.pgxConn.Close()
		err = nil
		return err
	}
	return nil
}

func (c *Config) Copy() cfg.ModelCfg {
	cfg := *c
	return &cfg
}

func (c *Config) Init(uuid string, di di.DI) error {
	mdi, ok := di.(ModelDI)
	if !ok {
		return errors.New("no ModelDI")
	}
	var err error

	c.di = mdi
	c.Log, err = mdi.NewLogger(di.GetService(), uuid)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) NewPgxConn(ctx context.Context) (conn.PgxConn, error) {
	if c.pgxConn != nil {
		return c.pgxConn, nil
	}
	var err error
	c.pgxConn, err = c.di.NewPgxConn(ctx)
	if err != nil {
		return nil, err
	}
	return c.pgxConn, nil
}

func GetConfigFromEnv() (*Config, error) {
	var mycfg Config
	err := cfg.GetFromEnv(&mycfg)
	if err != nil {
		return nil, err
	}

	return &mycfg, nil
}
