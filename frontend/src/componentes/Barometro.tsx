// Barómetro (DESIGN §8.9, U3 B-1, K54): una razón contra dos umbrales en un semicírculo con tramos verde, amarillo y rojo.
// Todo llega de C1 (valor, umbrales, escala y estado): el componente no fija 30 ni 40 ni calcula el estado. SVG en línea con
// colores solo de tokens (clases stroke-/fill- de C2, que cambian solas con el tema); trazos con `non-scaling-stroke` para que
// los grosores sean los de los tokens en píxeles (`size.meter.height`, `size.viz.line`) a cualquier ancho. Las marcas y el
// dígito son HTML (no texto SVG) para que su tamaño sea el de la escala tipográfica y no crezca con el dibujo. Sin animación.
import type { ReactNode } from 'react'
import { CircleAlert, CircleCheck, TriangleAlert } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { cn } from './utils'

export type EstadoBarometro = 'normal' | 'warning' | 'exceeded'
export interface PropsBarometro {
  /** Valor entero en % (C1 `ratio_pct`); ausente = sin dato (sin aguja ni valor, dígito «—»). */
  valor?: number
  /** Límite del tramo verde y del amarillo (C1 `warning_pct`, `limit_pct`). */
  aviso: number
  tope: number
  /** Fin de la escala (C1 `scale_max_pct`). */
  escala: number
  /** Estado que manda C1; ausente sin valor. */
  estado?: EstadoBarometro
  /** Nombre accesible del medidor (título del bloque). */
  titulo: string
  /** `aria-valuetext` completo (título, valor, tope y estado), armado por la pantalla con sus textos; sin dato es el nombre de la imagen. */
  textoAccesible: string
  /** Contenido debajo de la insignia (línea secundaria, avisos). */
  children?: ReactNode
  className?: string
}
const ANCHO = 200, CX = 100, CY = 100, R_ESCALA = 86, R_PISTA = 72
const fraccion = (v: number, escala: number) => Math.min(1, Math.max(0, escala > 0 ? v / escala : 0))
/** Punto del semicírculo para la fracción f (0 = izquierda, 1 = derecha), pasando por arriba. */
const punto = (f: number, r: number) => [CX - r * Math.cos(Math.PI * f), CY - r * Math.sin(Math.PI * f)] as const
const arco = (a: number, b: number, r: number) => {
  const [x1, y1] = punto(a, r), [x2, y2] = punto(b, r)
  return `M ${x1.toFixed(2)} ${y1.toFixed(2)} A ${r} ${r} 0 0 1 ${x2.toFixed(2)} ${y2.toFixed(2)}`
}
const TRAMO: Record<EstadoBarometro, string> = { normal: 'stroke-role-success', warning: 'stroke-role-warning', exceeded: 'stroke-role-error' }
const INSIGNIA: Record<EstadoBarometro, string> = { normal: 'bg-role-success-surface text-role-success', warning: 'bg-role-warning-surface text-role-warning', exceeded: 'bg-role-error-surface text-role-error' }
const ICONO: Record<EstadoBarometro, ReactNode> = { normal: <CircleCheck aria-hidden />, warning: <TriangleAlert aria-hidden />, exceeded: <CircleAlert aria-hidden /> }
const GRUESO = '[vector-effect:non-scaling-stroke] [stroke-width:var(--size-meter-height)]'
const FINO = '[vector-effect:non-scaling-stroke]'

export function Barometro({ titulo, valor, aviso, tope, escala, estado, textoAccesible, children, className }: PropsBarometro) {
  const t = textos.barometro
  const fa = fraccion(aviso, escala), ft = fraccion(tope, escala)
  const conValor = valor !== undefined && estado !== undefined
  const fv = conValor ? fraccion(valor, escala) : 0
  const [ax, ay] = punto(fv, R_ESCALA - 4)
  // Accesibilidad (DESIGN §8.9): con valor, `role="meter"` con `aria-valuenow`; sin dato no hay valor que dar y un `meter` sin
  // `aria-valuenow` es inválido (axe `aria-required-attr`), así que el dibujo es una imagen cuyo nombre es el texto accesible
  // («… sin dato, tope 40 %»). PROVISIONAL: bloqueo docs/estado/bloqueos/WF3c-barometro-sin-dato.md (DESIGN §8.9 no dice el rol sin dato).
  /** Marca de umbral: posición del texto «[n] %» fuera del arco, en % del lienzo (HTML encima del SVG). */
  const marca = (f: number) => { const [x, y] = punto(f, R_ESCALA + 12); return { left: `${(x / ANCHO) * 100}%`, top: `${(y / 110) * 100}%` } }
  return (
    <div className={cn('flex flex-col items-center gap-2', className)}>
      <div {...(conValor
        ? { role: 'meter', 'aria-label': titulo, 'aria-valuemin': 0, 'aria-valuemax': escala, 'aria-valuenow': valor, 'aria-valuetext': textoAccesible }
        : { role: 'img', 'aria-label': textoAccesible })}
        data-prueba="barometro" className="relative w-full max-w-80">
        <svg aria-hidden viewBox={`0 0 ${ANCHO} 110`} className="block h-auto w-full overflow-visible forced-colors:[&_path]:stroke-[CanvasText]">
          {/* Arco de escala: tres tramos; separación de 2 px del color de la superficie en cada límite (§8.1). */}
          <path d={arco(0, fa, R_ESCALA)} fill="none" className={cn('stroke-role-success', GRUESO)} data-tramo="normal" />
          <path d={arco(fa, ft, R_ESCALA)} fill="none" className={cn('stroke-role-warning', GRUESO)} data-tramo="warning" />
          <path d={arco(ft, 1, R_ESCALA)} fill="none" className={cn('stroke-role-error', GRUESO)} data-tramo="exceeded" />
          {[fa, ft].map((f) => { const [x1, y1] = punto(f, R_ESCALA - 7), [x2, y2] = punto(f, R_ESCALA + 7); return <line key={`s${f}`} x1={x1} y1={y1} x2={x2} y2={y2} className={cn('stroke-surface-card [stroke-width:var(--size-viz-gap)]', FINO)} /> })}
          {/* Marcas de umbral: 1 px en ui.viz.baseline, hacia afuera del arco. */}
          {[fa, ft].map((f) => { const [x1, y1] = punto(f, R_ESCALA + 5), [x2, y2] = punto(f, R_ESCALA + 9); return <line key={`m${f}`} data-marca={f} x1={x1} y1={y1} x2={x2} y2={y2} className={cn('stroke-viz-baseline [stroke-width:1px]', FINO)} /> })}
          {/* Pista del valor y valor en el color del tramo donde cae. */}
          <path d={arco(0, 1, R_PISTA)} fill="none" className={cn('stroke-viz-track', GRUESO)} />
          {conValor && fv > 0 && <path d={arco(0, fv, R_PISTA)} fill="none" className={cn(TRAMO[estado], GRUESO)} data-prueba="barometro-valor" />}
          {/* Aguja con extremo redondeado y eje circular con anillo de 2 px de la superficie. */}
          {conValor && <>
            <line x1={CX} y1={CY} x2={ax} y2={ay} strokeLinecap="round" data-prueba="barometro-aguja" className={cn('stroke-text-default [stroke-width:var(--size-viz-line)]', FINO)} />
            <circle cx={CX} cy={CY} r={5} className={cn('fill-text-default stroke-surface-card [stroke-width:var(--size-viz-gap)]', FINO)} />
          </>}
        </svg>
        {[aviso, tope].map((n, i) => (
          <span key={`${i}-${n}`} data-prueba="barometro-marca" style={marca(i ? ft : fa)} className="absolute -translate-x-1/2 -translate-y-1/2 whitespace-nowrap text-caption text-text-secondary tabular-nums">{`${n} %`}</span>
        ))}
      </div>
      <div aria-hidden className="flex flex-col items-center">
        <span data-prueba="barometro-digito" className="whitespace-nowrap text-title-1 font-bold text-text-default tabular-nums">{conValor ? `${valor} %` : t.sinDato}</span>
        <span className="text-small text-text-secondary tabular-nums">{t.deTope(tope)}</span>
      </div>
      {conValor && <span className={cn('inline-flex h-6 w-max items-center gap-1 rounded-sm px-2 text-small font-medium [&_svg]:size-5', INSIGNIA[estado])}>{ICONO[estado]}{t.estados[estado]}</span>}
      {children}
    </div>
  )
}
