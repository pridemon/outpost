package internal

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlModule = di.Options(
	di.Provide(MysqlConfigProvider),
	di.Provide(MysqlDatabaseProvider),
	di.Invoke(GormMigrate),
)

type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func MysqlConfigProvider(v *viper.Viper) (*MysqlConfig, error) {
	var conf MysqlConfig
	err := v.UnmarshalKey("mysql", &conf)
	return &conf, err
}

func MysqlDatabaseProvider(conf *MysqlConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GormMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.RefreshInfo{},
	)
}
