package chug

import (
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	GenesisTx struct {
		Version   string
		Timestamp int64
		Hash      string
		PubKey    string
		Lock      string
		Data      string
	}
	Paths struct {
		DataDir  string
		BlockDir string
		HugsDir  string
	}
}

const configDir = "./config/"
const configPath = configDir + "main.yaml"

func configCreateDefaultIfNotExists(fs *fs.FileSystem) error {
	if fs.FileExists(configPath) {
		return nil
	}

	c := new(Config)
	c.GenesisTx.Version = "1.0.0"
	c.GenesisTx.Timestamp = 1562764857
	c.GenesisTx.PubKey = "4e1BUTgGBfqVWUkx6dR1NYi9u3GwGPS5uwuo44Bsp5MuJjEGGWC891PknRQBofGkS6MbsPbTLWmt2BZdQKYDgQUwmGqWsfpBHHVGXUfL6uNEiVuBq3AGKsN6uvcvsfjV8fnTiN2zLNwHYzwFEWBVPdcuzSKiXEoGKXB88YYs3RtYV4cubFDcgu4QR2GNLBu4YxvmuyUgBYFYuGRGgD388oM9VLfXVijsAJpUvA9YwRAN2DU4QVb7m78vhLwwhT5QJjiT2L7dVLS8tfY9uLuuSq15TFwdiNS11f1rAUde2nobLEqXp4DAVaUw2BbNvgCr7JmKgLMxVQPUQHwNb5unnHQsZuomgYwocNDhr2UyktfbmyLt8"
	c.GenesisTx.Hash = "5VJfiQNmqcKn7Wzyo25LEUfcPihtwLLNpiDzopWue9Wj"
	c.GenesisTx.Lock = "8TDhCNb535YtS7AhrG7mfbkvCXEmer2329DkzZJt7VXBgjW2HmwJAg1eMeqfiics74FPQD8s9CNcbAXj7Y8CSAwkBD3desuG1ZoC8uNjXU9nVsZXHJ61sPTo5ZvQCpZrxsb5iyWYRt4Qas6Bc67FKtY24jV6WAUZZAsnveNoiiS1vo8kjrU2vDiRFStub9YHvNcnaE5UgACxFGuRQDHnJYqauYK8GgP6pNsoh4wVJxaQoVHn2Lgz6NQBMsRWnnnYrhbHE5iV21NxzxTjNuB4rzjXETK5nwSdPUDMPgUsfrg1hfXyeHvnvmpBKjouBF8DBvG22wM8G6WqUs3JN5wAHmmhE2bNA2"
	c.GenesisTx.Data = utils.Base58ToStr([]byte("hug the planed"))

	c.Paths.DataDir = "./data/"
	c.Paths.BlockDir = c.Paths.DataDir + "transactions/"
	c.Paths.HugsDir = c.Paths.DataDir + "hugs/"

	content, err := yaml.Marshal(&c)
	if err != nil {
		return errors.Wrap(err, "could not serialize config")
	}

	err = fs.WriteFile(configPath, content)
	return errors.Wrap(err, "could not write config file")
}

func configLoad(fs *fs.FileSystem) (*Config, error) {
	data, err := fs.ReadFile(configPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conf := new(Config)
	err = yaml.Unmarshal(data, conf)
	return conf, err
}
