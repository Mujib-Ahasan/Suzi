package report

import (
	"encoding/json"
	"os"
)

func WriteJSON(result []LoadTestResult, path string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func WriteJSONToStdout(result []LoadTestResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
