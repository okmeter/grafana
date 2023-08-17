package azcredentials

import (
	"fmt"

	"github.com/grafana/grafana-azure-sdk-go/azsettings"
)

func GetAzureCloud(settings *azsettings.AzureSettings, credentials AzureCredentials) (string, error) {
	switch c := credentials.(type) {
	case *AadCurrentUserCredentials:
		// In case of user identity, the cloud is always same as where Grafana is hosted
		return settings.GetDefaultCloud(), nil
	case *AzureManagedIdentityCredentials:
		// In case of managed identity, the cloud is always same as where Grafana is hosted
		return settings.GetDefaultCloud(), nil
	case *AzureClientSecretCredentials:
		return c.AzureCloud, nil
	case *AzureClientSecretOboCredentials:
		return c.ClientSecretCredentials.AzureCloud, nil
	default:
		err := fmt.Errorf("the Azure credentials of type '%s' not supported", c.AzureAuthType())
		return "", err
	}
}
