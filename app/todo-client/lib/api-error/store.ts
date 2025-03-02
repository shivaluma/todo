import { create } from "zustand"

interface ApiErrorTrackerData {
  error: string | null
  requestId: string | null
}

interface State {
  data: ApiErrorTrackerData
}

interface Actions {
  setData: (data: ApiErrorTrackerData) => void
  clearData: () => void
}

export const useApiErrorTrackerStore = create<State & Actions>((set) => ({
  data: {
    error: null,
    requestId: null,
  },
  setData: (data) => set({ data }),
  clearData: () => set({ data: { error: null, requestId: null } }),
}))
