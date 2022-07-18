package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fluxcd/flux2/pkg/manifestgen/install"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/weaveworks/weave-gitops/cmd/gitops/cmderrors"
	"github.com/weaveworks/weave-gitops/cmd/gitops/config"
	"github.com/weaveworks/weave-gitops/cmd/internal"
	"github.com/weaveworks/weave-gitops/pkg/kube"
	"github.com/weaveworks/weave-gitops/pkg/run"
	"github.com/weaveworks/weave-gitops/pkg/version"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type runCommandFlags struct {
	FluxVersion     string
	AllowK8sContext string
	Components      []string
	ComponentsExtra []string
	Timeout         time.Duration
	// Global flags.
	Namespace  string
	KubeConfig string
	// Flags, created by genericclioptions.
	Context string
}

var flags runCommandFlags

var kubeConfigArgs *genericclioptions.ConfigFlags

func RunCommand(opts *config.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Set up an interactive sync between your cluster and your local file system",
		Long:  "This will set up a sync between the cluster in your kubeconfig and the path that you specify on your local filesystem.  If you do not have Flux installed on the cluster then this will add it to the cluster automatically.  This is a requirement so we can sync the files successfully from your local system onto the cluster.  Flux will take care of producing the objects for you.",
		Example: `
# Run the sync on the current working directory
gitops beta run . [flags]

# Run the sync against the dev overlay path
gitops beta run ./deploy/overlays/dev [flags]`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		PreRunE:           betaRunCommandPreRunE(&opts.Endpoint),
		RunE:              betaRunCommandRunE(opts),
		DisableAutoGenTag: true,
	}

	cmdFlags := cmd.Flags()

	cmdFlags.StringVar(&flags.FluxVersion, "flux-version", version.FluxVersion, "The version of Flux to install")
	cmdFlags.StringVar(&flags.AllowK8sContext, "allow-k8s-context", "", "The name of the kubeconfig context to allow explicitly")
	cmdFlags.StringSliceVar(&flags.Components, "components", []string{"source-controller", "kustomize-controller", "helm-controller", "notification-controller"}, "The Flux components to install, as a comma-separated list: --components=component1,component2,component3")
	cmdFlags.StringSliceVar(&flags.ComponentsExtra, "components-extra", []string{}, "Additional Flux components to install, as a comma-separated list: --components=component1,component2,component3")
	cmdFlags.DurationVar(&flags.Timeout, "timeout", 5*time.Second, "The timeout for Flux installation")

	kubeConfigArgs = run.GetKubeConfigArgs()

	kubeConfigArgs.AddFlags(cmd.Flags())

	return cmd
}

func betaRunCommandPreRunE(endpoint *string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		numArgs := len(args)

		if numArgs == 0 {
			return cmderrors.ErrNoFilePath
		}

		if numArgs > 1 {
			return cmderrors.ErrMultipleFilePaths
		}

		return nil
	}
}

// isLocalCluster checks if it's a local cluster. See https://skaffold.dev/docs/environment/local-cluster/
func isLocalCluster(kubeClient *kube.KubeHTTP) bool {
	const (
		// cluster context prefixes
		kindPrefix = "kind-"
		k3dPrefix  = "k3d-"

		// cluster context names
		minikube         = "minikube"
		dockerForDesktop = "docker-for-desktop"
		dockerDesktop    = "docker-desktop"
	)

	if strings.HasPrefix(kubeClient.ClusterName, kindPrefix) ||
		strings.HasPrefix(kubeClient.ClusterName, k3dPrefix) ||
		kubeClient.ClusterName == minikube ||
		kubeClient.ClusterName == dockerForDesktop ||
		kubeClient.ClusterName == dockerDesktop {
		return true
	} else {
		return false
	}
}

func betaRunCommandRunE(opts *config.Options) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var err error

		if flags.Namespace, err = cmd.Flags().GetString("namespace"); err != nil {
			return err
		}

		kubeConfigArgs.Namespace = &flags.Namespace

		if flags.KubeConfig, err = cmd.Flags().GetString("kubeconfig"); err != nil {
			return err
		}

		if flags.Context, err = cmd.Flags().GetString("context"); err != nil {
			return err
		}

		log := internal.NewCLILogger(os.Stdout)

		if flags.KubeConfig != "" {
			kubeConfigArgs.KubeConfig = &flags.KubeConfig

			if flags.Context == "" {
				log.Failuref("A context should be provided if a kubeconfig is provided")
				return cmderrors.ErrNoContextForKubeConfig
			}
		}

		log.Actionf("Checking for a cluster in the kube config ...")

		var contextName string

		if flags.Context != "" {
			contextName = flags.Context
		} else {
			_, contextName, err = kube.RestConfig()
			if err != nil {
				log.Failuref("Error getting a restconfig: %v", err.Error())
				return cmderrors.ErrNoCluster
			}
		}

		cfg, err := kubeConfigArgs.ToRESTConfig()
		if err != nil {
			return fmt.Errorf("error getting a restconfig from kube config args: %w", err)
		}

		kubeClientOpts := run.GetKubeClientOptions()
		kubeClientOpts.BindFlags(cmd.Flags())

		kubeClient, err := run.GetKubeClient(log, contextName, cfg, kubeClientOpts)
		if err != nil {
			return cmderrors.ErrGetKubeClient
		}

		contextName = kubeClient.ClusterName
		if flags.AllowK8sContext == contextName {
			log.Actionf("Explicitly allow GitOps Run on %s context", contextName)
		} else if !isLocalCluster(kubeClient) {
			return errors.New("allowed to run against a local cluster only")
		}

		ctx := context.Background()

		log.Actionf("Checking if Flux is already installed ...")

		if fluxVersion, err := run.GetFluxVersion(log, ctx, kubeClient); err != nil {
			log.Warningf("Flux is not found: %v", err.Error())

			installOpts := install.Options{
				BaseURL:         install.MakeDefaultOptions().BaseURL,
				Version:         flags.FluxVersion,
				Namespace:       flags.Namespace,
				Components:      flags.Components,
				ComponentsExtra: flags.ComponentsExtra,
				ManifestFile:    "flux-system.yaml",
				Timeout:         flags.Timeout,
			}

			if err := run.InstallFlux(log, ctx, kubeClient, installOpts, kubeConfigArgs); err != nil {
				return fmt.Errorf("flux installation failed: %w", err)
			} else {
				log.Successf("Flux has been installed")
			}
		} else {
			log.Successf("Flux version %s is found", fluxVersion)
		}

		prompt := promptui.Prompt{
			Label:     "Would you like to install the GitOps Dashboard",
			IsConfirm: true,
		}
		_, err = prompt.Run()

		if err == nil {
			err = run.InstallDashboard(log, ctx, kubeClient, kubeConfigArgs)
			if err != nil {
				return fmt.Errorf("gitops dashboard installation failed: %w", err)
			}
		}

		const fluxSystemNS = "flux-system"
		for _, controllerName := range []string{"source-controller", "kustomize-controller", "helm-controller", "notification-controller"} {
			log.Actionf("Waiting for %s/%s to be ready ...", fluxSystemNS, controllerName)

			if err := run.WaitForDeploymentToBeReady(log, kubeClient, controllerName, fluxSystemNS); err != nil {
				return err
			}

			log.Successf("%s/%s is now ready ...", fluxSystemNS, controllerName)
		}

		if err := run.InstallBucketServer(log, kubeClient, cfg); err != nil {
			return err
		}

		return nil
	}
}
