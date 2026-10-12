// Variantes del botón (DESIGN.md §7.2): variantes, estados, tamaños y superficies, en claro y oscuro (barra de herramientas «Tema»).
// Hover, foco y presionado se ven con el puntero y el teclado; scripts/contraste-variantes.mjs los provoca y los mide.
import type { Meta, StoryObj } from '@storybook/react-vite'
import { LoaderCircle, Save, Settings2, Trash2 } from 'lucide-react'
import { Boton } from './Boton'
import { textos } from '../textos/es-MX'
import { Etiquetada, Fila, Grupo, Superficies } from '../../.storybook/apoyo'

const meta = {
  id: 'boton',
  title: 'Componentes/Botón',
  component: Boton,
  args: { children: 'Guardar gasto', variante: 'primario', tam: 'normal' },
  argTypes: {
    variante: { control: 'select', options: ['primario', 'secundario', 'fantasma', 'destructivo', 'enlace'] },
    tam: { control: 'inline-radio', options: ['normal', 'hoja', 'icono'] },
    deshabilitadoPor: { control: 'text' },
  },
} satisfies Meta<typeof Boton>
export default meta
type Historia = StoryObj<typeof meta>

const VARIANTES = [
  ['primario', 'Guardar gasto'],
  ['secundario', 'Ver detalle'],
  ['fantasma', textos.app.cancelar],
  ['destructivo', 'Borrar gasto'],
  ['enlace', 'Ir a Créditos'],
] as const

/** Editable con los controles de Storybook. */
export const Editable: Historia = {}

export const Variantes: Historia = {
  name: 'Variantes y estados',
  render: () => (
    <div className="flex flex-col gap-6">
      {VARIANTES.map(([v, texto]) => (
        <Grupo key={v} titulo={v[0].toUpperCase() + v.slice(1)}>
          <Fila>
            <Etiquetada etiqueta="Reposo"><Boton variante={v}>{texto}</Boton></Etiquetada>
            <Etiquetada etiqueta="Con icono">
              <Boton variante={v}>{v === 'destructivo' ? <Trash2 aria-hidden className="size-5" /> : <Save aria-hidden className="size-5" />}{texto}</Boton>
            </Etiquetada>
            <Etiquetada etiqueta="Deshabilitado">
              <Boton variante={v} deshabilitadoPor={textos.conexion.requiere}>{texto}</Boton>
            </Etiquetada>
            <Etiquetada etiqueta="Ocupado">
              <Boton variante={v} aria-busy="true"><LoaderCircle aria-hidden className="size-5" />{texto}</Boton>
            </Etiquetada>
          </Fila>
        </Grupo>
      ))}
    </div>
  ),
}

export const Tamanos: Historia = {
  name: 'Tamaños',
  render: () => (
    <div className="flex max-w-sm flex-col gap-4">
      <Grupo titulo="Normal" nota="36 px con puntero fino y 44 px con puntero grueso (§4.6).">
        <Fila><Boton variante="primario">Guardar gasto</Boton><Boton variante="secundario">Cancelar</Boton></Fila>
      </Grupo>
      <Grupo titulo="Hoja" nota="48 px, a todo el ancho: el primario de una hoja inferior.">
        <Boton variante="primario" tam="hoja">Guardar gasto</Boton>
      </Grupo>
      <Grupo titulo="Icono" nota="Cuadrado de 36 px (fino) o 44 px (grueso); lleva aria-label.">
        <Fila>
          <Boton variante="fantasma" tam="icono" aria-label={textos.app.masOpciones}><Settings2 aria-hidden className="size-5" /></Boton>
          <Boton variante="secundario" tam="icono" aria-label={textos.vistas.filtrar}><Settings2 aria-hidden className="size-5" /></Boton>
          <Boton variante="fantasma" tam="icono" deshabilitadoPor={textos.conexion.requiere} aria-label={textos.app.masOpciones}><Settings2 aria-hidden className="size-5" /></Boton>
        </Fila>
      </Grupo>
    </div>
  ),
}

export const Superficie: Historia = {
  name: 'Sobre página, tarjeta y superficie flotante',
  render: () => (
    <Superficies>
      {() => (
        <Fila>
          {VARIANTES.map(([v, texto]) => <Boton key={v} variante={v}>{texto}</Boton>)}
          <Boton variante="secundario" deshabilitadoPor={textos.conexion.requiere}>Subir ticket</Boton>
        </Fila>
      )}
    </Superficies>
  ),
}
