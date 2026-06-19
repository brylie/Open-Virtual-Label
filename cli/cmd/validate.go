package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/schema"
	"github.com/spf13/cobra"
)

var validateAll bool

var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate workspace records against their schemas",
	Long: `Validate a single JSON file or all files in the workspace.

Schema is inferred from the file's location within workspace/.`,
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&validateAll, "all", false, "validate every JSON record in the workspace")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}

	var files []string

	switch {
	case len(args) > 0:
		abs, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		files = append(files, abs)

	case validateAll:
		if err := filepath.WalkDir(wsPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(path, ".json") {
				files = append(files, path)
			}
			return nil
		}); err != nil {
			return err
		}

	default:
		return errors.New("provide a file path or use --all")
	}

	totalErrors := 0
	for _, f := range files {
		schemaName, err := schema.InferSchema(wsPath, f)
		if err != nil {
			output.Info("  skipping %s (%v)", relPath(wsPath, f), err)
			continue
		}
		errs, err := schema.ValidateFile(schemaName, f)
		if err != nil {
			output.Fail("%s\n    error reading file: %v", relPath(wsPath, f), err)
			totalErrors++
			continue
		}
		if len(errs) == 0 {
			output.Success("%s", relPath(wsPath, f))
		} else {
			output.Fail("%s", relPath(wsPath, f))
			for _, e := range errs {
				fmt.Println("   ", e)
			}
			totalErrors += len(errs)
		}
	}

	if totalErrors > 0 {
		fmt.Printf("\n%d error(s) found.\n", totalErrors)
		return &ExitError{Code: 3, Msg: fmt.Sprintf("%d validation error(s)", totalErrors)}
	}
	if len(files) > 0 {
		fmt.Printf("\nAll %d file(s) valid.\n", len(files))
	}
	return nil
}

func relPath(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}
	return rel
}
