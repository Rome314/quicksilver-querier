package grpcclient

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
)

type PaginationReq interface {
	GetPagination() *query.PageRequest
}
type PaginationResp[entityType any] interface {
	GetPagination() *query.PageResponse
}

type paginator[request PaginationReq, response PaginationResp[entity], entity any] struct {
	req         request
	fn          func(context.Context, request) (response, error)
	getEntities func(response) []entity
}

// TODO: Benchmark to find optimal batch size
func (p *paginator[request, response, entity]) All(ctx context.Context) ([]entity, error) {
	pagesCount := 1
	resp, err := p.fn(ctx, p.req)
	if err != nil {
		return nil, fmt.Errorf("failed to get page %d: %w", pagesCount, err)
	}

	batchSize := p.req.GetPagination().Limit
	var totalCount uint64
	if resp.GetPagination() != nil {
		totalCount = resp.GetPagination().Total
	}

	entities := p.getEntities(resp)
	usedCount := len(entities)

	result := make([]entity, 0, totalCount)
	result = append(result, entities...)

	for usedCount < int(totalCount) {
		pagesCount++
		p.req.GetPagination().Offset += batchSize
		resp, err = p.fn(ctx, p.req)
		if err != nil {
			return nil, fmt.Errorf("failed to get page %d: %w", pagesCount, err)
		}
		entities = p.getEntities(resp)
		usedCount += len(entities)
		result = append(result, entities...)
	}
	return result, nil
}
