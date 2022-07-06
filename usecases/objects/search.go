package objects

import (
	"context"
	"fmt"

	"github.com/semi-technologies/weaviate/entities/additional"
	"github.com/semi-technologies/weaviate/entities/filters"
	"github.com/semi-technologies/weaviate/entities/models"
)

type QueryInput struct {
	Class      string
	Offset     int
	Limit      int
	Filters    filters.LocalFilter
	Sort       []filters.Sort
	Additional additional.Properties
}

func (m *Manager) Query(ctx context.Context, principal *models.Principal,
	class string,
	offset, limit *int64,
	sort, order *string,
	additional additional.Properties) ([]*models.Object, *Error) {
	path := fmt.Sprintf("objects/%s", class)
	if err := m.authorizer.Authorize(principal, "list", path); err != nil {
		return nil, &Error{path, StatusForbidden, err}
	}
	unlock, err := m.locks.LockConnector()
	if err != nil {
		return nil, &Error{"cannot lock", StatusInternalServerError, err}
	}
	defer unlock()

	smartOffset, smartLimit, err := m.localOffsetLimit(offset, limit)
	if err != nil {
		return nil, &Error{"offset or limit", StatusBadRequest, err}
	}
	q := QueryInput{
		Class:      class,
		Offset:     smartOffset,
		Limit:      smartLimit,
		Sort:       m.getSort(sort, order),
		Additional: additional}
	res, rerr := m.vectorRepo.Query(ctx, &q)
	if err != nil {
		return nil, rerr
	}

	if m.modulesProvider != nil {
		res, err = m.modulesProvider.ListObjectsAdditionalExtend(ctx, res, additional.ModuleParams)
		if err != nil {
			return nil, &Error{"extend results", StatusInternalServerError, err}
		}
	}

	return res.ObjectsWithVector(additional.Vector), nil
}
