package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "View and update the label profile",
}

var labelShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the current label profile",
	RunE:  runLabelShow,
}

var labelSetNameCmd = &cobra.Command{
	Use:   "set-name <name>",
	Short: "Rename the label",
	Args:  cobra.ExactArgs(1),
	RunE:  runLabelSetName,
}

var labelSetDescriptionCmd = &cobra.Command{
	Use:   "set-description <text>",
	Short: "Set the label description",
	Args:  cobra.ExactArgs(1),
	RunE:  runLabelSetDescription,
}

var labelSetLicenseCmd = &cobra.Command{
	Use:   "set-license <license>",
	Short: "Set the label-wide default license",
	Args:  cobra.ExactArgs(1),
	RunE:  runLabelSetLicense,
}

var labelSetDistributorCmd = &cobra.Command{
	Use:   "set-distributor <distributor>",
	Short: "Set the label-wide default distributor",
	Args:  cobra.ExactArgs(1),
	RunE:  runLabelSetDistributor,
}

var (
	labelSetContactEmail    string
	labelSetContactWebsite  string
	labelSetContactLocation string
)

var labelSetContactCmd = &cobra.Command{
	Use:   "set-contact",
	Short: "Update label contact details",
	RunE:  runLabelSetContact,
}


func init() {
	labelSetContactCmd.Flags().StringVar(&labelSetContactEmail, "email", "", "public contact email for the label")
	labelSetContactCmd.Flags().StringVar(&labelSetContactWebsite, "website", "", "label website URL")
	labelSetContactCmd.Flags().StringVar(&labelSetContactLocation, "location", "", "city and country (e.g. Helsinki, Finland)")

	labelCmd.AddCommand(labelShowCmd)
	labelCmd.AddCommand(labelSetNameCmd)
	labelCmd.AddCommand(labelSetDescriptionCmd)
	labelCmd.AddCommand(labelSetLicenseCmd)
	labelCmd.AddCommand(labelSetDistributorCmd)
	labelCmd.AddCommand(labelSetContactCmd)
	rootCmd.AddCommand(labelCmd)
}

func readLabel(wsPath string) (models.Label, error) {
	var l models.Label
	if err := ws.ReadJSON(ws.LabelFile(wsPath), &l); err != nil {
		return l, errors.New("WORKSPACE_NOT_FOUND: label profile not found; run 'ovl init' first")
	}
	return l, nil
}

func writeLabel(wsPath string, l *models.Label) error {
	raw, err := json.Marshal(l)
	if err != nil {
		return err
	}
	errs, err := schema.Validate("label", raw)
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: "VALIDATION_FAILED: " + strings.Join(errs, "; ")}
	}
	return ws.WriteJSON(ws.LabelFile(wsPath), l)
}

func runLabelShow(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	l, err := readLabel(wsPath)
	if err != nil {
		return &ExitError{Code: 2, Msg: err.Error()}
	}

	if output.IsJSON() {
		return output.JSON(l)
	}

	fmt.Printf("Label:        %s\n", l.Name)
	fmt.Printf("ID:           %s\n", l.ID)
	if l.Description != "" {
		fmt.Printf("Description:  %s\n", l.Description)
	}
	fmt.Printf("License:      %s\n", l.DefaultLicense)
	if l.DefaultDistributor != "" {
		fmt.Printf("Distributor:  %s\n", l.DefaultDistributor)
	}
	if l.Contact != nil {
		if l.Contact.Email != "" {
			fmt.Printf("Contact:      %s\n", l.Contact.Email)
		}
		if l.Contact.Website != "" {
			fmt.Printf("Website:      %s\n", l.Contact.Website)
		}
		if l.Contact.Location != "" {
			fmt.Printf("Location:     %s\n", l.Contact.Location)
		}
	}
	if len(l.Sites) > 0 {
		fmt.Printf("Sites:\n")
		for _, s := range l.Sites {
			scope := "all artists"
			if s.ArtistID != "" {
				scope = s.ArtistID
			}
			desc := ""
			if s.Description != "" {
				desc = "  " + s.Description
			}
			fmt.Printf("  %-20s %s  [%s]%s\n", s.ID, s.Path, scope, desc)
		}
	}
	if l.CreatedDate != "" {
		fmt.Printf("Created:      %s\n", l.CreatedDate)
	}
	return nil
}

func runLabelSetName(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	l, err := readLabel(wsPath)
	if err != nil {
		return err
	}
	name := args[0]
	if name == "" {
		return errors.New("VALIDATION_FAILED: name is required")
	}
	oldName := l.Name
	l.Name = name
	if err := writeLabel(wsPath, &l); err != nil {
		return err
	}
	output.Success("Renamed label %q to %q", oldName, name)
	return nil
}

func runLabelSetDescription(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	l, err := readLabel(wsPath)
	if err != nil {
		return err
	}
	l.Description = args[0]
	if err := writeLabel(wsPath, &l); err != nil {
		return err
	}
	output.Success("Updated label description")
	return nil
}

func runLabelSetLicense(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	l, err := readLabel(wsPath)
	if err != nil {
		return err
	}
	oldLicense := l.DefaultLicense
	l.DefaultLicense = args[0]
	if err := writeLabel(wsPath, &l); err != nil {
		return err
	}
	output.Success("Changed default license from %q to %q", oldLicense, l.DefaultLicense)
	return nil
}

func runLabelSetDistributor(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	l, err := readLabel(wsPath)
	if err != nil {
		return err
	}
	l.DefaultDistributor = args[0]
	if err := writeLabel(wsPath, &l); err != nil {
		return err
	}
	output.Success("Set default distributor to %q", l.DefaultDistributor)
	return nil
}

func runLabelSetContact(cmd *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	l, err := readLabel(wsPath)
	if err != nil {
		return err
	}

	if l.Contact == nil {
		l.Contact = &models.LabelContact{}
	}
	if cmd.Flags().Changed("email") {
		l.Contact.Email = labelSetContactEmail
	}
	if cmd.Flags().Changed("website") {
		l.Contact.Website = labelSetContactWebsite
	}
	if cmd.Flags().Changed("location") {
		l.Contact.Location = labelSetContactLocation
	}

	if err := writeLabel(wsPath, &l); err != nil {
		return err
	}
	output.Success("Updated label contact details")
	return nil
}

