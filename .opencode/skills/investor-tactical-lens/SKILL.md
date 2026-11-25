name: investor-tactical-lens
description: Short-term (1–3 month) entry/exit and FOMO-risk assessment. Optional tactical lens that never changes the 3–5 year profile-based recommendation or the 1–5 fundamental valuation level.

# Lente táctico (investor-tactical) — uso opcional post-veredicto

## Cuándo cargarlo

Carga esta skill **solo después** de que el evaluador haya emitido su recomendación de largo plazo (`Strong Buy`, `Buy`, `Hold`, `Sell`, `Strong Sell`, o `N/D`) y haya registrado su lectura fundamental del precio. Es un paso opcional: úsalo únicamente cuando el usuario pregunte explícitamente por condiciones de entrada o salida a corto plazo, o cuando quiera evaluar el riesgo de FOMO para un ticker ya evaluado.

## Qué hace

`investor-tactical` evalúa el setup de corto plazo de un único ticker:

- Riesgo de FOMO (`elevated`, `mixed`, `low`, o `insufficient evidence`)
- Acción táctica (`enter`, `wait`, `avoid`, o `N/D`)
- Catalizador verificado (evento con fecha, fuente y timestamp)
- Horizonte (1–3 meses por defecto; más corto solo si hay un evento verificable cercano)
- Condición de entrada, salida e invalidación (observables y con fuente)
- Evidencia: cotización, precio-volumen, catalizador y contexto de mercado, con enlaces y timestamps
- Confianza (`high`, `medium`, `low`, o `N/D`)

## Ejecución

La skill define el límite entre horizontes; la investigación táctica la hace el subagente `investor-tactical`. Cuando el usuario active esta lente, el evaluador lo invoca **después** de cerrar el veredicto fundamental y le pasa el ticker/listing exacto, la pregunta táctica concreta y el horizonte solicitado. El subagente devuelve su salida separada; el evaluador no cambia ni reabre la recomendación 3–5 años por una señal táctica.

## Reglas de integración

- **Nunca cambia** la recomendación de largo plazo ni la lectura fundamental del precio (incluida su puntuación interna).
- **No justifica** ignorar una tesis rota ni comprar una acción que el evaluador marcó como `Sell` o `Strong Sell`.
- Si la recomendación principal es `Hold` o `N/D`, el lente táctico puede decir `wait` o `avoid`, pero no convierte eso en `Buy`.
- La acción táctica es análisis general, no recomendación personalizada de tamaño de posición ni de orden.
- Separar hechos de interpretación. No inventar precios, movimientos, volumen, catalizadores ni fuentes.

## Salida

Un bloque separado, claramente etiquetado como **Lente táctico (1–3 meses)**, distinto del veredicto principal. Incluye FOMO risk, acción táctica, catalizador, horizonte, condiciones de entrada/salida/invalidación, evidencia, confianza y limitaciones.
