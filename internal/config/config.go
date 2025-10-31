package config

import "github.com/spf13/viper"

type Config struct {
	DBPath    string
	JWTSecret string
	Port      string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return &Config{
		DBPath:    viper.GetString("DB_PATH"),
		JWTSecret: viper.GetString("JWT_SECRET"),
		Port:      viper.GetString("PORT"),
	}, nil
}
