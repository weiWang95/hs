package dao

import (
	"hs/pkg/config"
	"hs/repository"
	"hs/repository/entity"

	"github.com/sirupsen/logrus"
)

var ConfigRepo = NewConfigDao()

type configDao struct {
	server entity.ServerConfig
	shop   entity.ShopConfig
}

func NewConfigDao() repository.ConfigRepo {
	d := configDao{}
	d.init()
	return &d
}

func (d *configDao) init() {
	if err := config.LoadJsonConfig("server.json", &d.server); err != nil {
		logrus.Errorf("load server config error: %v", err)
	}

	if err := config.LoadJsonConfig("shop.json", &d.shop); err != nil {
		logrus.Errorf("load shop config error: %v", err)
	}
}

func (d *configDao) GetServerConfig() entity.ServerConfig {
	return d.server
}

func (d *configDao) GetShopConfig() entity.ShopConfig {
	return d.shop
}
