# Bloqueo criterios-cobertura
Contexto: S2-15, S3-18, S4-15, S5-13, S6-02: plan.md §3.4 fila S6: «`go test` ≥ 85 % en dominio y motores».
Qué falta: ¿Qué paquetes son «motores»? §3.2 lista `internal/{dominio,motor_fin,motor_cred,recetas,precios,api,auth,jobs,reportes,telegram,respaldo}`; §5 llama «motor» a recetas/comidas/salud y «comparador» a precios. ¿El 85 % aplica a recetas y precios, y a api/auth/jobs/reportes?
Opciones:
- (a) 85 % en dominio, motor_fin, motor_cred, recetas y precios (los 5 de lógica de negocio).
- (b) 85 % solo en dominio, motor_fin y motor_cred; los demás se miden y reportan sin umbral.
- (c) 85 % en todo `internal/`.
Impacto: Los criterios exigen 85 % en dominio, motor_fin y motor_cred (seguro) y miden los demás; recetas y precios quedan con 85 % sujeto a confirmación.
Cómo seguí: S4-15 y S5-13 piden medir y reportar con umbral 85 % sujeto a esta confirmación.
