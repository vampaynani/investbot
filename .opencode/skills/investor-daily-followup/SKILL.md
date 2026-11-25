name: investor-daily-followup
description: Daily follow-up for portfolio positions: read positions.json, check evaluations, output compact table.

# Seguimiento diario

Cuando se solicite un seguimiento diario:

## Paso 0 — leer cartera

Lea `.investor_portfolio/positions.json` primero. Si no existe, diga que no hay posiciones registradas y termine sin analizar ningún ticker; no trate la ausencia como cartera vacía ni analice `.investor_evaluator/`.

## Paso 1 — escanear solo posiciones

Escanear solo los tickers listados en el snapshot de cartera; no escanee cada archivo `.investor_evaluator/*.md`. Para cada posición:

- Revise si existe una evaluación guardada en `.investor_evaluator/<TICKER>.md`.
- Si no existe evaluación, reporte que falta evaluación y no asigne recomendación.
- Revise si un desarrollo sourcado que cambie la tesis o un precio verificado requiere una revisión más temprana, y si las fuentes discrepan.
- Ejecute los cuatro agentes especialistas solo cuando la evaluación guardada del ticker exista y esté vencida o haya cambiado materialmente; manténgalos en una sola ronda paralela cuando no haya dependencia entre ellos.
- No copie detalles de cartera a los archivos de ticker.
- Re-ejecutar el seguimiento diario el mismo día UTC no debe crear historial duplicado ni empujar silenciosamente la fecha de revisión.

## Salida

Una tabla compacta por cada posición en `.investor_portfolio/positions.json` con:

- Ticker
- Recomendación basada en el perfil (`Strong Buy`, `Buy`, `Hold`, `Sell`, `Strong Sell`, o `N/D`)
- Contexto de cartera (`held`, `not held` solo para un snapshot completo, o `unknown`)
- Rating agregado de analistas con proveedor, número/desglose y fecha de actualización cuando se verifique (si no, `N/D`)
- Flag de narrativa (`compra potencial contraria`, `señal claramente desfavorable`, `sin flag` o `evidencia insuficiente`)
- Razón de la recomendación
- Cambio desde la última revisión
- Próxima fecha de revisión completa y razón

Incluya los pesos reportados solo con su fecha de snapshot. Mantenga la recomendación distinta del rating de analistas y la comparación de precio objetivo. No muestre la puntuación interna de valoración 1–5; explique la lectura del precio en el motivo de recomendación solo cuando aporte información. Distinga `sin cambios comprobados` de `sin datos nuevos`. La recomendación se basa en el perfil de inversión y la evidencia disponible, no en un tamaño de posición personalizado. Nunca afirme que ocurrió un chequeo en segundo plano automáticamente; este agente se ejecuta solo cuando se invoca.

## Contexto de cartera

- Use `weight_pct` solo cuando el usuario lo haya suministrado, etiquetando su fecha de snapshot. No calcule concentración de cartera a partir de recuentos de acciones por sí solos.
- Trate `coverage: partial` como incompleto. Un ticker faltante no es evidencia de no tener posición. Solo diga que un ticker no se tiene cuando el snapshot está completo y actual. Un snapshot completo vacío significa que el usuario confirmó no tener activos relevantes.
