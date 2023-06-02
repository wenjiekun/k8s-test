package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"k8s-test/gin-mysql/initialize"
	"k8s-test/gin-mysql/model"
)

var (
	DB *gorm.DB
)

func InitDB() {
	//10.104.171.81:3306
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("user"),
		viper.GetString("password"),
		viper.GetString("MysqlHost"),
		viper.GetInt64("MysqlPort"),
		viper.GetString("MysqlName"))

	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
	//	logger.Config{
	//		SlowThreshold: time.Second, //慢sql阀值
	//		LogLevel:      logger.Info, //Log level
	//		Colorful:      true,
	//	},
	//)

	//全局模式
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	//设置全局logger，这个logger在我们执行每个sql语句的时候打印每天一行sql
	//sql才是最重要的，本着这个原则生成对应的表 -migrations
	//迁移 schema
	//_ = db.AutoMigrate(&model.User{})
}

// 初始化配置文件
func init() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	//viper.AddConfigPath("/etc/appname/") // 查找配置文件所在路径
	viper.AddConfigPath("/Users/wenjiekun/Documents/go/k8s-test/gin-mysql/conf/")
	viper.AddConfigPath(".")      // 还可以在工作目录中搜索配置文件
	viper.AddConfigPath("./conf") // 还可以在工作目录中搜索配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("read config failed: %v \r", err)
	}
	//if err := viper.Unmarshal(&conf); err != nil {
	//	fmt.Println(err)
	//}
}

func pong(c *gin.Context) {
	var demo1 model.Demo1
	DB.Where(&model.Demo1{Name: "wjk"}).First(&demo1)
	c.JSON(200, gin.H{
		"message": demo1.Name,
	})
}

func main() {
	//1.初始化logger
	initialize.InitLogger()
	InitDB()
	r := gin.Default()
	r.GET("/ping", pong)

	r.Run(":8080")

}
