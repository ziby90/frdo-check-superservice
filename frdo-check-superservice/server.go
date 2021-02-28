package frdo_check_superservice

import (
	"context"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func InitConfig(pathDir string) error {
	viper.AddConfigPath(pathDir)
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func RootDir() string {
	var d string
	_, exist := os.LookupEnv("GO_DEVELOPER")
	if exist {
		_, b, _, _ := runtime.Caller(0)
		d = path.Join(path.Dir(b))
	} else {
		e, err := os.Executable()
		if err != nil {
			panic(err)
		}
		d = path.Dir(e)
	}
	return d
}
