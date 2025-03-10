package dbrepo

import (
	"github.com/DanielChungYi/puna/internal/config"
	"github.com/DanielChungYi/puna/internal/repository"
	"gorm.io/gorm"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *gorm.DB
}

func NewPostgresRepo(a *config.AppConfig, conn *gorm.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
