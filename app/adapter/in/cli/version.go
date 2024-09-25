package cli

import (
	"fmt"

	"github.com/Ignaciojeria/einar/app/shared/archetype/cmd"

	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "retrieve einar cli version",
	Run:   runversion,
}

func runversion(cmd *cobra.Command, args []string) {
	fmt.Println("1.45.0")
}
