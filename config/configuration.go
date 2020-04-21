package config

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
	"time"
)

var (
	Conf = Configuration{}
	Db   = DbConnection{}
)

type DbConnection struct {
	ConnSQLx    sqlx.DB
	ConnGORM    gorm.DB
	Transaction *gorm.DB
}

// Configuration struct of config
type Configuration struct {
	Port     string `json:"port"`
	DbString string `json:"dbstring"`
	Dblog    bool   `json:"dblog"`
	Salt     string `json:"salt"`
	Mailer   Mailer `json:"mailer"`

	//Sms      Sms    `json:"sms"`
	JwtRequestDeamon DeamonContext `json:"jwtRequestDeamon"`
	JwtParseDeamon   DeamonContext `json:"jwtParseDeamon"`
	Crypta           Crypta        `json:"crypta"`
	RootDir          string        `json:"rootdir"`
	Redis            Redis         `json:"redis"`
	Jwtgost          jwtgost       `json:"jwtgost"`
	Log              LogConf
}

// DeamonContext Configuration struct of deamon
type DeamonContext struct {
	PidFileName  string `json:"pidfilename"`
	PidFilePerm  string `json:"pidfileperm"`
	PidFilePermI uint32 `json:"-"`
	LogFileName  string `json:"logfilename"`
	LogFilePerm  string `json:"logfileperm"`
	LogFilePermI uint32 `json:"-"`
	WorkDir      string `json:"workdir"`
	Umask        string `json:"umask"`
	UmaskI       uint32 `json:"-"`
	Args         string `json:"args"`
	Workers      int    `json:"workers"`
}

// Configuration struct of crypta
type Crypta struct {
	Tcp        string `json:"tcp"`
	Decrypt    string `json:"decrypt"`
	Aeneas     string `json:"aeneas"`
	OpenCert64 string `json:"openCert64"`
}

// Configuration struct of redis
type Redis struct {
	Ip          string        `json:"ip"`
	Port        string        `json:"port"`
	MaxIdle     int           `json:"max_idle"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	Type        string        `json:"type"`
	Pool        *redis.Pool   `json:"pool"`
}

type jwtgost struct {
	JwtHeaderValSchemaPath string `json:"headervalschemapath"`
}

type LogConf struct {
	LogLevel string `json:"loglevel"`
}

// Configuration struct of mailer
type Mailer struct {
	IsSMTP   bool   `json:"smtp"`
	Hostname string `json:"hostname"`
	Host     string `json:"host"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (C *Configuration) GetJwtHeaderValSchemaPath() string {
	return C.Jwtgost.JwtHeaderValSchemaPath

}
func (C *Configuration) GetRootDir() string {
	return C.RootDir
}

func (C *Configuration) GetAeneasAddres() string {
	return C.Crypta.Aeneas
}

func GetAeneasAddres() string {
	return Conf.Crypta.Aeneas
}

func GetConfiguration(url string) *Configuration {
	Conf.Redis.Pool = &redis.Pool{
		MaxIdle:     Conf.Redis.MaxIdle,
		IdleTimeout: Conf.Redis.IdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(Conf.Redis.Type, Conf.Redis.Ip+`:`+Conf.Redis.Port)
		},
	}
	file, err := os.Open(url)
	if err != nil {
		fmt.Println("При открытии файла конфигурации что-то не так:", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Conf)
	if err != nil {
		fmt.Println("Ошибка декодирования конфига error: ", err)
	}
	deamonFilesPemToUint()
	//fmt.Println(Conf.Crypta.OpenCert64)
	return &Conf
}

func InitLoger() {
	//hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	//if err != nil {
	//	log.Error("Unable to connect to local syslog daemon")
	//} else {
	//	log.AddHook(hook)
	//}
	//switch Conf.Log.LogLevel {
	//case "Fatal":
	//	log.SetReportCaller(true)
	//	log.SetLevel(log.FatalLevel)
	//case "Error":
	//	log.SetReportCaller(true)
	//	log.SetLevel(log.ErrorLevel)
	//case "Warn":
	//	log.SetLevel(log.WarnLevel)
	//case "Info":
	//	log.SetLevel(log.InfoLevel)
	//case "Debug":
	//	//log.SetReportCaller(true)
	//	log.SetLevel(log.DebugLevel)
	//case "Trace":
	//	//log.SetReportCaller(true)
	//	log.SetLevel(log.TraceLevel)
	//default:
	//	log.SetReportCaller(true)
	//	log.WithFields(log.Fields{
	//		"package": "Config",
	//	}).Fatal("Unknown LOG Level, check conf.json ")
	//
	//}
}

func deamonFilesPemToUint() {
	if Conf.JwtRequestDeamon.LogFilePerm != "" {
		Conf.JwtRequestDeamon.LogFilePermI = StrFilePem(Conf.JwtRequestDeamon.LogFilePerm)

	} else {
		Conf.JwtRequestDeamon.LogFilePermI = StrFilePem("0644")
	}

	if Conf.JwtRequestDeamon.PidFilePerm != "" {
		Conf.JwtRequestDeamon.PidFilePermI = StrFilePem(Conf.JwtRequestDeamon.PidFilePerm)

	} else {
		Conf.JwtRequestDeamon.PidFilePermI = StrFilePem("0640")
	}

	if Conf.JwtRequestDeamon.Umask != "" {
		Conf.JwtRequestDeamon.UmaskI = StrFilePem(Conf.JwtRequestDeamon.Umask)

	} else {
		Conf.JwtRequestDeamon.UmaskI = StrFilePem("027")
	}

	if Conf.JwtParseDeamon.LogFilePerm != "" {
		Conf.JwtParseDeamon.LogFilePermI = StrFilePem(Conf.JwtParseDeamon.LogFilePerm)

	} else {
		Conf.JwtParseDeamon.LogFilePermI = StrFilePem("0644")
	}

	if Conf.JwtParseDeamon.PidFilePerm != "" {
		Conf.JwtParseDeamon.PidFilePermI = StrFilePem(Conf.JwtParseDeamon.PidFilePerm)

	} else {
		Conf.JwtParseDeamon.PidFilePermI = StrFilePem("0640")
	}

	if Conf.JwtParseDeamon.Umask != "" {
		Conf.JwtParseDeamon.UmaskI = StrFilePem(Conf.JwtParseDeamon.Umask)

	} else {
		Conf.JwtParseDeamon.UmaskI = StrFilePem("027")
	}

}

func StrFilePem(parseString string) uint32 {
	//fmt.Println(parseString)

	ui32, err := strconv.ParseUint(parseString, 8, 32)
	if err != nil {
		fmt.Println("error:", err)
	}
	return uint32(ui32)
}

func GetRedisConn() redis.Conn {
	conn := Conf.Redis.Pool.Get()
	return conn
}

func GetDbConnection() {
	db, err := sqlx.Connect("postgres", Conf.DbString)

	if err == nil {
		Db.ConnSQLx = *db
	} else {
		fmt.Println("Ошибка подключения sqlx")
		fmt.Println(err)
	}
	gdb, err := gorm.Open("postgres", Conf.DbString)
	gdb.DB().SetMaxIdleConns(0)
	if err == nil {
		Db.ConnGORM = *gdb
	} else {
		fmt.Println("Ошибка подключения gorm")
		fmt.Println(err)
	}
}
