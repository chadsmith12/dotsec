package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chadsmith12/dotsec/cmdcontext"
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
		listMembers(cmd, projectConfig)
	} else {
		fmt.Println(cmd.Help())
	}
}

func listMembers(cmd *cobra.Command, cmdConfig *config.ProjectConfig) {
	fmt.Printf("Listing members for %s: \n", cmdConfig.Team)
	cmdContext, err := cmdcontext.NewCommandContext(cmd, cmdConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create command context: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := cmdContext.UserClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create passbolt client: %v\n", err)
		os.Exit(1)
	}
	members, err := client.GetGroupMembers(cmdConfig.Team)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get group members for %s: %v\n", cmdConfig.Team, err)
		os.Exit(1)
	}
	for i, member := range members {
		fmt.Printf("%d: %s %s\n", i+1, member.UserFirstName, member.UserLastName)
	}
}
