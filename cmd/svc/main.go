package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/piusalfred/kitsvc"
	"github.com/piusalfred/kitsvc/svc"
	"github.com/piusalfred/kitsvc/svc/api"
	svchttp "github.com/piusalfred/kitsvc/svc/api/svc/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defLogLevel   = "error"
	defHTTPPort   = "9021"
	defJaegerURL  = ""
	defServerCert = ""
	defServerKey  = ""
	defSecret     = "secret"

	envLogLevel   = "MF_MFXKIT_LOG_LEVEL"
	envHTTPPort   = "MF_MFXKIT_HTTP_PORT"
	envJaegerURL  = "MF_JAEGER_URL"
	envServerCert = "MF_MFXKIT_SERVER_CERT"
	envServerKey  = "MF_MFXKIT_SERVER_KEY"
	envSecret     = "MF_MFXKIT_SECRET"
)

type config struct {
	logLevel     string
	httpPort     string
	authHTTPPort string
	authGRPCPort string
	jaegerURL    string
	serverCert   string
	serverKey    string
	secret       string
}


func main() {

	cfg := loadConfig()

	logger:= log.NewLogfmtLogger(os.Stderr)


	service := newService("piusalfred",logger)

	errs := make(chan error, 2)

	go startHTTPServer(svchttp.MakeHandler(service), cfg.httpPort, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Log(fmt.Sprintf("svc service terminated: %s", err))
}

func newService(username string, logger log.Logger) svc.Service {
	service := svc.New(username)

	service = api.LoggingMiddleware(service, logger)
	return service
}

func loadConfig() config {
	return config{
		logLevel:   kitsvc.Env(envLogLevel, defLogLevel),
		httpPort:   kitsvc.Env(envHTTPPort, defHTTPPort),
		serverCert: kitsvc.Env(envServerCert, defServerCert),
		serverKey:  kitsvc.Env(envServerKey, defServerKey),
		jaegerURL:  kitsvc.Env(envJaegerURL, defJaegerURL),
		secret:     kitsvc.Env(envSecret, defSecret),
	}
}

func startHTTPServer(handler http.Handler, port string, cfg config, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	/*if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info(fmt.Sprintf("Mfxkit service started using https on port %s with cert %s key %s",
			port, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, handler)
		return
	}*/
	logger.Log(fmt.Sprintf("svc service started using http on port %s", cfg.httpPort))
	errs <- http.ListenAndServe(p, handler)
}
