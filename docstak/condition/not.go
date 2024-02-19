package condition

import (
	"context"

	"github.com/kasaikou/docstak/docstak/model"
)

type Not struct{ Internal model.Condition }

func (cond *Not) IsEnable(ctx context.Context) (bool, error) {
	results, err := cond.Internal.IsEnable(ctx)
	if err != nil {
		return false, err
	}

	return !results, nil
}
