package auth

import (
	"context"
	"fmt"

	"conformitea/infrastructure/gateway/hydra"
	"conformitea/server/types"
)

func (a *Auth) ProcessConsent(ctx context.Context, req types.ConsentRequest) (types.ConsentResult, error) {
	consentSession, err := a.hydraClient.GetConsentSession(req.ConsentChallenge)
	if err != nil {
		return types.ConsentResult{}, fmt.Errorf("failed to get consent session: %w", err)
	}

	// TODO: check if skip_consent is true

	acceptReq := hydra.HydraPutAcceptConsentRequest{
		GrantScope:               consentSession.RequestedScope,
		GrantAccessTokenAudience: consentSession.RequestedAccessTokenAudience,
		Remember:                 true,
		RememberFor:              3600, // Remember for 1 hour
		Session:                  a.buildConsentSession(consentSession),
	}

	result, err := a.hydraClient.AcceptConsentSession(req.ConsentChallenge, acceptReq)
	if err != nil {
		return types.ConsentResult{}, fmt.Errorf("failed to accept consent session: %w", err)
	}

	return types.ConsentResult{
		RedirectTo: result.RedirectTo,
	}, nil
}

func (a *Auth) buildConsentSession(consentSession *hydra.HydraGetConsentResponse) hydra.HydraConsentSessionTokens {
	// TODO: In the future, fetch user details from database and add proper claims
	// For now, we'll add basic claims based on the subject

	accessTokenClaims := map[string]any{
		"roles":       []string{"user"}, // Default role
		"permissions": []string{"read:profile", "write:profile"},
	}

	idTokenClaims := map[string]interface{}{
		// Email and name will be populated from database in future
		// For now, we'll use the subject as a placeholder
		"sub": consentSession.Subject,
	}

	return hydra.HydraConsentSessionTokens{
		AccessToken: accessTokenClaims,
		IDToken:     idTokenClaims,
	}
}
