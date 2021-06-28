package list

import (
	"fmt"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdList(f *cmdutils.Factory) *cobra.Command {
	var projectsListCmd = &cobra.Command{
		Use:   "list [flags]",
		Short: `Get the list of projects`,
		Example: heredoc.Doc(`
	$ glab repo list
	`),
		Long: ``,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var titleQualifier string

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			l := &gitlab.ListProjectsOptions{}
			l.Page = 1
			l.PerPage = 30

			if m, _ := cmd.Flags().GetString("orderBy"); m != "" {
				l.OrderBy = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetString("sort"); m != "" {
				l.Sort = gitlab.String(m)
			}
			if p, _ := cmd.Flags().GetInt("page"); p != 0 {
				l.Page = p
			}
			if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
				l.PerPage = p
			}
			repos, err := api.ListProjects(apiClient, l)
			//	repo, err := api.ListProjects(apiClient, repo.FullName(), l)
			if err != nil {
				return err
			}

			title := utils.NewListTitle(fmt.Sprintf("%s project", titleQualifier))
			title.RepoName = repo.FullName()
			title.Page = l.Page
			title.CurrentPageTotal = len(repos)

			fmt.Fprintf(f.IO.StdOut, "%s\n%s\n", title.Describe())

			return nil
		},
	}
	projectsListCmd.Flags().StringP("orderBy", "o", "", "Order pipeline by <string>")
	projectsListCmd.Flags().StringP("sort", "", "desc", "Sort pipeline by {asc|desc}. (Defaults to desc)")
	projectsListCmd.Flags().IntP("page", "p", 1, "Page number")
	projectsListCmd.Flags().IntP("per-page", "P", 30, "Number of items to list per page. (default 30)")

	return projectsListCmd
}
