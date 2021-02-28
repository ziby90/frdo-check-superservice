package main

import (
	"fmt"
	"frdo-check-superservice/pkg/entities"
	"frdo-check-superservice/pkg/handler"
	"frdo-check-superservice/pkg/repository"
	"frdo-check-superservice/pkg/repository/rabbit"
	"frdo-check-superservice/pkg/repository/sql"
	"frdo-check-superservice/pkg/service"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "15:04:05", // the "time" field configuration
		FullTimestamp:          true,
		DisableLevelTruncation: true, // log level field configuration
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// this function is required when you want to introduce your custom format.
			// In my case I wanted file and line to look like this `file="engine.go:141`
			// but f.File provides a full path along with the file name.
			// So in `formatFilePath()` function I just trimmed everything before the file name
			// and added a line number in the end
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
	logrus.SetFormatter(formatter)

	//logrus.SetFormatter(new(logrus.JSONFormatter))
	t := time.Now()
	pathLogName := fmt.Sprintf(`./logs/%d-%d-%d_log.txt`, t.Day(), t.Month(), t.Year())
	f, errLog := os.OpenFile(pathLogName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if errLog != nil {
		msg := fmt.Sprintf("failed to initialization log file : %s", errLog.Error())
		fmt.Println(msg)
		logrus.Error(msg)
		os.Exit(1)
	}
	logrus.SetOutput(f)
	DBConnects, err := sql.InitDB()
	if err != nil {
		msg := fmt.Sprintf("failed to initialization db : %s", err.Error())
		fmt.Println(msg)
		logrus.Error(msg)
		os.Exit(2)
	}

	rabbitCh, errInit := rabbit.InitRabbit()
	if errInit != nil {
		msg := fmt.Sprintf("failed to initialization rabbit : %s", errInit.Error())
		fmt.Println(msg)
		logrus.Error(msg)
		os.Exit(3)
	}

	repos := repository.NewRepository(DBConnects, rabbitCh)
	services := service.NewService(repos)
	ents := entities.NewEntities(services)
	hand := handler.NewHandler(services, ents)

	hand.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	//
	//logrus.Print("Todo Shutting Down")
	//if err := srv.Shutdown(context.Background()); err != nil {
	//	logrus.Printf("error occured on server shutting down : %s", err.Error())
	//}
	//if err := db.Close(); err != nil {
	//	logrus.Printf("error occured on db connection close : %s", err.Error())
	//}
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
