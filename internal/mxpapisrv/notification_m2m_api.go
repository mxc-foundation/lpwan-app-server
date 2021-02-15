package mxpapisrv

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

// NotificationAPI keeps variables required for NotificationAPI service
type NotificationAPI struct {
	st     *pgstore.PgStore
	mailer *email.Mailer
}

// NewNotificationAPI creates a new notification API service
func NewNotificationAPI(h *pgstore.PgStore, mailer *email.Mailer) *NotificationAPI {
	return &NotificationAPI{
		st:     h,
		mailer: mailer,
	}
}

// SendStakeIncomeNotification is called to send email to user when stake revenue is applied
func (a *NotificationAPI) SendStakeIncomeNotification(ctx context.Context, req *pb.SendStakeIncomeNotificationRequest) (*pb.SendStakeIncomeNotificationResponse, error) {
	resp := pb.SendStakeIncomeNotificationResponse{}
	// get user id from organization id
	users, err := a.st.GetOrganizationUsers(ctx, req.OrganizationId, 999, 0)
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

		_ = a.mailer.SendStakeIncomeNotification(v.Email, a.mailer.S.DefaultLanguage, email.Param{
			Amount: amountMap,
			ItemID: itemIDMap,
			Date:   dateMap,
		})

	}

	return &resp, nil
}
