// Secciones 01 a 07 de Ajustes (U4 parte a). P26Ajustes junta estas con las de las demás áreas (08 a 14).
import { lazy, type ComponentType } from 'react'
import type { PropsPantalla } from '../../app/registro'
export const secciones: Record<string, ComponentType<PropsPantalla>> = {
  '01': lazy(() => import('./S01Cuenta')),
  '02': lazy(() => import('./S02Presupuesto')),
  '03': lazy(() => import('./S03Ingresos')),
  '04': lazy(() => import('./S04Categorias')),
  '05': lazy(() => import('./S05Creditos')),
  '06': lazy(() => import('./S06Tiendas')),
  '07': lazy(() => import('./S07Recetas')),
}
