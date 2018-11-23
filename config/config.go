package config

import (
	"time"
	"os"
	"github.com/sirupsen/logrus"
	"path"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/thinmonkey/user-manager/utils/log"
	"github.com/thinmonkey/user-manager/utils"
	"github.com/spf13/viper"
	"github.com/thinmonkey/user-manager/models/db"
)

const (
	LOG_MAX_AGE    = 7 * 24 * time.Hour //日志最长保存一周
	LOG_ROTATETIME = 1 * 24 * time.Hour //一天切割一次文件

	LOGINFOFILE  = "user-manager.info.%Y%m%d%H%M.log"
	LOGERRORFILE = "user-manager.error.%Y%m%d%H%M.log"

	LOCAL_FILEPATH = "./locallog"
)

var logFilePath = LOCAL_FILEPATH

//配置初始化入口
func Init() {
	initViper()
	initLog()
	initDb()
}

func initViper() {
	viper.SetConfigFile("config.json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
}

//初始化数据库连接配置
func initDb() {
	mysqlURl := viper.GetString("mysql.url")
	mysqlPort := viper.GetString("mysql.port")
	mysqlUser := viper.GetString("mysql.user")
	mysqlPassword := viper.GetString("mysql.password")

	dbName := viper.GetString("mysql.databaseName")

	db.ConfigDbEngine(mysqlURl, mysqlPort, mysqlUser, mysqlPassword, dbName)

}

//初始化log配置
func initLog() {
	log.SetFormatter(&logrus.JSONFormatter{})
	//如果日志路径不存在则创建
	if ok, _ := utils.Exists(logFilePath); !ok {
		err := os.MkdirAll(logFilePath, 777)
		utils.CheckErr(err)
	}
	InfoPath := path.Join(logFilePath, LOGINFOFILE)
	ErrorPath := path.Join(logFilePath, LOGERRORFILE)
	InfoWriter, infoerr := rotatelogs.New(
		InfoPath,
		rotatelogs.WithLinkName(InfoPath),           // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(LOG_MAX_AGE),          // 文件最大保存时间
		rotatelogs.WithRotationTime(LOG_ROTATETIME), // 日志切割时间间隔
	)
	utils.CheckErr(infoerr)
	ErrorWriter, errorerr := rotatelogs.New(
		ErrorPath,
		rotatelogs.WithLinkName(ErrorPath),          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(LOG_MAX_AGE),          // 文件最大保存时间
		rotatelogs.WithRotationTime(LOG_ROTATETIME), // 日志切割时间间隔
	)
	utils.CheckErr(errorerr)
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: os.Stdout, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  InfoWriter,
		logrus.WarnLevel:  InfoWriter,
		logrus.ErrorLevel: ErrorWriter,
		logrus.FatalLevel: ErrorWriter,
		logrus.PanicLevel: ErrorWriter,
	}, &logrus.JSONFormatter{})
	log.AddHook(lfHook)
}
