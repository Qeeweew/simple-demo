package dbinit

import (
	"fmt"
	"simple-demo/common/config"
	"simple-demo/common/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func MySQLInit() *gorm.DB {
	var err error
	MySQL, err := createDB(struct {
		Addr string
		User string
		Pass string
		DB   string
	}{
		Addr: config.MysqlCfg.Address,
		User: config.MysqlCfg.User,
		Pass: config.MysqlCfg.Password,
		DB:   config.MysqlCfg.DBName,
	})

	if err != nil {
		logrus.Panic("connect mysql error: ", err.Error())
	}

	logrus.Infof("Connected mysql success")

	db, _ := MySQL.DB()
	db.SetMaxIdleConns(config.MysqlCfg.MaxIdle)        // 设置最大空闲连接数
	db.SetMaxOpenConns(config.MysqlCfg.MaxOpen)        // 设置最大连接数
	db.SetConnMaxLifetime(config.MysqlCfg.MaxLifetime) // 设置连接最大存活时间

	// TODO: model层设计好之后放开注释

	// 自动建表
	AutoCreateTable(MySQL)
	return MySQL
}

// AutoCreateTable TODO: model层设计好之后修改自动建表的逻辑
func AutoCreateTable(db *gorm.DB) {
	_ = db.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.Favorite{},
		// &model.Friend{},
		&model.Message{},
	)

}

func createDSN(dbInfo struct {
	Addr string
	User string
	Pass string
	DB   string
}) string {
	//user:password@/dbname?charset=utf8&parseTime=True&loc=Local
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbInfo.User, dbInfo.Pass, dbInfo.Addr, dbInfo.DB)
}

func createDB(dbInfo struct {
	Addr string
	User string
	Pass string
	DB   string
}) (*gorm.DB, error) {
	cfg := struct {
		Addr string
		User string
		Pass string
		DB   string
	}{
		Addr: dbInfo.Addr,
		User: dbInfo.User,
		Pass: dbInfo.Pass,
		DB:   dbInfo.DB,
	}
	DB, err := gorm.Open(mysql.Open(createDSN(cfg)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		PrepareStmt: true,           // 预处理语句
		Logger:      logger.Default, // 日志级别
	})
	return DB, err
}
