package store

import (
	"context"

	"github.com/gofrs/uuid"
)

type MulticastGroupStore interface {
	// validator
	CheckCreateMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckReadMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
	CheckUpdateDeleteMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
	CheckMulticastGroupQueueAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
}

func (h *Handler) CheckCreateMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckCreateMulticastGroupsAccess(ctx, username, organizationID, userID)
}
func (h *Handler) CheckListMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckListMulticastGroupsAccess(ctx, username, organizationID, userID)
}
func (h *Handler) CheckReadMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckReadMulticastGroupAccess(ctx, username, multicastGroupID, userID)
}
func (h *Handler) CheckUpdateDeleteMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckUpdateDeleteMulticastGroupAccess(ctx, username, multicastGroupID, userID)
}
func (h *Handler) CheckMulticastGroupQueueAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckMulticastGroupQueueAccess(ctx, username, multicastGroupID, userID)
}
