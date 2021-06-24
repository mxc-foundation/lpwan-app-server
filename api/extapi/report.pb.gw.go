// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: report.proto

/*
Package extapi is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package extapi

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = descriptor.ForMessage
var _ = metadata.Join

func request_ReportService_GetFiatCurrencyList_0(ctx context.Context, marshaler runtime.Marshaler, client ReportServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetFiatCurrencyListRequest
	var metadata runtime.ServerMetadata

	msg, err := client.GetFiatCurrencyList(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ReportService_GetFiatCurrencyList_0(ctx context.Context, marshaler runtime.Marshaler, server ReportServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetFiatCurrencyListRequest
	var metadata runtime.ServerMetadata

	msg, err := server.GetFiatCurrencyList(ctx, &protoReq)
	return msg, metadata, err

}

var (
	filter_ReportService_MiningReportCSV_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_ReportService_MiningReportCSV_0(ctx context.Context, marshaler runtime.Marshaler, client ReportServiceClient, req *http.Request, pathParams map[string]string) (ReportService_MiningReportCSVClient, runtime.ServerMetadata, error) {
	var protoReq MiningReportRequest
	var metadata runtime.ServerMetadata

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_ReportService_MiningReportCSV_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	stream, err := client.MiningReportCSV(ctx, &protoReq)
	if err != nil {
		return nil, metadata, err
	}
	header, err := stream.Header()
	if err != nil {
		return nil, metadata, err
	}
	metadata.HeaderMD = header
	return stream, metadata, nil

}

var (
	filter_ReportService_MiningReportPDF_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_ReportService_MiningReportPDF_0(ctx context.Context, marshaler runtime.Marshaler, client ReportServiceClient, req *http.Request, pathParams map[string]string) (ReportService_MiningReportPDFClient, runtime.ServerMetadata, error) {
	var protoReq MiningReportRequest
	var metadata runtime.ServerMetadata

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_ReportService_MiningReportPDF_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	stream, err := client.MiningReportPDF(ctx, &protoReq)
	if err != nil {
		return nil, metadata, err
	}
	header, err := stream.Header()
	if err != nil {
		return nil, metadata, err
	}
	metadata.HeaderMD = header
	return stream, metadata, nil

}

// RegisterReportServiceHandlerServer registers the http handlers for service ReportService to "mux".
// UnaryRPC     :call ReportServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterReportServiceHandlerFromEndpoint instead.
func RegisterReportServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ReportServiceServer) error {

	mux.Handle("GET", pattern_ReportService_GetFiatCurrencyList_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ReportService_GetFiatCurrencyList_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ReportService_GetFiatCurrencyList_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ReportService_MiningReportCSV_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	})

	mux.Handle("GET", pattern_ReportService_MiningReportPDF_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	})

	return nil
}

// RegisterReportServiceHandlerFromEndpoint is same as RegisterReportServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterReportServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
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

	return RegisterReportServiceHandler(ctx, mux, conn)
}

// RegisterReportServiceHandler registers the http handlers for service ReportService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterReportServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterReportServiceHandlerClient(ctx, mux, NewReportServiceClient(conn))
}

// RegisterReportServiceHandlerClient registers the http handlers for service ReportService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ReportServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ReportServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ReportServiceClient" to call the correct interceptors.
func RegisterReportServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ReportServiceClient) error {

	mux.Handle("GET", pattern_ReportService_GetFiatCurrencyList_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ReportService_GetFiatCurrencyList_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ReportService_GetFiatCurrencyList_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ReportService_MiningReportCSV_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ReportService_MiningReportCSV_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ReportService_MiningReportCSV_0(ctx, mux, outboundMarshaler, w, req, func() (proto.Message, error) { return resp.Recv() }, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ReportService_MiningReportPDF_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ReportService_MiningReportPDF_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ReportService_MiningReportPDF_0(ctx, mux, outboundMarshaler, w, req, func() (proto.Message, error) { return resp.Recv() }, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_ReportService_GetFiatCurrencyList_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"api", "report", "supported-fiat-currencies"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_ReportService_MiningReportCSV_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "report", "mining-income", "csv"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_ReportService_MiningReportPDF_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "report", "mining-income", "pdf"}, "", runtime.AssumeColonVerbOpt(true)))
)

var (
	forward_ReportService_GetFiatCurrencyList_0 = runtime.ForwardResponseMessage

	forward_ReportService_MiningReportCSV_0 = runtime.ForwardResponseStream

	forward_ReportService_MiningReportPDF_0 = runtime.ForwardResponseStream
)