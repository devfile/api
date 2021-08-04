package common

import (
	"fmt"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// CreateParentEndpoints creates and returns a random number of endpoints in a schema structure
func (devfile *TestDevfile) CreateParentEndpoints() []schema.EndpointParentOverride {

	numEndpoints := GetRandomNumber(1, 5)
	endpoints := make([]schema.EndpointParentOverride, numEndpoints)

	commonPort := devfile.getUniquePort()

	for i := 0; i < numEndpoints; i++ {

		endpoint := schema.EndpointParentOverride{}

		endpoint.Name = GetRandomUniqueString(GetRandomNumber(5, 24), true)
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d name  : %s", i, endpoint.Name))

		if GetBinaryDecision() {
			endpoint.TargetPort = devfile.getUniquePort()
		} else {
			endpoint.TargetPort = commonPort
		}
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d targetPort: %d", i, endpoint.TargetPort))

		if GetBinaryDecision() {
			endpoint.Exposure = schema.EndpointExposureParentOverride(getRandomExposure())
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d exposure: %s", i, endpoint.Exposure))
		}

		if GetBinaryDecision() {
			endpoint.Protocol = schema.EndpointProtocolParentOverride(getRandomProtocol())
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d protocol: %s", i, endpoint.Protocol))
		}

		endpoint.Secure = GetBinaryDecision()
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d secure: %t", i, endpoint.Secure))

		if GetBinaryDecision() {
			endpoint.Path = "/Path_" + GetRandomString(GetRandomNumber(3, 15), false)
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d path: %s", i, endpoint.Path))
		}

		endpoints[i] = endpoint

	}

	return endpoints
}
