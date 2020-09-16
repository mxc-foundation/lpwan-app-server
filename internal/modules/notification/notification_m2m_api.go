package notification

import (
	"context"
	"fmt"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
)

type NotificationAPI struct{}

func NewNotificationAPI() *NotificationAPI {
	return &NotificationAPI{}
}

// SendStakeIncomeNotification is called to send email to user when stake revenue is applied
func (a *NotificationAPI) SendStakeIncomeNotification(ctx context.Context, req *pb.SendStakeIncomeNotificationRequest) (*pb.SendStakeIncomeNotificationResponse, error) {
	resp := pb.SendStakeIncomeNotificationResponse{}
	// get user id from organization id
	users, err := organization.Service.St.GetOrganizationUsers(ctx, req.OrganizationId, 999, 0)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to get users for organization")
	}

	for _, v := range users {
		amountMap := make(map[string]string)
		amountMap[email.StakeIncomeAmount] = req.StakeIncomeAmount
		amountMap[email.StakeAmount] = req.StakeAmount
		amountMap[email.StakeIncomeInterest] = req.StakeIncomeInterest

		itemIDMap := make(map[string]string)
		itemIDMap[email.UserID] = fmt.Sprintf("%d", v.UserID)
		itemIDMap[email.StakeID] = req.StakeId
		itemIDMap[email.StakeRevenueID] = req.StakeRevenueId

		dateMap := make(map[string]string)
		dateMap[email.StakeStartDate] = req.StakeStartDate
		dateMap[email.StakeRevenueDate] = req.StakeRevenueDate

		_ = email.SendInvite(v.Email, email.Param{
			Amount: amountMap,
			ItemID: itemIDMap,
			Date:   dateMap,
		}, email.EmailLanguage(serverinfo.GetSettings().DefaultLanguage), email.StakingIncome)

	}

	return &resp, nil
}
