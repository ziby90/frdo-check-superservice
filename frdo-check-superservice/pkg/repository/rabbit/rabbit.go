package rabbit

import (
	"10.10.11.220/ursgis/rabbitMQ.git"
	frdo_check_superservice "frdo-check-superservice"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func NewConnectRabbit(config *rabbitMQ.RabbitConfig) (*rabbitMQ.RabbitChan, error) {
	cfg := rabbitMQ.GetConfig()
	*cfg = *config
	ch, err := rabbitMQ.GetChannel(1)
	if err != nil {
		return nil, err
	}
	errConsume := ch.SetConsumerSettings(viper.GetString("rabbitMq.queuesIn.alias"))
	if errConsume != nil {
		return ch, err
	}

	logrus.Println(`connect to rabbit for reading`)
	return ch, nil
}

func InitRabbit() (*rabbitMQ.RabbitChan, error) {
	if err := frdo_check_superservice.InitConfig(frdo_check_superservice.RootDir() + `/configs`); err != nil {
		logrus.Fatalf("error initialization config : %s", err.Error())
		return nil, err
	}
	if err := godotenv.Load(frdo_check_superservice.RootDir() + `/.env`); err != nil {
		logrus.Fatalf("error loading variables: %s", err.Error())
		return nil, err
	}

	cfg := rabbitMQ.GetConfig()
	config := rabbitMQ.RabbitConfig{
		Login:       viper.GetString("rabbitMq.login"),
		Pass:        os.Getenv("RABBIT_PASSWORD"),
		Host:        viper.GetString("rabbitMq.host"),
		VirtualHost: viper.GetString("rabbitMq.vHost"),
		QueuesIn: []rabbitMQ.QueueIn{
			rabbitMQ.QueueIn{
				Alias:      viper.GetString("rabbitMq.queuesIn.alias"),
				RemoteKoef: viper.GetInt("rabbitMq.queuesIn.remoteKoef"),
				QueueName:  viper.GetString("rabbitMq.queuesIn.queue"),
			},
		},
		QueuesOut: []rabbitMQ.QueueOut{
			rabbitMQ.QueueOut{
				Alias:      viper.GetString("rabbitMq.queuesOut.alias"),
				Exchange:   viper.GetString("rabbitMq.queuesOut.exchange"),
				RoutingKey: viper.GetString("rabbitMq.queuesOut.routingKey"),
			},
		},
	}
	*cfg = config

	rabbitChannel, err := NewConnectRabbit(cfg)

	return rabbitChannel, err
}
