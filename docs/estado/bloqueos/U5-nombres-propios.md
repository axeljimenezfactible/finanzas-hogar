# Bloqueo U5-nombres-propios

Contexto: WF-3b, `revisor-u5-r1` (12-oct-2026). U5-07 permite en pantalla solo estos nombres propios: Walmart, Bodega Aurrera, Sam's Club, Kueski, Telegram, USDA, Open Food Facts, PROFECO, Notion. En las capturas aparecen además «Claude» y «Ollama» (P-26.11, IA opcional), «Tesseract» (P-26.06, OCR local), las siglas «API», «CSV», «PDF» y «OCR», y las claves internas de categoría «groceries» y «eating_out» (P-26.04, columna «Clave que usan los cálculos»).

Qué falta: saber cuáles se aceptan.
- A: ampliar la lista con Claude, Ollama, Tesseract y las siglas API, CSV, PDF, OCR; las claves internas no se muestran o se muestran en español (lo que la ronda 1 da por hecho: hallazgo alto A-4).
- B: ampliar la lista y además aceptar las claves internas como identificadores técnicos.
- C: no ampliar: cambiar los textos («la IA de Anthropic», «modelo local», «archivo de hoja de cálculo»).

Consecuencia: con A o C, A-4 sigue alto hasta quitar «groceries»/«eating_out» y el «Choose File» del selector de archivo; con B solo queda el selector. Lo decide Axel.
