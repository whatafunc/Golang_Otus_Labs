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
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}
