package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/chadsmith12/dotsec/cmdcontext"
	"github.com/chadsmith12/dotsec/colors"
	"github.com/chadsmith12/dotsec/config"
	"github.com/chadsmith12/dotsec/input"
	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/passbolt/go-passbolt/helper"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates a project to a new version of dotsec",
	Long: `Migrates a project that was using the old structure of dotsec to a newer one.
	This will move your folder to the "dotsec" folder in Passbolt and will also share
	all secrets with the group set to be used`,
	Run: runMigrateCmd,
}

type migrationProcess struct {
	group       string
	groupId     string
	createGroup bool
	usersToMove []helper.GroupMembership
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrateCmd(cmd *cobra.Command, args []string) {
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
	ctx := context.Background()
	client, err := cmdContext.UserClient(ctx)
	_ = client
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create passbolt client: %v\n", err)
		os.Exit(1)
	}

	if projectConfig.Team == "" {
		projectConfig.Team, _ = input.PromptUser(colors.Yellow("No team configured!\n What team/group would you like to use in Passbolt? "), false)
	}

	migration := migrationProcess{createGroup: false}
	foundGroup, err := client.GetGroup(projectConfig.Team)
	if err != nil {
		if _, ok := err.(*passbolt.GroupNotFoundErr); !ok {
			fmt.Fprintf(os.Stderr, "failed to get group: %v\n", err)
			os.Exit(1)
		}
		confirm, _ := input.PromptYesOrNo(colors.Yellow("No group found. Would you like to create it? [Y/n]: "))
		if !confirm {
			fmt.Printf("Must provide an existing group or create a new group. Exiting...")
			return
		}

		migration.createGroup = true
		migration.group = projectConfig.Team
	} else {
		migration.groupId = foundGroup.ID
	}

	folder, err := client.GetFolderWithResources(projectConfig.Folder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get folder: %v\n", err)
		os.Exit(1)
	}

	userPermissions, err := client.GetUsersFromPermissions(folder.Permissions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get users: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Migrate %s to dotsec/%s\n", projectConfig.Folder, projectConfig.Folder)
	if migration.createGroup {
		fmt.Printf("Create new group: %s\n", projectConfig.Team)
	} else {
		fmt.Printf("Update group %s\n", projectConfig.Team)
	}
	fmt.Printf("Share %d resources with %s (replacing individual permissions)\n", len(folder.ChildrenResources), projectConfig.Team)
	fmt.Printf("  With Members: \n")
	for i, user := range userPermissions {
		permissionType := "Member"
		if user.Type == passbolt.OwnerPermission {
			permissionType = "Group Owner"
		}
		fmt.Printf("  %d: %s %s (%s)\n", i+1, user.User.Profile.FirstName, user.User.Profile.LastName, permissionType)
	}

	confirm, _ := input.PromptYesOrNo(colors.Yellow("Continue with migration? [Y/n]: "))
	if !confirm {
		return
	}

	if migration.createGroup {
		id, err := client.CreateGroup(projectConfig.Team, userPermissions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get create group: %v\n", err)
			os.Exit(1)
		}
		migration.groupId = id
	} else {
		results := client.UpdateGroup(migration.groupId, projectConfig.Team, userPermissions)
		if len(results.Errors) > 0 {
			joinedErr := errors.Join(results.Errors...)
			fmt.Printf("%d errors occured updating the group: %s", len(results.Errors), joinedErr)
		}
	}

	dotsecFolder, err := client.GetFolderWithResources("dotsec")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get parent folder to move to: %v\n", err)
		os.Exit(1)
	}
	oldFolder, err := client.GetFolderWithResources(projectConfig.Folder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get folder: %v\n", err)
		os.Exit(1)
	}

	err = client.MoveFolder(oldFolder.ID, dotsecFolder.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to move %s to dotsec/%s. Try moving it in Passbolt. %v\n", projectConfig.Folder, projectConfig.Folder, err)
		os.Exit(1)
	}

	err = client.ShareFolderWithGroup(oldFolder.ID, migration.groupId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to share folder: %v\n", err)
		os.Exit(1)
	}

	resourceIds := make([]string, len(oldFolder.ChildrenResources))
	for i, resource := range oldFolder.ChildrenResources {
		resourceIds[i] = resource.ID
	}

	shareErrors := client.ShareResourcesWithGroup(resourceIds, migration.groupId)
	if len(shareErrors) > 0 {
		joinedErr := errors.Join(shareErrors...)
		fmt.Fprintf(os.Stderr, "failed to share some resources: %v\n", joinedErr)
		os.Exit(1)
	}
}
