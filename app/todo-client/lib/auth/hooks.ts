import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useSession } from "next-auth/react"

import { UserRole } from "@/lib/auth/types"

declare module "next-auth" {
  interface User {
    id: string
    email: string
    name: string
    roles: string[]
  }

  interface Session {
    user: User
  }
}

/**
 * Custom hook to protect routes and handle auth state
 */
export function useAuth(requiredRoles?: UserRole[]) {
  const { data: session, status } = useSession()
  const router = useRouter()

  useEffect(() => {
    // If not authenticated, redirect to login
    if (status === "unauthenticated") {
      router.push("/login")
      return
    }

    // If authenticated but missing required roles, redirect to home
    if (status === "authenticated" && requiredRoles?.length) {
      const userRoles = session?.user?.roles || []
      const hasRequiredRole = requiredRoles.some((role) =>
        userRoles.includes(role.toString())
      )

      if (!hasRequiredRole) {
        router.push("/")
      }
    }
  }, [status, session, requiredRoles, router])

  return {
    session,
    status,
    isLoading: status === "loading",
    isAuthenticated: status === "authenticated",
  }
}

/**
 * Custom hook to handle auth-protected API requests
 */
export function useAuthenticatedFetch() {
  const { data: session } = useSession()
  const router = useRouter()

  return async (input: RequestInfo | URL, init?: RequestInit) => {
    if (!session) {
      throw new Error("Not authenticated")
    }

    const response = await fetch(input, {
      ...init,
      headers: {
        ...init?.headers,
        "Content-Type": "application/json",
      },
    })

    if (response.status === 401) {
      // Handle token expiration
      router.push("/login")
      throw new Error("Session expired")
    }

    return response
  }
}

/**
 * Custom hook to check if user has specific roles
 */
export function useHasRole(roles: UserRole[]) {
  const { data: session } = useSession()
  const userRoles = session?.user?.roles || []

  return roles.some((role) => userRoles.includes(role.toString()))
}
