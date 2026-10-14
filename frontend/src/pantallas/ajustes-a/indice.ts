// Área Ajustes (parte a): P-26 con la columna de las 14 secciones y las secciones 01 a 07 (U4 parte a).
// Las secciones 08 a 14 las aporta otra área desde src/pantallas/<area>/secciones-ajustes.ts (export `secciones`).
import { lazy } from 'react'
export const pantallas = { 'P-26': lazy(() => import('./P26Ajustes')) }
