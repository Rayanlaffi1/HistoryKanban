const etiquettesTypes: { prefixe: string; libelle: string }[] = [
  { prefixe: "image/", libelle: "IMG" },
  { prefixe: "application/pdf", libelle: "PDF" },
  { prefixe: "text/csv", libelle: "CSV" },
  { prefixe: "text/markdown", libelle: "MD" },
  { prefixe: "text/", libelle: "TXT" },
  { prefixe: "application/json", libelle: "JSON" },
  { prefixe: "application/zip", libelle: "ZIP" },
  { prefixe: "application/gzip", libelle: "GZ" },
  { prefixe: "application/x-7z-compressed", libelle: "7Z" },
  { prefixe: "application/vnd.rar", libelle: "RAR" },
  { prefixe: "application/vnd.openxmlformats-officedocument.wordprocessingml", libelle: "DOCX" },
  { prefixe: "application/vnd.openxmlformats-officedocument.spreadsheetml", libelle: "XLSX" },
  { prefixe: "application/vnd.openxmlformats-officedocument.presentationml", libelle: "PPTX" },
  { prefixe: "application/vnd.oasis.opendocument.text", libelle: "ODT" },
  { prefixe: "application/vnd.oasis.opendocument.spreadsheet", libelle: "ODS" },
  { prefixe: "application/msword", libelle: "DOC" },
  { prefixe: "application/vnd.ms-excel", libelle: "XLS" },
]

export function estImage(typecontenu: string) {
  return typecontenu.startsWith("image/")
}

export function libelleType(typecontenu: string, nom: string) {
  for (const entree of etiquettesTypes) {
    if (typecontenu.startsWith(entree.prefixe)) return entree.libelle
  }
  const point = nom.lastIndexOf(".")
  if (point > 0 && point < nom.length - 1) return nom.slice(point + 1).toUpperCase().slice(0, 4)
  return "FIC"
}

export function formaterTaille(octets: number) {
  if (octets < 1024) return `${octets} o`
  if (octets < 1024 * 1024) return `${(octets / 1024).toFixed(0)} ko`
  return `${(octets / (1024 * 1024)).toFixed(1)} Mo`
}

export const typesAcceptes =
  "image/png,image/jpeg,image/gif,image/webp,application/pdf,text/plain,text/csv,text/markdown," +
  "application/json,application/zip,application/gzip,application/x-7z-compressed,application/vnd.rar," +
  "application/msword,application/vnd.ms-excel," +
  "application/vnd.openxmlformats-officedocument.wordprocessingml.document," +
  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet," +
  "application/vnd.openxmlformats-officedocument.presentationml.presentation," +
  "application/vnd.oasis.opendocument.text,application/vnd.oasis.opendocument.spreadsheet"
