import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

import { parseCookies } from "@/lib/auth/cookies"

// Define public routes that don't require authentication
const publicRoutes = ["/login", "/register"]

/**
 * Middleware to protect routes and handle authentication
 */
export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // Allow public routes
  if (publicRoutes.some((route) => pathname.startsWith(route))) {
    return NextResponse.next()
  }

  // Check for auth token
  const { authToken } = parseCookies(request.cookies)

  // If no auth token and trying to access protected route, redirect to login
  if (!authToken) {
    const loginUrl = new URL("/login", request.url)
    loginUrl.searchParams.set("from", pathname)
    return NextResponse.redirect(loginUrl)
  }

  // Add auth token to request headers for API routes
  const requestHeaders = new Headers(request.headers)
  requestHeaders.set("Authorization", `Bearer ${authToken}`)

  // Continue with modified headers
  return NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  })
}

// Configure which routes to run middleware on
export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     * - public routes (/login, /register)
     * - api routes (/api/*)
     */
    "/((?!_next/static|_next/image|favicon.ico|public|login|register|api).*)",
  ],
}
