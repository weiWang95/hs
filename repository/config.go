package repository

import "hs/repository/entity"

type ConfigRepo interface {
	GetServerConfig() entity.ServerConfig
	GetShopConfig() entity.ShopConfig
}
