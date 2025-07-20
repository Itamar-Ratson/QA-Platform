package yaml

import (
    "gopkg.in/yaml.v3"
    "io/ioutil"
)

type TestCase struct {
    Metadata struct {
        Name           string `yaml:"name"`
        Type           string `yaml:"type"`
        Priority       string `yaml:"priority"`
        Severity       string `yaml:"severity"`
        ExpectedResult string `yaml:"expected_result"`
        Description    string `yaml:"description"`
    } `yaml:"metadata"`
    Terraform struct {
        TfVars map[string]interface{} `yaml:"tfvars"`
    } `yaml:"terraform"`
    TestFunctions []string `yaml:"test_functions"`
}

func ParseTestCase(filepath string) (*TestCase, error) {
    data, err := ioutil.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

    var tc TestCase
    err = yaml.Unmarshal(data, &tc)
    return &tc, err
}
