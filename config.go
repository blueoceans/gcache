package gcache

var (
	clientID     string
	clientSecret string
	refreshToken string
)

// SetConfig sets `ClientID`, `ClientSecret` and `RefreshToken` for Google OAuth2.
func SetConfig(
	cid string,
	cs string,
	rt string,
) {
	clientID = cid
	clientSecret = cs
	refreshToken = rt
}
