package config

import(
	"os"
)


type AppConfig struct {
	AppEnv     string // env: APP_ENV, default: "local"
	MongoDBURI string // env: MONGODB_URI, default: "mongodb://localhost:27017"
	AWSRegion  string // env: AWS_REGION, default: "us-east-1"
	Port       string // env: PORT, default: "8080"
	SvcName    string // env: SVC_NAME, default: "my-service"
}


// LoadAppConfig populates AppConfig from environment variables.
func LoadAppConfig() AppConfig {
	return AppConfig{
		AppEnv:     GetString("APP_ENV", "local"),
		MongoDBURI: GetString("MONGODB_URI", "mongodb://localhost:27017"),
		AWSRegion:  GetString("AWS_REGION", "us-east-1"),
		Port:       GetString("PORT", "8080"),
		SvcName:    GetString("SVC_NAME", "my-service"),
	}
}

// GetString retrieves the value of an environment variable or returns a default.
func GetString(env, defVal string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return defVal
}