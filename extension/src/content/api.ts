const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080/api/v1"

type UserResponse = {
  id: string
  email: string
  name: string
  has_notion_integration: boolean
  notion_workspace_id?: string
}

export async function loadJwt(): Promise<string | undefined> {
  const { pn_jwt } = await chrome.storage.local.get("pn_jwt")
  return pn_jwt
}

export async function getCurrentUser(): Promise<UserResponse | null> {
  const token = await loadJwt()
  if (!token) {
    return null
  }

  const resp = await fetch(`${API_BASE}/users/me`, {
    headers: {
      Authorization: `Bearer ${token}`
    }
  })

  if (!resp.ok) {
    return null
  }

  return (await resp.json()) as UserResponse
}
