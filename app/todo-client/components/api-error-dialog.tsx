"use client"

import { useApiErrorTrackerStore } from "@/lib/api-error/store"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"

export function RequestIdDialog() {
  const data = useApiErrorTrackerStore((state) => state.data)
  const clearData = useApiErrorTrackerStore((state) => state.clearData)

  return (
    <Dialog
      open={data.error !== null}
      onOpenChange={(open) => {
        if (!open) {
          clearData()
        }
      }}
    >
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Sorry, something went wrong</DialogTitle>
          <DialogDescription>
            <ul className="list-none space-y-2">
              <li>
                <p>Error: {data.error}</p>
              </li>
              <li>
                <p>Request ID: {data.requestId}</p>
              </li>
            </ul>
          </DialogDescription>
        </DialogHeader>

        <DialogFooter className="sm:justify-start">
          <DialogClose asChild>
            <Button type="button" variant="secondary">
              Close
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
