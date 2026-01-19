package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/goravaa/apache-polaris-cli/pkg/api"
	"github.com/goravaa/apache-polaris-cli/pkg/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	clientID     string
	clientSecret string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Commands for authenticating with the Apache Polaris server.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Polaris server",
	Long: `Authenticate with the Apache Polaris server using client credentials.

You can provide credentials via flags:
  polaris auth login --client-id <id> --client-secret <secret>

Or you can run interactively (will prompt for credentials):
  polaris auth login

Environment variables are also supported:
  POLARIS_CLIENT_ID and POLARIS_CLIENT_SECRET`,
	RunE: runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the Polaris server",
	Long:  `Clear stored authentication credentials.`,
	RunE:  runLogout,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Check if you are currently authenticated with the Polaris server.`,
	RunE:  runStatus,
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the access token",
	Long:  `Refresh the current access token using the stored credentials.`,
	RunE:  runRefresh,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(refreshCmd)

	loginCmd.Flags().StringVar(&clientID, "client-id", "", "OAuth client ID")
	loginCmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth client secret")
}

func runLogin(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Host == "" {
		return fmt.Errorf("Polaris host not configured. Run 'polaris config set --host <url>' first")
	}

	id := clientID
	if id == "" {
		id = os.Getenv("POLARIS_CLIENT_ID")
	}
	if id == "" {
		fmt.Print("Client ID: ")
		reader := bufio.NewReader(os.Stdin)
		id, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read client ID: %w", err)
		}
		id = strings.TrimSpace(id)
	}

	if id == "" {
		return fmt.Errorf("client ID is required")
	}

	secret := clientSecret
	if secret == "" {
		secret = os.Getenv("POLARIS_CLIENT_SECRET")
	}
	if secret == "" {
		fmt.Print("Client Secret: ")
		secretBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read client secret: %w", err)
		}
		fmt.Println()
		secret = string(secretBytes)
	}

	if secret == "" {
		return fmt.Errorf("client secret is required")
	}

	fmt.Printf("Authenticating with %s...\n", cfg.Host)
	authClient := api.NewAuthClient(cfg)
	credentials, err := authClient.Login(id, secret)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := config.SaveCredentials(credentials); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Println("✓ Successfully authenticated!")
	if credentials.ExpiresIn > 0 {
		fmt.Printf("  Token expires in: %d seconds\n", credentials.ExpiresIn)
	}
	if credentials.Scope != "" {
		fmt.Printf("  Scope: %s\n", credentials.Scope)
	}

	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	authClient := api.NewAuthClient(cfg)
	if err := authClient.Logout(); err != nil {
		return fmt.Errorf("failed to log out: %w", err)
	}

	fmt.Println("✓ Successfully logged out!")
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Host: %s\n", cfg.Host)
	if cfg.Realm != "" {
		fmt.Printf("Realm: %s\n", cfg.Realm)
	}

	creds, err := config.LoadCredentials()
	if err != nil {
		fmt.Println("Status: Not authenticated")
		fmt.Println("\nRun 'polaris auth login' to authenticate.")
		return nil
	}

	fmt.Println("Status: Authenticated ✓")
	if creds.TokenType != "" {
		fmt.Printf("Token Type: %s\n", creds.TokenType)
	}
	if creds.Scope != "" {
		fmt.Printf("Scope: %s\n", creds.Scope)
	}
	if creds.ClientID != "" {
		fmt.Printf("Client ID: %s\n", creds.ClientID)
	}

	if len(creds.AccessToken) > 20 {
		fmt.Printf("Access Token: %s...%s\n", creds.AccessToken[:10], creds.AccessToken[len(creds.AccessToken)-5:])
	}

	return nil
}

func runRefresh(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	creds, err := config.LoadCredentials()
	if err != nil {
		return fmt.Errorf("not authenticated: %w", err)
	}

	fmt.Println("Refreshing access token...")
	authClient := api.NewAuthClient(cfg)
	newCreds, err := authClient.RefreshToken(creds.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	if err := config.SaveCredentials(newCreds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Println("✓ Token refreshed successfully!")
	if newCreds.ExpiresIn > 0 {
		fmt.Printf("  New token expires in: %d seconds\n", newCreds.ExpiresIn)
	}

	return nil
}
