package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)


// Conf 全局变量 用来保存程序的所有储存信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64 `mapstructure:"machine_id"`
	Port int `mapstructure:"port"`

	*LogConfig `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	Filename string `mapstructure:"filename"`
	MaxSize int `mapstructure:"max_size"`
	MaxAge int `mapstructure:"max_age"`
	MaxBackups int `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host string `mapstructure:"host"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName string `mapstructure:"dbname"`
	Port int `mapstructure:"port"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	MaxIdleConns int `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port int `mapstructure:"port"`
	DB int `mapstructure:"db"`
	PoolSize int `mapstructure:"pool_size"`
}

func Init() (err error) {

	//法一：直接指定配置文件路径（相对/绝对）
	//相对路径：相对于可执行文件的路径
	//绝对路径：系统中实际的文件路径

	//法二：指定配置文件名和位置  viper自行查找可用的配置文件
	//文件名不需要带后缀
	//位置可配置多个
	viper.SetConfigFile("config.yaml")
	//viper.SetConfigName("config") // 指定配置文件路径（无后缀）
	//viper.SetConfigType("yaml")   // 指定配置文件类型（专用于从远程获取配置信息时，告诉程序 指定配置文件的类型[不会与上面的SetConfigName对应结合]）
	viper.AddConfigPath(".")      // 指定为当前路径
	//viper.AddConfigPath("./conf")  //配置多个查找路径


	err = viper.ReadInConfig()    // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return
	}

	//把读取到的配置信息反序列化到Conf变量中
	if err:=viper.Unmarshal(Conf);err!=nil {
		fmt.Printf("viper.Unmarshal failed,err:%v\n",err)
	}

	// 监控配置文件变化，实时更新配置信息
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		//把最新的配置信息反序列化到Conf中
		if err:=viper.Unmarshal(Conf);err!=nil {
			fmt.Printf("viper.Unmarshal failed,err:%v\n",err)
		}
	})
	return
}

//func Init()(err error){
//	viper.SetConfigFile("./conf/config.yaml") // 指定配置文件路径
//	err = viper.ReadInConfig()        // 读取配置信息
//	if err != nil {                    // 读取配置信息失败
//		panic(fmt.Errorf("Fatal error config file: %s \n", err))
//	}
//
//	// 监控配置文件变化
//	viper.WatchConfig()
//
//	r := gin.Default()
//	// 访问/version的返回值会随配置文件的变化而变化
//	r.GET("/version", func(c *gin.Context) {
//		c.String(http.StatusOK, viper.GetString("version"))
//	})
//
//	if err := r.Run(
//		fmt.Sprintf(":%d", viper.GetInt("port"))); err != nil {
//		panic(err)
//	}
//}
