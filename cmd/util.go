package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/starlingbank/vault-unsealer/pkg/kv"
	"github.com/starlingbank/vault-unsealer/pkg/kv/aws_kms"
	"github.com/starlingbank/vault-unsealer/pkg/kv/aws_ssm"
	"github.com/starlingbank/vault-unsealer/pkg/kv/env_file"
	"github.com/starlingbank/vault-unsealer/pkg/kv/cloudkms"
	"github.com/starlingbank/vault-unsealer/pkg/kv/gcs"
	"github.com/starlingbank/vault-unsealer/pkg/kv/local"
	"github.com/starlingbank/vault-unsealer/pkg/vault"
)

func vaultConfigForConfig(cfg *viper.Viper) (vault.Config, error) {

	return vault.Config{
		KeyPrefix: "vault",

		SecretShares:    appConfig.GetInt(cfgSecretShares),
		SecretThreshold: appConfig.GetInt(cfgSecretThreshold),

		InitRootToken:  appConfig.GetString(cfgInitRootToken),
		StoreRootToken: appConfig.GetBool(cfgStoreRootToken),

		OverwriteExisting: appConfig.GetBool(cfgOverwriteExisting),
	}, nil
}

func kvStoreForConfig(cfg *viper.Viper) (kv.Service, error) {

	switch cfg.GetString(cfgMode) {
	case cfgModeValueGoogleCloudKMSGCS:

		g, err := gcs.New(
			cfg.GetString(cfgGoogleCloudStorageBucket),
			cfg.GetString(cfgGoogleCloudStoragePrefix),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating google cloud storage kv store: %s", err.Error())
		}

		kms, err := cloudkms.New(g,
			cfg.GetString(cfgGoogleCloudKMSProject),
			cfg.GetString(cfgGoogleCloudKMSLocation),
			cfg.GetString(cfgGoogleCloudKMSKeyRing),
			cfg.GetString(cfgGoogleCloudKMSCryptoKey),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating google cloud kms kv store: %s", err.Error())
		}

		return kms, nil

	case cfgModeValueAWSKMSSSM:
		ssm, err := aws_ssm.New(cfg.GetString(cfgAWSSSMKeyPrefix))
		if err != nil {
			return nil, fmt.Errorf("error creating AWS SSM kv store: %s", err.Error())
		}

		kms, err := aws_kms.New(ssm, cfg.GetString(cfgAWSKMSKeyID))
		if err != nil {
			return nil, fmt.Errorf("error creating AWS KMS ID kv store: %s", err.Error())
		}

		return kms, nil

	case cfgModeValueLocal:
		return local.New(cfg.GetString(cfgLocalKeyDir))

	case cfgModeValueAWSKMSParamFile:
		envf, err := env_file.New(cfg.GetString(cfgEnvFileName))
		if err != nil {
			return nil, fmt.Errorf("error creating Env file kv store: %s", err.Error())
		}

		kms, err := aws_kms.New(envf, cfg.GetString(cfgAWSKMSKeyID))
		if err != nil {
			return nil, fmt.Errorf("error creating AWS KMS ID kv store: %s", err.Error())
		}

		return kms, nil

	default:
		return nil, fmt.Errorf("Unsupported backend mode: '%s'", cfg.GetString(cfgMode))
	}
}
