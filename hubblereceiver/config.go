package hubblereceiver

import (
	"context"
	"time"

	"github.com/cilium/hubble-otel/common"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/receiver"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	receiver.CreateSettings       `mapstructure:",squash"`
	configgrpc.GRPCClientSettings `mapstructure:",squash"`

	// BufferSize is the size of the buffer to build before exporting.
	BufferSize int `mapstructure:"buffer_size"`

	FlowEncodingOptions FlowEncodingOptions `mapstructure:"flow_encoding_options"`
	IncludeFlowTypes    IncludeFlowTypes    `mapstructure:"include_flow_types"`

	FallbackServiceNamePrefix string        `mapstructure:"fallback_service_name_prefix"`
	TraceCacheWindow          time.Duration `mapstructure:"trace_cache_window"`
	ParseTraceHeaders         bool          `mapstructure:"parse_trace_headers"`
}

type FlowEncodingOptions struct {
	Traces common.EncodingOptions `mapstructure:"traces"`
	Logs   common.EncodingOptions `mapstructure:"logs"`
}

type IncludeFlowTypes struct {
	Traces common.IncludeFlowTypes `mapstructure:"traces"`
	Logs   common.IncludeFlowTypes `mapstructure:"logs"`
}

func (cfg *Config) NewOutgoingContext(ctx context.Context) context.Context {
	if cfg.GRPCClientSettings.Headers == nil {
		return ctx
	}

	stringConv := make(map[string]string, len(cfg.GRPCClientSettings.Headers))
	for k, v := range cfg.GRPCClientSettings.Headers {
		stringConv[k] = string(v)
	}

	return metadata.NewOutgoingContext(ctx, metadata.New(stringConv))
}
