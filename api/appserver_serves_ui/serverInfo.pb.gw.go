// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: serverInfo.proto

/*
Package appserver_serves_ui is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package appserver_serves_ui

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = descriptor.ForMessage

func request_ServerInfoService_GetAppserverVersion_0(ctx context.Context, marshaler runtime.Marshaler, client ServerInfoServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq empty.Empty
	var metadata runtime.ServerMetadata

	msg, err := client.GetAppserverVersion(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ServerInfoService_GetAppserverVersion_0(ctx context.Context, marshaler runtime.Marshaler, server ServerInfoServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq empty.Empty
	var metadata runtime.ServerMetadata

	msg, err := server.GetAppserverVersion(ctx, &protoReq)
	return msg, metadata, err

}

func request_ServerInfoService_GetServerRegion_0(ctx context.Context, marshaler runtime.Marshaler, client ServerInfoServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq empty.Empty
	var metadata runtime.ServerMetadata

	msg, err := client.GetServerRegion(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ServerInfoService_GetServerRegion_0(ctx context.Context, marshaler runtime.Marshaler, server ServerInfoServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq empty.Empty
	var metadata runtime.ServerMetadata

	msg, err := server.GetServerRegion(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterServerInfoServiceHandlerServer registers the http handlers for service ServerInfoService to "mux".
// UnaryRPC     :call ServerInfoServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
func RegisterServerInfoServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ServerInfoServiceServer) error {

	mux.Handle("GET", pattern_ServerInfoService_GetAppserverVersion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ServerInfoService_GetAppserverVersion_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ServerInfoService_GetAppserverVersion_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ServerInfoService_GetServerRegion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ServerInfoService_GetServerRegion_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ServerInfoService_GetServerRegion_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterServerInfoServiceHandlerFromEndpoint is same as RegisterServerInfoServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterServerInfoServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterServerInfoServiceHandler(ctx, mux, conn)
}

// RegisterServerInfoServiceHandler registers the http handlers for service ServerInfoService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterServerInfoServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterServerInfoServiceHandlerClient(ctx, mux, NewServerInfoServiceClient(conn))
}

// RegisterServerInfoServiceHandlerClient registers the http handlers for service ServerInfoService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ServerInfoServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ServerInfoServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ServerInfoServiceClient" to call the correct interceptors.
func RegisterServerInfoServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ServerInfoServiceClient) error {

	mux.Handle("GET", pattern_ServerInfoService_GetAppserverVersion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ServerInfoService_GetAppserverVersion_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ServerInfoService_GetAppserverVersion_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ServerInfoService_GetServerRegion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ServerInfoService_GetServerRegion_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ServerInfoService_GetServerRegion_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_ServerInfoService_GetAppserverVersion_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"api", "server-info", "appserver-version"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_ServerInfoService_GetServerRegion_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"api", "server-info", "server-region"}, "", runtime.AssumeColonVerbOpt(true)))
)

var (
	forward_ServerInfoService_GetAppserverVersion_0 = runtime.ForwardResponseMessage

	forward_ServerInfoService_GetServerRegion_0 = runtime.ForwardResponseMessage
)
