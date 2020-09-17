package storage

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// SearchResult defines a search result.
type SearchResult store.SearchResult

// GlobalSearch performs a search on organizations, applications, gateways
// and devices.
func GlobalSearch(ctx context.Context, handler *store.Handler, userID int64, globalAdmin bool, search string, limit, offset int) ([]SearchResult, error) {
	res, err := handler.GlobalSearch(ctx, userID, globalAdmin, search, limit, offset)
	if err != nil {
		return nil, err
	}

	var resList []SearchResult
	for _, v := range res {
		resItem := SearchResult(v)
		resList = append(resList, resItem)
	}

	return resList, nil
}
