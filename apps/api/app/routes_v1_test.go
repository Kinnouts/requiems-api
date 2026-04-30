package app

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/config"
)

func TestServiceEnabled(t *testing.T) {
	t.Run("returns true when enabled services is blank", func(t *testing.T) {
		if !serviceEnabled(config.Config{}, "text") {
			t.Fatal("expected service to be enabled when config is blank")
		}
	})

	t.Run("matches trimmed service names from config", func(t *testing.T) {
		cfg := config.Config{EnabledServices: "text, validation,networking"}

		if !serviceEnabled(cfg, "text") {
			t.Fatal("expected text service to be enabled")
		}
		if serviceEnabled(cfg, "finance") {
			t.Fatal("expected finance service to be disabled")
		}
	})
}

func TestRegisterV1Routes(t *testing.T) {
	t.Run("mounts all services when enabled services is blank", func(t *testing.T) {
		r := chi.NewRouter()

		registerV1Routes(context.Background(), r, nil, nil, config.Config{})

		routes := walkRoutes(t, r)
		for _, prefix := range []string{
			"/entertainment", "/finance", "/health", "/networking",
			"/places", "/technology", "/text", "/validation",
		} {
			if !hasRoutePrefix(routes, prefix) {
				t.Fatalf("expected mounted routes to include prefix %s; got %v", prefix, routes)
			}
		}
	})

	t.Run("mounts only explicitly enabled services", func(t *testing.T) {
		r := chi.NewRouter()

		registerV1Routes(context.Background(), r, nil, nil, config.Config{EnabledServices: "validation,text"})

		routes := walkRoutes(t, r)
		if !hasRoutePrefix(routes, "/validation") {
			t.Fatalf("expected validation routes to be mounted; got %v", routes)
		}
		if !hasRoutePrefix(routes, "/text") {
			t.Fatalf("expected text routes to be mounted; got %v", routes)
		}
		if hasRoutePrefix(routes, "/finance") {
			t.Fatalf("expected finance routes to be absent; got %v", routes)
		}
		if hasRoutePrefix(routes, "/technology") {
			t.Fatalf("expected technology routes to be absent; got %v", routes)
		}
	})
}

func walkRoutes(t *testing.T, r chi.Router) []string {
	t.Helper()

	var routes []string
	if err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routes = append(routes, route)
		return nil
	}); err != nil {
		t.Fatalf("walk routes: %v", err)
	}

	return routes
}

func hasRoutePrefix(routes []string, prefix string) bool {
	for _, route := range routes {
		if strings.HasPrefix(route, prefix) {
			return true
		}
	}
	return false
}
