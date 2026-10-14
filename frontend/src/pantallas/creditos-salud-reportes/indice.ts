// Área Créditos, Salud y Reportes (agente front-creditos-salud-reportes): P-17..P-21 y P-23..P-25.
// El cascarón las descubre solo (src/app/registro.ts); P-22 (registrar peso) es de la base.
import { lazy } from 'react'
export const pantallas = {
  'P-17': lazy(() => import('./P17Lineas')),
  'P-18': lazy(() => import('./P18Credito')),
  'P-19': lazy(() => import('./P19Compra')),
  'P-20': lazy(() => import('./P20Comparador')),
  'P-21': lazy(() => import('./P21Salud')),
  'P-23': lazy(() => import('./P23Centro')),
  'P-24': lazy(() => import('./P24Mensual')),
  'P-25': lazy(() => import('./P25Karen')),
}
