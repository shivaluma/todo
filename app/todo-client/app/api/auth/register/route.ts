import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

import { RegisterCredentials } from "@/lib/auth/types"

export async function POST(request: NextRequest) {
  try {
    const credentials: RegisterCredentials = await request.json()

    // Call the backend API
    const response = await fetch(`${process.env.API_URL}/auth/register`, {
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

    // Return user data without sensitive information
    return NextResponse.json({
      user: {
        id: data.id,
        email: data.email,
        name: data.name,
        roles: ["USER"], // Default role for new users
      },
    })
  } catch (error) {
    console.error("Registration error:", error)
    return NextResponse.json(
      { message: "Internal server error" },
      { status: 500 }
    )
  }
}
