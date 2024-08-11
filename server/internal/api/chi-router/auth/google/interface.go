package googleapi

type GoogleUserInfo struct {
	LocalId           string
	Email             string
	DisplayName       string
	Language          string
	PhotoUrl          string
	TimeZone          string
	DateOfBirth       string
	PasswordHash      string
	Salt              string
	Version           int
	EmailVerified     bool
	PasswordUpdatedAt int
	ProviderUserInfo  []ProviderUserInfo
	ValidSince        string
	Disabled          bool
	LastLoginAt       string
	CreatedAt         string
	ScreenName        string
	CustomAuth        bool
	RawPassword       string
	PhoneNumber       string
	CustomAttributes  string
	EmailLinkSignin   bool
	TenantId          string
	MfaInfo           []MfaEnrollment
	InitialEmail      string
	LastRefreshAt     string
}

type MfaEnrollment struct {
	MfaEnrollmentId       string
	DisplayName           string
	EnrolledAt            string
	PhoneInfo             string
	TotpInfo              TotpInfo
	EmailInfo             EmailInfo
	UnobfuscatedPhoneInfo string
}

type TotpInfo struct {
	EmailAddress string
}

type EmailInfo struct {
	EmailAddress string
}

type ProviderUserInfo struct {
	ProviderId  string
	DisplayName string
	PhotoUrl    string
	FederatedId string
	Email       string
	RawId       string
	ScreenName  string
	PhoneNumber string
}
