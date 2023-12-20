package mysql

import (
	"fmt"
	"web_app/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	//parseTime=True:sqlx库可以做时间格式的转换（数据库里存的是时间戳）  时区是local本地时间
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		//viper.GetString("mysql.user"),
		//viper.GetString("mysql.password"),
		//viper.GetString("mysql.host"),
		//viper.GetInt("mysql.port"),
		//viper.GetString("mysql.dbname"),
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn) //open+ping
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err)) //日志 存一个message信息和error信息
		return
	}
	//db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	//db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

func Close() {
	_ = db.Close()
}
