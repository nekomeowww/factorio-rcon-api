package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/nekomeowww/xo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	input        string
	inputFormat  string
	output       string
	outputFormat string
)

func main() {
	root := &cobra.Command{
		Use: "openapiv2conv",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasPrefix(input, "/") {
				input = xo.RelativePathBasedOnPwdOf(input)
			}

			openapiV2DocContent, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("failed to read input file %s: %w", input, err)
			}

			var openapiV2Doc openapi2.T

			switch inputFormat {
			case "json":
				err = json.Unmarshal(openapiV2DocContent, &openapiV2Doc)
				if err != nil {
					return fmt.Errorf("failed to unmarshal input file %s: %w", input, err)
				}
			case "yaml":
				err = yaml.Unmarshal(openapiV2DocContent, &openapiV2Doc)
				if err != nil {
					return fmt.Errorf("failed to unmarshal input file %s: %w", input, err)
				}
			default:
				return fmt.Errorf("unsupported input format: %s", inputFormat)
			}

			openapiV3Doc, err := openapi2conv.ToV3(&openapiV2Doc)
			if err != nil {
				return fmt.Errorf("failed to convert openapi v2 to v3: %w", err)
			}

			openapiV3DocBuffer := new(bytes.Buffer)

			switch outputFormat {
			case "json":
				encoder := json.NewEncoder(openapiV3DocBuffer)
				encoder.SetIndent("", "  ")

				err = encoder.Encode(openapiV3Doc)
				if err != nil {
					return fmt.Errorf("failed to encode openapi v3 doc: %w", err)
				}
			case "yaml":
				encoder := yaml.NewEncoder(openapiV3DocBuffer)
				encoder.SetIndent(2)

				err = encoder.Encode(openapiV3Doc)
				if err != nil {
					return fmt.Errorf("failed to encode openapi v3 doc: %w", err)
				}
			default:
				return fmt.Errorf("unsupported output format: %s", outputFormat)
			}

			if !strings.HasPrefix(output, "/") {
				output = xo.RelativePathBasedOnPwdOf(output)
			}

			err = os.WriteFile(output, openapiV3DocBuffer.Bytes(), 0644) //nolint
			if err != nil {
				return fmt.Errorf("failed to write output file %s: %w", output, err)
			}

			return nil
		},
	}

	root.Flags().StringVarP(&input, "input", "i", "", "input file path")
	root.Flags().StringVar(&inputFormat, "input-format", "json", "input file format")
	root.Flags().StringVarP(&output, "output", "o", "", "output file path")
	root.Flags().StringVar(&outputFormat, "output-format", "yaml", "output file format")

	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
