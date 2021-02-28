package sql

import (
	"errors"
	"fmt"
	frdo_check_superservice "frdo-check-superservice"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDbNamesList() []string {
	return []string{`CHECK_SMEV`, `DPO`}
}

func InitDB() (map[string]*sqlx.DB, error) {
	DBConnects := make(map[string]*sqlx.DB)
	logrus.Println(`rootdir`, frdo_check_superservice.RootDir())
	if err := frdo_check_superservice.InitConfig(frdo_check_superservice.RootDir() + `/configs`); err != nil {
		msg := fmt.Sprintf(`error for initialization cinfig: %s`, err.Error())
		logrus.Println(msg)
		return DBConnects, errors.New(msg)
	}
	if err := godotenv.Load(frdo_check_superservice.RootDir() + `/.env`); err != nil {
		msg := fmt.Sprintf("error loading variables: %s", err.Error())
		logrus.Println(msg)
		return DBConnects, errors.New(msg)
	}
	DBList := GetDbNamesList()
	for _, name := range DBList {
		db, err := NewPostgresDB(Config{
			Host:     viper.GetString(fmt.Sprintf("db.%s.host", name)),
			Port:     viper.GetString(fmt.Sprintf("db.%s.port", name)),
			Username: viper.GetString(fmt.Sprintf("db.%s.user", name)),
			Password: os.Getenv(fmt.Sprintf("DB_%s_PASSWORD", name)),
			DBName:   viper.GetString(fmt.Sprintf("db.%s.dbname", name)),
			SSLMode:  viper.GetString(fmt.Sprintf("db.%s.sslmode", name)),
		})
		if err != nil {
			logrus.Printf(`No connection to %s db, error: %s.`, name, err.Error())
			return DBConnects, err
		}
		logrus.Printf(`Success connection to %s db.`, name)
		DBConnects[name] = db
	}
	return DBConnects, nil
}
