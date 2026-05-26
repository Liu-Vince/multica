package handler

import (
	"net/http"
	"os"

	"github.com/multica-ai/multica/server/internal/analytics"
)

type AppConfig struct {
	CdnDomain string `json:"cdn_domain"`
	// Public auth config consumed by the web app at runtime so self-hosted
	// deployments do not need to rebuild the frontend image when operators
	// toggle signup or wire Google OAuth.
	AllowSignup    bool   `json:"allow_signup"`
	GoogleClientID string `json:"google_client_id,omitempty"`

	// PostHog public config for the frontend. The key is the same Project
	// API Key the backend uses; returning it here (instead of baking it
	// into the frontend bundle via NEXT_PUBLIC_*) means self-hosted
	// instances — whose server returns an empty key — automatically
	// disable frontend event shipping too.
	PosthogKey           string `json:"posthog_key"`
	PosthogHost          string `json:"posthog_host"`
	AnalyticsEnvironment string `json:"analytics_environment"`

	// CLI / install URLs that frontend surfaces show to users when adding
	// remote machines.  Empty strings are omitted from the JSON response so
	// the web app can detect "not configured" and skip rendering commands
	// that would reference an unreachable host.
	InstallScriptURL  string `json:"install_script_url,omitempty"`
	CLIDownloadBaseURL string `json:"cli_download_base_url,omitempty"`
	CLIServerURL       string `json:"cli_server_url,omitempty"`
	CLIAppURL          string `json:"cli_app_url,omitempty"`
}

const (
	defaultInstallScriptURL   = "https://raw.githubusercontent.com/multica-ai/multica/main/scripts/install.sh"
	defaultCLIDownloadBaseURL = "https://github.com/multica-ai/multica/releases/latest/download"
)

// GetConfig is mounted on the public (unauthenticated) route group because
// the web app calls it before login to decide whether to render the Google
// sign-in button and signup UI. Only add fields here that are safe to expose
// to anonymous callers — never user- or tenant-scoped data.
func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	config := AppConfig{
		AllowSignup:    os.Getenv("ALLOW_SIGNUP") != "false",
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),

		InstallScriptURL:  envOrDefault("INSTALL_SCRIPT_URL", defaultInstallScriptURL),
		CLIDownloadBaseURL: envOrDefault("CLI_DOWNLOAD_BASE_URL", defaultCLIDownloadBaseURL),
		CLIServerURL:       os.Getenv("CLI_SERVER_URL"),
		CLIAppURL:          os.Getenv("CLI_APP_URL"),
	}
	if h.Storage != nil {
		config.CdnDomain = h.Storage.CdnDomain()
	}

	// Re-read from env on every request so operators can rotate keys via
	// secret refresh without a server restart.
	if v := os.Getenv("ANALYTICS_DISABLED"); v != "true" && v != "1" {
		config.PosthogKey = os.Getenv("POSTHOG_API_KEY")
		config.PosthogHost = os.Getenv("POSTHOG_HOST")
		config.AnalyticsEnvironment = analytics.EnvironmentFromEnv()
		if config.PosthogHost == "" && config.PosthogKey != "" {
			config.PosthogHost = "https://us.i.posthog.com"
		}
	}

	writeJSON(w, http.StatusOK, config)
}

// envOrDefault returns os.Getenv(key) if set and non-empty, otherwise fallback.
func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}