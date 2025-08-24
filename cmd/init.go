package cmd

import (
	"log"

	"github.com/chadsmith12/dotsec/config"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a .dotsecrc file",
	Long: `Initializes a new .dotsecrc file that can be used for project configuration to allow using
	dotsec without additional arguments or flags.`,

	Run: initRun,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initRun(cmd *cobra.Command, args []string) {
	err := config.WriteProjectConfig()
	if err != nil {
		log.Fatalf("%v", err)
	}
	// server, err := input.PromptUser("Server (https://passbolt.example.com): ", false)
	// if err != nil {
	// 	log.Fatalf("Error getting server: %v", err)
	// }
	//
	// privateKey, err := input.PromptUser("Path to Private Key: ", false)
	// if err != nil {
	// 	log.Fatalf("Error getting path to the private key: %v", err)
	// }
	//
	// password, err := input.PromptUser("Master Password (leave blank to ask each time): ", true)
	// if err != nil {
	// 	log.Fatalf("Error getting users master password: %v", err)
	// }
	//
	// fmt.Println("")
	// viper.Set("server", server)
	// viper.Set("privateKey", privateKey)
	// viper.Set("password", password)
	//
	// saveConfigFile()
}

