package env

import (
	"os"
	"path"
)

func SetSecrets(project, env, key, value string) error {
	folder := ""
	if project != "" {
		folder = project
	}
	file, err := os.Open(folder)
	
	filePath := path.Join(folder, env)
}
