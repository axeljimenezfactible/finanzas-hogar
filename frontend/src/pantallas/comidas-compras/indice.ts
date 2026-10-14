// Área Comidas y Compras: P-10 a P-16 (U1 §2.1). El cascarón las descubre solo (src/app/registro.ts).
import { lazy } from 'react'
export const pantallas = {
  'P-10': lazy(() => import('./P10Plan')),
  'P-11': lazy(() => import('./P11Recetario')),
  'P-12': lazy(() => import('./P12Receta')),
  'P-13': lazy(() => import('./P13Lista')),
  'P-14': lazy(() => import('./P14Canasta')),
  'P-15': lazy(() => import('./P15Productos')),
  'P-16': lazy(() => import('./P16Inventario')),
}
