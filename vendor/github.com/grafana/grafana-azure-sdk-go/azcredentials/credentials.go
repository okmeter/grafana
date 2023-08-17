package azcredentials

const (
	AzureAuthCurrentUserIdentity = "currentuser"
	AzureAuthManagedIdentity     = "msi"
	AzureAuthClientSecret        = "clientsecret"
	AzureAuthClientSecretObo     = "clientsecret-obo"
)

type AzureCredentials interface {
	AzureAuthType() string
}

// AadCurrentUserCredentials "Current User" user identity credentials of the current Grafana user.
type AadCurrentUserCredentials struct {
}

// AzureManagedIdentityCredentials "Managed Identity" service managed identity credentials configured
// for the current Grafana instance.
type AzureManagedIdentityCredentials struct {
	ClientId string
}

// AzureClientSecretCredentials "App Registration" AAD service identity credentials configured in the datasource.
type AzureClientSecretCredentials struct {
	AzureCloud   string
	Authority    string
	TenantId     string
	ClientId     string
	ClientSecret string
}

// AzureClientSecretOboCredentials "App Registration (On-Behalf-Of)" user identity credentials obtained using
// service identity configured in the datasource.
type AzureClientSecretOboCredentials struct {
	ClientSecretCredentials AzureClientSecretCredentials
}

func (credentials *AadCurrentUserCredentials) AzureAuthType() string {
	return AzureAuthCurrentUserIdentity
}

func (credentials *AzureManagedIdentityCredentials) AzureAuthType() string {
	return AzureAuthManagedIdentity
}

func (credentials *AzureClientSecretCredentials) AzureAuthType() string {
	return AzureAuthClientSecret
}

func (credentials *AzureClientSecretOboCredentials) AzureAuthType() string {
	return AzureAuthClientSecretObo
}
