package config

import "testing"

func TestGetString(t *testing.T) {
	t.Setenv("APP_ENV", "prod")

	got := GetString("APP_ENV", "local")

	if got != "prod" {
		t.Fatalf(
			"expected %q but got %q",
			"prod",
			got,
		)
	}
}

func TestGetString_DefaultValue(t *testing.T) {
	got := GetString("UNKNOWN_ENV", "default")

	if got != "default" {
		t.Fatalf(
			"expected %q but got %q",
			"default",
			got,
		)
	}
}

func TestLoadAppConfig(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("MONGODB_URI", "mongodb://prod:27017")
	t.Setenv("AWS_REGION", "us-west-2")
	t.Setenv("PORT", "9090")
	t.Setenv("SVC_NAME", "weather-api")

	cfg := LoadAppConfig()

	if cfg.AppEnv != "production" {
		t.Fatalf(
			"expected %q but got %q",
			"production",
			cfg.AppEnv,
		)
	}

	if cfg.MongoDBURI != "mongodb://prod:27017" {
		t.Fatalf(
			"expected %q but got %q",
			"mongodb://prod:27017",
			cfg.MongoDBURI,
		)
	}

	if cfg.AWSRegion != "us-west-2" {
		t.Fatalf(
			"expected %q but got %q",
			"us-west-2",
			cfg.AWSRegion,
		)
	}

	if cfg.Port != "9090" {
		t.Fatalf(
			"expected %q but got %q",
			"9090",
			cfg.Port,
		)
	}

	if cfg.SvcName != "weather-api" {
		t.Fatalf(
			"expected %q but got %q",
			"weather-api",
			cfg.SvcName,
		)
	}
}
