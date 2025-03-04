import type { Metadata, Viewport } from "next"
import { Be_Vietnam_Pro, Geist_Mono } from "next/font/google"

import "@/app/globals.css"

import { preconnect } from "react-dom"
import { Toaster } from "sonner"

import { RequestIdDialog } from "@/components/api-error-dialog"
import { TailwindIndicator } from "@/components/tailwind-indicator"
import config from "@/app/config"

const beVietnamPro = Be_Vietnam_Pro({
  variable: "--font-be-vietnam-pro",
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  subsets: ["latin", "latin-ext", "vietnamese"],
})

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
})

export const metadata: Metadata = {
  metadataBase: new URL("https://todo.shiro.fit"),
  title: config.seo.todo.title,
  description: config.seo.todo.description,
  openGraph: {
    title: config.seo.todo.title,
    description: config.seo.todo.description,
    images: [
      {
        url: config.seo.todo.og_image,
        width: 1200,
        height: 630,
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: config.seo.todo.title,
    description: config.seo.todo.description,
    images: [config.seo.todo.og_image],
  },
}

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="vi">
      <head>
        <link rel="preconnect" href={config.api.url} />
      </head>
      <body
        className={`${beVietnamPro.className} ${geistMono.variable} antialiased overflow-x-hidden`}
      >
        <RequestIdDialog />
        {children}
        <TailwindIndicator />
        <Toaster />
      </body>
    </html>
  )
}
