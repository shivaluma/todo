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
  data: {
    token: string
  }
}

export interface RegisterResponse {
  status: string
  message: string
  data: UserInformation
  timestamp: string
  request_id: string
}

interface UserInformation {
  id: string
  username: string
  email: string
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
  email: string
  password: string
}

/**
 * Registration credentials interface
 */
export interface RegisterCredentials {
  fullname: string
  email: string
  password: string
}

/**
 * Auth API response interfaces
 */
export interface AuthResponse {
  data: {
    token: string
  }
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
