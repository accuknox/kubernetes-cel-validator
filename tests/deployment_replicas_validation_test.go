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

func TestReplicaCountValidation(t *testing.T) {
	featureTrueValidation := features.New("Check Deployment Validation").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			deployment, err := util.GetDeployment("./scenarios/check-deployment-replicas/deployment.yaml")
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
			return context.WithValue(ctx, "nginx-deployment", deployment)
		}).
		Assess("Check Number of Replicas", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationPrecondition, err := util.GetKubernetesResourcePrecondition("./scenarios/check-deployment-replicas/validation.yaml")
			if err != nil {
				t.Fatal(err)
			}
			result, selectionPreconditionValidationFailures, _, err :=
				celvalidator.GetKubernetesResourcePreconditionResult(validationPrecondition, config.Client().RESTConfig())
			if err != nil {
				t.Fatal(err)
			}
			if result != true {
				t.Fatal("Validation Failed, result:", result)
			}
			if len(selectionPreconditionValidationFailures) != 0 {
				t.Fatal("Selection Precondition Validation Failures Should be 0, but got:", len(selectionPreconditionValidationFailures))
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			deployment := ctx.Value("nginx-deployment").(*appsv1.Deployment)
			if err := config.Client().Resources().Delete(ctx, deployment); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	featureFalseValidation := features.New("Check Deployment Validation With 2 Replicas").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			deployment, err := util.GetDeployment("./scenarios/check-deployment-replicas/deployment.yaml")
			if err != nil {
				t.Fatal(err)
			}

			// Change the replicas of deployment to 2
			replicas := int32(2)
			deployment.Spec.Replicas = &replicas

			client, err := config.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			if err = client.Resources().Create(ctx, deployment); err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(client.Resources()).DeploymentConditionMatch(deployment, appsv1.DeploymentAvailable, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, "nginx-deployment", deployment)
		}).
		Assess("Check Number of Replicas", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			client, err := config.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			validationPrecondition, err := util.GetKubernetesResourcePrecondition("./scenarios/check-deployment-replicas/validation.yaml")
			if err != nil {
				t.Fatal(err)
			}
			result, selectionPreconditionValidationFailures, validationFailure, err :=
				celvalidator.GetKubernetesResourcePreconditionResult(validationPrecondition, client.RESTConfig())
			if err != nil {
				t.Fatal(err)
			}
			if result != false {
				t.Fatal("Validation Failed, result:", result)
			}
			if len(selectionPreconditionValidationFailures) != 0 {
				t.Fatal("Selection Precondition Validation Failures Should be 0, but got:", len(selectionPreconditionValidationFailures))
			}

			validationFailureExpectedMessage := "resource: nginx does not have 1 replica, it has 2"
			if validationFailure.Message != validationFailureExpectedMessage {
				t.Fatal("Validation Failure Message Should be:", validationFailureExpectedMessage, "but got:", validationFailure.Message)
			}

			validationFailureExpectedKey := "single-replica"
			if validationFailure.Key != validationFailureExpectedKey {
				t.Fatal("Validation Failure Key Should be:", validationFailureExpectedKey, "but got:", validationFailure.Key)
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			deployment := ctx.Value("nginx-deployment").(*appsv1.Deployment)
			if err := config.Client().Resources().Delete(ctx, deployment); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()
	testEnv.Test(t, featureTrueValidation, featureFalseValidation)
}
