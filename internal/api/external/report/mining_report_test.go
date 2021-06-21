package report

import (
	"context"
	"google.golang.org/grpc/metadata"

	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
)

type testMXPCli struct {
	pb.FinanceReportServiceClient
}

func (tc *testMXPCli) GetMXCMiningReportByDate(ctx context.Context, in *pb.GetMXCMiningReportByDateRequest,
	opts ...grpc.CallOption) (*pb.GetMXCMiningReportByDateResponse, error) {
	return &pb.GetMXCMiningReportByDateResponse{}, nil
}

type testAuth struct {
	auth.Authenticator
	validator func(opts *auth.Options) (*auth.Credentials, error)
}

func (ta *testAuth) GetCredentials(ctx context.Context, opts *auth.Options) (*auth.Credentials, error) {
	if ta.validator != nil {
		return ta.validator(opts)
	}
	return nil, fmt.Errorf("validator is not defined")
}

type streamServer struct {
}

func (s *streamServer) Send(response *api.MiningReportResponse) error {
	return s.SendMsg(response)
}

func (s *streamServer) SetHeader(md metadata.MD) error {
	return nil
}

func (s *streamServer) SendHeader(md metadata.MD) error {
	return nil
}

func (s *streamServer) SetTrailer(md metadata.MD) {
}

func (s *streamServer) Context() context.Context {
	return context.Background()
}

func (s *streamServer) SendMsg(m interface{}) error {
	return nil
}

func (s *streamServer) RecvMsg(m interface{}) error {
	return nil
}

func TestGetMiningReportCSVFileURI(t *testing.T) {
	mxpCli := &testMXPCli{}
	ta := &testAuth{
		validator: func(opts *auth.Options) (*auth.Credentials, error) {
			if opts.OrgID != 1 {
				return &auth.Credentials{
					UserID:     2,
					Username:   "user2",
					IsExisting: true,
					OrgID:      2,
				}, nil
			}
			return &auth.Credentials{
				UserID:     1,
				Username:   "user1",
				IsExisting: true,
				OrgID:      1,
				IsOrgAdmin: true,
			}, nil
		},
	}
	server := NewServer(mxpCli, ta, "local_test_server")
	request := api.MiningReportRequest{
		OrganizationId: 1,
		Currency:       []string{"ETH_MXC"},
		FiatCurrency:   "usd",
		Start:          timestamppb.New(time.Now().AddDate(-1, 0, 0)),
		End:            timestamppb.New(time.Now()),
		Decimals:       6,
	}

	err := server.MiningReportCSV(&request, &streamServer{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMiningReportPDFFileURI(t *testing.T) {
	mxpCli := &testMXPCli{}
	ta := &testAuth{
		validator: func(opts *auth.Options) (*auth.Credentials, error) {
			if opts.OrgID != 1 {
				return &auth.Credentials{
					UserID:     2,
					Username:   "user2",
					IsExisting: true,
					OrgID:      2,
				}, nil
			}
			return &auth.Credentials{
				UserID:     1,
				Username:   "user1",
				IsExisting: true,
				OrgID:      1,
				IsOrgAdmin: true,
			}, nil
		},
	}
	server := NewServer(mxpCli, ta, "local_test_server")
	request := api.MiningReportRequest{
		OrganizationId: 1,
		Currency:       []string{"ETH_MXC"},
		FiatCurrency:   "usd",
		Start:          timestamppb.New(time.Now().AddDate(-1, 0, 0)),
		End:            timestamppb.New(time.Now()),
		Decimals:       6,
	}
	err := server.MiningReportPDF(&request, &streamServer{})
	if err != nil {
		t.Fatal(err)
	}
}
