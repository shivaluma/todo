import { NextResponse } from "next/server"

import { getAuthToken } from "@/lib/auth/server"

export async function GET() {
  try {
    const token = await getAuthToken()

    if (!token) {
      return NextResponse.json({ message: "Unauthorized" }, { status: 401 })
    }

    // Call the backend API with the token
    const response = await fetch(`${process.env.API_URL}/users/me`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      return NextResponse.json(
        { message: "Failed to fetch user data" },
        { status: response.status }
      )
    }

    const data = await response.json()

    // Return user data
    return NextResponse.json({
      user: {
        id: data.id,
        email: data.email,
        name: data.name,
        roles: data.roles || ["USER"],
      },
    })
  } catch (error) {
    console.error("Error fetching user data:", error)
    return NextResponse.json(
      { message: "Internal server error" },
      { status: 500 }
    )
  }
}
