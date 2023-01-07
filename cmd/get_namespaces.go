package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getNamespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "Show information about kubernetes namespaces.",
	Long: `Show information about Kubernetes namespaces such as names, status (percentage of pods in a running status),
managed or not with the current playbook, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runGetNamespaces()
		if err != nil {
			logrus.Fatal(err.Error())
		}

	},
}

func NewGetNamespacesCommand() *cobra.Command {
	return getNamespacesCmd
}

func runGetNamespaces() error {

	api := newAPI(newFileClient(playbookDir), newKubernetesClient())

	namespaces, err := api.ListNamespaces()
	if err != nil {
		return errors.New(fmt.Sprintf("an error occured when getting information about namespaces : %v", err))
	}

	x := new(tabwriter.Writer)
	x.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(x, "Namespace\tPhase\tStatus\tManaged\t")
	for _, namespace := range namespaces {
		fmt.Fprint(x, fmt.Sprintf("%s\t%s\t%d%%\t%t\t\n", namespace.Name, namespace.Phase, namespace.Status, namespace.Managed))
	}
	fmt.Fprintln(x)
	x.Flush()

	return nil

}
