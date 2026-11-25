# Plantilla de reporte de evaluación bursátil

Plantilla de presentación para evaluaciones completas. La respuesta final al usuario usa primero el formato compacto de resultados de abajo. El informe completo conserva después el detalle necesario para revisar las fuentes, las hipótesis y la consistencia de la recomendación. Las instrucciones del agente y el registro vigente por ticker siguen siendo la fuente de las reglas y de los datos.

## Formato de resultados para la respuesta final

Abre con una frase que indique que la evaluación terminó, el ticker, la recomendación, la calidad, la confianza y la lectura del precio. Menciona dónde quedó guardado el informe. Resume los factores decisivos en viñetas, sin repetir el análisis completo:

```markdown
Evaluación de **[TICKER]** completada. Veredicto: **[acción]**, calidad [alta/media/baja], confianza [alta/media/baja], lectura del precio [barata/atractiva/razonable/cara/muy cara/N/D]. Resumen del informe (guardado en `.investor_evaluator/[TICKER].md`):

- **Precio verificado:** [moneda, precio, listing, fecha/hora y fuente; o N/D].
- **Valoración:** [base financiera normalizada, múltiplo y caveat; estado de control de alto múltiplo cuando corresponda; lectura del margen de seguridad].
- **Escenarios (3y / 5y):** [retorno total y anualizado para pesimista/conservador, base y optimista; ratios de asimetría por horizonte]. Si un resultado no se puede calcular, marca N/D y di por qué.
- **Consenso (lente separado):** [proveedor, rating, número y desglose de analistas, objetivo, diferencia frente al precio, horizonte, fecha y fuente; aclara que es expectativa de mercado, no valor intrínseco].
- **Narrativas:** [narrativa general y narrativa experta por separado, evidencia y flag aplicable; marca evidencia insuficiente si corresponde].
- **Killer switch:** [condición verificable que invalidaría la tesis]. **Próxima prueba:** [evento y fecha, si se conocen].
- **Cartera:** [qué permite afirmar el snapshot sobre posición y concentración, con fecha y cobertura; si falta o es parcial, dilo sin inferir ausencia de posición]. Explica por separado qué implica para acciones existentes y para capital nuevo.

**Próxima revisión completa:** [fecha UTC], después de [evento/razón]. [Añade los disparadores de revisión anticipada.]
```

Usa las cifras verificadas y conservadas en el registro. No rellenes huecos con estimaciones. La respuesta compacta debe incluir los tres escenarios y los ratios de asimetría, pero no necesita duplicar la tabla completa de supuestos del informe. Mantén las recomendaciones, la valoración fundamental, el consenso y las narrativas como lentes distintas. No repitas cantidades, coste base ni pesos de cartera.

El informe completo de referencia debe seguir el formato detallado siguiente; el resumen compacto no lo sustituye.

## Reglas de presentación

- Escribe en español claro y profesional. Abre con una conclusión directa; evita esconder la acción al final.
- Separa explícitamente **Hechos** de **Interpretación**. Cita la fuente y fecha cerca de cada dato material.
- Identifica la listing exacta, moneda, fecha UTC del análisis y la última cotización verificable. Si no puede verificarse un dato material, escribe `N/D`; no lo estimes ni presentes una cifra antigua como actual.
- Mantén separados calidad empresarial, lectura fundamental del precio, recomendación propia, consenso de analistas y precio objetivo de consenso. **No muestres la puntuación interna 1–5 de valoración en el reporte al usuario**; expresa la lectura del precio en lenguaje claro y conserva la puntuación solo en el registro técnico.
- Incluye siempre los tres escenarios (**pesimista (conservador), base y optimista**) para 3 y 5 años, calculados con `tools/investor-metrics` a partir de inputs verificados. La selección de supuestos es juicio analítico; la aritmética debe ser reproducible. Muestra retornos totales acumulados y anualizados, y ratios de asimetría por horizonte cuando correspondan. Escenarios condicionales, no pronósticos. Solo asigna probabilidades cuando haya base defendible y explica su origen; si no, marca el retorno ponderado `N/D`.
- Expón fuentes en conflicto, limitaciones y evidencia que cambiaría la conclusión. No prometas retornos.
- El reporte completo es el análisis de referencia, no una adaptación editorial para publicar. No añadas versiones para redes, newsletter u otros medios salvo que se soliciten aparte.

---

# [Empresa] ([TICKER] — [bolsa/listing])

**Fecha de análisis (UTC):** [AAAA-MM-DD]

**Cotización verificada:** [moneda y precio] al [fecha/hora y zona horaria] — [fuente enlazada]. Si no hay precio vigente: `N/D` y conclusión provisional.

## Veredicto

**[Strong Buy / Buy / Hold / Sell / Strong Sell / N/D] — calidad [alta/media/baja/N/D].**

[En 2–4 frases: acción relevante para el perfil y horizonte de 3–5 años, razón principal, confianza y condición que acota el veredicto. Si faltan datos críticos, señalar que es provisional.]

### En una frase

[Qué ofrece el negocio y qué tiene que salir bien al precio actual.]

### Resumen de escenarios y asimetría (retorno total)

| Caso | Supuestos clave | Múltiplo de salida | 3 años: acumulado / anualizado | 5 años: acumulado / anualizado |
|---|---|---:|---:|---:|
| Pesimista (conservador) | [crecimiento EPS/FCF por acción, margen, dilución/recompras] | [x] | [% / %] | [% / %] |
| Base | [crecimiento EPS/FCF por acción, margen, dilución/recompras] | [x] | [% / %] | [% / %] |
| Optimista | [crecimiento EPS/FCF por acción, margen, dilución/recompras] | [x] | [% / %] | [% / %] |

| Horizonte | Pérdida pesimista / retorno base `|R_p| / R_b` | Retorno optimista / pérdida pesimista `R_o / |R_p|` |
|---|---:|---:|
| 3 años | [x / N/D] | [x / N/D] |
| 5 años | [x / N/D] | [x / N/D] |

**Probabilidad/retorno ponderado:** [solo si hay probabilidades defendibles, con fuente y base; de otro modo `N/D`]. Los ratios comparan magnitudes de escenarios, no probabilidades ni pérdida máxima. Si el precio o los inputs no se pueden verificar, conserva las filas y marca los resultados afectados `N/D`, explicando por qué.

**Base de retorno:** total return = retorno de precio + rendimiento de dividendos acumulado simple, sin reinversión; impuestos, comisiones y FX se excluyen. Si no se verificó el dividendo, indica que los resultados son price-only.

## 1. Calidad del negocio

**Hechos:** [modelo de negocio, segmentos y tendencias relevantes; ventaja competitiva observable; resultados o evidencia de retención/distribución; gestión y asignación de capital. Enlazar fuentes.]

**Interpretación:** [durabilidad del moat, retornos de reinversión, previsibilidad, capacidad de adaptación y límites de la ventaja.]

**Confianza:** [alta/media/baja] — [razón o dato que falta].

## 2. Valoración: base financiera y supuestos

**Precio y base financiera:** [precio verificado, moneda, fecha/fuente; beneficio y FCF normalizados; deuda/caja si son materiales; múltiplos actuales y método. Distinguir GAAP/ajustes y explicar normalizaciones.]

### Contexto de precio (solo si cambia la tesis)

**Hechos:** [Incluya máximo histórico, drawdown, recuperación o comparación con benchmark solo si aporta evidencia material a la tesis y puede verificarse con datos comparables; de lo contrario, omita esta subsección.]

**Interpretación:** [Si se incluye, separe corrección de mercado/compresión de múltiplo de deterioro del negocio. Una caída desde máximos no demuestra por sí sola que la acción esté barata, ni una recuperación pasada garantiza otra futura.]

**Lectura del precio:** [explica en lenguaje claro si el precio ofrece margen de seguridad, refleja razonablemente los fundamentos o exige resultados especialmente buenos. No muestres una puntuación numérica interna.]

**Supuestos y sensibilidad:** [explica brevemente qué mueve más los resultados —crecimiento, margen, capex, dilución, múltiplo, moneda—, qué escenario está más expuesto y qué no cubre el modelo. La tabla de escenarios y ratios aparece una sola vez, en el resumen del veredicto.]

## 3. Consenso de analistas y precio objetivo — lentes separados

**Rating agregado:** [proveedor y categoría exacta; número de analistas; reparto Buy/Strong Buy/Hold/Sell; método y fecha de actualización; fuente.] [Si no es verificable: `N/D`, no reutilizar datos viejos.]

**Objetivo agregado:** [moneda, objetivo, horizonte, media/mediana u otro método, número/rango de analistas, fuente y actualización].

**Comparación con la cotización:** [porcentaje implícito y clasificación descuento/paridad/prima/N/D] solo si listing, moneda, fechas y horizonte son comparables. Explicar que el objetivo es una expectativa de mercado, no valor intrínseco ni retorno total. No usar rating u objetivo para cambiar la lectura fundamental del precio o rescatar una recomendación débil.

## 4. Macro, industria y narrativas

**Contexto material:** [estructura competitiva, demanda, regulación, geopolítica, cadena de suministro y otros factores que puedan cambiar la economía del negocio. Separar hechos de interpretación y enlazar fuentes.]

- **Narrativa general del mercado:** [favorable / mixta / contraria / evidencia insuficiente]. [Evidencia reciente, fuente/fecha y límites.]
- **Narrativa experta:** [favorable / mixta / contraria / evidencia insuficiente]. [Quién sostiene la tesis, expertise relevante verificable, razonamiento, fecha y fuentes; señalar incentivos/conflictos conocidos y desacuerdos creíbles. Un consenso de precio objetivo no es narrativa experta.]
- **Flag narrativa:** [compra potencial contraria / señal claramente desfavorable / sin flag / evidencia insuficiente]. [Por qué corresponde —o no—; no mezclar este flag con la valoración ni convertirlo automáticamente en acción.]

## 5. Riesgos y prueba de tesis

**Riesgos principales:**

1. [Riesgo operativo/competitivo y evidencia.]
2. [Riesgo financiero, valoración o asignación de capital.]
3. [Riesgo de regulación, concentración, suministro, ciberseguridad u otro material.]

**Killer switch:** [hecho o combinación de evidencia que invalidaría la tesis; distinguir señal práctica de umbral oficial o guía empresarial.]

**Qué cambiaría la recomendación:** [condición verificable para mejorarla y condición para degradarla; no usar solo movimiento de precio sin reevaluar supuestos.]

**Concentración:** [contexto de cartera conocido; si no hay datos, decir que no se pudo evaluar y no inventar una asignación.]

## Consistencia de la recomendación

Antes de cerrar, concilia explícitamente, **remitiendo al resumen de escenarios de arriba sin repetir sus cifras**:

- acción con precio verificado y lectura fundamental del precio;
- si los casos base y pesimista (conservador) respaldan la acción elegida y si los ratios de asimetría cambian la lectura del riesgo;
- retorno ponderado solo cuando las probabilidades estén justificadas;
- pruebas de tesis pendientes y concentración disponible;
- diferencias de horizonte entre la recomendación propia (3–5 años), el consenso/objetivo (a menudo 12 meses) y las narrativas.

**Conclusión de consistencia:** [por qué el retorno base compensa —o no— el downside y la incertidumbre, remitiendo a la tabla en vez de repetir resultados. Si el precio parece caro, el retorno base es solo medio/alto de un dígito, o hay una prueba material sin resolver, normalmente `Hold` salvo evidencia específica que justifique `Buy`. No aplicar un umbral mecánico; explicar el juicio.]

## Conclusión acotada y próxima revisión

**Acción:** [qué significa para una posición existente y para capital nuevo; no asumir que el usuario tiene acciones.]

**Cambiaría de opinión si:** [una o dos condiciones concretas y observables.]

**Tesis invalidada si:** [condición más fuerte, si aplica.]

**Próxima revisión completa (UTC):** [AAAA-MM-DD]. **Motivo:** [evento fechado y fuente; o cadencia provisional de 7 días por riesgo/evidencia poco fiable, o 30 días para tesis estable. Si falta información fresca, retry en 1 día y nivel N/D.]

## Fuentes principales

- [Fuente primaria o regulatoria, fecha — enlace]
- [Fuente de cotización / consenso, fecha — enlace]
- [Fuente independiente relevante, fecha — enlace]

---

## Lista de verificación para quien redacta

- [ ] Listing, moneda, precio y fecha comprobados; datos viejos marcados como históricos.
- [ ] Contexto de máximo histórico/drawdown incluido solo si es material y verificable; no hacer investigación histórica rutinaria.
- [ ] Hechos e interpretación diferenciados y fuentes trazables.
- [ ] Consenso y precio objetivo no mezclados con valor intrínseco ni con la recomendación propia.
- [ ] Tres escenarios calculados con la herramienta para 3 y 5 años; ratios de asimetría calculados por horizonte cuando sean matemáticamente aplicables; probabilidades solo con evidencia justificable; resultados comparados en la prueba de consistencia y mostrados una sola vez en el resumen del veredicto.
- [ ] Narrativas general y experta separadas; flag aplicado solo si se cumplen sus condiciones.
- [ ] Respuesta final presentada primero en el formato compacto de resultados; acción directa, confianza, cartera conocida y revisión futura explicitadas.
- [ ] El reporte termina con la conclusión y fuentes; cualquier adaptación para publicar se solicita por separado mediante el skill editorial.
