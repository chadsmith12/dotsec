package cmd

import (
	"fmt"
	"os"

	"github.com/chadsmith12/dotsec/config"
	"github.com/spf13/cobra"
)

// teamCmd represents the team command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Provides commands to work with teams in dotsec",
	Long: `Provides commands to work with teams in dotsec.
Team in dotsec are a way to manage your team that you are sharing secrets
with in your secrets manager.`,
	Run: runTeamsCmd,
}

func init() {
	rootCmd.AddCommand(teamCmd)

	teamCmd.Flags().BoolP("list", "l", true, "List all the members in the team")
}

func runTeamsCmd(cmd *cobra.Command, args []string) {
	flags := cmd.Flags()
	list, _ := flags.GetBool("list")
	projectConfig, err := config.LoadProjectConfig(cmd, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}
	if list {
		listMembers(projectConfig)
	} else {
		fmt.Println(cmd.Help())
	}
}

func listMembers(cmdConfig *config.ProjectConfig) {
	fmt.Printf("Team: %s\n", cmdConfig.Team)
}
