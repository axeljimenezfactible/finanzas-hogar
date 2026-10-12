// Tipos derivados del contrato C1 (src/api/esquema.gen.ts lo genera openapi-typescript; no se edita).
import type { paths, components } from './esquema.gen'
export type { paths, components }
export type Esquemas = components['schemas']
export type Metodo = 'get' | 'post' | 'put' | 'patch' | 'delete'
type Op<P extends keyof paths, M extends Metodo> = NonNullable<paths[P][M]>
export type RutasCon<M extends Metodo> = { [P in keyof paths]: [Op<P, M>] extends [never] ? never : P }[keyof paths]
type Json<R> = R extends { content: { 'application/json': infer J } } ? J : undefined
export type Respuesta<P extends keyof paths, M extends Metodo> = Op<P, M> extends { responses: infer R }
  ? R extends { 200: infer A } ? Json<A> : R extends { 201: infer A } ? Json<A> : R extends { 202: infer A } ? Json<A> : undefined
  : never
export type Cuerpo<P extends keyof paths, M extends Metodo> = Op<P, M> extends { requestBody?: infer B }
  ? NonNullable<B> extends { content: { 'application/json': infer J } } ? J : NonNullable<B> extends { content: { 'multipart/form-data': unknown } } ? FormData : never
  : never
export type Consulta<P extends keyof paths, M extends Metodo> = Op<P, M> extends { parameters: { query?: infer Q } } ? NonNullable<Q> : never
export type ParamsRuta<P extends keyof paths, M extends Metodo> = Op<P, M> extends { parameters: { path?: infer Q } } ? NonNullable<Q> : never
type Encabezados<P extends keyof paths, M extends Metodo> = Op<P, M> extends { parameters: { header?: infer H } } ? NonNullable<H> : never
export type RequiereClave<P extends keyof paths, M extends Metodo> = [Encabezados<P, M>] extends [never] ? false : Encabezados<P, M> extends { 'Idempotency-Key': string } ? true : false
export type ItemsDe<R> = R extends { items: (infer I)[] } ? I : never
