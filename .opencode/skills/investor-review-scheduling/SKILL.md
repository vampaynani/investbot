name: investor-review-scheduling
description: Review scheduling for stock evaluations: event-based dates and provisional cadence.

# Calendario de revisión

Establezca una fecha explícita de próxima revisión completa y una razón para cada ticker evaluado.

## Revisión basada en evento

Base la fecha en el evento que cambie la tesis más cercano y verificable: un reporte de ganancias con fecha, decisión regulatoria, vencimiento de deuda, o una comprobación del killer-switch acercándose. Revise en la fecha del evento o al siguiente día UTC si los resultados llegan después de la sesión. Indique fechas de evento solo cuando estén sourcadas.

## Cadencia provisional sin evento fechado

Sin un evento fechado, use una cadencia provisional basada en la conclusión:

- Dentro de 7 días: riesgo de killer-switch activo o evidencia poco confiable.
- Dentro de 30 días: tesis estable.

Son valores por defecto de flujo de trabajo, no pronósticos. Si no hay datos disponibles, reintente dentro de 1 día y deje la valoración en `N/D` en lugar de afirmar el nivel anterior.

## Revisión anticipada

Un evento material o un movimiento de precio que cambie los supuestos implícitos dispara una revisión completa anticipada. No establezca un umbral arbitrario de movimiento de precio ni infiera un cambio de nivel solo por el precio. Si la fecha de la próxima revisión pasa, marque la tesis como vencida (overdue); no presente el nivel anterior como actual.
