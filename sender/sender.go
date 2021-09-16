package sender

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Exporter interface {
	Export(context.Context, <-chan protoreflect.Message) error
}

type NullExporter struct{}

func (s *NullExporter) Export(ctx context.Context, flows <-chan protoreflect.Message) error {
	for {
		<-flows
	}
}

func Run(ctx context.Context, s Exporter, flows <-chan protoreflect.Message, errs chan<- error) {
	for {
		switch err := s.Export(ctx, flows); err {
		case io.EOF, context.Canceled:
			return
		case nil:
			// fmt.Printf("wrote %d entries to the OTLP receiver\n", logBufferSize)
			continue
		default:
			if status.Code(err) == codes.Canceled {
				return
			}
			errs <- err
			return
		}
	}
}