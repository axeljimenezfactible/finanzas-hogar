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
/** Propiedad que se puede editar a la vez en varias filas (U4 §0.2 punto 2 y §0.7, `Ctrl+/`). */
export interface PropiedadVarias {
  id: string
  nombre: string
  tipo: 'opcion' | 'texto' | 'dinero' | 'numero' | 'fecha'
  opciones?: { valor: string; texto: string }[]
}
/** Edición en grupo: la pantalla calcula la vista previa SUMADA sobre las filas elegidas y guarda todas o ninguna. */
export interface EditarVarias<T> {
  propiedades: PropiedadVarias[]
  /** «Antes → Después» sumado (p. ej. «Gasto fijo mensual · $1,200.00 → $1,450.00», «Pasan a Reponer · 0 → 3»).
   *  Si falla, el bloque dice «No se pudo calcular el efecto · Reintentar» y no aplica. */
  vistaPrevia: (filas: T[], propiedad: string, valor: string) => Promise<Cambio[]>
  /** Aplica el valor a todas las filas. Si lanza, se muestra el error y el bloque sigue abierto. */
  guardar: (filas: T[], propiedad: string, valor: string) => Promise<void>
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
  /** Con esto, la tabla muestra casillas de selección y `Ctrl+/` (o «Editar seleccionadas») abre la edición en grupo. */
  editarVarias?: EditarVarias<T>
}
