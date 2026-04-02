package configs

import "github.com/spf13/viper"

type conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    int    `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort int    `mapstructure:"GRAPHQL_SERVER_PORT"`
	RabbitMQDSN       string `mapstructure:"RABBITMQ_DSN"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	// ReadInConfig is optional: absent in Docker where env vars are injected directly.
	_ = viper.ReadInConfig()
	// BindEnv registers each key so that Unmarshal picks up env vars even
	// when no config file is present (AutomaticEnv alone is not enough for Unmarshal).
	for _, key := range []string{
		"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"WEB_SERVER_PORT", "GRPC_SERVER_PORT", "GRAPHQL_SERVER_PORT", "RABBITMQ_DSN",
	} {
		_ = viper.BindEnv(key)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
