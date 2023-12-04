package config

import (
	"github.com/spf13/viper"
)

type S3Bucket struct {
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret"`
	Region          string `mapstructure:"Region"`
	BucketName      string `mapstructure:"BucketName"`
}
type DataBase struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBName     string `mapstructure:"DBNAME"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost     string `mapstructure:"DBHOST"`
	DBPort     string `mapstructure:"DBPORT"`
}
type OTP struct {
	AccountSid string `mapstructure:"AccountSid"`
	AuthToken  string `mapstructure:"AuthToken"`
	ServiceSid string `mapstructure:"ServiceSid"`
}
type Razopay struct {
	RazopayKey    string `mapstructure:"RAZOPAYKEY"`
	RazopaySecret string `mapstructure:"RAZOPAYSECRET"`
}

type Config struct {
	S3aws S3Bucket
	DB DataBase
	Otp OTP
	Razopay Razopay
}

func LoadConfig() (*Config, error) {
	var (
		s3 S3Bucket
		db DataBase
		otp OTP
		razorpay Razopay
	)

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&s3)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&db)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&otp)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&razorpay)
	if err != nil {
		return nil, err
	}
	config := Config{S3aws: s3,DB: db,Razopay: razorpay,Otp:otp}
	return &config, nil
}
