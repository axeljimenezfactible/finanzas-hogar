// Vistas de base de datos tipo Notion (U1 §2.5–§2.6, DESIGN §7.5). Una «propiedad» describe una columna.
import type { ReactNode } from 'react'
export type TipoVista = 'tabla' | 'tablero' | 'calendario' | 'galeria'
export interface Cambio { etiqueta: string; antes: string; despues: string }
export interface Propiedad<T> {
  id: string
  nombre: string
  /** Valor crudo para ordenar, filtrar y agrupar (centavos, fecha ISO, texto…). */
  valor: (fila: T) => string | number | boolean | null | undefined
  /** Cómo se pinta (por omisión, el valor como texto). Dinero con dinero(), fechas con fecha(). */
  pintar?: (fila: T) => ReactNode
  alinear?: 'derecha'
  /** Texto del grupo para agrupar o para las columnas del tablero. */
  grupo?: (fila: T) => string
  /** Edición en el lugar (U4 §0.2): guarda al salir o con Enter; Esc descarta. */
  editar?: {
    tipo: 'texto' | 'dinero' | 'fecha' | 'numero' | 'opcion'
    opciones?: { valor: string; texto: string }[]
    /** Vista previa «Antes → Después» para campos con efecto numérico; si falla, no guarda. */
    vistaPrevia?: (fila: T, nuevo: string) => Promise<Cambio[]>
    guardar: (fila: T, nuevo: string) => Promise<void>
    textoInicial: (fila: T) => string
  }
}
export interface DefVista<T> {
  /** Clave Local para recordar la vista elegida (U1 §5.2: /tabla, /tablero… son preferencias del dispositivo). */
  clave: string
  filas: T[]
  idFila: (fila: T) => string | number
  propiedades: Propiedad<T>[]
  vistas: TipoVista[]
  /** Propiedad de fecha (ISO) para el calendario. */
  fecha?: (fila: T) => string | null | undefined
  /** Propiedad que da las columnas del tablero; mover una tarjeta llama a moverA. */
  columnasTablero?: { propiedad: string; columnas: string[]; moverA?: (fila: T, columna: string) => Promise<void> }
  titulo: (fila: T) => ReactNode
  alAbrir?: (fila: T) => void
  vacio?: ReactNode
}
