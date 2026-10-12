// P-26 Ajustes (U1 §2.1 y §2.3, U4 §0.7): columna interna con las 14 secciones (escritorio) o lista (móvil);
// cada sección es una página con su propósito, propiedades y bases de datos en línea.
import { Suspense, type ComponentType } from 'react'
import { ChevronRight } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { SECCIONES_AJUSTES } from '../../app/rutas'
import { Enlace } from '../../app/enrutador'
import { Cargando, Vacio } from '../../componentes/Estados'
import { cn } from '../../componentes/utils'
import { textos } from '../../textos/es-MX'
import { t } from './textos'

type Mod = { secciones: Record<string, ComponentType<PropsPantalla>> }
const SECCIONES: Record<string, ComponentType<PropsPantalla>> = Object.assign({}, ...Object.values(import.meta.glob<Mod>('../*/secciones-ajustes.ts', { eager: true })).map((m) => m.secciones))

export default function P26Ajustes(props: PropsPantalla) {
  const id = props.params.seccion && SECCIONES_AJUSTES.some((s) => s.id === props.params.seccion) ? props.params.seccion : null
  const actual = SECCIONES_AJUSTES.find((s) => s.id === id)
  const Seccion = id ? SECCIONES[id] : undefined
  return (
    <div className="flex flex-col gap-4 desktop:flex-row desktop:items-start desktop:gap-6">
      <nav aria-label={t.indice.titulo} className={cn("flex flex-col desktop:w-60 desktop:shrink-0", id && "hidden desktop:flex")}>
        <ul className="flex flex-col gap-1 rounded-lg border border-border-divider bg-surface-card p-2 desktop:sticky desktop:top-4">
          {SECCIONES_AJUSTES.map((s) => (
            <li key={s.id}>
              <Enlace a={`/ajustes/${s.id}`} aria-current={s.id === id ? 'page' : undefined}
                className={cn('flex min-h-11 items-center gap-2 rounded-md px-3 text-body text-text-default hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus desktop:min-h-9 desktop:text-small', s.id === id && 'bg-surface-selected font-medium')}>
                <span className="w-6 text-caption tabular-nums text-text-secondary">{s.id}</span><span className="flex-1">{s.nombre}</span><ChevronRight aria-hidden className="size-4 text-text-secondary desktop:hidden" />
              </Enlace>
            </li>
          ))}
        </ul>
      </nav>
      {actual ? (
        <article aria-labelledby="ajustes-seccion" className="flex min-w-0 flex-1 flex-col gap-4">
          <Enlace a="/ajustes" className="flex min-h-11 w-fit items-center text-small text-text-link desktop:hidden">{t.indice.volver}</Enlace>
          <header className="flex flex-col gap-1">
            <h2 id="ajustes-seccion" className="text-title-2 text-text-default">{actual.nombre}</h2>
            {t.secciones[actual.id] && <p className="text-body text-text-secondary">{t.secciones[actual.id]}</p>}
          </header>
          {Seccion ? <Suspense fallback={<Cargando />}><Seccion {...props} /></Suspense> : <Vacio texto={textos.pendiente.titulo} />}
        </article>
      ) : <p className="hidden text-body text-text-secondary desktop:block">{t.indice.elige}</p>}
    </div>
  )
}
