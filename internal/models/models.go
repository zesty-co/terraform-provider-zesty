package models

import "time"

type (
	OnboardingStatus string
	CloudProvider    string
	Product          string
)

const (
	AWS   CloudProvider = "AWS"
	Azure CloudProvider = "Azure"
	GCP   CloudProvider = "GCP"

	Kompass   Product = "Kompass"
	CM        Product = "CM"
	ZestyDisk Product = "ZestyDisk"

	DefaultHostURL string = "https://api.cloudvisor.io/kompass-platform"
)

type ProductDetails struct {
	Active bool `json:"active" dynamodbav:"active"`
}

type Payload struct {
	AccountID     string                     `json:"accountID"`
	CloudProvider CloudProvider              `json:"cloudProvider"`
	AWSRegion     *string                    `json:"awsRegion,omitempty"`
	RoleARN       string                     `json:"roleARN"`
	ExternalID    string                     `json:"externalID"`
	Products      map[Product]ProductDetails `json:"products"`
}

type Account struct {
	OrganizationID   int64
	OnboardingStatus OnboardingStatus
	AccountID        string
	AWSRegion        *string
	CloudProvider    CloudProvider
	Products         map[Product]ProductDetails

	CreatedAt      time.Time
	UpdatedAt      time.Time
	AdditionalData map[string]any
}
