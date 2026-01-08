package application

import (
	"context"

	"src/internal/modules/users/domain"
)

type ExtensionOAuthUseCase struct {
	authURLUC *GetAuthorizationURLUseCase
	oauthUC   *NotionOAuthUseCase
	stateSvc  domain.ExtensionOAuthStateService
}

func NewExtensionOAuthUseCase(
	authURLUC *GetAuthorizationURLUseCase,
	oauthUC *NotionOAuthUseCase,
	stateSvc domain.ExtensionOAuthStateService,
) *ExtensionOAuthUseCase {
	return &ExtensionOAuthUseCase{
		authURLUC: authURLUC,
		oauthUC:   oauthUC,
		stateSvc:  stateSvc,
	}
}

func (uc *ExtensionOAuthUseCase) Start(ctx context.Context, redirectURI string) (GetAuthorizationURLResponse, string, error) {
	state, err := uc.stateSvc.Generate(redirectURI)
	if err != nil {
		return GetAuthorizationURLResponse{}, "", err
	}

	resp, err := uc.authURLUC.Execute(ctx, GetAuthorizationURLRequest{
		State: state,
	})
	if err != nil {
		return GetAuthorizationURLResponse{}, "", err
	}

	return resp, state, nil
}

func (uc *ExtensionOAuthUseCase) Exchange(ctx context.Context, code, state string) (NotionOAuthResponse, string, error) {
	redirectURI, err := uc.stateSvc.Validate(state)
	if err != nil {
		return NotionOAuthResponse{}, "", err
	}

	resp, err := uc.oauthUC.Execute(ctx, NotionOAuthRequest{
		Code:  code,
		State: state,
	})
	if err != nil {
		return NotionOAuthResponse{}, "", err
	}

	return resp, redirectURI, nil
}

func (uc *ExtensionOAuthUseCase) IsExtensionState(state string) bool {
	_, err := uc.stateSvc.Validate(state)
	return err == nil
}

func (uc *ExtensionOAuthUseCase) ResolveRedirect(state string) (string, error) {
	return uc.stateSvc.Validate(state)
}
