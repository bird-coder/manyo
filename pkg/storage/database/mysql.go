/*
 * @Author: yujiajie
 * @Date: 2025-05-30 11:44:18
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-06-04 17:33:20
 * @FilePath: /manyo/pkg/storage/database/mysql.go
 * @Description:
 */
package database

import (
	"fmt"

	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func NewMysql(host string, cfg *config.MysqlConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(cfg.Dsn, cfg.Password),
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt: true,
	})
	if err != nil {
		logger.Error("mysql init failed, host: %s, err: %v", host, err)
		return nil, err
	}
	if cfg.Cluster {
		initCluster(db, cfg)
	}
	logger.Info("mysql init success, host: %s", host)
	return db, nil
}

func initCluster(db *gorm.DB, cfg *config.MysqlConfig) {
	if len(cfg.Sources) == 0 && len(cfg.Replicas) == 0 {
		return
	}

	config := dbresolver.Config{}
	config.Sources = []gorm.Dialector{}
	config.Replicas = []gorm.Dialector{}

	for _, sourceConfig := range cfg.Sources {
		config.Sources = append(config.Sources, mysql.Open(fmt.Sprintf(sourceConfig.Dsn, sourceConfig.Password)))
	}
	for _, replicaConfig := range cfg.Replicas {
		config.Replicas = append(config.Replicas, mysql.Open(fmt.Sprintf(replicaConfig.Dsn, replicaConfig.Password)))
	}
	db.Use(dbresolver.Register(config))
}
