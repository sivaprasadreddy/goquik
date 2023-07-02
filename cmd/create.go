package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sivaprasadreddy/goquik/internal/project"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Example: "goquik create myapp",
	Short:   "Create a new project",
	Long:    `Create a new Go project`,
	Run:     run,
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, args []string) {
	p, err := prompt(args)
	if errors.Is(err, ErrAbort) {
		fmt.Println("Aborting project generation")
		os.Exit(1)
	} else {
		handleError(err)
	}

	fmt.Printf("Generating a project with name: %s\n", p.ProjectName)

	err = os.Mkdir(p.ProjectName, os.ModePerm)
	handleError(err)

	err = project.Generate(p)
	handleError(err)
}

var ErrAbort = errors.New("abort project generation")

func NewProject() *project.Project {
	return &project.Project{}
}

func prompt(args []string) (*project.Project, error) {
	p := NewProject()
	if len(args) == 0 {
		err := survey.AskOne(&survey.Input{
			Message: "What is your project name?",
			Help:    "project name.",
			Suggest: nil,
		}, &p.ProjectName, survey.WithValidator(survey.Required))
		if err != nil {
			return nil, err
		}
	} else {
		p.ProjectName = args[0]
	}

	err := overrideProject(p.ProjectName)
	handleError(err)

	err = survey.AskOne(&survey.Input{
		Message: "What is your module path? (Ex: github.com/username/project)",
		Help:    "Go Module path",
	}, &p.ModulePath, survey.WithValidator(survey.Required))

	if err != nil {
		return nil, err
	}
	return p, err
}

func overrideProject(projectName string) error {
	stat, _ := os.Stat(projectName)
	if stat != nil {
		var overwrite = false

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Folder %s already exists, do you want to overwrite it?", projectName),
			Help:    "Remove old project and create new project.",
		}
		err := survey.AskOne(prompt, &overwrite)
		handleError(err)
		if !overwrite {
			fmt.Println("Aborting project generation")
			os.Exit(1)
		}
		err = os.RemoveAll(projectName)
		if err != nil {
			fmt.Println("Remove old project error: ", err)
			return err
		}
	}
	return nil
}

func handleError(err error) {
	if err != nil {
		_ = fmt.Errorf("project generation failed. error: %v", err)
		os.Exit(1)
	}
}
