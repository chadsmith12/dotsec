package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultName = ".dotsecrc"

type ProjectConfig struct {
	Folder string `json:"folder"`
	Type string `json:"type"`
	Path string `json:"path"`
}

func defaultProjectConfig() ProjectConfig {
	return ProjectConfig{
		Folder: "",
		Type: "dotnet",
		Path: "",	
	}
}

func WriteProjectConfig() error {
	file, err := os.Create(defaultName)
	if err != nil {
		return fmt.Errorf("error creating .dotsecrc file: %w", err)
	}
	defer file.Close()
	config := defaultProjectConfig()	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating json for .dotsecrc file: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("error writing json to .dotsecrc file: %w", err)	
	}

	return nil
}

func LoadProjectConfig(cmd *cobra.Command, folder string) (*ProjectConfig, error) {
	config := &ProjectConfig{}
	fileConfig, err := loadFromFile()
	if err == nil {
		config = fileConfig
	}

	overrideFromFlags(cmd, config)
	
	if folder != "" {
		config.Folder = folder
	}

	if config.Folder == "" {
		return nil, fmt.Errorf("folder is required. Provide from argument or a .dotsecrc file")
	}

	return config, nil
}

func loadFromFile() (*ProjectConfig, error) {
	data, err := os.ReadFile(defaultName)
	if err != nil {
		return &ProjectConfig{}, err
	}
	
	projectConfig := &ProjectConfig{}
	if err := json.Unmarshal(data, projectConfig); err != nil {
		return &ProjectConfig{}, err
	}

	return projectConfig, nil
}

func overrideFromFlags(cmd *cobra.Command, config *ProjectConfig) {
	flags := cmd.Flags()
	if secretType, _ := flags.GetString("type"); secretType != "" {
		fmt.Println("setting secretType from the flag...")
		config.Type = secretType
	}

	if config.Type == "" {
		config.Type = "dotnet"
	}

	if config.Type == "dotnet" {
		if project, _ := flags.GetString("project"); project != "" {
			config.Path = project
		}
	} else if config.Type == "env" {
		if envFile, _ := flags.GetString("file"); envFile != "" {
			config.Path = envFile
		}
	}
}

func loadFromFlags(cmd *cobra.Command) ProjectConfig {
	flagConfig := ProjectConfig{}
	flags := cmd.Flags()
	if !flags.HasFlags() {
		return ProjectConfig{}
	}

	secretType, err := flags.GetString("type")
	if err != nil || secretType == "" {
		flagConfig.Type = ""
	}

	projectPath, _ := flags.GetString("project")
	envFile, _ := flags.GetString("file")
	
	if secretType == "env" {
		flagConfig.Path = envFile
	} else if secretType == "dotnet" {
		flagConfig.Path = projectPath
	}

	return flagConfig
}
