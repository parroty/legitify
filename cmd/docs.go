package cmd

import (
	"os"
	"sort"
	"strings"

	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/opa"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/open-policy-agent/opa/ast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(newDocsCommand())
}

const (
	cmdDocs           = "generate-docs"
	argDocsOutputFile = "output-file"
)

func newDocsCommand() *cobra.Command {
	docsCmd := &cobra.Command{
		Use:   cmdDocs,
		Short: `Generate policies documentation (as a yaml)`,
		RunE:  executeDocsCommand,
	}
	flags := docsCmd.Flags()
	flags.StringP(ArgOutputFile, "o", "", "output file, defaults to stdout")

	return docsCmd
}

func executeDocsCommand(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	outputFile, err := flags.GetString(argDocsOutputFile)
	if err != nil {
		return err
	} else if outputFile != "" {
		if err = setOutputFile(outputFile); err != nil {
			return err
		}
	}

	// loading only built-in policies
	engine, err := opa.Load([]string{})
	if err != nil {
		return err
	}

	result := generateDocs(&engine)
	data, err := yaml.Marshal(result)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		return err
	}

	return nil
}

type PolicyDoc struct {
	PolicyName  string `yaml:"policy_name"`
	Title       string
	Description string
	Severity    string
	Remediation []string
	Threat      []string
}

func newPolicyDoc(policy *ast.Rule, ref *ast.AnnotationsRef) PolicyDoc {
	return PolicyDoc{
		PolicyName:  policy.Head.Name.String(),
		Title:       ref.Annotations.Title,
		Description: ref.Annotations.Description,
		Severity:    ref.Annotations.Custom["severity"].(string),
		Remediation: resolveStringArray(ref.Annotations.Custom["remediationSteps"]),
		Threat:      resolveStringArray(ref.Annotations.Custom["threat"]),
	}
}

func resolveStringArray(customField interface{}) []string {
	retval := make([]string, 0)
	switch t := customField.(type) {
	case []interface{}:
		for _, enricherString := range t {
			switch ts := enricherString.(type) {
			case string:
				retval = append(retval, ts)
			}
		}
	}
	return retval
}

func generateDocs(engine *opa_engine.Enginer) map[string][]PolicyDoc {
	result := make(map[string][]PolicyDoc)
	annotations := (*engine).Annotations().Flatten()

	for _, a := range annotations {
		policy := a.GetRule()
		module := policy.Module.Package.Path.String()
		module = strings.Replace(module, "data.", "", 1)

		if _, ok := result[module]; !ok {
			result[module] = []PolicyDoc{}
		}
		val := result[module]
		result[module] = append(val, newPolicyDoc(policy, a))
	}

	for _, v := range result {
		sort.Slice(v, func(i, j int) bool {
			return severity.Less(v[i].Severity, v[j].Severity)
		})
	}

	return result
}
