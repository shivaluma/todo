import { ResponseCookie } from "next/dist/compiled/@edge-runtime/cookies"
import { RequestCookies } from "next/dist/server/web/spec-extension/cookies"
import { cookies } from "next/headers"

const AUTH_TOKEN_NAME = "auth_token"
const CSRF_TOKEN_NAME = "csrf_token"

interface CookieOptions extends Omit<ResponseCookie, "name" | "value"> {
  maxAge?: number
  expires?: Date
  path?: string
  domain?: string
  secure?: boolean
  httpOnly?: boolean
  sameSite?: "strict" | "lax" | "none"
}

const defaultOptions: CookieOptions = {
  httpOnly: true,
  secure: process.env.NODE_ENV === "production",
  sameSite: "lax",
  path: "/",
  maxAge: 30 * 24 * 60 * 60, // 30 days
}

/**
 * Server-side cookie utilities
 */
export const serverCookies = {
  /**
   * Set an HTTP-only cookie (server-side only)
   */
  async set(name: string, value: string, options: CookieOptions = {}) {
    const cookieStore = await Promise.resolve(cookies())
    const opts = { ...defaultOptions, ...options }
    cookieStore.set({
      name,
      value,
      ...opts,
    })
  },

  /**
   * Get a cookie value (server-side only)
   */
  async get(name: string): Promise<string | undefined> {
    const cookieStore = await Promise.resolve(cookies())
    const cookie = cookieStore.get(name)
    return cookie?.value
  },

  /**
   * Delete a cookie (server-side only)
   */
  async delete(name: string) {
    const cookieStore = await Promise.resolve(cookies())
    cookieStore.delete(name)
  },

  /**
   * Set the authentication token in an HTTP-only cookie
   */
  async setAuthToken(token: string) {
    await this.set(AUTH_TOKEN_NAME, token, {
      ...defaultOptions,
      sameSite: "strict",
    })
  },

  /**
   * Get the authentication token from cookies
   */
  async getAuthToken(): Promise<string | undefined> {
    return this.get(AUTH_TOKEN_NAME)
  },

  /**
   * Delete the authentication token cookie
   */
  async deleteAuthToken() {
    await this.delete(AUTH_TOKEN_NAME)
  },

  /**
   * Set CSRF token
   */
  async setCsrfToken(token: string) {
    await this.set(CSRF_TOKEN_NAME, token, {
      ...defaultOptions,
      httpOnly: false, // Needs to be accessible by JavaScript
    })
  },

  /**
   * Get CSRF token
   */
  async getCsrfToken(): Promise<string | undefined> {
    return this.get(CSRF_TOKEN_NAME)
  },
}

/**
 * Parse cookies from a request
 */
export const parseCookies = (cookies: RequestCookies) => {
  const authToken = cookies.get(AUTH_TOKEN_NAME)?.value
  const csrfToken = cookies.get(CSRF_TOKEN_NAME)?.value

  return {
    authToken,
    csrfToken,
  }
}

/**
 * Generate a random CSRF token
 */
export const generateCsrfToken = (): string => {
  return crypto.randomUUID()
}
