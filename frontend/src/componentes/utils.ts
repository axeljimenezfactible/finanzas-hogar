import { clsx, type ClassValue } from 'clsx'
import { extendTailwindMerge } from 'tailwind-merge'
// Los tamaños de texto de C2 (text-body, text-small…) no deben chocar con los colores (text-text-default…).
const unir = extendTailwindMerge({ extend: { classGroups: { 'font-size': [{ text: ['caption', 'small', 'body', 'title-1', 'title-2', 'title-3'] }] } } })
export const cn = (...c: ClassValue[]) => unir(clsx(c))
