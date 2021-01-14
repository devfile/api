package validation

import "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"

// validateEndpoints validates that the endpoints

func validateEndpoints(endpoints []v1alpha2.Endpoint, processedEndPointPort map[int]bool,  processedEndPointName map[string]bool) error {
	currentComponentEndPointPort := make(map[int]bool)

	for _, endPoint := range endpoints {
		if isInt(endPoint.Name) {
			return &InvalidNameOrIdError{name: endPoint.Name, resourceType: "endpoint"}
		}
		if _, ok := processedEndPointName[endPoint.Name]; ok {
			return &InvalidEndpointError{name: endPoint.Name}
		}
		processedEndPointName[endPoint.Name] = true
		currentComponentEndPointPort[endPoint.TargetPort] = true
	}

	for targetPort := range currentComponentEndPointPort {
		if _, ok :=processedEndPointPort[targetPort]; ok {
			return &InvalidEndpointError{port: targetPort}
		}
		processedEndPointPort[targetPort] = true
	}
	return nil
}
