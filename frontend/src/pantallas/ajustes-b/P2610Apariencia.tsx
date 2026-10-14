// P-26.10 Apariencia (U4 parte b; decisión G, Q2-82 A, K51). Todo es del dispositivo (U1 §0.3): no usa C1 y funciona sin conexión.
import { useState, useSyncExternalStore } from 'react'
import { RadioGroup } from 'radix-ui'
import { Monitor, Moon, Sun } from 'lucide-react'
import { ponerTema, useTema, type Tema } from '../../app/tema'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { dinero } from '../../api/formatos'
import { cn } from '../../componentes/utils'
import { Bloque, Interruptor, Marco } from './comun'
import { CircleCheck } from 'lucide-react'
import { textosAjustesB as T } from './textos'

const t = T.a
/** Preferencia del dispositivo «Atajos de una tecla» (K51). PROPUESTA: el cascarón la lee con leerAtajosUnaTecla() (solicitud a frontend-base). */
export const CLAVE_ATAJOS = 'fh-atajos-una-tecla'
const oyentes = new Set<() => void>()
export function leerAtajosUnaTecla(): boolean { try { return localStorage.getItem(CLAVE_ATAJOS) !== 'no' } catch { return true } }
function ponerAtajos(v: boolean) {
  try { if (v) localStorage.removeItem(CLAVE_ATAJOS); else localStorage.setItem(CLAVE_ATAJOS, 'no') } catch { /* sin almacenamiento */ }
  window.dispatchEvent(new CustomEvent('fh:atajos-una-tecla', { detail: { activos: v } }))
  oyentes.forEach((f) => f())
}
const OPCIONES: { valor: Tema; texto: string; Icono: typeof Sun }[] = [{ valor: 'sistema', texto: t.sistema, Icono: Monitor }, { valor: 'claro', texto: t.claro, Icono: Sun }, { valor: 'oscuro', texto: t.oscuro, Icono: Moon }]

export default function P2610Apariencia() {
  const tema = useTema()
  const atajos = useSyncExternalStore((f) => { oyentes.add(f); return () => oyentes.delete(f) }, leerAtajosUnaTecla)
  const [guardado, setGuardado] = useState(false)
  const marcar = () => { setGuardado(true); setTimeout(() => setGuardado(false), 2000) }
  return (
    <Marco id="10">
      <p className="text-body text-text-secondary">{t.proposito}</p>
      <Bloque titulo={t.tema}>
        <div className="flex flex-col gap-4 desktop:flex-row desktop:items-start">
          <RadioGroup.Root aria-label={t.tema} value={tema} onValueChange={(v) => { ponerTema(v as Tema); marcar() }} className="grid grid-cols-3 gap-2 desktop:w-96">
            {OPCIONES.map(({ valor, texto, Icono }) => (
              <RadioGroup.Item key={valor} value={valor} className={cn('flex min-h-11 items-center justify-center gap-2 rounded-md border border-border-control bg-surface-card px-3 text-small text-text-default hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus data-[state=checked]:border-border-selected data-[state=checked]:bg-surface-selected data-[state=checked]:font-semibold')}>
                <Icono aria-hidden className="size-5" />{texto}
              </RadioGroup.Item>
            ))}
          </RadioGroup.Root>
          <div className="flex-1" aria-label={t.muestra}>
            <TarjetaMetrica titulo={t.muestra} valor={dinero(123456)} nota={tema === 'oscuro' ? t.muestraOscuro : tema === 'claro' ? t.muestraClaro : t.sistema} />
          </div>
        </div>
        <p className="text-caption text-text-secondary">{t.fijos}</p>
      </Bloque>
      <Bloque titulo={t.atajos} ayuda={t.atajosAyuda}>
        <Interruptor etiqueta={t.atajos} activo={atajos} alCambiar={(v) => { ponerAtajos(v); marcar() }} />
      </Bloque>
      <span aria-live="polite" className="inline-flex min-h-5 items-center gap-1 text-caption text-role-success">{guardado && <><CircleCheck aria-hidden className="size-5" />{T.comun.guardado}</>}</span>
    </Marco>
  )
}
