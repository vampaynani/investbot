name: investor-record-keeping
description: Persistent follow-up record for a stock evaluation under .investor_evaluator/<TICKER>.md.

# Registro de seguimiento persistente

Cuando el evaluador complete una evaluación (inicial o actualización), cree o actualice un único archivo Markdown por ticker bajo `.investor_evaluator/`, relativo a la raíz del proyecto.

## Reglas del archivo

- **Clave del archivo:** el exchange symbol, no el nombre de la empresa. Normalizar a mayúsculas y permitir solo letras, dígitos, `.` y `-`. Si el símbolo o el listing es ambiguo, resuélvalo antes de escribir. No mezcle listings ni clases de acciones.
- **Leer primero:** lea el archivo existente antes de cada evaluación. No reemplace silenciosamente una tesis previa ni una fuente conflictiva. Cree el directorio si falta.
- **No guarde secretos:** nunca almacene API keys ni detalles personales de cartera.
- **No cantidades de cartera:** no guarde quantities, cost basis ni weights en los archivos de ticker bajo `.investor_evaluator/`. No repita detalles exactos de posición en los informes a menos que sea necesario para responder la pregunta del usuario.

## Estructura del snapshot (sección actual)

Cada archivo tiene un snapshot actual y un historial fechado. El snapshot incluye, como mínimo:

- Ticker y exchange
- Fecha de análisis en UTC
- Último chequeo diario en UTC
- Último precio verificado con moneda, fuente URL y timestamp, o `sin dato vigente`
- Calidad empresarial (quality assessment)
- Nivel de valoración interno 1–5 (o `N/D`); consérvelo en el snapshot para consistencia y seguimiento, pero no lo muestre en el reporte dirigido al usuario.
- Recomendación basada en el perfil (profile-based recommendation)
- Rationale y confianza
- Acción acotada (bounded action)
- Thesis
- Assumptions (supuestos de valoración a 3–5 años, incluyendo múltiplos actuales y de escenario cuando sea significativo)
- Catalyst
- Killer switch
- Qué cambiaría la recomendación
- Próxima revisión completa en UTC y por qué se eligió esa fecha

También registre:
- El panorama de negocio a 3–5 años y los supuestos de valoración, incluyendo múltiplos de ganancias actuales y de escenario cuando sea significativo.
- Los retornos de escenario, los ratios de asimetría de escenario, y cualquier probabilidad justificada y retorno esperado ponderado con sus fuentes/supuestos.
- El consenso de analistas por separado de la recomendación: rating definido por el proveedor, reparto de recomendaciones y número de analistas si está disponible, fuente/método y hora de actualización.
- El precio objetivo de consenso por separado: objetivo, moneda, horizonte, método de agregación, número/rango de analistas si está disponible, fuente y hora de actualización, diferencia implícita vs. el precio verificado, interpretación (descuento/paridad/prima/N/D) y confianza.
- Mantenga todos los campos actuales solo cuando sean reverificados.
- La evaluación de narrativas en el snapshot y en cada entrada material del historial: narrativa general del mercado, narrativa de expertos, evidencia/enlaces de fuentes, y confianza. Mantenga esas evaluaciones distintas del nivel de valoración.

## Historial

- Comience el historial con la evaluación inicial.
- Ante un cambio material, apunte la conclusión anterior y la razón de la nueva sin perder la entrada original. Incluya cambios en la recomendación o en cualquiera de los lentes de consenso cuando sea material, y preservar fuentes y supuestos.
- Mantenga el historial en orden cronológico.
- Una segunda ejecución en el mismo día UTC actualiza el chequeo diario en lugar de añadir una entrada duplicada.
- No reutilice un precio caducado como precio de hoy. Si el precio o la evidencia financiera no se pueden verificar, use `N/D` para el nivel de valoración, marque la conclusión como provisional y establezca una fecha de reintento concreta. Las valoraciones existentes permanecen como históricas, no actuales, hasta que se refresquen.

## Verificación de fuentes

- No fabrique precios, cotizaciones, fechas de ganancias, fuentes, catalizadores ni calendarios de revisión.
- Cite la fuente y fecha cerca de cada dato material.
- Separe hechos de interpretación en cada salida.
