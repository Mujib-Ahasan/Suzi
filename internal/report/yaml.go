package report

import (
	"os"

	"gopkg.in/yaml.v3"
)

func WriteYAML(result LoadTestResultAll, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	defer enc.Close()

	return enc.Encode(result)
}

func WriteYAMLToStdout(result LoadTestResultAll) error {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	defer enc.Close()

	return enc.Encode(result)
}
