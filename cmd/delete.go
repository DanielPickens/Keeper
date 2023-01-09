package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [command]",
	Short: "Delete an object",
	Long: `Delete resources by namespace or names.

Deletion of a namespace will delete the namespace and remove all his attached object including the intentory attached to it. While removing an object will only supress it's metadata from the namespace but keep everything else.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDelete()
	},
}

func NewDeleteCommand() *cobra.Command {
	deleteCmd.AddCommand(NewDeleteJobCommand())
	deleteCmd.AddCommand(NewDeleteNamespaceCommand())

	return deleteCmd
}

func runDelete() {
	tpl := template.Must(template.New("deleteCmd").Parse(`
Using the delete command with a sub-command is helpful. Please use one of the following sub-command :
{{range . -}}
- {{.}}
{{end -}}
`))

	data := []string{"delete namespace", "delete job"}

	contents := bytes.Buffer{}
	if err := tpl.Execute(&contents, data); err != nil {
		logrus.Fatalf("error while executing template : %v", err)
	}

	fmt.Println(contents.String())
}
