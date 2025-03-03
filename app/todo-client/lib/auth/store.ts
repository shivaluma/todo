import { create } from "zustand"
import { persist, PersistOptions } from "zustand/middleware"

import authService from "@/lib/api/auth-service"
import {
  AuthState,
  LoginCredentials,
  RegisterCredentials,
  UserSession,
} from "@/lib/auth/types"

// Define the store actions
interface AuthStoreActions {
  setUser: (user: UserSession | null) => void
  setError: (error: string | null) => void
  clearState: () => void
  login: (credentials: LoginCredentials) => Promise<void>
  register: (credentials: RegisterCredentials) => Promise<void>
  logout: () => Promise<void>
  fetchCurrentUser: () => Promise<void>
}

// Combine state and actions
type AuthStore = AuthState & AuthStoreActions

const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
}

type PersistedState = Pick<AuthState, "user" | "isAuthenticated">

const persistConfig: PersistOptions<AuthStore, PersistedState> = {
  name: "auth-storage",
  partialize: (state) => ({
    user: state.user,
    isAuthenticated: state.isAuthenticated,
  }),
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      ...initialState,

      setUser: (user) => {
        set({
          user,
          isAuthenticated: !!user,
          error: null,
        })
      },

      setError: (error) => {
        set({ error })
      },

      clearState: () => {
        set({ ...initialState })
      },

      login: async (credentials) => {
        try {
          set({ isLoading: true, error: null })

          const user = await authService.login(credentials)

          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          })
        } catch (error) {
          set({
            isLoading: false,
            error: error instanceof Error ? error.message : "Login failed",
          })
          throw error
        }
      },

      register: async (credentials) => {
        try {
          set({ isLoading: true, error: null })

          await authService.register(credentials)
        } catch (error) {
          set({
            isLoading: false,
            error:
              error instanceof Error ? error.message : "Registration failed",
          })
          throw error
        }
      },

      logout: async () => {
        try {
          set({ isLoading: true })

          await authService.logout()

          set({ ...initialState })
        } catch (error) {
          set({
            isLoading: false,
            error: error instanceof Error ? error.message : "Logout failed",
          })
          throw error
        }
      },

      fetchCurrentUser: async () => {
        try {
          set({ isLoading: true, error: null })

          const user = await authService.getCurrentUser()

          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          })
        } catch {
          set({
            isLoading: false,
            error: null, // Don't set error for this case
          })
          // Don't throw error for this case
        }
      },
    }),
    persistConfig
  )
)
