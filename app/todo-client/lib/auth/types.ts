/**
 * User roles for authorization
 */
export enum UserRole {
  USER = "USER",
  ADMIN = "ADMIN",
}

/**
 * User session data (non-sensitive information for client-side)
 */
export interface UserSession {
  id: string
  email: string
  name: string
  roles: UserRole[]
}

/**
 * Authentication state interface
 */
export interface AuthState {
  user: UserSession | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}

/**
 * Login credentials interface
 */
export interface LoginCredentials {
  username_or_email: string
  password: string
}

/**
 * Registration credentials interface
 */
export interface RegisterCredentials {
  username: string
  email: string
  password: string
}

/**
 * Auth API response interfaces
 */
export interface AuthResponse {
  user: UserSession
  token?: string
}

/**
 * Route authorization configuration
 */
export interface RouteConfig {
  requireAuth?: boolean
  roles?: UserRole[]
}

/**
 * JWT payload interface
 */
export interface JWTPayload {
  sub?: string
  email?: string
  name?: string
  roles?: UserRole[]
  exp?: number
  iat?: number
  [key: string]: unknown
}
