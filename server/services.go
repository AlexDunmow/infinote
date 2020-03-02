package infinote

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/caddyserver/caddy"
	// http driver for caddy
	_ "github.com/caddyserver/caddy/caddyhttp"
	"go.uber.org/zap"
)

const caddyfileTemplate = `
{{ .caddyAddr}} {
	tls off
    proxy /api/ localhost{{ .apiAddr }} {
		transparent
		websocket
		timeout 10m
    }
    root {{ .rootPath }}
    rewrite { 
        if {path} not_match ^/api
        to {path} /
    }
}
`

// LoadbalancerService for caddy style load balancing
type LoadbalancerService struct {
	Addr string
	Log  *zap.SugaredLogger
}

// Run the API service
func (s *LoadbalancerService) Run(ctx context.Context, caddyAddr, apiAddr, rootPath string) error {
	s.Log.Infow("Starting caddy")
	caddy.AppName = "Boilerplate"
	caddy.AppVersion = "0.0.1"
	caddy.Quiet = true
	t := template.Must(template.New("CaddyFile").Parse(caddyfileTemplate))
	data := map[string]string{
		"caddyAddr": caddyAddr,
		"apiAddr":   apiAddr,
		"rootPath":  rootPath,
	}

	result := &bytes.Buffer{}
	err := t.Execute(result, data)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	caddyfile := &caddy.CaddyfileInput{
		Contents:       result.Bytes(),
		Filepath:       "Caddyfile",
		ServerTypeName: "http",
	}

	instance, err := caddy.Start(caddyfile)
	if err != nil {
		return fmt.Errorf("start caddy: %w", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			s.Log.Info("Stopping caddy")
			err := instance.Stop()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	instance.Wait()
	return nil
}

// APIService for long running
type APIService struct {
	Addr string
	Log  *zap.SugaredLogger
}

// Run the API service
func (s *APIService) Run(ctx context.Context, controller http.Handler) error {
	s.Log.Infow("Starting API")

	server := &http.Server{
		Addr:    s.Addr,
		Handler: controller,
	}

	go func() {
		select {
		case <-ctx.Done():
			s.Log.Info("Stopping API")
			err := server.Shutdown(ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	return server.ListenAndServe()
}
