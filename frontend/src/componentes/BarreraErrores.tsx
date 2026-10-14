// Barrera de errores por pantalla (U1 §2.4: «nunca una pantalla en blanco»). Si una pantalla lanza al pintarse, el
// cascarón (barras, franja «Sin conexión», cola y hojas) sigue vivo y se muestra «No se pudo mostrar esta página ·
// Reintentar». Se reinicia sola al cambiar de pantalla (key en Layout) o de subsección o consulta (`clave`).
import { Component, type ErrorInfo, type ReactNode } from 'react'
import { RefreshCw } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { Franja } from './Estados'
import { Boton } from './Boton'
interface Estado { error: unknown }
export class BarreraErrores extends Component<{ children: ReactNode; clave?: string }, Estado> {
  state: Estado = { error: null }
  componentDidUpdate(previas: { clave?: string }) { if (this.state.error && previas.clave !== this.props.clave) this.setState({ error: null }) }
  static getDerivedStateFromError(error: unknown): Estado { return { error: error ?? new Error('sin detalle') } }
  componentDidCatch(error: unknown, info: ErrorInfo) { console.error('Pantalla con error', error, info.componentStack) }
  render() {
    if (!this.state.error) return this.props.children
    return (
      <div data-prueba="barrera-errores">
        <Franja rol="error" alerta accion={<Boton variante="enlace" onClick={() => this.setState({ error: null })}><RefreshCw aria-hidden className="size-5" />{textos.app.reintentar}</Boton>}>{textos.errores.pantalla}</Franja>
      </div>
    )
  }
}
