---
description: Consulta o actualiza el registro local de posiciones para dar contexto a las evaluaciones.
agent: investor-evaluator
---

Administra el snapshot local de posiciones en `.investor_portfolio/positions.json`.

- Si no proporciono posiciones ni una acción concreta, explica qué datos necesito compartir para iniciar el registro: ticker, bolsa/listing, cantidad de acciones, moneda, costo promedio por acción si lo conozco, peso actual si lo conozco, fecha de cada posición y si la lista está completa o parcial. No inventes posiciones ni crees un snapshot vacío.
- Si proporciono posiciones, crea o actualiza el archivo con `schema_version: 1`, `updated_at` en `YYYY-MM-DD`, `coverage` (`complete` o `partial`) y un array `positions`. Cada posición debe tener `ticker`, `exchange`, `currency`, `shares` como decimal, `average_cost` como decimal o `null`, `weight_pct` como porcentaje o `null`, y `as_of` en `YYYY-MM-DD`. Omite `average_cost` o `weight_pct` como `null` si no los conozco; no los estimes.
- Normaliza tickers a mayúsculas y confirma el listing exacto. No mezcles ADRs, clases de acciones ni bolsas. Si la identidad del listing no está clara, pregunta antes de guardar.
- `coverage: complete` significa que incluí todos los activos relevantes para evaluar concentración. Usa `partial` si falta alguno. Un snapshot completo con `positions: []` solo significa que confirmé no tener activos que incluir; no lo infieras por una lista ausente.
- `weight_pct` y `average_cost` son valores que yo proporciono. No los calcules ni los actualices con cotizaciones. No guardes brokers, números de cuenta, credenciales ni otra información personal innecesaria.
- Si pido agregar o actualizar un ticker, conserva los demás registros, actualiza el `as_of` solo de la posición modificada y refresca `updated_at`. Si pido quitarlo, elimina solo ese ticker/listing. Si pido reemplazar el snapshot, confirma la cobertura indicada y reemplaza únicamente el archivo de cartera. No alteres archivos de `.investor_evaluator/`.
- Si pido mostrar el registro, presenta lo guardado y su fecha/cobertura; no consultes cotizaciones ni sugieras operaciones.

El archivo está excluido de Git, pero es texto local sin cifrar. Cuando uses sus datos para evaluar una empresa, su contenido entra en el contexto del modelo. No afirmes que existe sincronización con el broker ni seguimiento automático.
