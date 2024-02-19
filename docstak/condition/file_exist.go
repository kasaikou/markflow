package condition

import (
	"context"

	"github.com/kasaikou/docstak/docstak/resolver"
)

type FileIsExisted struct {
	Config resolver.FileGlobConfig
}

func (cond *FileIsExisted) IsEnable(ctx context.Context) (bool, error) {
	results, err := resolver.ResolveFileGlob(cond.Config)
	if err != nil {
		return false, err
	}

	return len(results) > 0, err
}
