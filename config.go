package chug

import (
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	GenesisTx struct {
		Version   string
		Timestamp int64
		Hash      string
		PubKey    string
		Address   string
		Lock      string
		Data      string
	}
	Paths struct {
		BlockDir    string
		StatesDir   string
		HugsDir     string
		TxStagePath string
	}
	Blocks struct {
		Size int
	}
}

const configDir = "./_config/"
const configPath = configDir + "main.yaml"

func NewDefaultConfig() *Config {
	c := new(Config)
	c.GenesisTx.Version = "1.0.0"
	c.GenesisTx.Timestamp = 1563731742
	// priv key: 69xUcVszy2udRNngci9EWMopFZMhsRsKW8vcYz19drZZ
	c.GenesisTx.PubKey = "2yU5NpDEB3ed9RSL84jXswGCVDW3Mydj9oaGb1MhhLUo7YuL1QMzbBYF3XdMkwZZUBL4mo2czDQeDuNagUifmufn"
	c.GenesisTx.Address = "i6hAWT3BrRBCBehfCJdncLwp3HpSCSpPF"
	c.GenesisTx.Hash = "9zXkjGdPTHupSieTqx11GF67TLZ6fgmBi4y7sqQCnvgJ"
	c.GenesisTx.Lock = "iF1g9zPvp5wKSnhqde2zWiKYE7sbD7sSoTtej2WbBRQAvXhPAsNPrg5n5tgdS1Wvw96cjZ3SEJdcAP2qnQCUDCx"
	c.GenesisTx.Data = "9hdjRw5sudFHrEtyUE" //utils.Base58ToStr([]byte("hug the planed"))

	c.Paths.BlockDir = "./blocks/"
	c.Paths.StatesDir = "./states/"
	c.Paths.HugsDir = c.Paths.StatesDir + "hugs/"
	c.Paths.TxStagePath = c.Paths.StatesDir + "tx-stage/"

	c.Blocks.Size = 100

	return c
}
func NewConfigFromFileOrDefault(fs *fs.FileSystem) (*Config, error) {
	if ConfigExists(fs) {
		result, err := NewConfigFromFile(fs)
		return result, err
	}

	result := NewDefaultConfig()
	err := result.WriteToFile(fs)
	return result, err
}

func ConfigExists(fs *fs.FileSystem) bool {
	return fs.FileExists(configPath)
}

func (c *Config) WriteToFile(fs *fs.FileSystem) error {
	content, err := yaml.Marshal(&c)
	if err != nil {
		return errors.Wrap(err, "could not serialize config")
	}

	err = fs.WriteFile(configPath, content)
	return errors.Wrap(err, "could not write config file")
}

func NewConfigFromFile(fs *fs.FileSystem) (*Config, error) {
	data, err := fs.ReadFile(configPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conf := new(Config)
	err = yaml.Unmarshal(data, conf)
	return conf, err
}
