export const LINKEDIN_CONFIG = {
  clientId: process.env.VITE_LINKEDIN_CLIENT_ID || '',
  redirectUri: process.env.VITE_LINKEDIN_REDIRECT_URI || 'http://localhost:5173/auth/linkedin/callback',
  scope: 'r_liteprofile r_emailaddress',
  responseType: 'code',
  state: 'random_state_string', // In production, generate this dynamically
};

export const getLinkedInAuthUrl = () => {
  const params = new URLSearchParams({
    response_type: LINKEDIN_CONFIG.responseType,
    client_id: LINKEDIN_CONFIG.clientId,
    redirect_uri: LINKEDIN_CONFIG.redirectUri,
    scope: LINKEDIN_CONFIG.scope,
    state: LINKEDIN_CONFIG.state,
  });

  return `https://www.linkedin.com/oauth/v2/authorization?${params.toString()}`;
}; 