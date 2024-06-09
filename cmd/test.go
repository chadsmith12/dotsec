package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the credentials that dotsec is about to use",
	Long: `This is a quick command that will read and set the credentials that dotsec is about to use,`,
	Run: runTestCmd,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func runTestCmd(cmd *cobra.Command, args []string) {
	server := viper.GetString("server")
	password := viper.GetString("password")
	privateKey := viper.GetString("privateKey")

	fmt.Printf("Using following data: \n Server: %v\n Password: %v\n PrivateKey: %v\n", server, password, privateKey)
}
