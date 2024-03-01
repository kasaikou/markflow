/*
Copyright 2024 Kasai Kou

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package condition

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"sort"

	"github.com/kasaikou/docstak/docstak/resolver"
)

type FileNotChanged struct {
	Config resolver.FileGlobConfig
	MD5    string
}

func (cond *FileNotChanged) CurrentMD5(ctx context.Context) (string, error) {
	results, err := resolver.ResolveFileGlob(cond.Config)
	if err != nil {
		return "", err
	}

	sort.Strings(results)
	hash := md5.New()

	for i := range results {
		err := func(ctx context.Context) error {
			if err := ctx.Err(); err != nil {
				return err
			}

			file, err := os.Open(results[i])
			if err != nil {
				return err
			}

			defer file.Close()
			_, err = io.Copy(hash, file)
			if err != nil {
				return err
			}

			return nil
		}(ctx)

		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (cond *FileNotChanged) IsEnable(ctx context.Context) (bool, error) {

	if cond.MD5 == "" {
		return false, nil
	}

	current, err := cond.CurrentMD5(ctx)
	if err != nil {
		return false, err
	}

	return cond.MD5 == current, nil
}
