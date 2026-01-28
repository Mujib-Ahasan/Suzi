package cli

import (
	"os"

	"gopkg.in/yaml.v3"
)

func WriteToYAML(currentAttack Attack, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	defer enc.Close()

	return enc.Encode(currentAttack)
}

func WriteYAMLToStdout(currentAttack Attack) error {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	defer enc.Close()

	return enc.Encode(currentAttack)
}
