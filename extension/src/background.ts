const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080/api/v1"

type AuthResult = {
  ok: boolean
  token?: string
  error?: string
}

async function startNotionAuth(): Promise<AuthResult> {
  const redirectUrl = chrome.identity.getRedirectURL("notion")
  const authorizeUrl = new URL(`${API_BASE}/auth/notion/authorize-extension`)
  authorizeUrl.searchParams.set("redirect_uri", redirectUrl)

  const authResp = await fetch(authorizeUrl.toString(), { credentials: "include" })
  if (!authResp.ok) {
    return { ok: false, error: "failed_to_get_authorize_url" }
  }

  const { authorization_url: authorizationUrl } = (await authResp.json()) as {
    authorization_url?: string
  }

  if (!authorizationUrl) {
    return { ok: false, error: "missing_authorization_url" }
  }

  const resultUrl = await chrome.identity.launchWebAuthFlow({
    url: authorizationUrl,
    interactive: true
  })

  if (!resultUrl) {
    return { ok: false, error: "missing_result_url" }
  }

  const parsed = new URL(resultUrl)
  const errorParam = parsed.searchParams.get("error")
  if (errorParam) {
    return { ok: false, error: errorParam }
  }

  const jwtToken = parsed.searchParams.get("jwt_token")
  if (!jwtToken) {
    return { ok: false, error: "missing_jwt_token" }
  }

  await chrome.storage.local.set({ pn_jwt: jwtToken })
  return { ok: true, token: jwtToken }
}

chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  if (message?.type === "pn-auth-start") {
    startNotionAuth()
      .then((result) => sendResponse(result))
      .catch((error) => {
        sendResponse({
          ok: false,
          error: error instanceof Error ? error.message : "auth_failed"
        })
      })
    return true
  }
  return false
})
