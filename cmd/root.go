package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile           string
	playbookDir       string
	kubectlConfigPath string
	v                 string
	namespace         string
	cors              bool
	wait              bool
	timeout           time.Duration
	port              int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "keeper",
	Short: "Keeper is a tool that let you create and manage multiple version of the same stack using Kubernetes and namespace",
	Long: `Keeper allows you apply a bunch of configuration file templates and combine them into different namespaces using some provided values.

Keeper is made to be executed using a directory containing configuration files and directories called a Playbook.

Using Keeper and a Playbook, you can easily create a namespace by using the "create" command.
This command will generate an inventory file containing the default configuration for the namespace you are creating.

You can update your specific inventory file manually.

Then Keeper configures your namespace using a auto-generated Kubernetes config using the specified inventory file.
This action can be done using the "apply" command.
	`,
}

func NewKeeperCommand() *cobra.Command {

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, v); err != nil {
			return err
		}
		return nil
	}

	rootCmd.AddCommand(NewServeCommand())
	rootCmd.AddCommand(NewApplyCommand())
	rootCmd.AddCommand(NewCreateCommand())
	rootCmd.AddCommand(NewDeleteCommand())
	rootCmd.AddCommand(NewGetCommand())
	rootCmd.AddCommand(NewResetCommand())
	rootCmd.AddCommand(NewVersionCommand())

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.keeper.yaml)")
	rootCmd.PersistentFlags().StringVar(&playbookDir, "dir", "", "Use the specified directory as root path to execute commands. Default is the current directory.")
	rootCmd.PersistentFlags().StringVar(&kubectlConfigPath, "kube-config-path", kubernetes.KubeConfigDefaultPath(), "kubectl config file")
	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	viper.BindPFlag("working-dir", rootCmd.PersistentFlags().Lookup("dir"))

	initConfig()

	return rootCmd

}

func addCommonNamespaceCommandFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace where to apply configuration")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".keeper") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	//Define current working dir as default value
	currentDir, err := os.Getwd()
	if err != nil {
		logrus.Fatal("Error when getting the working dir : ", err)
	}
	viper.SetDefault("working-dir", currentDir)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())

	}

	playbookDir = viper.GetString("working-dir")

}

func askForConfirmation(message string, reader io.Reader) bool {

	r := bufio.NewReader(reader)

	for {
		fmt.Printf("%s [y/n]: ", message)

		response, err := r.ReadString('\n')
		if err != nil {
			logrus.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else {
			return false
		}
	}
}

func newKubernetesClient() *kubernetes.Client {
	kube, err := kubernetes.NewClient(kubectlConfigPath)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	return kube
}

func newFileClient(dir string) *files.Client {
	f, err := files.NewClient(dir)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	return f

}

func newAPI(files *files.Client, kube *kubernetes.Client) api.Api {
	return api.NewApi(
		files.Inventories(),
		files.Configs(),
		files.Playbooks(),
		kube.Namespaces(),
		kube.Pods(),
		kube.Deployments(),
	)
}

func setUpLogs(out io.Writer, level string) error {

	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return nil
}
