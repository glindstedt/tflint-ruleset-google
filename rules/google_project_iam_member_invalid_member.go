package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-google/project"
)

// GoogleProjectIamMemberInvalidMemberRule checks whether member value is invalid
type GoogleProjectIamMemberInvalidMemberRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewGoogleProjectIamMemberInvalidMemberRule returns new rule with default attributes
func NewGoogleProjectIamMemberInvalidMemberRule() *GoogleProjectIamMemberInvalidMemberRule {
	return &GoogleProjectIamMemberInvalidMemberRule{
		resourceType:  "google_project_iam_member",
		attributeName: "member",
	}
}

// Name returns the rule name
func (r *GoogleProjectIamMemberInvalidMemberRule) Name() string {
	return "google_project_iam_member_invalid_member"
}

// Enabled returns whether the rule is enabled by default
func (r *GoogleProjectIamMemberInvalidMemberRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *GoogleProjectIamMemberInvalidMemberRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *GoogleProjectIamMemberInvalidMemberRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether member format is invalid
func (r *GoogleProjectIamMemberInvalidMemberRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		var member string
		err := runner.EvaluateExpr(attribute.Expr, &member, nil)

		err = runner.EnsureNoError(err, func() error {
			if isValidIAMMemberFormat(member) {
				return nil
			}
			return runner.EmitIssue(
				r,
				fmt.Sprintf("%s is an invalid member format", member),
				attribute.Expr.Range(),
			)
		})
		if err != nil {
			return err
		}
	}

	return nil
}
