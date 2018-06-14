package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/tidwall/sjson"

	"github.com/Peripli/service-manager/pkg/filter"

	"os"
	"os/signal"

	cfenv "github.com/Peripli/service-manager/cf/env"
	"github.com/Peripli/service-manager/env"
	"github.com/Peripli/service-manager/rest"
	"github.com/Peripli/service-manager/server"
	"github.com/Peripli/service-manager/sm"
	"github.com/sirupsen/logrus"
)

func main() {
	flags := initFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handleInterrupts(ctx, cancel)

	api := rest.API{}
	api.RegisterFilters(filter.Filter{
		RequestMatcher: filter.RequestMatcher{
			Methods:     []string{"GET"},
			PathPattern: "/v1/service_brokers",
		},
		Middleware: func(request *filter.Request, next filter.Handler) (*filter.Response, error) {
			res, err := next(request)
			if err == nil {
				res.Body, err = sjson.SetBytes(res.Body, "extra", "value")
			}
			fmt.Println(">>>>>", string(res.Body))
			return res, err
		},
	})

	config := &sm.Parameters{
		Context:     ctx,
		Environment: getEnvironment(flags),
		API:         &api,
	}
	srv, err := sm.NewServer(config)
	if err != nil {
		logrus.Fatal("Error creating the server: ", err)
	}

	srv.Run(ctx)
}

func handleInterrupts(ctx context.Context, cancelFunc context.CancelFunc) {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt)
	go func() {
		select {
		case <-term:
			logrus.Error("Received OS interrupt, exiting gracefully...")
			cancelFunc()
		case <-ctx.Done():
			return
		}
	}()
}

func initFlags() map[string]interface{} {
	configFileLocation := flag.String("config_location", ".", "Location of the application.yaml file")
	flag.Parse()
	return map[string]interface{}{"config_location": *configFileLocation}
}

func getEnvironment(flags map[string]interface{}) server.Environment {
	configFileLocation := flags["config_location"].(string)
	logrus.Infof("config_location: %s", configFileLocation)

	runEnvironment := env.New(&env.ConfigFile{
		Path:   configFileLocation,
		Name:   "application",
		Format: "yaml",
	}, "SM")

	if _, exists := os.LookupEnv("VCAP_APPLICATION"); exists {
		return cfenv.New(runEnvironment)
	}
	return runEnvironment
}
