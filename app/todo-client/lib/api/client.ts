import axios from "axios"

import { useApiErrorTrackerStore } from "@/lib/api-error/store"

/**
 * Axios instance for API requests
 */
const apiClient = axios.create({
  baseURL: "/api",
  headers: {
    "Content-Type": "application/json",
  },
})

/**
 * Request interceptor
 * - Add auth token if available
 * - Add CSRF token if available
 */
apiClient.interceptors.request.use(
  (config) => {
    // Add any request interceptors here
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

/**
 * Response interceptor
 * - Handle common errors
 * - Transform response data
 */
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // Handle common errors
    const errorMessage = error.response?.data?.message || "Something went wrong"

    // Handle specific status codes
    if (error.response?.status === 401) {
      // Handle unauthorized (e.g., redirect to login)
      console.error("Unauthorized access")
    }

    if (error.response?.status >= 500) {
      // Handle server errors
      useApiErrorTrackerStore.getState().setData({
        error: error.response?.data?.message || "Something went wrong",
        requestId: error.response?.data?.request_id || "none",
      })
    }

    return Promise.reject({
      ...error,
      message: errorMessage,
    })
  }
)

export default apiClient
