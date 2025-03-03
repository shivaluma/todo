import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

import { serverCookies } from "@/lib/auth/cookies"
import { LoginCredentials } from "@/lib/auth/types"

export async function POST(request: NextRequest) {
  try {
    const credentials: LoginCredentials = await request.json()

    // Call the backend API
    const response = await fetch(`${process.env.API_URL}/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(credentials),
    })

    if (!response.ok) {
      const error = await response.json()
      return NextResponse.json(error, { status: response.status })
    }

    const data = await response.json()

    // Set the JWT token in an HTTP-only cookie
    await serverCookies.setAuthToken(data.data.token)

    // Return user data without sensitive information
    return NextResponse.json(data)
  } catch (error) {
    return NextResponse.json(error, { status: 500 })
  }
}
