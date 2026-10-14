// Secciones 08 a 14 de Ajustes (U4 parte b). P26Ajustes (área ajustes-a) junta estas con las 01 a 07.
import { lazy, type ComponentType } from 'react'
import type { PropsPantalla } from '../../app/registro'
export const secciones: Record<string, ComponentType<PropsPantalla>> = {
  '08': lazy(() => import('./P2608Reportes')),
  '09': lazy(() => import('./P2609Notificaciones')),
  '10': lazy(() => import('./P2610Apariencia')),
  '11': lazy(() => import('./P2611Datos')),
  '12': lazy(() => import('./P2612Personas')),
  '13': lazy(() => import('./P2613Salud')),
  '14': lazy(() => import('./P2614Inventario')),
}
