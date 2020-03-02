package api

import (
	infinote "boilerplate"
	"boilerplate/canlog"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	reqsName    = "chi_requests_total"
	latencyName = "chi_request_duration_milliseconds"
)

// PromMiddleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type PromMiddleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

var (
	dflBuckets = []float64{300, 1200, 5000}
)

// NewPromMiddleware returns a new prometheus Middleware handler.
func NewPromMiddleware() func(next http.Handler) http.Handler {
	var m PromMiddleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "boilerplate",
			Subsystem: "api",
			Name:      reqsName,
			Help:      "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
		},
		[]string{"code", "method", "path", "ip"},
	)
	prometheus.MustRegister(m.reqs)

	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "boilerplate",
		Subsystem: "api",
		Name:      latencyName,
		Help:      "How long it took to process the request, partitioned by status code, method and HTTP path.",
	},
		[]string{"code", "method", "path", "ip"},
	)
	prometheus.MustRegister(m.latency)
	return m.handler
}
func (m *PromMiddleware) handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		ip := r.RemoteAddr
		if strings.Contains(ip, "127.0.0.1") || strings.Contains(ip, "[::1]") {
			ip = "127.0.0.1"
		}

		m.reqs.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, r.URL.Path, ip).Inc()
		m.latency.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, r.URL.Path, ip).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}

// NewGQLMetrics registers metrics and returns a middleware struct for GraphQL purposes
func NewGQLMetrics() *GraphQLMetrics {
	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "boilerplate",
			Subsystem: "api",
			Name:      "graphql_error_total",
			Help:      "Total number of errors returned from the graphql server.",
		},
		[]string{"object", "field"},
	)

	requestStartedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "boilerplate",
			Subsystem: "api",
			Name:      "graphql_request_started_total",
			Help:      "Total number of requests started on the graphql server.",
		},
		[]string{},
	)

	requestCompletedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "boilerplate",
			Subsystem: "api",
			Name:      "graphql_request_completed_total",
			Help:      "Total number of requests completed on the graphql server.",
		},
		[]string{},
	)

	resolverStartedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "boilerplate",
			Subsystem: "api",
			Name:      "graphql_resolver_started_total",
			Help:      "Total number of resolver started on the graphql server.",
		},
		[]string{"object", "field"},
	)

	resolverCompletedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "boilerplate",
			Subsystem: "api",
			Name:      "graphql_resolver_completed_total",
			Help:      "Total number of resolver completed on the graphql server.",
		},
		[]string{"object", "field"},
	)

	timeToResolveField := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "boilerplate",
		Subsystem: "api",
		Name:      "graphql_resolver_duration_millseconds",
		Help:      "The time taken to resolve a field by graphql server.",
	}, []string{"object", "field"})

	timeToHandleRequest := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "boilerplate",
		Subsystem: "api",
		Name:      "graphql_request_duration_millseconds",
		Help:      "The time taken to handle a request by graphql server.",
	}, []string{})

	prometheus.MustRegister(
		requestStartedCounter,
		requestCompletedCounter,
		resolverStartedCounter,
		resolverCompletedCounter,
		timeToResolveField,
		timeToHandleRequest,
	)
	return &GraphQLMetrics{
		ErrorCounter:             errorCounter,
		RequestStartedCounter:    requestStartedCounter,
		RequestCompletedCounter:  requestCompletedCounter,
		ResolverStartedCounter:   resolverStartedCounter,
		ResolverCompletedCounter: resolverCompletedCounter,
		TimeToResolveField:       timeToResolveField,
		TimeToHandleRequest:      timeToHandleRequest,
	}
}

// GraphQLMetrics contains the graphql middleware
type GraphQLMetrics struct {
	ErrorCounter             *prometheus.CounterVec
	RequestStartedCounter    *prometheus.CounterVec
	RequestCompletedCounter  *prometheus.CounterVec
	ResolverStartedCounter   *prometheus.CounterVec
	ResolverCompletedCounter *prometheus.CounterVec
	TimeToResolveField       *prometheus.HistogramVec
	TimeToHandleRequest      *prometheus.HistogramVec
}

// ResolverMiddleware runs before each resolver
func (m *GraphQLMetrics) ResolverMiddleware() graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		rctx := graphql.GetResolverContext(ctx)
		if strings.Contains(rctx.Object, "__") {
			return next(ctx)
		}
		m.ResolverStartedCounter.WithLabelValues(rctx.Object, rctx.Field.Name).Inc()

		observerStart := time.Now()
		res, err := next(ctx)

		if err != nil {
			m.ErrorCounter.WithLabelValues(rctx.Object, rctx.Field.Name).Inc()
		}

		m.TimeToResolveField.WithLabelValues(rctx.Object, rctx.Field.Name).Observe(float64(time.Since(observerStart).Nanoseconds()) / 1000000)

		m.ResolverCompletedCounter.WithLabelValues(rctx.Object, rctx.Field.Name).Inc()

		return res, err
	}
}

// RequestMiddleware runs before each request
func (m *GraphQLMetrics) RequestMiddleware() graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		m.RequestStartedCounter.WithLabelValues().Inc()
		observerStart := time.Now()

		res := next(ctx)

		rctx := graphql.GetResolverContext(ctx)
		reqCtx := graphql.GetRequestContext(ctx)
		errList := reqCtx.GetErrors(rctx)
		if len(errList) > 0 {
			m.ErrorCounter.WithLabelValues(rctx.Object, rctx.Field.Name).Inc()
		}

		m.TimeToHandleRequest.WithLabelValues().Observe(float64(time.Since(observerStart).Nanoseconds()) / 1000000)

		m.RequestCompletedCounter.WithLabelValues().Inc()

		return res
	}
}

func canonicalLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			ctx := infinote.WithCanonicalLogger(r.Context(), logger)
			canlog.Set(ctx, "start", now.Unix())
			canlog.Set(ctx, "ip", r.RemoteAddr)
			canlog.Set(ctx, "reqId", middleware.GetReqID(ctx))
			next.ServeHTTP(w, r.WithContext(ctx))
			canlog.Set(ctx, "duration_ms", fmt.Sprintf("%.02f", float64(time.Since(now).Microseconds())/1000))

			if infinote.ClaimExistsInContext(ctx, infinote.ClaimUserName) {
				username, _ := infinote.ClaimValueFromContext(ctx, infinote.ClaimUserName)
				canlog.Set(ctx, "username", username)
			}
			if infinote.ClaimExistsInContext(ctx, infinote.ClaimRoles) {
				roles, _ := infinote.ClaimValueFromContext(ctx, infinote.ClaimRoles)
				canlog.Set(ctx, "roles", roles)
			}
			if infinote.ClaimExistsInContext(ctx, infinote.ClaimUserID) {
				userID, _ := infinote.ClaimValueFromContext(ctx, infinote.ClaimUserID)
				canlog.Set(ctx, "user_id", userID)
			}
			canlog.Log(ctx, "request")
		}
		return http.HandlerFunc(fn)
	}

}
