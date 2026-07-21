import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function classes(...entrees: ClassValue[]) {
  return twMerge(clsx(entrees))
}
