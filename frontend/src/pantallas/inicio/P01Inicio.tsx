// P-01 Inicio (U1 §2.1–§2.2, §3.3; U3 2.0): Resumen (= tablero de control, K54) · Detalle · Tendencias · Calendario. El cascarón pinta título, pestañas y «Registrar gasto».
import { useEffect } from 'react'
import { Camera, Banknote } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { Boton } from '../../componentes/Boton'
import { ProveedorInicio, useInicio } from './comun'
import { esDiaDePago } from './fechas'
import { textosInicio as t } from './textos'
import Tablero from './Tablero'
import Detalle from './Detalle'
import Tendencias from './Tendencias'
import Calendario from './Calendario'
const CLAVE = 'fh:inicio:pestana'
const SUBS = ['resumen', 'detalle', 'tendencias', 'calendario']

export default function P01Inicio({ sub, consulta }: PropsPantalla) {
  // La pestaña elegida se recuerda en el dispositivo (U3 §1.1, frecuencia «Local»).
  useEffect(() => {
    try {
      if (!consulta.has('sub')) { const g = localStorage.getItem(CLAVE); if (g && g !== 'resumen' && SUBS.includes(g)) navegar(conConsulta({ sub: g }), { reemplazar: true }) }
      else if (sub) localStorage.setItem(CLAVE, sub)
    } catch { /* sin almacenamiento local: abre en Resumen */ }
  }, [sub, consulta])
  return (
    <ProveedorInicio>
      <div className="flex flex-col gap-4">
        <Acciones />
        {sub === 'detalle' ? <Detalle consulta={consulta} /> : sub === 'tendencias' ? <Tendencias /> : sub === 'calendario' ? <Calendario /> : <Tablero />}
      </div>
    </ProveedorInicio>
  )
}
/** Acciones secundarias de P-01 (U1 §2.2): «Foto de ticket»; «Capturar quincena» solo el día de pago si falta (ingreso aún esperado). */
function Acciones() {
  const c = useInicio()
  const falta = esDiaDePago(c.hoy) && c.flujo.datos?.income_is_expected === true
  return (
    <div className="flex flex-wrap justify-end gap-2">
      <Boton variante="secundario" onClick={() => navegar(rutaDe('P-05'))}><Camera aria-hidden className="size-5" />{t.acciones.fotoTicket}</Boton>
      {falta && <Boton variante="secundario" onClick={() => navegar(`${rutaDe('P-08')}?hoja=quincena`)}><Banknote aria-hidden className="size-5" />{t.acciones.capturarQuincena}</Boton>}
    </div>
  )
}
