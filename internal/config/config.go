package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cors     CorsConfig
}

type ServerConfig struct {
	Url            string   `env:"APP_URL" envDefault:"localhost:8000"`
	TrustedProxies []string `env:"TRUSTED_PROXIES" envDefault:""`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"3306"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type CorsConfig struct {
	AllowOrigins     []string `env:"CORS_ALLOW_ORIGINS" envDefault:"*"`
	AllowMethods     []string `env:"CORS_ALLOW_METHODS" envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Url:            GetEnv("APP_URL", "localhost:8000"),
			TrustedProxies: GetEnvSlice("TRUSTED_PROXIES", []string{}),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnvInt("DB_PORT", 3306),
			Username: GetEnv("DB_USERNAME", ""),
			Password: GetEnv("DB_PASSWORD", ""),
			Name:     GetEnv("DB_NAME", "golang"),
		},
		Cors: CorsConfig{
			AllowOrigins:     GetEnvSlice("CORS_ALLOW_ORIGINS", []string{"*"}),
			AllowMethods:     GetEnvSlice("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowCredentials: GetEnvBool("CORS_ALLOW_CREDENTIALS", true),
		},
	}

	return cfg, nil
}
