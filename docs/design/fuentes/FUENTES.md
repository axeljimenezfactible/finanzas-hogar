# Fuentes e iconos autoalojados · procedencia

Descargados el 11-oct-2026 por el agente U2b-1-design (WF-1), cada archivo comprimido en una carpeta nueva y temporal de `datos/tmp/` (borrada al terminar). Solo se conservan los archivos que usa el diseño y sus licencias; nada se carga desde servicios de terceros en tiempo de ejecución (protocolo §8). Verificación: `python3 -I docs/design/fuentes/scripts/verificar.py .`

## Inter 4.1 (SIL Open Font License 1.1)

- Origen oficial: publicación v4.1 del repositorio `rsms/inter` en GitHub (16-nov-2024), archivo `Inter-4.1.zip`, 33,707,794 B, sha256 `9883fdd4a49d4fb66bd8177ba6625ef9a64aa45899767dde3d36aa425756b11e`.
- Se conservan `web/InterVariable.woff2`, `web/InterVariable-Italic.woff2`, `InterVariable.ttf`, `InterVariable-Italic.ttf` y `LICENSE.txt`, sin modificar.
- Ejes medidos: `wght` 100–900 y `opsz` 14–32; rasgos `tnum`, `case`, `zero`, `cv01`–`cv14`, `ss01`–`ss08`, entre otros; glifos `+`, `−` (U+2212), `$`, `,` y `.` presentes.

## Lucide 1.55.0, paquete estático `lucide-static` (ISC; iconos derivados de Feather también MIT)

- Origen oficial: registro npm, `lucide-static-1.55.0.tgz` (publicado el 10-oct-2026), 6.6 MB, integridad npm `sha512-h18lunTFsIJWnY0FMRPjXeO9KcAtNi7fDqtqgWHOFbssfgUWsctn/pA7bpiyrd4g1DuBOJ24CsRnPfkYblwm/A==` (coincide con la que publica el registro), sha256 `809e2e759aeecd05ef3abc231395a0ccb6a58169823cbd26a425960b4d282797`.
- Se conservan `sprite.svg` (1,877 símbolos, `viewBox` 24, sin trazo fijo), `LICENSE` y `package.json` (como `package-1.55.0.json`), sin modificar. Los 2,141 SVG sueltos, la fuente de iconos y el código del paquete no se conservan: el sprite los contiene.

## Archivos conservados

| Archivo | Bytes | sha256 |
|---|---|---|
| `inter/web/InterVariable.woff2` | 352240 | `693b77d4f32ee9b8bfc995589b5fad5e99adf2832738661f5402f9978429a8e3` |
| `inter/web/InterVariable-Italic.woff2` | 387976 | `e564f652916db6c139570fefb9524a77c4d48f30c92928de9db19b6b5c7a262a` |
| `inter/ttf/InterVariable.ttf` | 879708 | `4989b125924991b90d05b2d16e0e388c48f7d5bb8b30539bbf9c755278d0ccaf` |
| `inter/ttf/InterVariable-Italic.ttf` | 910252 | `d6f1f6a172d9e588438db9f986fd5cfad7b30f644374080a8a9d4d91e344586f` |
| `inter/LICENSE.txt` | 4380 | `262481e844521b326f5ecd053e59b98c8b2da78c8ee1bdbb6e8174305e54935a` |
| `lucide/sprite.svg` | 523536 | `e52e460b86eb0f3c31c197b504cc33f44d4918fc903fa1cda0cd2446f0b56843` |
| `lucide/LICENSE` | 3208 | `b495047bd93a9b06913511076f504daba17d5bbeb3e0650f3bb53a4220329c57` |
| `lucide/package-1.55.0.json` | 1699 | `9ba8d309b1b3274b45cf1f13f289b0afb0c0d5a1fe69705e1f3606ceabb171cf` |
