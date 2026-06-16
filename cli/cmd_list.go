package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) listCmd() *cobra.Command {
	var conf string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List papers from a USENIX conference",
		RunE: func(cmd *cobra.Command, _ []string) error {
			limit := a.effectiveLimit(0)
			papers, err := a.client.List(cmd.Context(), conf, limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(papers, len(papers))
		},
	}
	cmd.Flags().StringVar(&conf, "conf", "usenixsecurity24", "conference ID (e.g. usenixsecurity24, nsdi24)")
	return cmd
}
