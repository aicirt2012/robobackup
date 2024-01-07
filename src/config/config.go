package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aicirt2012/robobackup/src/sys/fs"
	"github.com/aicirt2012/robobackup/src/sys/prompt"
	"sigs.k8s.io/yaml"
)

func Import(version string) ConfigDto {
	cfg, err := performImport(version)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func performImport(version string) (ConfigDto, error) {
	fileName, err := selectFile()
	if err != nil {
		return ConfigDto{}, err
	}
	return readFile(version, fileName)
}

func selectFile() (string, error) {
	fileNames, _ := fs.FileNamesWithSuffix(fileSuffix)
	switch len(fileNames) {
	case 0:
		return "", errors.New("config file '*" + fileSuffix + "' does not exist")
	case 1:
		return fileNames[0], nil
	default:
		i := prompt.Select("Please select configuration:", fileNames)
		return fileNames[i], nil
	}
}

func readFile(version string, fileName string) (ConfigDto, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return ConfigDto{}, errors.New("could not import config file '" + fileName + "'.")
	}
	return unmarshal(content, version)
}

func unmarshal(content []byte, version string) (ConfigDto, error) {
	config := ConfigDto{}
	if err := yaml.UnmarshalStrict(content, &config); err != nil {
		prefix := "error unmarshaling JSON: while decoding JSON: json: "
		msg := strings.TrimPrefix(err.Error(), prefix)
		return config, errors.New("could not deserialize config file: " + msg)
	}
	if version != config.Version {
		return config, errors.New("config version does not match the execution version")
	}
	return config, nil
}
