import {
  AuthResponse,
  LoginCredentials,
  RegisterCredentials,
  UserSession,
  type RegisterResponse,
} from "@/lib/auth/types"

import apiClient from "./client"

/**
 * Authentication service
 * Handles all authentication-related API calls
 */
const authService = {
  /**
   * Login user
   * @param credentials User login credentials
   * @returns User session data
   */
  async login(credentials: LoginCredentials): Promise<UserSession> {
    const response = await apiClient.post<AuthResponse>(
      "/auth/login",
      credentials
    )
    return response.data
  },

  /**
   * Register new user
   * @param credentials User registration data
   * @returns User session data
   */
  async register(credentials: RegisterCredentials): Promise<RegisterResponse> {
    const response = await apiClient.post<RegisterResponse>(
      "/auth/register",
      credentials
    )
    return response.data
  },

  /**
   * Logout user
   */
  async logout(): Promise<void> {
    await apiClient.post("/auth/logout")
  },

  /**
   * Get current user data
   * @returns User session data
   */
  async getCurrentUser(): Promise<UserSession> {
    const response = await apiClient.get<{ user: UserSession }>("/auth/me")
    return response.data.user
  },
}

export default authService
