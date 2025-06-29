package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	DatabaseURL    string
	SessionSecret  string
	PearlPoints    int
	TargetScore    int
	MaxGameTime    time.Duration
	MoveCooldown   time.Duration
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "boba_vim.db"),
		SessionSecret: getEnv("SESSION_SECRET", "your-secret-key-change-in-production"),
		PearlPoints:   getEnvInt("PEARL_POINTS", 100),
		TargetScore:   getEnvInt("TARGET_SCORE", 1000),
		MaxGameTime:   time.Duration(getEnvInt("MAX_GAME_TIME", 1800)) * time.Second, // 30 minutes
		MoveCooldown:  time.Duration(getEnvInt("MOVE_COOLDOWN", 100)) * time.Millisecond, // 100ms cooldown
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
