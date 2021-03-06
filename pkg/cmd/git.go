package cmd

import (
	"github.com/puppetlabs/relay-sdk-go/pkg/task"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/spf13/cobra"
)

func NewGitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "git",
		Short:                 "Manage git data",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewGitCloneCommand())

	return cmd
}

func NewGitCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "clone",
		Short:                 "Clone git repository",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, _ := cmd.Flags().GetString("directory")
			revision, _ := cmd.Flags().GetString("revision")

			u, err := taskutil.MetadataSpecURL()
			if err != nil {
				return err
			}
			planOpts := taskutil.DefaultPlanOptions{SpecURL: u}
			task := task.NewTaskInterface(planOpts)
			return task.CloneRepository(revision, directory)
		},
	}

	cmd.Flags().StringP("revision", "r", "", "git revision")
	cmd.Flags().StringP("directory", "d", "", "git clone output directory")

	return cmd
}
