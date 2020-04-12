package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pinzolo/xdgdir"
)

type config struct {
	StorageDir string
}

func Read() config {
	configPath := configPath()
	if configJSON, err := ioutil.ReadFile(configPath); os.IsNotExist(err) {
		return writeDefaultConfiguration(configPath)
	} else {
		var conf config
		if err := json.Unmarshal(configJSON, &conf); err == nil {
			return conf
		} else {
			panic(err)
		}
	}
}

func writeDefaultConfiguration(configurationPath string) config {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	conf := config{StorageDir: filepath.Join(usr.HomeDir, "backpocket")}
	if err := os.MkdirAll(filepath.Dir(configurationPath), os.ModePerm); err != nil {
		panic(err)
	}

	if configJSON, err := json.MarshalIndent(conf, "", " "); err != nil {
		panic(err)
	} else if err := ioutil.WriteFile(configurationPath, configJSON, 0640); err != nil {
		panic(err)
	}
	return conf
}

func configPath() string {
	xdgApp := xdgdir.NewApp("backpocket")
	xdgConfigurationFilePath, err := xdgApp.ConfigFile("config.json")
	if err != nil {
		panic(err)
	}

	return xdgConfigurationFilePath
}
