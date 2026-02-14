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
	Type   string `json:"type"`
	Path   string `json:"path"`
	Team   string `json:"team"`
}

func defaultProjectConfig() ProjectConfig {
	return ProjectConfig{
		Folder: "",
		Type:   "dotnet",
		Path:   "",
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

func WriteProjectConfigWithData(folder, secretType, path, team string) error {
	file, err := os.Create(defaultName)
	if err != nil {
		return fmt.Errorf("error creating .dotsecrc file: %w", err)
	}
	defer file.Close()

	config := ProjectConfig{
		Folder: folder,
		Type:   secretType,
		Path:   path,
		Team:   team,
	}

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

	if team, _ := flags.GetString("team"); team != "" {
		config.Team = team
	}

	if secretType, _ := flags.GetString("type"); secretType != "" {
		config.Type = secretType
	}

	if config.Type == "" {
		config.Type = "dotnet"
	}

	switch config.Type {
	case "dotnet":
		if project, _ := flags.GetString("project"); project != "" {
			config.Path = project
		}
	case "env":
		if envFile, _ := flags.GetString("file"); envFile != "" {
			config.Path = envFile
		}
	}
}
