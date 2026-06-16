package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) searchCmd() *cobra.Command {
	var conf string
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search papers by title (client-side filter)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			if query == "" {
				return codeError(exitUsage, fmt.Errorf("query cannot be empty"))
			}
			limit := a.effectiveLimit(0)
			papers, err := a.client.Search(cmd.Context(), query, conf, limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(papers, len(papers))
		},
	}
	cmd.Flags().StringVar(&conf, "conf", "usenixsecurity24", "conference ID (e.g. usenixsecurity24, nsdi24)")
	return cmd
}
