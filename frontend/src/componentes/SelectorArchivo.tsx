// Selector de archivo con botón propio (U5-07): el <input type="file"> nativo pinta textos del navegador en inglés
// («Choose File», «No file chosen»), así que va oculto y lo abre un Boton con texto en español; el nombre del archivo
// elegido (o «Ningún archivo elegido») se muestra al lado. El input conserva su nombre accesible para las pruebas.
import { useId, useRef, type ReactNode } from 'react'
import { FileUp } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { Boton } from './Boton'
interface Props {
  etiqueta: string
  /** Texto del botón; por omisión «Elegir archivo». */
  boton?: string
  accept: string
  archivo: File | null
  alElegir: (f: File | null) => void
  /** Motivo visible para no dejar elegir (sin conexión, importando…). */
  deshabilitadoPor?: string
  /** Lo que se muestra con un archivo elegido; por omisión, su nombre. */
  detalle?: ReactNode
}
export function SelectorArchivo({ etiqueta, boton = textos.archivo.elegir, accept, archivo, alElegir, deshabilitadoPor, detalle }: Props) {
  const ref = useRef<HTMLInputElement>(null)
  const id = useId()
  return (
    <div role="group" aria-labelledby={id} className="flex flex-col gap-1">
      <span id={id} className="text-small font-medium text-text-default">{etiqueta}</span>
      <div className="flex flex-wrap items-center gap-2">
        <Boton deshabilitadoPor={deshabilitadoPor} onClick={() => ref.current?.click()}><FileUp aria-hidden />{boton}</Boton>
        <span className="min-w-0 break-words text-small text-text-secondary tabular-nums">{archivo ? (detalle ?? archivo.name) : textos.archivo.ninguno}</span>
      </div>
      <input ref={ref} type="file" hidden accept={accept} aria-label={etiqueta} onChange={(e) => { alElegir(e.target.files?.[0] ?? null); e.target.value = '' }} />
    </div>
  )
}
