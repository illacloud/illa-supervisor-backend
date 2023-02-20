package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

const DEPLOY_MODE_SELF_HOST = "self-host"
const DEPLOY_MODE_CLOUD = "cloud"
const DRIVE_TYPE_AWS = "aws"
const DRIVE_TYPE_MINIO = "minio"

type Config struct {
	// server config
	ServerHost         string `env:"ILLA_SERVER_HOST"              envDefault:"0.0.0.0"`
	ServerPort         string `env:"ILLA_SERVER_PORT"              envDefault:"8003"`
	InternalServerPort string `env:"ILLA_SERVER_INTERNAL_PORT"     envDefault:"9001"`
	ServerMode         string `env:"ILLA_SERVER_MODE"              envDefault:"debug"`
	DeployMode         string `env:"ILLA_DEPLOY_MODE"              envDefault:"self-host"`
	ServeHTTPS         string `env:"ILLA_DEPLOY_SERVE_HTTPS"       envDefault:"false"`

	// drive config
	DriveType             string `env:"ILLA_DRIVE_TYPE"               envDefault:""`
	DriveAccessKeyID      string `env:"ILLA_DRIVE_ACCESS_KEY_ID"      envDefault:"minioadmin"`
	DriveAccessKeySecret  string `env:"ILLA_DRIVE_ACCESS_KEY_SECRET"  envDefault:"minioadmin"`
	DriveRegion           string `env:"ILLA_DRIVE_REGION"             envDefault:""`
	DriveEndpoint         string `env:"ILLA_DRIVE_ENDPOINT"           envDefault:"127.0.0.1:9000"`
	DriveSystemBucketName string `env:"ILLA_DRIVE_SYSTEM_BUCKET_NAME" envDefault:"illa-supervisor"`
	DriveTeamBucketName   string `env:"ILLA_DRIVE_TEAM_BUCKET_NAME"   envDefault:"illa-supervisor-team"`
	DriveUploadTimeoutRaw string `env:"ILLA_DRIVE_UPLOAD_TIMEOUT"     envDefault:"30s"`
	DriveUploadTimeout    time.Duration
}

func GetConfig() (*Config, error) {
	// fetch
	cfg := &Config{}
	err := env.Parse(cfg)

	// process data
	var errInParseDuration error
	cfg.DriveUploadTimeout, errInParseDuration = time.ParseDuration(cfg.DriveUploadTimeoutRaw)
	if errInParseDuration != nil {
		return nil, errInParseDuration
	}

	// ok
	fmt.Printf("----------------\n")
	fmt.Printf("%+v\n", cfg)
	fmt.Printf("%+v\n", err)

	return cfg, err
}

func (c *Config) IsSelfHostMode() bool {
	if c.DeployMode == DEPLOY_MODE_SELF_HOST {
		return true
	}
	return false
}

func (c *Config) IsCloudMode() bool {
	if c.DeployMode == DEPLOY_MODE_CLOUD {
		return true
	}
	return false
}

func (c *Config) IsServeHTTPS() bool {
	if c.ServeHTTPS == "true" {
		return true
	}
	return false
}

func (c *Config) GetServeHTTPAddress() string {
	if c.ServeHTTPS == "true" {
		return "https://" + c.ServerHost
	}
	return "http://" + c.ServerHost
}

func (c *Config) GetDriveType() string {
	return c.DriveType
}

func (c *Config) IsAWSDrive() bool {
	if c.DriveType == DRIVE_TYPE_AWS {
		return true
	}
	return false
}

func (c *Config) IsMINIODrive() bool {
	if c.DriveType == DRIVE_TYPE_MINIO {
		return true
	}
	return false
}

func (c *Config) GetAWSS3AccessKeyID() string {
	return c.DriveAccessKeyID
}

func (c *Config) GetAWSS3AccessKeySecret() string {
	return c.DriveAccessKeySecret
}

func (c *Config) GetAWSS3Region() string {
	return c.DriveRegion
}

func (c *Config) GetAWSS3SystemBucketName() string {
	return c.DriveSystemBucketName
}

func (c *Config) GetAWSS3TeamBucketName() string {
	return c.DriveTeamBucketName
}

func (c *Config) GetAWSS3Timeout() time.Duration {
	return c.DriveUploadTimeout
}

func (c *Config) GetMINIOAccessKeyID() string {
	return c.DriveAccessKeyID
}

func (c *Config) GetMINIOAccessKeySecret() string {
	return c.DriveAccessKeySecret
}

func (c *Config) GetMINIOEndpoint() string {
	return c.DriveEndpoint
}

func (c *Config) GetMINIOSystemBucketName() string {
	return c.DriveSystemBucketName
}

func (c *Config) GetMINIOTeamBucketName() string {
	return c.DriveTeamBucketName
}

func (c *Config) GetMINIOTimeout() time.Duration {
	return c.DriveUploadTimeout
}
