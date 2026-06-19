package auth

import "context"

// DeviceLoginOptions configures dev/dogfood login without full GitHub App flow.
type DeviceLoginOptions struct {
	Token  string
	APIURL string
	OrgID  string
}

// DeviceLogin stores upload credentials for capture server access.
// Full GitHub device flow lands in a later slice; V1 dogfood uses --token.
func DeviceLogin(_ context.Context, opts DeviceLoginOptions) error {
	creds := UploadCredentials{
		Token:  opts.Token,
		APIURL: opts.APIURL,
		OrgID:  opts.OrgID,
	}
	return SaveUpload(creds)
}
