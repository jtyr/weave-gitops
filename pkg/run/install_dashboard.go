package run

import (
	"context"
	"fmt"

	"github.com/weaveworks/weave-gitops/pkg/kube"
	"github.com/weaveworks/weave-gitops/pkg/logger"
	"github.com/weaveworks/weave-gitops/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func InstallDashboard(log logger.Logger, ctx context.Context, kubeClient *kube.KubeHTTP, kubeConfigArgs *genericclioptions.ConfigFlags) error {
	password, err := utils.ReadPasswordFromStdin("Please enter your password to generate your secret: ")
	if err != nil {
		return fmt.Errorf("could not read password: %w", err)
	}

	secret, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Successf("Secret has been generated:")
	fmt.Println(string(secret))

	log.Actionf("Installing GitOps Dashboard ...")

	manifests := createManifests(string(secret))

	applyOutput, err := apply(log, ctx, kubeClient, kubeConfigArgs, manifests)
	if err != nil {
		log.Failuref("GitOps Dashboard install failed")
		return err
	}

	fmt.Println(applyOutput)

	return nil
}

func createManifests(secret string /*, helmRepository *sourcev1.HelmRepository, helmRelease *helmv2.HelmRelease, values string*/) []byte {
	contentString := fmt.Sprintf(`
---
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: HelmRepository
metadata:
  name: ww-gitops
  namespace: flux-system
spec:
  interval: 1m0s
  url: https://helm.gitops.weave.works
---
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: ww-gitops
  namespace: flux-system
spec:
  chart:
    spec:
      chart: weave-gitops
      reconcileStrategy: ChartVersion
      sourceRef:
        kind: HelmRepository
        name: ww-gitops
      version: 2.0.6
  interval: 1m0s
  values:
    adminUser:
      create: true
      passwordHash: %s
      username: admin
`, secret)

	fmt.Println(contentString)

	content := []byte(contentString)

	return content
}

// func createHelmRepository() *sourcev1.HelmRepository {
// 	helmRepository := &sourcev1.HelmRepository{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "ww-gitops",
// 			Namespace: "flux-system",
// 		},
// 		Spec: sourcev1.HelmRepositorySpec{
// 			URL: "https://helm.gitops.weave.works",
// 			Interval: metav1.Duration{
// 				Duration: time.Minute,
// 			},
// 		},
// 	}

// 	// helmRepository := &sourcev1.HelmRepository{
// 	// 	ObjectMeta: metav1.ObjectMeta{
// 	// 		Name:      name,
// 	// 		Namespace: "",
// 	// 		Labels:    map[string]string{},
// 	// 	},
// 	// 	Spec: sourcev1.HelmRepositorySpec{
// 	// 		URL: sourceHelmArgs.url,
// 	// 		Interval: metav1.Duration{
// 	// 			Duration: createArgs.interval,
// 	// 		},
// 	// 	},
// 	// }

// 	return helmRepository
// }

// func createHelmRelease(passwordHash string) *helmv2.HelmRelease {
// 	helmRelease := &helmv2.HelmRelease{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "ww-gitops",
// 			Namespace: "flux-system",
// 		},
// 		Spec: helmv2.HelmReleaseSpec{
// 			Interval: metav1.Duration{
// 				Duration: time.Minute,
// 			},
// 			Chart: helmv2.HelmChartTemplate{
// 				Spec: helmv2.HelmChartTemplateSpec{
// 					Chart:   "weave-gitops",
// 					Version: "2.0.6",
// 					SourceRef: helmv2.CrossNamespaceObjectReference{
// 						Kind: "HelmRepository",
// 						Name: "ww-gitops",
// 					},
// 					ReconcileStrategy: "ChartVersion",
// 				},
// 			},
// 			Suspend: false,
// 		},
// 	}

// 	// helmRelease := helmv2.HelmRelease{
// 	// 	ObjectMeta: metav1.ObjectMeta{
// 	// 		Name:      name,
// 	// 		Namespace: *kubeconfigArgs.Namespace,
// 	// 		Labels:    sourceLabels,
// 	// 	},
// 	// 	Spec: helmv2.HelmReleaseSpec{
// 	// 		ReleaseName: helmReleaseArgs.name,
// 	// 		DependsOn:   utils.MakeDependsOn(helmReleaseArgs.dependsOn),
// 	// 		Interval: metav1.Duration{
// 	// 			Duration: createArgs.interval,
// 	// 		},
// 	// 		TargetNamespace: helmReleaseArgs.targetNamespace,

// 	// 		Chart: helmv2.HelmChartTemplate{
// 	// 			Spec: helmv2.HelmChartTemplateSpec{
// 	// 				Chart:   helmReleaseArgs.chart,
// 	// 				Version: helmReleaseArgs.chartVersion,
// 	// 				SourceRef: helmv2.CrossNamespaceObjectReference{
// 	// 					Kind:      helmReleaseArgs.source.Kind,
// 	// 					Name:      helmReleaseArgs.source.Name,
// 	// 					Namespace: helmReleaseArgs.source.Namespace,
// 	// 				},
// 	// 				ReconcileStrategy: helmReleaseArgs.reconcileStrategy,
// 	// 			},
// 	// 		},
// 	// 		Suspend: false,
// 	// 	},
// 	// }

// 	return helmRelease
// }

// func createValues(secret string) string {
// 	return fmt.Sprintf(`
//   values:
//     adminUser:
//       create: true
//       passwordHash: "%s"
//       username: admin
// 	`, secret)
// }
