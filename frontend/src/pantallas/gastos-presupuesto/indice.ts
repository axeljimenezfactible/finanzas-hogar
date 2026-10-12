// Área Gastos y Presupuesto (WF-3a, agente front-gastos-presupuesto): P-02, P-04, P-05, P-06, P-07, P-08, P-09.
import { lazy } from 'react'
export const pantallas = {
  'P-02': lazy(() => import('./P02Gastos')),
  'P-04': lazy(() => import('./P04Gasto')),
  'P-05': lazy(() => import('./P05Ticket')),
  'P-06': lazy(() => import('./P06Presupuesto')),
  'P-07': lazy(() => import('./P07Flujo')),
  'P-08': lazy(() => import('./P08Ingresos')),
  'P-09': lazy(() => import('./P09Ahorro')),
}
