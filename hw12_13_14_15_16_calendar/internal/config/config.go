package config

type Config struct {
	Logger         LoggerConf    `yaml:"logger"`
	HTTP           HTTPConf      `yaml:"http"`
	Storage        StorageConfig `yaml:"storage"`
	MigrationsPath string        `yaml:"migrationsPath"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type HTTPConf struct {
	Listen string `yaml:"listen"`
}

type StorageConfig struct {
	Type     string         `yaml:"type"`
	Redis    RedisConfig    `yaml:"redis"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}
