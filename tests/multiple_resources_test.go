package tests

import (
	"context"
	"testing"
	"time"

	"github.com/accuknox/kubernetes-cel-validator/tests/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	celvalidator "github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg"
)

func TestMultipleResourcesValidation(t *testing.T) {
	feature := features.New("Check Multiple Resources Validation").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			deployment, err := util.GetDeployment("./scenarios/check-multiple-resources/deployment.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, deployment); err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(config.Client().Resources()).DeploymentConditionMatch(deployment, appsv1.DeploymentAvailable, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			pod, err := util.GetPod("./scenarios/check-multiple-resources/pod.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(config.Client().Resources()).PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			var resources []any
			resources = append(resources, deployment, pod)

			return context.WithValue(ctx, "resources", resources)
		}).
		Assess("Check The Container Image", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationPrecondition, err := util.GetKubernetesResourcePrecondition("./scenarios/check-multiple-resources/validation.yaml")
			if err != nil {
				t.Fatal(err)
			}
			result, selectionPreconditionValidationFailures, _, err := celvalidator.GetKubernetesResourcePreconditionResult(validationPrecondition, config.Client().RESTConfig())
			if err != nil {
				t.Fatal(err)
			}
			if result != true {
				t.Fatal("Validation Failed, result:", result)
			}
			if len(selectionPreconditionValidationFailures) != 2 {
				t.Fatal("Selection Precondition Validation Failures Should be 2, but got:", len(selectionPreconditionValidationFailures))
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			resources := ctx.Value("resources").([]any)
			deployment := resources[0].(*appsv1.Deployment)
			if err := config.Client().Resources().Delete(ctx, deployment); err != nil {
				t.Fatal(err)
			}
			pod := resources[1].(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, feature)
}
