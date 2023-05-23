package tests

import (
	"context"
	"testing"

	celvalidator "github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg"
	"github.com/accuknox/kubernetes-cel-validator/tests/util"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestStarClusterRole(t *testing.T) {
	featureTrueValidation := features.New("Check ClusterRole Validation").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			clusterRole, err := util.GetClusterRole("./scenarios/check-star-cluster-role/cluster-role.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, clusterRole); err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, "cluster-role", clusterRole)
		}).
		Assess("Check Contains * Resource", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationPrecondition, err := util.GetKubernetesResourcePrecondition("./scenarios/check-star-cluster-role/validation.yaml")
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
			if len(selectionPreconditionValidationFailures) == 0 {
				t.Fatal("Selection Precondition Validation Failures Should be greater than 0, but got:", len(selectionPreconditionValidationFailures))
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			clusterRole := ctx.Value("cluster-role").(*rbacv1.ClusterRole)
			if err := config.Client().Resources().Delete(ctx, clusterRole); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	featureTrueValidationExistsMacro := features.New("Check ClusterRole Validation").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			clusterRole, err := util.GetClusterRole("./scenarios/check-star-cluster-role/cluster-role.yaml")
			if err != nil {
				t.Fatal(err)
			}
			// Change the resource to pods, the test will still pass as exists macro is used.
			clusterRole.Rules[0].Resources = []string{"pods", "pods/*"}
			if err = config.Client().Resources().Create(ctx, clusterRole); err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, "cluster-role", clusterRole)
		}).
		Assess("Check Contains * Resource With One Resource Changed", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationPrecondition, err := util.GetKubernetesResourcePrecondition("./scenarios/check-star-cluster-role/validation.yaml")
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
			if len(selectionPreconditionValidationFailures) == 0 {
				t.Fatal("Selection Precondition Validation Failures Should be greater than 0, but got:", len(selectionPreconditionValidationFailures))
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			clusterRole := ctx.Value("cluster-role").(*rbacv1.ClusterRole)
			if err := config.Client().Resources().Delete(ctx, clusterRole); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	featureFalseValidationExistsMacro := features.New("Check ClusterRole Validation").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			clusterRole, err := util.GetClusterRole("./scenarios/check-star-cluster-role/cluster-role.yaml")
			if err != nil {
				t.Fatal(err)
			}
			// Change the resource to pods for validation to fail
			clusterRole.Rules[0].Resources = []string{"pods", "pods/*"}
			clusterRole.Rules[1].Resources = []string{"deployments", "deployment/*"}
			if err = config.Client().Resources().Create(ctx, clusterRole); err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, "cluster-role", clusterRole)
		}).
		Assess("Check Does Not Contain * Resource", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationPrecondition, err := util.GetKubernetesResourcePrecondition("./scenarios/check-star-cluster-role/validation.yaml")
			if err != nil {
				t.Fatal(err)
			}
			result, selectionPreconditionValidationFailures, validationFailure, err := celvalidator.GetKubernetesResourcePreconditionResult(validationPrecondition, config.Client().RESTConfig())
			if err != nil {
				t.Fatal(err)
			}
			if result != false {
				t.Fatal("Validation Failed, result:", result)
			}
			if len(selectionPreconditionValidationFailures) == 0 {
				t.Fatal("Selection Precondition Validation Failures Should be greater than 0, but got:", len(selectionPreconditionValidationFailures))
			}

			validationFailureExpectedMessage := "superuser does not have access to all resources"
			if validationFailure.Message != validationFailureExpectedMessage {
				t.Fatal("Validation Failure Message Should be:", validationFailureExpectedMessage, "but got:", validationFailure.Message)
			}

			validationFailureExpectedKey := "star-cluster-role"
			if validationFailure.Key != validationFailureExpectedKey {
				t.Fatal("Validation Failure Key Should be:", validationFailureExpectedKey, "but got:", validationFailure.Key)
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			clusterRole := ctx.Value("cluster-role").(*rbacv1.ClusterRole)
			if err := config.Client().Resources().Delete(ctx, clusterRole); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()
	testEnv.Test(t, featureTrueValidation, featureTrueValidationExistsMacro, featureFalseValidationExistsMacro)
}
