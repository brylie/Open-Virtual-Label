package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

// SetReader replaces the default stdin reader. Intended for use in tests.
func SetReader(r io.Reader) {
	scanner = bufio.NewScanner(r)
}

// Ask prompts for a string value. Shows defaultVal in brackets; returns it if user presses Enter.
func Ask(label, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, defaultVal)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}
	if scanner.Scan() {
		v := strings.TrimSpace(scanner.Text())
		if v == "" {
			return defaultVal, nil
		}
		return v, nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return defaultVal, nil
}

// Confirm asks a yes/no question. Returns true for y/yes.
func Confirm(message string) (bool, error) {
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", message)
	if scanner.Scan() {
		v := strings.TrimSpace(strings.ToLower(scanner.Text()))
		return v == "y" || v == "yes", nil
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}

// Select prompts the user to choose one item from a list. Returns the chosen string.
func Select(label string, options []string) (string, error) {
	fmt.Fprintf(os.Stderr, "%s:\n", label)
	for i, opt := range options {
		fmt.Fprintf(os.Stderr, "  %d) %s\n", i+1, opt)
	}
	fmt.Fprintf(os.Stderr, "Choice [1-%d]: ", len(options))
	if scanner.Scan() {
		v := strings.TrimSpace(scanner.Text())
		for i, opt := range options {
			if v == strconv.Itoa(i+1) || strings.EqualFold(v, opt) {
				return opt, nil
			}
		}
		return "", fmt.Errorf("invalid selection: %q", v)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("no selection made")
}
