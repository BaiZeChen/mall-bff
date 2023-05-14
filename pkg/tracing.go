package pkg

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
)

var TracingCloser io.Closer

func initJaeger() (opentracing.Tracer, io.Closer, error) {
	var cfg = jaegercfg.Configuration{
		ServiceName: "mall-bff",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://39.96.82.208:14268/api/traces",
		},
	}
	return cfg.NewTracer()
}

func InitTracing() {
	tracer, tracingCloser, err := initJaeger()
	if err != nil {
		panic(err)
	}
	TracingCloser = tracingCloser
	opentracing.SetGlobalTracer(tracer)
}
