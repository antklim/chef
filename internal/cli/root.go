package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// TODO: add crush report capability

// Execute is the primary entrypoint of the CLI app.
func Execute() {
	rootCmd := &cobra.Command{
		Use:           "chef",
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "Supercharge your development workflow.",
		Long: "Supercharge your development workflow.\n" +
			"Bootstrap a new project using predefined categories or bring your own layout.\n" +
			"Add new components to an existing project.\n",
		Version: "v0.1.0", // TODO (feat): add build info and version
	}

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(componentsCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
