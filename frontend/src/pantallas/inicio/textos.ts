// Textos visibles de P-01 Inicio (U3 v1.5), en español de México (U5-07). Mismo estilo que src/textos/es-MX.ts.
const red = 'Revisa que estés en la red de casa.'
const dias = (d: number) => (d === 1 ? '1 día' : `${d} días`)
export const MESES_LARGOS = ['enero', 'febrero', 'marzo', 'abril', 'mayo', 'junio', 'julio', 'agosto', 'septiembre', 'octubre', 'noviembre', 'diciembre']
export const textosInicio = {
  acciones: { fotoTicket: 'Foto de ticket', capturarQuincena: 'Capturar quincena', reintentar: 'Reintentar', verTodo: 'Ver todo', verComoTabla: 'Ver como tabla', verComoGrafica: 'Ver como gráfica' },
  comun: { cargadoA: (h: string) => `Cargado a las ${h}`, nuncaCargada: 'Sin conexión · Esta página aún no se ha cargado en este teléfono', sinDato: '—', hogar: 'Hogar', karen: 'Karen', hogarYKaren: (h: string, k?: string) => (k ? `Hogar ${h} · Karen ${k}` : `Hogar ${h}`), activa: 'Activa', lectura: (l: string) => `Lectura ${l}`, provisional: 'Provisional' },
  pago: { pagado: 'Pagado', vencido: 'Vencido', sinMarcar: 'Sin marcar como pagado', venceHoy: 'Vence hoy', venceManana: 'Vence mañana', venceEn: (d: number) => `Vence en ${dias(d)}`, faltaKaren: 'Falta la parte de Karen', faltaHogar: 'Falta la parte del hogar', pendiente: 'Pendiente', esperado: 'Esperado' },
  t1: {
    titulo: 'Mes contra plan', teQuedan: 'te quedan', rebasado: 'rebasado', deDe: (r: string, p: string, n: string) => `${r} de ${p} · ${n}`, faltan: (d: number) => (d === 1 ? 'Falta 1 día del mes' : `Faltan ${d} días del mes`),
    vasEn: (n: number) => `Vas en ${n} %`, rebasaste: 'Rebasaste el plan', accion: 'Ver presupuesto', sinPlan: (m: string) => `Aún no hay presupuesto para ${m}`, llevas: (r: string) => `Llevas ${r} gastado`, accionVacio: 'Capturar presupuesto',
    falla: `No se pudo cargar el presupuesto del mes. ${red}`,
  },
  t2: {
    titulo: 'Flujo de la quincena', delAl: (a: number, b: number, m: string) => `Del ${a} al ${b} ${m}`, alCierre: 'al cierre', faltante: 'faltante', entra: 'Entra', sale: 'Sale', yaSalio: 'Ya salió', porPagar: 'Por pagar',
    esperado: 'Ingreso esperado', accion: 'Ver flujo', faltaIngreso: 'Falta tu ingreso quincenal', accionVacio: 'Capturar ingreso', falla: `No se pudo cargar el flujo de la quincena. ${red}`,
  },
  t3: {
    titulo: 'Próximo pago', variosPagos: (n: number) => `${n} pagos este día`, faltan: (d: number) => (d === 0 ? 'Hoy' : d === 1 ? 'Mañana' : `Faltan ${dias(d)}`), vencioHace: (d: number) => `Sin marcar como pagado · venció hace ${dias(d)}`,
    disponibleHoy: (m: string, l: string) => `Disponible hoy ${m} · Lectura ${l}`, pendienteConfirmar: 'Pendiente confirmar', reposicion: (n: string) => `Reposición de capital pendiente de confirmar con ${n}`, masLineas: (n: number) => `+${n} líneas`,
    accion: 'Ver crédito', sinPagos: 'Sin pagos pendientes', sinCreditos: 'Aún no hay créditos', accionVacio: 'Agregar crédito', falla: `No se pudieron cargar los pagos. ${red}`,
  },
  t4: {
    titulo: 'Comidas de la semana', costo: 'costo de la semana', sinArmar: (n: number) => (n === 1 ? '1 casilla sin armar' : `${n} casillas sin armar`), completa: 'Semana completa',
    persona: (nombre: string, kcal: string, meta: string | null, d: number) => `${nombre} · ${kcal} kcal al día ${meta ? `de ${meta}` : '· sin meta'} · ${d} de 7 días`, personaSinDias: (nombre: string) => `${nombre} · —`,
    accion: 'Ver plan', sinPlan: 'Aún no hay plan para esta semana', accionVacio: 'Armar plan', falla: `No se pudo cargar el plan de comidas. ${red}`,
  },
  t5: {
    titulo: 'Mejor canasta', menosQue: (d: string, t: string) => `${d} menos que ${t}`, menosQueSiguiente: (d: string) => `${d} menos que la siguiente tienda`, cobertura: (c: number, i: number) => `${c} de ${i} artículos con precio`,
    viejos: (n: number) => (n === 1 ? '1 precio viejo' : `${n} precios viejos`), membresia: 'Como si tuvieran membresía', listaVacia: (f: string) => `Tu lista está vacía; se compara la compra del ${f}`,
    accion: 'Comparar canasta', sinPrecios: 'Aún no hay precios para comparar', accionVacio: 'Registrar precios', falla: `No se pudo comparar la canasta. ${red}`,
  },
  t6: {
    titulo: 'Ahorro y metas', ahorrado: 'ahorrado', esteMes: (a: string, p: string) => `Este mes ${a} de ${p}`, meta: (n: string, p: string, f: string) => `${n} · ${p} · ${f}`,
    estados: { in_progress: 'En curso', overdue: 'Vencida', achieved: 'Cumplida' }, accion: 'Ver ahorro', sinAhorro: 'Aún no tienes fondo ni metas de ahorro', accionVacio: 'Crear meta', falla: `No se pudo cargar el ahorro. ${red}`,
  },
  /** Tablero de control (U3 2.0 §2, K54). */
  tablero: {
    region: 'Tablero de control', masDelHogar: 'Más del hogar',
    servidorViejo: 'Disponible cuando el servidor se actualice',
    comparar: (cifra: string, dif: number, m: string) => (dif > 0 ? `${cifra} más que ${m}` : dif < 0 ? `${cifra} menos que ${m}` : `Igual que ${m}`),
    sinDatosDe: (m: string) => `Sin datos de ${m}`,
    k1: {
      titulo: 'Gasto del mes', gastadoEn: (m: string) => `gastado en ${m}`, dePlan: (n: string, p: string, r: string) => `${n} de ${p} · te quedan ${r}`, rebasado: (m: string) => `${m} rebasado`,
      delPlan: (n: string) => `${n} del plan`, efectivo: (m: string) => `Pagado en efectivo ${m}`, gastado: 'Gastado', enEfectivo: 'Efectivo', barra: 'Avance del plan',
      vasEn: (n: number) => `Vas en ${n} %`, rebasaste: 'Rebasaste el plan', rebasadas: (n: number) => (n === 1 ? '1 categoría rebasada' : `${n} categorías rebasadas`),
      sinPlan: (m: string) => `Aún no hay presupuesto para ${m}`, llevas: (r: string) => `Llevas ${r} gastado`,
      accion: 'Ver presupuesto', accionVacio: 'Capturar presupuesto', falla: `No se pudo cargar el gasto del mes. ${red}`,
    },
    k2: {
      titulo: 'Ingreso del mes', deQuincenas: 'de quincenas', extra: (m: string) => `Extra este mes ${m}`, sinExtra: 'Sin extra este mes',
      capturadas: (n: number, m: number) => `${n} de ${m} quincenas capturadas`, esperado: 'Ingreso esperado', faltaIngreso: 'Falta tu ingreso quincenal',
      accion: 'Ver ingresos', accionVacio: 'Capturar ingreso', falla: `No se pudo cargar el ingreso del mes. ${red}`,
    },
    k3: {
      titulo: 'Créditos', deudaViva: 'deuda viva total', pagosDe: (m: string, total: string, n?: number) => `Pagos de ${m}: ${total}${n === undefined ? '' : n === 1 ? ' · 1 pago' : ` · ${n} pagos`}`,
      faltaPagar: (m: string) => `Falta pagar ${m}`, sinPagos: (m: string) => `Sin pagos en ${m}`, sinCreditos: 'Aún no hay créditos', disponibleSinDato: 'Disponible hoy —',
      accion: 'Ver créditos', accionVacio: 'Agregar crédito', falla: `No se pudieron cargar los créditos. ${red}`,
    },
    k4: {
      titulo: 'Pagos próximos', hayMas: 'Hay más pagos en Calendario', sinPagos: 'Sin pagos en los próximos 30 días', sinCreditos: 'Aún no hay créditos',
      fila: (f: string, l: string, m: string) => `${f} · ${l} · ${m}`, accion: 'Ver calendario', accionVacio: 'Agregar crédito', falla: `No se pudieron cargar los pagos. ${red}`,
    },
    b1: {
      titulo: 'Capacidad de endeudamiento', valorTexto: (n: string, tope: number, color: string) => `Capacidad de endeudamiento ${n} de ${tope} %: ${color}`,
      sinValorTexto: (tope: number) => `Capacidad de endeudamiento sin dato, tope ${tope} %`,
      pagosDe: (p: string, i: string) => `Pagos ${p} de un ingreso de ${i}`, teCaben: (m: string) => `Te caben ${m} más al mes antes del rojo`, arriba: (m: string) => `${m} arriba del tope`,
      conKaren: (n: number) => `Con la parte de Karen: ${n} %`, sinCreditos: 'Sin créditos activos', esperado: 'Ingreso esperado', faltaIngreso: 'Falta tu ingreso del mes para calcular la capacidad',
      accion: 'Ver créditos', accionVacio: 'Capturar ingreso', falla: `No se pudo calcular la capacidad de endeudamiento. ${red}`,
    },
    g1: {
      titulo: 'Ingreso, gasto y pagos por mes', subtitulo: (a: string, b: string) => `De ${a} a ${b}`, ingreso: 'Ingreso', gasto: 'Gasto', pagos: 'Pagos de crédito', quincenas: 'Quincenas', extras: 'Extras',
      nota: 'Los pagos de crédito pagan compras que ya cuentan en el gasto de su mes; no se suman al gasto',
      resumen: (m: string, i: string, g: string, p: string) => `En ${m}: ingreso ${i}, gasto ${g} y pagos de crédito ${p}`,
      vacio: 'Aún no hay meses con ingreso ni gasto', accionVacio: 'Capturar ingreso',
    },
    g2: {
      titulo: 'Deuda viva por línea', subtitulo: (f: string) => `Al ${f}`, linea: 'Línea', monto: 'Deuda viva', resumen: (t: string, l: string, m: string) => `Deuda viva hoy: ${t}; la mayor, ${l} con ${m}`, vacio: 'No hay deudas vivas',
    },
    g3: {
      titulo: 'Presupuesto por categoría', subtitulo: (m: string) => `Mes de ${m}`, etiqueta: (c: string, r: string, p: string, n: string) => `${c} · ${r} de ${p} · ${n}`, vasEn: (n: number) => `Vas en ${n} %`, rebasada: 'Rebasada',
      subir: 'Sugerencia: subir presupuesto', recortar: 'Sugerencia: recortar presupuesto',
      motivo: (c: string, n: number, arriba: boolean, m: number) => `${c}: ${n} meses seguidos ${arriba ? 'arriba' : 'abajo'} del plan por más de ${m} %`,
      marca: (n: number) => `${n} %`, verTodas: 'Ver todas', sinPlan: (n: number) => (n === 1 ? '1 categoría con gasto y sin plan' : `${n} categorías con gasto y sin plan`),
      resumen: (m: string, n: number, w: number, x: number) => `En ${m}: ${n} ${n === 1 ? 'categoría' : 'categorías'} arriba de ${w} % y ${x} ${x === 1 ? 'rebasada' : 'rebasadas'}`,
      categoria: 'Categoría', real: 'Real', plan: 'Plan', usado: '% del plan', sugerencia: 'Sugerencia', ninguna: 'Sin sugerencia', abrir: (c: string) => `Ver gastos de ${c}`,
      vacio: 'Aún no hay presupuesto por categoría', accionVacio: 'Capturar presupuesto',
    },
    g6: {
      titulo: 'Tendencia del barómetro', serie: 'Capacidad de endeudamiento', nota: 'Proyección con pagos programados e ingreso esperado', faltaIngreso: '— · Falta ingreso',
      estado: 'Estado', capacidad: 'Capacidad', proyeccion: 'Proyección', si: 'Sí', no: 'No', punto: (m: string, n: number, e: string) => `${m}, ${n} %: ${e}`,
      resumen: (m: string, n: string, e: string, mm: string, nn: string) => `${m}: ${n} (${e}); el mayor de la proyección es ${mm} con ${nn}`,
      resumenSinProyeccion: (m: string, n: string, e: string) => `${m}: ${n} (${e})`,
      vacio: 'Falta tu ingreso para calcular la capacidad', accionVacio: 'Capturar ingreso',
    },
  },
  detalle: {
    filtrar: 'Filtrar', filtros: 'Filtros', periodo: 'Periodo', mes: 'Mes', q1: '1.ª quincena', q2: '2.ª quincena', persona: 'Persona', familia: 'Familia', credito: 'Crédito', categoria: 'Categoría', tienda: 'Tienda',
    familias: { household: 'Hogar', external: 'Karen', joint: 'Conjunto' }, noUsa: (f: string) => `Este bloque no usa el filtro ${f}`, verGastos: 'Ver todos los gastos con estos filtros', listo: 'Listo', limpiar: 'Quitar filtros',
  },
  d1: { titulo: 'Presupuesto por categoría', fila: (r: string, p: string, n: string) => `${r} de ${p} · ${n}`, vacio: 'Aún no hay presupuesto por categoría', accionVacio: 'Capturar presupuesto', falla: `No se pudo cargar el presupuesto por categoría. ${red}` },
  d2: {
    titulo: 'Quincena en curso', ingresos: 'Ingresos', esperado: 'esperado', pagosCredito: 'Pagos de crédito (parte del hogar)', fijos: 'Fijos y venideros', variable: 'Gasto variable estimado', ahorro: 'Ahorro planificado', saldo: 'Saldo de la quincena', acumulado: 'Saldo acumulado',
    vacio: 'Falta tu ingreso quincenal', accionVacio: 'Capturar ingreso', falla: `No se pudo cargar la quincena. ${red}`,
  },
  d3: {
    titulo: 'Pagos y disponible', vencidos: 'Vencidos', proximos: 'Próximos pagos', masVencidos: 'Hay más pagos vencidos', disponible: 'Disponible por línea', minimo: (m: string) => `Mínimo del horizonte ${m}`,
    sinPagos: 'Sin pagos pendientes', sinCreditos: 'Aún no hay créditos', accionVacio: 'Agregar crédito', falla: `No se pudieron cargar los pagos. ${red}`, fallaDisponible: `No se pudo cargar el disponible. ${red}`,
  },
  d4: { titulo: 'Comidas de hoy', fuera: 'Fuera', kcal: (n: string) => `${n} kcal`, vacio: 'Hoy no hay comidas en el plan', accionVacio: 'Armar plan', falla: `No se pudieron cargar las comidas de hoy. ${red}` },
  d5: { titulo: 'Diferencias de precio', viejo: 'viejo', totales: 'Total por tienda', articulo: 'Artículo', gana: 'Gana', diferencia: 'Diferencia', vacio: 'Aún no hay precios para comparar', accionVacio: 'Registrar precios', falla: `No se pudo comparar la canasta. ${red}` },
  d6: {
    titulo: 'Metas', deDe: (a: string, m: string) => `${a} de ${m}`, aportacion: (m: string) => `Aportación necesaria ${m} al mes`, fechaMeta: (f: string) => `Fecha meta ${f}`, avance: (p: string) => `Avance ${p}`,
    vacio: 'Aún no tienes metas de ahorro', accionVacio: 'Crear meta', falla: `No se pudieron cargar las metas. ${red}`,
  },
  d7: {
    titulo: 'Participación de Karen', despensa: (m: string) => `Su parte de la despensa este mes: ${m}`, compras: 'Sus compras a crédito', pagos: 'Sus pagos de este mes', debe: (m: string) => `Lo que debe en la línea: ${m}`,
    nota: 'Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga', comercio: 'Comercio', fecha: 'Fecha', total: 'Total', pagosDe: (h: number, t: number) => `Pagos ${h} de ${t}`, saldo: 'Saldo', monto: 'Su parte', estado: 'Estado',
    accion: 'Estado de cuenta de Karen', vacio: 'Karen no tiene compras ni pagos este mes', falla: `No se pudo cargar la participación de Karen. ${red}`,
  },
  graficas: {
    falla: `No se pudo cargar la gráfica. ${red}`, hoy: 'Hoy', proyeccion: 'Proyección', mes: 'Mes', semana: 'Semana', total: 'Total', todas: 'Todas las categorías', ultimos12: (c: string) => `Últimos 12 meses · ${c}`, sinFiltro: 'Todas',
    t1: { titulo: 'Plan contra real por mes', real: 'Real', plan: 'Plan', resumen: (m: string, r: string, p: string) => `En ${m}: real ${r} de un plan de ${p}`, resumenSinPlan: (m: string, r: string) => `En ${m}: real ${r}; aún sin plan`, vacio: 'Aún no hay presupuesto capturado', accionVacio: 'Capturar presupuesto' },
    t2: {
      titulo: 'Devengado contra efectivo', subtitulo: (a: string, b: string) => `De ${a} a ${b}`, efectivo: 'Efectivo acumulado', diferencia: 'Diferencia acumulada', devengado: 'Devengado acumulado', adelantado: 'Pagado por adelantado', adelantadoEn: (m: string) => `Pagado por adelantado en ${m}`,
      resumen: (u: string, s: string, d: string, e: string) => `Hasta ${u}, con proyección desde ${s}: devengado acumulado ${d} y efectivo acumulado ${e}`,
      conRepetidas: 'Proyección con despensas repetidas', sinRepetidas: 'Proyección sin despensas repetidas',
      columnas: { devengadoMes: 'Devengado', efectivoMes: 'Efectivo', devHogar: 'Devengado del hogar', devKaren: 'Devengado de Karen', efeHogar: 'Efectivo del hogar', efeKaren: 'Efectivo de Karen' },
      vacio: 'Aún no hay gastos ni pagos para comparar', accionVacio: 'Registrar gasto',
    },
    t3: { titulo: 'Deuda viva proyectada', total: 'Deuda viva total', resumen: (a: string, ma: string, b: string, mb: string) => `Deuda viva al cierre de ${a}: ${ma}; al cierre de ${b}: ${mb}`, resumenFiltro: (m: string, partes: string) => `Al cierre de ${m}: ${partes}`, interesesPagados: 'Intereses pagados', interesesPorPagar: 'Intereses por pagar', vacio: 'No hay deudas vivas' },
    t4: {
      tituloA: 'Costo por semana', tituloB: 'kcal al día por persona', subtitulo: 'Últimas 12 semanas', costo: 'Costo', metaDe: (n: string) => `Meta de ${n}`, semana: (a: string) => `Semana ${a}`,
      resumenA: (s: string, m: string) => `Semana ${s}: costo ${m}`, resumenPersona: (n: string, k: string, m: string) => `${n} ${k} kcal al día con meta de ${m}`, resumenB: (s: string, partes: string) => `Semana ${s}: ${partes}`,
      vacio: 'Aún no hay semanas con plan de comidas', accionVacio: 'Armar plan',
    },
  },
  calendario: {
    titulo: 'Calendario', anterior: 'Mes anterior', siguiente: 'Mes siguiente', hoy: 'Hoy', tipos: 'Tipos', mas30: 'Ver 30 días más', masN: (n: number) => `+${n} más`, q1: '1.ª quincena · del 1 al 15', q2: (u: number) => `2.ª quincena · del 16 al ${u}`,
    delAl: (a: string, b: string) => `del ${a} al ${b}`, diaSeleccionado: 'Día seleccionado', sinEventosDia: 'Sin eventos este día', semana: ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom'],
    tiposEvento: { income: 'Ingresos', credit_payment: 'Pagos de crédito', fixed_expense: 'Gastos fijos', upcoming_expense: 'Gastos venideros', card_statement: 'Fechas de corte', planned_grocery: 'Despensa programada', market_visit: 'Visita a mercados', promotion: 'Promociones', savings_goal: 'Metas de ahorro', monthly_report: 'Reporte mensual' },
    vacio: 'No hay eventos en este periodo', sinProgramados: 'Aún no hay pagos ni ingresos programados', accionVacio: 'Capturar ingreso', falla: `No se pudo cargar el calendario. ${red}`,
  },
} as const
