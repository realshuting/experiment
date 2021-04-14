package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiservervalidation "k8s.io/apiextensions-apiserver/pkg/apiserver/validation"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
	policyPath := "policy.yaml"
	crdPath := "crd-old.yaml"

	policyBytes := convertToJSONbytes(policyPath)
	u := &unstructured.Unstructured{}
	err := u.UnmarshalJSON(policyBytes)
	if err != nil {
		fmt.Println("failed to decode policy", err)
	}

	var v1crd apiextensions.CustomResourceDefinitionSpec
	crdBytes := convertToJSONbytes(crdPath)
	if err := json.Unmarshal(crdBytes, &v1crd); err != nil {
		fmt.Println("failed to decode crd: ", err)
	}

	versions := v1crd.Versions
	for _, version := range versions {
		validator, _, err := apiservervalidation.NewSchemaValidator(&apiextensions.CustomResourceValidation{OpenAPIV3Schema: version.Schema.OpenAPIV3Schema})
		if err != nil {
			fmt.Println("failed to create schema validator", err)
		}

		errList := apiservervalidation.ValidateCustomResource(nil, u.UnstructuredContent(), validator)
		if errList != nil {
			fmt.Println(errList)
		}
	}
}

func convertToJSONbytes(path string) []byte {
	pathBytes, err := ioutil.ReadFile(path)
	if err != nil {
		println("error in extracting in bytes: ", err)
	}
	jsonBytes, err := yaml.ToJSON(pathBytes)
	if err != nil {
		fmt.Printf("failed to convert to JSON: %v\n", err)
	}
	return jsonBytes
}
