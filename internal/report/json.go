package report

import (
	"encoding/json"
	"os"
)

func WriteJSON(result LoadTestResultAll, path string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func WriteJSONToStdout(result LoadTestResultAll) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
