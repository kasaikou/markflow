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
