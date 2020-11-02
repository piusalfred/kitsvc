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
	defHTTPPort = "9021"
	defUsername = "secret"

	envHTTPPort = "KITSVC_HTTP_PORT"
	envUsername = "KITSVC_SECRET"
)

type config struct {
	httpPort string
	username string
}

func main() {

	cfg := loadConfig()

	logger := log.NewLogfmtLogger(os.Stderr)

	service := newService(cfg.username, logger)

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
		httpPort: kitsvc.Env(envHTTPPort, defHTTPPort),
		username: kitsvc.Env(envUsername, defUsername),
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
