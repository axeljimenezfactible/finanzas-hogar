// IndexedDB del teléfono: «cola» (capturas sin conexión, K9) y «lecturas» (última carga de cada GET).
const NOMBRE = 'finanzas-hogar'
let abierta: Promise<IDBDatabase> | null = null
export function bd(): Promise<IDBDatabase> {
  abierta ??= new Promise((ok, mal) => {
    const pet = indexedDB.open(NOMBRE, 1)
    pet.onupgradeneeded = () => {
      const db = pet.result
      if (!db.objectStoreNames.contains('cola')) db.createObjectStore('cola', { keyPath: 'clave' }).createIndex('creadoEn', 'creadoEn')
      if (!db.objectStoreNames.contains('lecturas')) db.createObjectStore('lecturas', { keyPath: 'url' })
    }
    pet.onsuccess = () => ok(pet.result)
    pet.onerror = () => mal(pet.error)
  })
  return abierta
}
const prom = <T>(p: IDBRequest<T>) => new Promise<T>((ok, mal) => { p.onsuccess = () => ok(p.result); p.onerror = () => mal(p.error) })
export async function leer<T>(almacen: 'cola' | 'lecturas', clave: string): Promise<T | undefined> {
  return prom((await bd()).transaction(almacen).objectStore(almacen).get(clave)) as Promise<T | undefined>
}
export async function todos<T>(almacen: 'cola' | 'lecturas'): Promise<T[]> {
  const st = (await bd()).transaction(almacen).objectStore(almacen)
  return prom((almacen === 'cola' ? st.index('creadoEn') : st).getAll()) as Promise<T[]>
}
export async function poner<T>(almacen: 'cola' | 'lecturas', valor: T): Promise<void> {
  const tx = (await bd()).transaction(almacen, 'readwrite')
  tx.objectStore(almacen).put(valor)
  await new Promise<void>((ok, mal) => { tx.oncomplete = () => ok(); tx.onerror = () => mal(tx.error) })
}
export async function quitar(almacen: 'cola' | 'lecturas', clave?: string): Promise<void> {
  const tx = (await bd()).transaction(almacen, 'readwrite')
  if (clave === undefined) tx.objectStore(almacen).clear()
  else tx.objectStore(almacen).delete(clave)
  await new Promise<void>((ok, mal) => { tx.oncomplete = () => ok(); tx.onerror = () => mal(tx.error) })
}
