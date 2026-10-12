// P-00 Iniciar sesión (U1 §3.2): usuario y contraseña; mensaje de bloqueo (429 con retry_at); aviso de sesión vencida.
import { useState, type FormEvent } from 'react'
import { textos } from '../textos/es-MX'
import { Campo } from '../componentes/Campo'
import { Boton } from '../componentes/Boton'
import { Franja } from '../componentes/Estados'
import { ErrorApi, esSinConexion } from '../api/errores'
import { horaCorta } from '../api/formatos'
import { useConexion } from '../api/conexion'
import { useSesion } from './sesion'
export function IniciarSesion() {
  const { entrar, vencida } = useSesion()
  const conectado = useConexion()
  const [usuario, setUsuario] = useState('')
  const [contrasena, setContrasena] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [enviando, setEnviando] = useState(false)
  const alEnviar = async (e: FormEvent) => {
    e.preventDefault()
    setEnviando(true); setError(null)
    try { await entrar(usuario.trim(), contrasena) } catch (err) {
      if (err instanceof ErrorApi && err.estado === 429 && err.reintentarEn) setError(textos.sesion.bloqueo(horaCorta(new Date(err.reintentarEn))))
      else setError(esSinConexion(err) ? textos.sesion.requiereConexion : err instanceof ErrorApi ? err.message : textos.errores.generico)
    } finally { setEnviando(false); setContrasena('') }
  }
  return (
    <main id="contenido" className="mx-auto flex min-h-dvh w-full max-w-sm flex-col justify-center gap-6 px-4">
      <h1 className="text-title-1 font-semibold text-text-default">{textos.app.nombre}</h1>
      {vencida && <Franja rol="info">{textos.sesion.vencida}</Franja>}
      {!conectado && <Franja rol="sinConexion">{textos.sesion.requiereConexion}</Franja>}
      <form onSubmit={alEnviar} className="flex flex-col gap-4" aria-label={textos.sesion.titulo}>
        <h2 className="text-title-3 font-semibold">{textos.sesion.titulo}</h2>
        <Campo etiqueta={textos.sesion.usuario} name="username" autoComplete="username" autoCapitalize="none" required minLength={3} maxLength={64} value={usuario} onChange={(e) => setUsuario(e.target.value)} />
        <Campo etiqueta={textos.sesion.contrasena} name="password" type="password" autoComplete="current-password" required value={contrasena} onChange={(e) => setContrasena(e.target.value)} />
        {error && <Franja rol="error" alerta>{error}</Franja>}
        <Boton type="submit" variante="primario" tam="hoja" disabled={enviando}>{enviando ? textos.sesion.entrando : textos.sesion.entrar}</Boton>
      </form>
    </main>
  )
}
