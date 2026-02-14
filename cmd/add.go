package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chadsmith12/dotsec/cmdcontext"
	"github.com/chadsmith12/dotsec/config"
	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Allows you to add a user to a team",
	Long:  `Adds you to add a new user to the team`,
	Run:   runAddCmd,
}

func init() {
	teamCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("manager", "m", false, "sets if this user should be added as a manager")
}

func runAddCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Invaild number args - must supply email\n")
		os.Exit(1)
	}

	email := args[0]
	flags := cmd.Flags()
	asManager, _ := flags.GetBool("manager")

	projectConfig, err := config.LoadProjectConfig(cmd, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}
	cmdContext, err := cmdcontext.NewCommandContext(cmd, projectConfig)
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
	group, err := client.GetGroup(projectConfig.Team)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to find group: %v\n", err)
		os.Exit(1)
	}
	user, err := client.GetUser(email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to find user: %v\n", err)
		os.Exit(1)
	}

	addOperation := passbolt.AddUserGroupOptions{
		User:    user,
		Group:   group,
		Manager: asManager,
	}
	err = client.AddUserToGroup(addOperation)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add %s to %s: %v\n", user.Username, group.Name, err)
		os.Exit(1)
	}

	fmt.Printf("Added %s to %s\n", user.Username, group.Name)
}
