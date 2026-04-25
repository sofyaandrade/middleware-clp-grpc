package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Env struct {
	AccessToken  string `mapstructure:"ACCESS_TOKEN"`
	RefreshToken string `mapstructure:"REFRESH_TOKEN"`
	ApiPort      string `mapstructure:"API_PORT"`
}

const (
	ENV = ".env"
)

func CarregarValoresEnv() (*Env, error) {
	env, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	return env, nil
}

func GetEnvFilePath() (string, error) {
	workingDir, _ := os.Getwd()
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	envPaths := []string{
		filepath.Join(workingDir, ENV),
		filepath.Join(execDir, ENV),
	}

	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("não foi possível localizar o arquivo %s", ENV)
}

func LoadEnv() (*Env, error) {
	envFilePath, err := GetEnvFilePath()
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(envFilePath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Não foi possível localizar o arquivo %s: %s", ENV, err)
		return nil, err
	}

	var env Env
	if err := viper.Unmarshal(&env); err != nil {
		fmt.Println("Não foi possível carregar o ambiente:", err)
		return nil, err
	}

	return &env, nil
}
