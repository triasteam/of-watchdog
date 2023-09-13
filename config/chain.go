package config

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mitchellh/mapstructure"
	"github.com/openfaas/of-watchdog/logger"
	"github.com/spf13/viper"
)

const EnvPrefix = "CHAIN"

type Chain struct {
	Id                 int64  `mapstructure:"id" json:"id"`
	Addr               string `mapstructure:"addr" json:"addr"`
	FunctionClientAddr string `mapstructure:"function_client_addr" json:"function_client_addr"`
	FunctionOracleAddr string `mapstructure:"function_oracle_addr" json:"function_oracle_addr"`
	VerifierScoreUrl   string `mapstructure:"verifier_score_url" json:"verifier_score_url"`
	FunctionName       string `mapstructure:"function_name" json:"function_name"`
	KeyFilePath        string `mapstructure:"key_file_path" json:"key_file_path"`
	KeyPassword        string `mapstructure:"key_password" json:"key_password"`
	KeyBass64          string `mapstructure:"key_base64" json:"key_base64"`
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
		return errors.New("chain id is wrong, please set the value of env CHAIN_ID")
	}
	if cfg.Addr == "" {
		return errors.New("not found chain addr, please set the value of env CHAIN_ADDR")
	}
	if cfg.FunctionClientAddr == "" {
		return errors.New("not found FunctionClientAddr, please set the value of env CHAIN_FUNCTION_CLIENT_ADDR")
	}
	if cfg.FunctionOracleAddr == "" {
		return errors.New("not found FunctionOracleAddr, please set the value of env CHAIN_FUNCTION_ORACLE_ADDR")
	}
	if cfg.KeyFilePath == "" || cfg.KeyBass64 == "" {
		return errors.New("not found KeyFilePath, please set the value of env CHAIN_KEY_FILE_PATH")
	}

	if cfg.KeyPassword == "" {
		return errors.New("not found chain id, please set the value of env CHAIN_KEY_PASSWORD")
	}

	if cfg.FunctionName == "" {
		return errors.New("not found function name, please set the value of env function name")
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
	var (
		keyBytes []byte
		err      error
	)
	if c.KeyBass64 != "" {
		keyBytes, err = base64.StdEncoding.DecodeString(c.KeyBass64)
		if err != nil {
			logger.Fatal("failed to read key base64 string", "err", err)
		}
	}
	if c.KeyFilePath != "" {
		keyBytes, err = os.ReadFile(c.KeyFilePath)
		if err != nil {
			logger.Fatal("failed to read key file", "err", err)
		}

	}
	if len(keyBytes) == 0 {
		logger.Fatal("not found key")
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

func (c Chain) FuncName() [32]byte {
	functionNameHash := crypto.Keccak256Hash([]byte(c.FunctionName))
	logger.Info("hash function name", "name", c.FunctionName, "hash", functionNameHash, "hex", hex.EncodeToString(functionNameHash[:]))
	return functionNameHash
}

func (c Chain) VerifierScoreAddr() string {
	_, err := url.Parse(c.VerifierScoreUrl)
	if err != nil {
		logger.Fatal("invalid VerifierScoreUrl", "raw", c.VerifierScoreUrl)
		return ""
	}
	return c.VerifierScoreUrl
}

func SetEnv(
	id int64,
	addr, functionClientAddr, functionOracleAddr,
	verifierScoreUrl, functionName,
	keyFilePath, keyPassword string,
) {

	_ = os.Setenv("CHAIN_ID", fmt.Sprintf("%d", id))
	_ = os.Setenv("CHAIN_ADDR", addr)
	_ = os.Setenv("CHAIN_FUNCTION_CLIENT_ADDR", functionClientAddr)
	_ = os.Setenv("CHAIN_FUNCTION_ORACLE_ADDR", functionOracleAddr)
	_ = os.Setenv("CHAIN_VERIFIER_SCORE_URL", verifierScoreUrl)
	_ = os.Setenv("CHAIN_FUNCTION_NAME", functionName)
	_ = os.Setenv("CHAIN_KEY_FILE_PATH", keyFilePath)
	_ = os.Setenv("CHAIN_KEY_PASSWORD", keyPassword)

}
