package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/usenix-cli/usenix"
)

func (a *App) conferencesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "conferences",
		Short: "List known USENIX conferences",
		RunE: func(cmd *cobra.Command, _ []string) error {
			confs := usenix.Conferences()
			return a.render(confs)
		},
	}
}
