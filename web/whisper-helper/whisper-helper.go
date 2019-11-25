package whisper_helper

import (
	"fmt"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper-client/client"
	"net/http"
	"time"
)

var (
	defaultPostLoginRoute  = "/whisper/post-login"
	defaultPostLogoutRoute = "/whisper/post-logout"
	
	accessTokenCookieName  = "WHISPER_ACCESS_TOKEN"
	openIdTokenCookieName  = "WHISPER_OPENID_TOKEN"
	refreshTokenCookieName = "WHISPER_REFRESH_TOKEN"
)

// WhisperHelper implements common/recommended usages of the whisper client API
type WhisperHelper struct {
	*client.WhisperClient
}

// InitFromParams set the whisper client and returns a helper
// If the login redirect URI or the logout redirect URI is not set, it will use the default route
func (helper *WhisperHelper) InitFromParams(whisperURL, clientID, clientSecret, loginRedirectURI, logoutRedirectURI  string, scopes []string) *WhisperHelper {
	helper.WhisperClient = new(client.WhisperClient).
		InitFromParams(whisperURL, clientID, clientSecret, loginRedirectURI, logoutRedirectURI, scopes)

	return helper
}

// GetWhisperToken connects with whisper, register in hydra and retrieve token to be used for authentication inside whisper
func (helper *WhisperHelper) GetWhisperToken() (tokenString string, err error) {
	token, err := helper.CheckCredentials()
	if err != nil {
		return
	}

	tokenString = helper.GetTokenAsJSONStr(token)
	if tokenString == "" {
		err = fmt.Errorf("Unable to extract token")
	}

	return
}

// GetPostLoginHandler exchange the code in the request for tokens, set them as cookies and redirect to the specified url
// The purpose of this is to make it seamless for the developer and user the setup of the cookies.
func (helper *WhisperHelper) GetPostLoginHandler(redirectURI string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		tokens, err := helper.ExchangeCodeForToken(code)
		gohtypes.PanicIfError("Unable to exchange code for token", http.StatusInternalServerError, err)

		SetCookies(w, tokens.AccessToken, tokens.OpenIdToken, tokens.RefreshToken)

		http.Redirect(w, r, redirectURI, http.StatusFound)
	}
}

// GetPostLogoutHandler overwrite cookies and redirect to the specified url
// The purpose of this is to make it seamless for the developer and user the removal of the cookies.
func (helper *WhisperHelper) GetPostLogoutHandler(redirectURI string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		UnsetCookies(w)

		http.Redirect(w, r, redirectURI, http.StatusFound)
	}
}

// GetUserToken extract tokens and introspect the access token with hydra
// From the request it will identify the cookies set on post login and extract them.
// If they are valid, it means that there is a logged user and the lone token returned contains its username.
func (helper *WhisperHelper) GetUserToken(w http.ResponseWriter, r *http.Request) (tokens client.Tokens, token client.Token, err error) {
	accessToken, openIdToken, refreshToken, err := ExtractTokensFromRequest(r)
	if err != nil {
		return
	}

	tokens = client.Tokens{
		AccessToken:  accessToken,
		RefreshToken: openIdToken,
		OpenIdToken:  refreshToken,
	}

	token, err = helper.IntrospectToken(tokens.AccessToken)
	if err != nil {
		return
	}

	if !token.Active {
		UnsetCookies(w)
		return
	}

	return
}

// IsUserLogged verify if the user has a valid access token
// From the request it will identify the cookies set on post login and extract the access token.
// If that token is valid, it means that there is a user logged.
func (helper *WhisperHelper) IsUserLogged(r *http.Request) bool {
	accessToken, _, _, err := ExtractTokensFromRequest(r)
	if err != nil {
		return false
	}

	token, err := helper.IntrospectToken(accessToken)
	gohtypes.PanicIfError("Unable to introspect token", http.StatusInternalServerError, err)

	if token.Active {
		return false
	}

	return true
}

// ExtractTokensFromRequest retrieves from certain cookies the access token, open id token and refresh token
func ExtractTokensFromRequest(r *http.Request) (accessToken, openIdToken, refreshToken string, err error) {
	cookie, err := r.Cookie(accessTokenCookieName)
	if err != nil {
		return
	}

	accessToken = cookie.Value

	cookie, err = r.Cookie(openIdTokenCookieName)
	if err != nil {
		return
	}

	openIdToken = cookie.Value

	cookie, err = r.Cookie(refreshTokenCookieName)
	if err != nil {
		return
	}

	refreshToken = cookie.Value

	return
}

// SetCookie set a simple cookie that expires in a week
func SetCookies(w http.ResponseWriter, accessToken, openIdToken, refreshToken string) {
	setCookie := func(w http.ResponseWriter, name, value string) {
		oneWeekFromNow := time.Now().Add(7 * 24 * time.Hour)

		http.SetCookie(w, &http.Cookie{
			Name:    name,
			Value:   value,
			Expires: oneWeekFromNow,
		})
	}

	setCookie(w, accessTokenCookieName, accessToken)
	setCookie(w, openIdTokenCookieName, openIdToken)
	setCookie(w, refreshTokenCookieName, refreshToken)
}

// UnsetCookie overwrite a cookie and make it expired
func UnsetCookies(w http.ResponseWriter) {
	unsetCookie := func(w http.ResponseWriter, name string) {
		past := time.Unix(0, 0)

		http.SetCookie(w, &http.Cookie{
			Name:    name,
			Value:   "",
			Expires: past,
		})
	}

	unsetCookie(w, accessTokenCookieName)
	unsetCookie(w, openIdTokenCookieName)
	unsetCookie(w, refreshTokenCookieName)
}

// GetPostLoginWhisper retrieves the post login route used
func GetDefaultPostLoginRoute() string {
	return defaultPostLoginRoute
}

// GetPostLogoutWhisper retrieves the post logout route used
func GetDefaultPostLogoutRoute() string {
	return defaultPostLogoutRoute
}