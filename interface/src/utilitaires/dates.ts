import dayjs from "dayjs"
import "dayjs/locale/fr"
import relativeTime from "dayjs/plugin/relativeTime"

dayjs.extend(relativeTime)
dayjs.locale("fr")

export function formaterDate(valeur: string | null | undefined): string {
  if (!valeur) return ""
  return dayjs(valeur).format("DD MMM YYYY")
}

export function formaterDateHeure(valeur: string | null | undefined): string {
  if (!valeur) return ""
  return dayjs(valeur).format("DD MMM YYYY HH:mm")
}

export function depuis(valeur: string): string {
  return dayjs(valeur).fromNow()
}

export function estDepassee(valeur: string | null | undefined): boolean {
  if (!valeur) return false
  return dayjs(valeur).isBefore(dayjs())
}

export function versChampDate(valeur: string | null | undefined): string {
  if (!valeur) return ""
  return dayjs(valeur).format("YYYY-MM-DDTHH:mm")
}

export function depuisChampDate(valeur: string): string | null {
  if (!valeur) return null
  return dayjs(valeur).toISOString()
}
