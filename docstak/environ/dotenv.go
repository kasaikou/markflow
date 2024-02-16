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

package environ

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadDotenv(filename string, fn func(key, value string)) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	kv, err := godotenv.Parse(file)
	if err != nil {
		return err
	}

	for key, value := range kv {
		fn(key, value)
	}

	return nil
}
