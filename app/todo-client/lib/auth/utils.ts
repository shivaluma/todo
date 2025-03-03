import { headers } from "next/headers"
import { redirect } from "next/navigation"

import { serverCookies } from "@/lib/auth/cookies"
import { JWTPayload, RouteConfig, UserRole } from "@/lib/auth/types"

/**
 * Get the current auth token from the request headers or cookies
 */
export async function getAuthToken(): Promise<string | null> {
  try {
    // Try to get token from headers first (for API routes)
    const headersList = await Promise.resolve(headers())
    const authHeader = headersList.get("Authorization")
    if (authHeader?.startsWith("Bearer ")) {
      return authHeader.substring(7)
    }

    // Fall back to cookies
    const cookieToken = await serverCookies.getAuthToken()
    return cookieToken || null
  } catch {
    return null
  }
}

/**
 * Protect a route based on authentication and authorization requirements
 */
export async function protectRoute(config?: RouteConfig) {
  const token = await getAuthToken()

  // Check authentication
  if (config?.requireAuth && !token) {
    redirect("/login")
  }

  return token
}

/**
 * Verify JWT token structure (does not verify signature)
 */
export function isValidJWT(token: string): boolean {
  try {
    const parts = token.split(".")
    if (parts.length !== 3) return false

    // Check if each part is base64 encoded
    return parts.every((part) => {
      try {
        atob(part)
        return true
      } catch {
        return false
      }
    })
  } catch {
    return false
  }
}

/**
 * Parse JWT payload (without verification)
 */
export function parseJWT(token: string): JWTPayload | null {
  try {
    const base64Payload = token.split(".")[1]
    const payload = atob(base64Payload)
    return JSON.parse(payload)
  } catch {
    return null
  }
}

/**
 * Get user roles from JWT token
 */
export function getRolesFromToken(token: string): UserRole[] {
  const payload = parseJWT(token)
  return payload?.roles || []
}
