package m2m_ui

import (
	"context"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
)

type ExtAccountServerAPI struct {
	validator auth.Validator
}

func NewMoneyServerAPI(validator auth.Validator) *ExtAccountServerAPI {
	return &ExtAccountServerAPI{
		validator: validator,
	}
}

func (s *ExtAccountServerAPI) ModifyMoneyAccount(ctx context.Context, req *api.ModifyMoneyAccountRequest) (*api.ModifyMoneyAccountResponse, error) {

}

func (s *ExtAccountServerAPI) GetChangeMoneyAccountHistory(ctx context.Context, req *api.GetMoneyAccountChangeHistoryRequest) (*api.GetMoneyAccountChangeHistoryResponse, error) {

}

func (s *ExtAccountServerAPI) GetActiveMoneyAccount(ctx context.Context, req *api.GetActiveMoneyAccountRequest) (*api.GetActiveMoneyAccountResponse, error) {

}