// Cierre global: junta las partes por captura en evidencia/manifiesto.json y hosts.json (U5-08) y borra los temporales.
import fs from 'node:fs'
import path from 'node:path'
import { EVID, PARTES, SALIDA } from './apoyo'
export default async function cierre() {
  if (!fs.existsSync(PARTES)) return
  fs.mkdirSync(EVID, { recursive: true })
  const leer = (f: string) => { try { return JSON.parse(fs.readFileSync(f, 'utf8')) } catch { return null } }
  const manifiesto: Record<string, unknown> = (leer(path.join(EVID, 'manifiesto.json')) as { capturas: Record<string, unknown> } | null)?.capturas ?? {}
  const previoHosts = leer(path.join(SALIDA, 'hosts.json')) as { por_captura: Record<string, string[]> } | null
  const hosts: Record<string, string[]> = previoHosts?.por_captura ?? {}
  for (const f of fs.readdirSync(PARTES).filter((n) => n.endsWith('.json'))) {
    const p = leer(path.join(PARTES, f)); if (!p) continue
    manifiesto[p.nombre] = p.manifiesto; hosts[p.nombre] = p.hosts
  }
  const unicos = [...new Set(Object.values(hosts).flat())].sort()
  fs.writeFileSync(path.join(EVID, 'manifiesto.json'), JSON.stringify({ capturas: manifiesto }, null, 1))
  fs.writeFileSync(path.join(SALIDA, 'hosts.json'), JSON.stringify({ criterio: 'U5-08', unicos, total_unicos: unicos.length, capturas_con_registro: Object.keys(hosts).length, por_captura: hosts }, null, 1))
  fs.rmSync(PARTES, { recursive: true, force: true })
}
