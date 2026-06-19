package cmd

import (
	"github.com/spf13/cobra"
)

var contentCmd = &cobra.Command{
	Use:   "content",
	Short: "Content campaign briefs and social copy",
}

var socialCmd = &cobra.Command{
	Use:   "social",
	Short: "Generate social media copy",
}

func init() {
	contentBriefCmd := &cobra.Command{
		Use:   "brief",
		Short: "Generate a content campaign brief (requires content-strategist agent)",
		RunE:  func(*cobra.Command, []string) error { return agentStub("content-strategist") },
	}
	contentBriefCmd.Flags().String("release", "", "release ID (required)")
	_ = contentBriefCmd.MarkFlagRequired("release")

	socialDraftCmd := &cobra.Command{
		Use:   cmdUseDraft,
		Short: "Generate social media copy for a release campaign (requires social-media agent)",
		RunE:  func(*cobra.Command, []string) error { return agentStub("social-media") },
	}
	socialDraftCmd.Flags().String("release", "", "release ID (required)")
	socialDraftCmd.Flags().String("platform", "", "instagram, youtube, or facebook (all if omitted)")
	_ = socialDraftCmd.MarkFlagRequired("release")

	socialCmd.AddCommand(socialDraftCmd)

	contentCmd.AddCommand(contentBriefCmd)
	rootCmd.AddCommand(contentCmd)
	rootCmd.AddCommand(socialCmd)
}
