package main

import (
	"testing"

	"github.com/hariom-pal/go-epp/epp"
)

func TestValidateUATSafetyBlocksTransformWithoutSafeEnvironment(t *testing.T) {
	plan := operationPlan{names: []string{"domain-create"}, transform: true, objects: []string{"example.in"}}

	_, err := validateUATSafety(&epp.Config{Environment: "production"}, cliOptions{ConfirmTransform: true}, plan)
	if err == nil {
		t.Fatal("expected production transform to be blocked")
	}
}

func TestValidateUATSafetyRequiresTransformConfirmation(t *testing.T) {
	plan := operationPlan{names: []string{"domain-create"}, transform: true, objects: []string{"example.in"}}

	_, err := validateUATSafety(&epp.Config{Environment: "ote"}, cliOptions{}, plan)
	if err == nil {
		t.Fatal("expected transform without confirmation to be blocked")
	}
}

func TestValidateUATSafetyAllowsConfirmedOTETransform(t *testing.T) {
	plan := operationPlan{names: []string{"domain-create"}, transform: true, objects: []string{"example.in"}}

	environment, err := validateUATSafety(&epp.Config{Environment: "ote"}, cliOptions{ConfirmTransform: true}, plan)
	if err != nil {
		t.Fatalf("expected confirmed OTE transform to pass: %v", err)
	}
	if environment != "ote" {
		t.Fatalf("environment = %q, want ote", environment)
	}
}

func TestValidateUATSafetyAllowsOneOperationInUATMode(t *testing.T) {
	plan := operationPlan{names: []string{"domain-check", "contact-check"}}

	_, err := validateUATSafety(&epp.Config{Environment: "ote"}, cliOptions{UAT: true}, plan)
	if err == nil {
		t.Fatal("expected multiple UAT operations to be blocked")
	}
}

func TestValidateUATSafetyRequiresSafeEnvironmentInUATMode(t *testing.T) {
	plan := operationPlan{names: []string{"connect"}}

	_, err := validateUATSafety(&epp.Config{Environment: "production"}, cliOptions{UAT: true}, plan)
	if err == nil {
		t.Fatal("expected UAT mode with production environment to be blocked")
	}
}
