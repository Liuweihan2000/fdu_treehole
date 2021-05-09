package dal

import (
	"GoProject/fudan_bbs/utils"
	"context"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func InitRedis() {
	redisClient = NewRedisClient()
	// ping 一下看能不能跑通，因为即使提供的地址和密码是非法的话也不会报错
	status := Redis().Ping(context.Background())
	if status.Err() != nil {
		utils.FatalErrorHandle(status.Err(), "error occurred while initializing redis")
	}
}

func InitMySQL() {
	logFile, err := os.OpenFile("sql.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	utils.FatalErrorHandle(err, "error occurred while creating sql log file")
	// 将 log 信息打印并写在文件中
	mw := io.MultiWriter(os.Stdout, logFile)
	// 生产环境下只打印 Warn 级别及以上的 log 信息
	logLevel := logger.Warn
	if viper.GetBool("debugging") {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(mw, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Millisecond * 500, // 慢日志阈值
			LogLevel:      logLevel,
			Colorful:      false,
		},
	)

	db, err = gorm.Open(mysql.Open(
		viper.GetString("sql_source")+"charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"), &gorm.Config{
		Logger: newLogger,
	})
	utils.FatalErrorHandle(err, "error occurred while opening mysql")
}

func InitDal() {
	InitRedis()
	InitMySQL()
}
