package application

import (
	"context"
	"fmt"
	"log"

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/domain"
	"src/internal/pkg/notion"
)

// NotionOAuthRequest contains the data needed to complete Notion OAuth
type NotionOAuthRequest struct {
	Code  string
	State string
}

// NotionOAuthResponse contains the OAuth completion result
type NotionOAuthResponse struct {
	User        domain.User
	AccessToken string
	WorkspaceID string
	BotID       string
}

// NotionOAuthUseCase handles Notion OAuth flow completion
type NotionOAuthUseCase struct {
	repo         domain.UserRepository
	clock        shared.Clock
	txMgr        shared.TransactionManager
	idGen        shared.IDGenerator
	notionClient *notion.Service
}

// NewNotionOAuthUseCase creates a new NotionOAuthUseCase
func NewNotionOAuthUseCase(
	repo domain.UserRepository,
	clock shared.Clock,
	txMgr shared.TransactionManager,
	idGen shared.IDGenerator,
	notionClient *notion.Service,
) *NotionOAuthUseCase {
	return &NotionOAuthUseCase{
		repo:         repo,
		clock:        clock,
		txMgr:        txMgr,
		idGen:        idGen,
		notionClient: notionClient,
	}
}

// Execute completes the Notion OAuth flow and creates or updates a user
func (uc *NotionOAuthUseCase) Execute(ctx context.Context, req NotionOAuthRequest) (NotionOAuthResponse, error) {
	var response NotionOAuthResponse

	err := uc.txMgr.WithinTransaction(ctx, func(ctx context.Context) error {
		// Exchange code for token
		tokenResp, err := uc.notionClient.ExchangeCodeForToken(req.Code)
		if err != nil {
			return fmt.Errorf("failed to exchange code for token ( us ): %w", err)
		}

		// Get user info from Notion
		log.Println("tokenResp", tokenResp)
		notionUser, err := uc.notionClient.GetCurrentUser(tokenResp.AccessToken)
		if err != nil {
			return fmt.Errorf("failed to get user info ( us ): %w", err)
		}

		// Extract email from person object
		email := ""
		if notionUser.Person != nil && notionUser.Person.Email != "" {
			email = notionUser.Person.Email
		} else if notionUser.Bot != nil && notionUser.Bot.Owner.User.Person != nil && notionUser.Bot.Owner.User.Person.Email != "" {
			email = notionUser.Bot.Owner.User.Person.Email
		} else {
			return fmt.Errorf("no email available from Notion user")
		}

		// Try to find existing user by email
		user, err := uc.repo.GetByEmail(ctx, email)
		if err == domain.ErrUserNotFound {
			// Create new user (ID and PublicID will be generated inside NewUser)
			user, err = domain.NewUser(email, notionUser.Name, uc.idGen, uc.clock)
			if err != nil {
				return fmt.Errorf("failed to create new user: %w", err)
			}

			// Set Notion token information
			user.SetNotionToken(
				tokenResp.AccessToken,
				tokenResp.WorkspaceID,
				tokenResp.BotID,
				nil, // Notion tokens don't typically expire
				uc.clock,
			)

			// Save new user
			user, err = uc.repo.Create(ctx, user)
			if err != nil {
				return fmt.Errorf("failed to save new user: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check existing user: %w", err)
		} else {
			// Update existing user with new token
			user.SetNotionToken(
				tokenResp.AccessToken,
				tokenResp.WorkspaceID,
				tokenResp.BotID,
				nil,
				uc.clock,
			)

			// Update user name if it changed
			if user.Name != notionUser.Name {
				user.Name = notionUser.Name
				user.UpdatedAt = uc.clock.Now()
			}

			// Save updated user
			user, err = uc.repo.Update(ctx, user)
			if err != nil {
				return fmt.Errorf("failed to update user: %w", err)
			}
		}

		response = NotionOAuthResponse{
			User:        user,
			AccessToken: tokenResp.AccessToken,
			WorkspaceID: tokenResp.WorkspaceID,
			BotID:       tokenResp.BotID,
		}

		return nil
	})

	if err != nil {
		return NotionOAuthResponse{}, err
	}

	return response, nil
}

// GetAuthorizationURLRequest contains the data needed to generate auth URL
type GetAuthorizationURLRequest struct {
	State string
}

// GetAuthorizationURLResponse contains the authorization URL
type GetAuthorizationURLResponse struct {
	AuthorizationURL string
}

// GetAuthorizationURLUseCase generates Notion OAuth authorization URL
type GetAuthorizationURLUseCase struct {
	notionClient *notion.Service
}

// NewGetAuthorizationURLUseCase creates a new GetAuthorizationURLUseCase
func NewGetAuthorizationURLUseCase(notionClient *notion.Service) *GetAuthorizationURLUseCase {
	return &GetAuthorizationURLUseCase{
		notionClient: notionClient,
	}
}

// Execute generates the Notion OAuth authorization URL
func (uc *GetAuthorizationURLUseCase) Execute(ctx context.Context, req GetAuthorizationURLRequest) (GetAuthorizationURLResponse, error) {
	authURL := uc.notionClient.GetAuthorizationURL(req.State)

	return GetAuthorizationURLResponse{
		AuthorizationURL: authURL,
	}, nil
}
