package dhx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
)

type testDHXCli struct {
	pb.DHXServiceClient
	lastMining *pb.DHXGetLastMiningResponse
}

func (td *testDHXCli) DHXGetLastMining(ctx context.Context, req *pb.DHXGetLastMiningRequest, opts ...grpc.CallOption) (*pb.DHXGetLastMiningResponse, error) {
	return td.lastMining, nil
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

func TestDHXGetLastMining(t *testing.T) {
	td := &testDHXCli{}
	ta := &testAuth{
		validator: func(opts *auth.Options) (*auth.Credentials, error) {
			if opts.OrgID == 2 {
				return &auth.Credentials{
					UserID:     3,
					Username:   "three",
					IsExisting: true,
					OrgID:      2,
					IsOrgUser:  true,
				}, nil
			}
			if opts.OrgID == 3 {
				return nil, fmt.Errorf("not authenticated")
			}
			if opts.OrgID == 4 {
				return &auth.Credentials{
					UserID:     4,
					Username:   "four",
					IsExisting: true,
					OrgID:      4,
				}, nil
			}
			return &auth.Credentials{
				UserID:     2,
				Username:   "two",
				IsExisting: true,
			}, nil
		},
	}
	srv := NewServer(td, ta, nil)
	ctx := context.Background()

	// not authenticated user can't get response
	resp1, err := srv.DHXGetLastMining(ctx, &api.DHXGetLastMiningRequest{OrgId: 3})
	if resp1 != nil {
		t.Errorf("got non-nil response for unauthenticated user")
	}
	if err == nil {
		t.Errorf("got no error for unauthenticated user")
	}
	if st := status.Convert(err); st.Code() != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated got %s", st.Code().String())
	}

	// non org user can't get org's mining
	resp2, err := srv.DHXGetLastMining(ctx, &api.DHXGetLastMiningRequest{OrgId: 4})
	if resp2 != nil {
		t.Errorf("got non-nil response for unauthenticated user")
	}
	if err == nil {
		t.Errorf("got no error for non-org member")
	}
	if st := status.Convert(err); st.Code() != codes.PermissionDenied {
		t.Errorf("expected PermissionDenied got %s", st.Code().String())
	}

	date := &timestamppb.Timestamp{Seconds: time.Now().Unix()}
	// just supernode totals
	td.lastMining = &pb.DHXGetLastMiningResponse{
		Date:        date,
		MiningPower: "1000000",
		DhxAmount:   "1000000",
	}
	resp3, err := srv.DHXGetLastMining(ctx, &api.DHXGetLastMiningRequest{})
	if err != nil {
		t.Errorf("expected success, got: %v", err)
	}
	if resp3 == nil || resp3.MiningPower != "1000000" || resp3.OrgId != 0 || resp3.CouncilId != 0 {
		t.Errorf("response is not as expected: %#v", resp3)
	}

	// with org and council
	td.lastMining.OrgId = 2
	td.lastMining.OrgMiningPower = "1000"
	td.lastMining.DhxAmount = "1000"
	td.lastMining.CouncilId = 9
	td.lastMining.CouncilName = "nine"
	td.lastMining.CouncilMiningPower = "3000"
	td.lastMining.CouncilDhxAmount = "3000"
	resp4, err := srv.DHXGetLastMining(ctx, &api.DHXGetLastMiningRequest{OrgId: 2})
	if err != nil {
		t.Errorf("expected success, got: %v", err)
	}
	if resp4 == nil || resp4.MiningPower != "1000000" || resp4.OrgId != 2 || resp4.CouncilId != 9 || resp4.CouncilDhxAmount != "3000" {
		t.Errorf("response is not as expected: %#v", resp4)
	}
}
