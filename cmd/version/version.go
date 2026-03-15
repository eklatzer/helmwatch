package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v0.2.0"

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(Version)
		},
	}
}
