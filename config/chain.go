package config

import (
	"errors"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/openfaas/of-watchdog/logger"
)

const EnvPrefix = "CHAIN"

type Chain struct {
	Id                 int64  `mapstructure:"id" json:"id"`
	Addr               string `mapstructure:"addr" json:"addr"`
	FunctionClientAddr string `mapstructure:"function_client_addr" json:"function_client_addr"`
	FunctionOracleAddr string `mapstructure:"function_oracle_addr" json:"function_oracle_addr"`
	KeyFilePath        string `mapstructure:"key_file_path" json:"key_file_path"`
	KeyPassword        string `mapstructure:"key_password" json:"key_password"`
}

func LoadChainConfig() *Chain {
	cfg := &Chain{}
	v := viper.New()

	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	logger.Debug("all configs", "v", v.AllSettings())

	var envMap map[string]interface{}
	err := mapstructure.Decode(cfg, &envMap)
	if err != nil {
		logger.Fatal("failed to load config", "err", err)
		return nil
	}

	for k, _ := range envMap {
		err := v.BindEnv(k)
		if err != nil {
			logger.Fatal("fail to bind env", "err", err)
		}
	}
	if err = v.Unmarshal(cfg); err != nil {
		logger.Fatal("fail to unmarshal chain Config", "err", err)
	}
	logger.Info("successfully unmarshal chain config", "value", cfg)
	if err = validateChainConfig(cfg); err != nil {
		logger.Fatal("fail to validate chain Config", "err", err)
	}

	logger.Info("successfully load to chain config", "value", cfg)
	return cfg
}

func validateChainConfig(cfg *Chain) error {

	if cfg.Id == 0 {
		return errors.New("chain id is wrong, please set the value of env chain id")
	}
	if cfg.Addr == "" {
		return errors.New("not found chain addr, please set the value of env chain id")
	}
	if cfg.FunctionClientAddr == "" {
		return errors.New("not found chain id, please set the value of env chain id")
	}
	if cfg.FunctionOracleAddr == "" {
		return errors.New("not found chain id, please set the value of env chain id")
	}
	if cfg.KeyFilePath == "" {
		return errors.New("not found chain id, please set the value of env chain id")
	}

	if cfg.KeyPassword == "" {
		return errors.New("not found chain id, please set the value of env chain id")
	}

	return nil
}

func (c Chain) ChainID() int64 {
	return c.Id
}

func (c Chain) ChainAddr() string {
	return c.Addr
}

func (c Chain) Key() *keystore.Key {

	keyBytes, err := os.ReadFile(c.KeyFilePath)
	if err != nil {
		logger.Fatal("failed to read key file", "err", err)
	}

	key, err := keystore.DecryptKey(keyBytes, c.KeyPassword)
	if err != nil {
		logger.Fatal("wrong key", "key", string(keyBytes), "pw", c.KeyPassword)
	}
	return key
}

func (c Chain) FuncClientAddr() string {
	return c.FunctionClientAddr
}

func (c Chain) FuncOracleClientAddr() string {
	return c.FunctionOracleAddr
}
