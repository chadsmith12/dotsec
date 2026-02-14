package cmd

import (
	"fmt"

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
	if list {
		listMembers()
	} else {
		fmt.Println(cmd.Help())
	}
}

func listMembers() {

}
