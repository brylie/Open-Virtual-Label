package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var financeCmd = &cobra.Command{
	Use:   "finance",
	Short: "Log and summarize revenue and expenses",
}

// finance add-revenue
var (
	finRevSource, finRevCurrency, finRevPeriod, finRevArtist, finRevRelease, finRevDesc string
	finRevAmount                                                                         float64
)
var financeAddRevenueCmd = &cobra.Command{
	Use:   "add-revenue",
	Short: "Log a revenue entry",
	RunE:  runFinanceAddRevenue,
}

// finance add-expense
var (
	finExpSource, finExpCurrency, finExpDate, finExpArtist, finExpDesc string
	finExpAmount                                                        float64
)
var financeAddExpenseCmd = &cobra.Command{
	Use:   "add-expense",
	Short: "Log an expense entry",
	RunE:  runFinanceAddExpense,
}

// finance summary
var finSummaryPeriod string
var finSummaryBrief bool
var financeSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Generate a financial summary for a period",
	RunE:  runFinanceSummary,
}

func init() {
	financeAddRevenueCmd.Flags().StringVar(&finRevSource, "source", "", "revenue source (required)")
	financeAddRevenueCmd.Flags().Float64Var(&finRevAmount, "amount", 0, "amount (required)")
	financeAddRevenueCmd.Flags().StringVar(&finRevCurrency, "currency", "EUR", "ISO 4217 currency code")
	financeAddRevenueCmd.Flags().StringVar(&finRevPeriod, "period", "", "period YYYY-MM (required)")
	financeAddRevenueCmd.Flags().StringVar(&finRevArtist, "artist", "", "artist ID")
	financeAddRevenueCmd.Flags().StringVar(&finRevRelease, "release", "", "release ID")
	financeAddRevenueCmd.Flags().StringVar(&finRevDesc, "description", "", "description")
	_ = financeAddRevenueCmd.MarkFlagRequired("source")
	_ = financeAddRevenueCmd.MarkFlagRequired("amount")
	_ = financeAddRevenueCmd.MarkFlagRequired("period")

	financeAddExpenseCmd.Flags().StringVar(&finExpSource, "source", "", "expense category (required)")
	financeAddExpenseCmd.Flags().Float64Var(&finExpAmount, "amount", 0, "amount (required)")
	financeAddExpenseCmd.Flags().StringVar(&finExpCurrency, "currency", "EUR", "ISO 4217 currency code")
	financeAddExpenseCmd.Flags().StringVar(&finExpDate, "date", "", "date YYYY-MM-DD (required)")
	financeAddExpenseCmd.Flags().StringVar(&finExpArtist, "artist", "", "artist ID")
	financeAddExpenseCmd.Flags().StringVar(&finExpDesc, "description", "", "description")
	_ = financeAddExpenseCmd.MarkFlagRequired("source")
	_ = financeAddExpenseCmd.MarkFlagRequired("amount")
	_ = financeAddExpenseCmd.MarkFlagRequired("date")

	financeSummaryCmd.Flags().StringVar(&finSummaryPeriod, "period", "", "period YYYY-MM (required)")
	financeSummaryCmd.Flags().BoolVar(&finSummaryBrief, "brief", false, "one-paragraph summary")
	_ = financeSummaryCmd.MarkFlagRequired("period")

	financeQuoteCmd := &cobra.Command{
		Use:   "quote",
		Short: "Generate a pricing quote for a commission opportunity (requires finance-manager agent)",
		RunE:  func(*cobra.Command, []string) error { return agentStub("finance-manager") },
	}
	financeQuoteCmd.Flags().String("opportunity", "", "opportunity ID (required)")
	_ = financeQuoteCmd.MarkFlagRequired("opportunity")

	financeCmd.AddCommand(financeAddRevenueCmd, financeAddExpenseCmd, financeSummaryCmd, financeQuoteCmd)
	rootCmd.AddCommand(financeCmd)
}

func runFinanceAddRevenue(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	// Convert YYYY-MM to YYYY-MM-01
	date := finRevPeriod + "-01"
	entryID := ws.ToSlug(finRevSource + "-" + finRevPeriod)

	entry := models.FinanceEntry{
		SchemaVersion: "1",
		ID:            entryID,
		Type:          financeTypeRevenue,
		Date:          date,
		Amount:        finRevAmount,
		Currency:      strings.ToUpper(finRevCurrency),
		Source:        finRevSource,
		Description:   finRevDesc,
	}
	if finRevArtist != "" {
		entry.ArtistID = &finRevArtist
	}
	if finRevRelease != "" {
		entry.ReleaseID = &finRevRelease
	}

	raw, _ := json.Marshal(entry)
	if errs, err := schema.Validate("finance-entry", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: strings.Join(errs, "\n")}
	}

	if err := appendFinanceEntry(wsPath, "revenue.json", &entry); err != nil {
		return err
	}
	output.Success("Revenue entry %s added (%.2f %s)", entryID, finRevAmount, entry.Currency)
	return nil
}

func runFinanceAddExpense(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	entryID := ws.ToSlug(finExpSource + "-" + finExpDate)

	entry := models.FinanceEntry{
		SchemaVersion: "1",
		ID:            entryID,
		Type:          "expense",
		Date:          finExpDate,
		Amount:        finExpAmount,
		Currency:      strings.ToUpper(finExpCurrency),
		Source:        finExpSource,
		Description:   finExpDesc,
	}
	if finExpArtist != "" {
		entry.ArtistID = &finExpArtist
	}

	raw, _ := json.Marshal(entry)
	if errs, err := schema.Validate("finance-entry", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: strings.Join(errs, "\n")}
	}

	if err := appendFinanceEntry(wsPath, "expenses.json", &entry); err != nil {
		return err
	}
	output.Success("Expense entry %s added (%.2f %s)", entryID, finExpAmount, entry.Currency)
	return nil
}

func runFinanceSummary(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}

	prefix := finSummaryPeriod + "-"
	var revenue, expenses []models.FinanceEntry

	readEntries(filepath.Join(wsPath, "finance", "revenue.json"), &revenue)
	readEntries(filepath.Join(wsPath, "finance", "expenses.json"), &expenses)

	var periodRevenue, periodExpenses []models.FinanceEntry
	for i := range revenue {
		if strings.HasPrefix(revenue[i].Date, prefix) {
			periodRevenue = append(periodRevenue, revenue[i])
		}
	}
	for i := range expenses {
		if strings.HasPrefix(expenses[i].Date, prefix) {
			periodExpenses = append(periodExpenses, expenses[i])
		}
	}

	totalRev := sumEntries(periodRevenue)
	totalExp := sumEntries(periodExpenses)

	if finSummaryBrief {
		fmt.Printf("Period %s: Revenue %.2f, Expenses %.2f, Net %.2f\n",
			finSummaryPeriod, totalRev, totalExp, totalRev-totalExp)
		return nil
	}

	fmt.Printf("Financial Summary — %s\n\n", finSummaryPeriod)
	fmt.Printf("Revenue (%d entries):\n", len(periodRevenue))
	bySource := map[string]float64{}
	for i := range periodRevenue {
		bySource[periodRevenue[i].Source] += periodRevenue[i].Amount
	}
	for src, amt := range bySource {
		fmt.Printf("  %-24s  %.2f\n", src, amt)
	}
	fmt.Printf("  %-24s  %.2f\n", "TOTAL", totalRev)
	fmt.Println()

	fmt.Printf("Expenses (%d entries):\n", len(periodExpenses))
	bySource = map[string]float64{}
	for i := range periodExpenses {
		bySource[periodExpenses[i].Source] += periodExpenses[i].Amount
	}
	for src, amt := range bySource {
		fmt.Printf("  %-24s  %.2f\n", src, amt)
	}
	fmt.Printf("  %-24s  %.2f\n", "TOTAL", totalExp)
	fmt.Println()

	net := totalRev - totalExp
	fmt.Printf("Net: %.2f\n", net)

	return nil
}

func appendFinanceEntry(wsPath, filename string, entry *models.FinanceEntry) error {
	path := filepath.Join(wsPath, "finance", filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	entries := make([]models.FinanceEntry, 0, 1)
	readEntries(path, &entries)
	entries = append(entries, *entry)
	data, err := json.MarshalIndent(entries, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func readEntries(path string, entries *[]models.FinanceEntry) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, entries)
}

func sumEntries(entries []models.FinanceEntry) float64 {
	var total float64
	for i := range entries {
		total += entries[i].Amount
	}
	return total
}
