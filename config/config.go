package config

import "github.com/ganeshdipdumbare/goenv"

type EnvVar struct {
	MongoUri           string `json:"mongo_uri"`
	MongoDb            string `json:"mongo_db"`
	Port               string `json:"port"`
	MigrationFilesPath string `json:"migration_files_path"`
	ApiKey             string `json:"api_key"`
}

var (
	envVars = &EnvVar{
		Port:               "8080",
		MongoDb:            "scootin-aboot-db",
		MigrationFilesPath: "file://migration",
		MongoUri:           "mongodb://localhost:27017",
		ApiKey:             "secretkey",
	}
)

func init() {
	goenv.SyncEnvVar(&envVars)
}

func Get() *EnvVar {
	return envVars
}
