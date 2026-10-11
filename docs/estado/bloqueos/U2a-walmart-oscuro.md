# Bloqueo U2a-walmart-oscuro
Contexto: Colorimetría 1.1 (K26). Gasto pasa a 45° y Karen y Bodega Aurrera salen de 50°. Con la L que le toca en su grupo (entre Mercado y Walmart), Bodega Aurrera solo puede separarse de los roles por matiz, y las únicas ventanas a ≥ 20° de todos los roles, de acción crítica y de la marca quedan entre deuda 310° y error 20° (rosa). En oscuro el rosa pierde unos 0.014 de luminancia en gris y, con Walmart en L 0.720, ningún matiz entero cumple el escalonamiento ΔL ≥ 0.10 (L OKLab y gris) de las tiendas: ni con la L 0.830 y croma 0.12 de la versión 1.0 ni con L 0.829 y croma 0.11. El script lo calcula y lo muestra en `docs/design/colorimetria.md`, «Cambios 1.1», puntos 2 y 3.
Qué falta: ¿Aceptas que Walmart en oscuro (`color.store.walmart.dark`) baje de L 0.720 a 0.716 (`#7BA4F0` → `#7AA3EE`, ΔE 0.4, imperceptible) para que Bodega Aurrera quepa, o prefieres otra salida?
Opciones:
1. Aceptar Walmart oscuro en L 0.716 (lo entregado, PROVISIONAL). Consecuencia: cambia, sin diferencia visible, un token que K26 no nombra; el escalonamiento de las tiendas en oscuro queda en ΔL mínimo 0.101 (OKLab y gris), tan justo como en la 1.0 (0.101 y 0.103).
2. Bajar también Mercado en oscuro (L 0.615 → 0.611) y Walmart a 0.714–0.715, con Bodega Aurrera en croma 0.105 y L 0.827. Consecuencia: cambian dos tokens fuera de K26 (ΔE < 1) y el margen sube a unos 0.002; Mercado queda con contraste cercano a 3.0:1, en el mínimo.
3. Relajar solo para las tiendas en oscuro el escalonamiento en gris de 0.10 a 0.095. Consecuencia: Walmart no cambia, pero se aparta de la regla de §7 (ΔL ≥ 0.10 entre todos los pares de cada grupo) y del criterio U2; habría que reescribir ambos.
4. Intercambiar en oscuro la L de Walmart y de Bodega Aurrera (Walmart arriba, en L 0.825). Consecuencia: Walmart pierde croma (0.088, bajo el piso de 0.10) y queda a ΔL 0.04 de la marca 400. No se recomienda.
Impacto: solo `color.store.walmart.dark`, marcado PROVISIONAL en la descripción del token en `docs/contratos/C2-colores.json` y en `docs/design/colorimetria.md` («Cambios 1.1», §7.2 y §12). No toca cifras ni esquema.
Cómo seguí: Entregué la opción 1 marcada PROVISIONAL. El script verifica el escalonamiento, la separación de los roles y la regla Karen–Bodega Aurrera. Cambiar de opción es editar la L de Walmart o de Mercado (TIENDAS) y el diccionario PROVISIONAL en `docs/design/tools/colorimetria.py` y regenerar.

## Respuesta de Axel (puerta 1b, 11-oct-2026)
- **Decisión:** opción **1**: `color.store.walmart.dark` queda en L 0.716 (`#7AA3EE`), FIRME; se quita la marca.
- **Estado:** RESUELTO (K42). Se aplica en: colorimetría 1.2, DESIGN.md (WF-1).
