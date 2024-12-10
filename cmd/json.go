package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// printJSON prints the given data as formatted JSON to the specified writer
func printJSON(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

// printJSONToStdout is a convenience function that prints JSON to stdout
func printJSONToStdout(v interface{}) error {
	return printJSON(os.Stdout, v)
}

// formatJSON returns the formatted JSON string
func formatJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}
