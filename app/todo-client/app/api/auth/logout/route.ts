import { NextResponse } from "next/server"

import { serverCookies } from "@/lib/auth/cookies"

export async function POST() {
  try {
    // Clear the auth token cookie
    await serverCookies.deleteAuthToken()

    return NextResponse.json({ message: "Logged out successfully" })
  } catch (error) {
    console.error("Logout error:", error)
    return NextResponse.json(
      { message: "Internal server error" },
      { status: 500 }
    )
  }
}
