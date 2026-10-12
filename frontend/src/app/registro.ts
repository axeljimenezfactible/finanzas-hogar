// Registro de pantallas por área: cada área exporta `pantallas` desde src/pantallas/<area>/indice.ts
// (componentes con lazy()). Nadie edita este archivo para agregar una pantalla.
import type { ComponentType } from 'react'
import type { Subseccion } from './rutas'
export interface PropsPantalla { params: Record<string, string>; sub: Subseccion | null; consulta: URLSearchParams }
type Indice = { pantallas: Record<string, ComponentType<PropsPantalla>> }
const modulos = import.meta.glob<Indice>('../pantallas/*/indice.ts', { eager: true })
export const PANTALLAS_REGISTRADAS: Record<string, ComponentType<PropsPantalla>> = Object.assign({}, ...Object.values(modulos).map((m) => m.pantallas))
