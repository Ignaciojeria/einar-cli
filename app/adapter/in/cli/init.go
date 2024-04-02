package cli

import (
	"fmt"

	"github.com/Ignaciojeria/einar/app/business"
	"github.com/Ignaciojeria/einar/app/domain"
	"github.com/Ignaciojeria/einar/app/shared/archetype/cmd"
	"github.com/Ignaciojeria/einar/app/shared/utils"

	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(initCmd)
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project name] [repository template]",
	Short: "Initialize a new Go module",
	Run:   runInitCmd,
}

func runInitCmd(cmd *cobra.Command, args []string) {
	_, err := utils.ReadEinarCli()

	if err == nil {
		fmt.Println("einar cli already initialized")
		return
	}

	var repositoryURL string
	var userCredentials string
	var invalidArgsQuantity bool = true
	if len(args) == 1 {
		repositoryURL = "https://github.com/Ignaciojeria/einar-cli-standard-template"
		userCredentials = "no-auth"
		invalidArgsQuantity = false
	}

	if len(args) == 3 {
		repositoryURL = args[1]
		userCredentials = args[2]
		invalidArgsQuantity = false
	}

	if invalidArgsQuantity {
		fmt.Println("accept 1 or 3 args only")
		return
	}

	templatePath, err := utils.GitCloneTemplateInBinaryPath(repositoryURL, userCredentials, "")
	if err != nil {
		fmt.Println("error getting template path")
		return
	}
	tag, err := utils.GetLatestTag(templatePath)
	if err != nil {
		fmt.Println("error getting tag from templateURL")
		return
	}

	project := args[0]
	if args[0] == "." {
		project, _ = utils.GetCurrentFolderName()
	}
	project = utils.ConvertStringCase(project, "kebab")
	business.EinarInit(cmd.Context(), templatePath, project)

	err = utils.CreateEinarCLIJSON(domain.EinarCli{
		Project: args[0],
		Template: domain.Template{
			URL: repositoryURL,
			Tag: tag,
		},
	})

	if err != nil {
		fmt.Println("error creating einar cli file")
		return
	}
}
