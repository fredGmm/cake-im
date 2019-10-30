package path

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func GetConfigPath(path string) (p string, error error) {

	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	}
	if os.IsNotExist(err) {
		abPath, _ := filepath.Abs(os.Args[0])
		dir := filepath.Dir(abPath)
		log.Printf(dir)
		separator := string(filepath.Separator)
		configPath := dir + separator + "config" + separator + "Config.ini"
		_, err2 := os.Stat(configPath)
		if err2 == nil {
			return configPath, nil
		}
	}
	return "", errors.New("find config path unknown error")
}
