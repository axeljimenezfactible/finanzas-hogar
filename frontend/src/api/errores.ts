// Errores del cliente (C1 §3: esquema único Error con code, message, fields, request_id, retry_at).
import type { Esquemas } from './tipos'
export type CampoConError = Esquemas['FieldError']
export class ErrorApi extends Error {
  constructor(
    readonly estado: number,
    readonly codigo: string,
    mensaje: string,
    readonly campos: CampoConError[] = [],
    readonly idPeticion: string | null = null,
    readonly reintentarEn: string | null = null,
  ) {
    super(mensaje)
    this.name = 'ErrorApi'
  }
  /** Mensaje de un campo (`amount_cents`, …) para pintarlo junto al control. */
  campo(nombre: string): string | undefined {
    return this.campos.find((c) => c.field === nombre)?.message
  }
}
/** Sin respuesta del servidor de casa (red cortada, servidor apagado o fuera de la LAN). */
export class ErrorSinConexion extends Error {
  constructor() {
    super('sin_conexion')
    this.name = 'ErrorSinConexion'
  }
}
export const esSinConexion = (e: unknown): e is ErrorSinConexion => e instanceof ErrorSinConexion
