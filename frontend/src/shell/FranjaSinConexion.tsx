// Franja «Sin conexión con el servidor de casa · n en cola» (G-07, N-317) y panel para revisar la cola.
import { useState } from 'react'
import { textos } from '../textos/es-MX'
import { useConexion } from '../api/conexion'
import { dinero, fechaHora, kilos } from '../api/formatos'
import { quitarDeCola, reintentar, useCola } from '../cola/cola'
import { Franja } from '../componentes/Estados'
import { Boton } from '../componentes/Boton'
import { Hoja } from '../componentes/Hoja'
import { cambiar, useShell } from './estado'
export function FranjaSinConexion() {
  const conectado = useConexion()
  const { total, fallidas, retenidas } = useCola()
  // Con conexión, la franja sigue mientras haya capturas que no se pudieron enviar (rechazadas o con falla del
  // servidor / sesión vencida), con el contador de todo lo que está en cola (U1 §2.4, U5-13 paso 7).
  if (conectado && !fallidas && !retenidas) return null
  const revisar = <Boton variante="enlace" onClick={() => cambiar({ panelCola: true })}>{textos.conexion.revisar}</Boton>
  return (
    <div className="sticky top-0 z-30 px-3 pt-2" data-prueba="franja-conexion">
      <Franja rol={conectado ? 'error' : 'sinConexion'} alerta accion={total > 0 ? revisar : undefined}>
        {conectado ? textos.conexion.noSeEnvio : textos.conexion.sinConexion}
        {total > 0 && <> · <span data-prueba="contador-cola">{textos.conexion.enCola(total)}</span></>}
      </Franja>
    </div>
  )
}
export function PanelCola() {
  const { panelCola } = useShell()
  const { capturas } = useCola()
  const conectado = useConexion()
  const [porDescartar, setPorDescartar] = useState<string | null>(null)
  return (
    <Hoja abierta={panelCola} alCambiar={(a) => cambiar({ panelCola: a })} titulo={textos.cola.titulo}>
      {capturas.length === 0 ? <p className="text-text-secondary">{textos.cola.vacia}</p> : (
        <ul className="flex flex-col gap-2">
          {capturas.map((c) => (
            <li key={c.clave} className="flex flex-col gap-1 rounded-md border border-border-divider p-3 text-small">
              <span className="font-medium tabular-nums">{c.tipo === 'gasto' ? `${textos.cola.gasto} ${dinero(c.cuerpo.amount_cents)}` : `${textos.cola.peso} ${kilos(c.cuerpo.weight_g)}`}</span>
              <span className="text-text-secondary">{textos.cola.capturadoEl(fechaHora(c.cuerpo.captured_at ?? new Date(c.creadoEn).toISOString()))} · {c.error ? `${textos.cola.fallo}: ${c.error}` : textos.cola.pendiente}</span>
              {porDescartar === c.clave ? (
                <div className="flex flex-wrap items-center gap-2"><span>{textos.cola.confirmarDescartar}</span>
                  <Boton variante="destructivo" onClick={() => { void quitarDeCola(c.clave); setPorDescartar(null) }}>{textos.cola.descartar}</Boton>
                  <Boton onClick={() => setPorDescartar(null)}>{textos.app.cancelar}</Boton></div>
              ) : (
                <div className="flex gap-2">
                  <Boton deshabilitadoPor={conectado ? undefined : textos.conexion.requiere} onClick={() => void reintentar(c.clave)}>{textos.cola.reintentar}</Boton>
                  <Boton variante="fantasma" onClick={() => setPorDescartar(c.clave)}>{textos.cola.descartar}</Boton>
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </Hoja>
  )
}
