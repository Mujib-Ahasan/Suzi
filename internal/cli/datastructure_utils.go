package cli

import (
	"os"
	"strings"

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

// LoadAndValidateAttack validate attack yaml file
func (currentAttack *Attack) LoadAndValidateAttack(path string) (*Attack, []byte, string, map[string]string, error) {

	err := DecodeStrictYAML(path, &currentAttack)
	if err != nil {
		return &Attack{}, nil, "", nil, err
	}

	err = currentAttack.SetValue(set)
	if err != nil {
		return &Attack{}, nil, "", nil, err
	}

	err = currentAttack.isValidURL()
	if err != nil {
		return &Attack{}, nil, "", nil, err
	}

	currentAttack.AttackMethod = strings.ToUpper(currentAttack.AttackMethod)
	if err := currentAttack.ValidateYaml(); err != nil {
		return &Attack{}, nil, "", nil, err
	}

	err = currentAttack.ValidateMethod()
	if err != nil {
		return &Attack{}, nil, "", nil, err
	}

	err = currentAttack.ValidateBeforeAttack()
	if err != nil {
		return &Attack{}, nil, "", nil, err
	}

	payload, err := currentAttack.ValidateBody()
	if err != nil {
		return &Attack{}, nil, "", nil, err
	}

	attackContentType, err := ValidateContentType(currentAttack.AttackContentType)
	if err != nil {
		return &Attack{}, nil, "", nil, err
	}

	header, err := currentAttack.ValidateHeader()
	if err != nil {
		return &Attack{}, nil, "", nil, err
	}

	return currentAttack, payload, attackContentType, header, nil
}
