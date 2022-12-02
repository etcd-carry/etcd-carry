package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type Probe interface {
	Probe(conn *grpc.ClientConn) (Result, string, error)
}

type grpcProbe struct{}

func New() Probe {
	return grpcProbe{}
}

type Result string

const (
	// Success Result
	Success Result = "success"
	// Failure Result
	Failure Result = "failure"
	// Unknown Result
	Unknown Result = "unknown"
)

func (p grpcProbe) Probe(conn *grpc.ClientConn) (Result, string, error) {
	cli := grpchealth.NewHealthClient(conn)
	resp, err := cli.Check(context.TODO(), &grpchealth.HealthCheckRequest{})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			switch s.Code() {
			case codes.Unimplemented:
				return Failure, fmt.Sprintf("grpc health protocol is not implemented or not supported/enabled in this service: %s", s.Message()), nil
			case codes.DeadlineExceeded:
				return Failure, fmt.Sprintf("timeout: health rpc timeout: %s", s.Message()), nil
			case codes.Unavailable:
				return Failure, fmt.Sprintf("unavailable: service %v is currently unavailable", conn.Target()), err
			default:
				fmt.Println("rpc probe failed: ", s.Code())
			}
		} else {
			fmt.Println("health rpc probe failed: ", err)
		}
		return Failure, fmt.Sprintf("error: health rpc probe failed: %+v", err), nil
	}

	if resp.GetStatus() != grpchealth.HealthCheckResponse_SERVING {
		return Failure, fmt.Sprintf("service unhealthy, responded with %q", resp.GetStatus().String()), nil
	}

	return Success, "service healthy", nil
}
