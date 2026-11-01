# Plantilla de reporte de evaluación bursátil

Plantilla de presentación para evaluaciones completas. Está pensada para que una persona u otro LLM pueda encontrar rápidamente la decisión, revisar de dónde salen los datos y comprobar que la recomendación concuerda con los escenarios. Las instrucciones del agente y el registro vigente por ticker siguen siendo la fuente de las reglas y de los datos.

## Reglas de presentación

- Escribe en español claro y profesional. Abre con una conclusión directa; evita esconder la acción al final.
- Separa explícitamente **Hechos** de **Interpretación**. Cita la fuente y fecha cerca de cada dato material.
- Identifica la listing exacta, moneda, fecha UTC del análisis y la última cotización verificable. Si no puede verificarse un dato material, escribe `N/D`; no lo estimes ni presentes una cifra antigua como actual.
- Mantén separados calidad empresarial, nivel fundamental de valoración, recomendación propia, consenso de analistas y precio objetivo de consenso.
- Escenarios son ilustraciones condicionales, no pronósticos ni probabilidades. Muestra supuestos y horizonte; no uses un único múltiplo como argumento de valoración.
- Expón fuentes en conflicto, limitaciones y evidencia que cambiaría la conclusión. No prometas retornos.
- El reporte completo es el análisis de referencia, no una adaptación editorial para publicar. No añadas versiones para redes, newsletter u otros medios salvo que se soliciten aparte.

---

# [Empresa] ([TICKER] — [bolsa/listing])

**Fecha de análisis (UTC):** [AAAA-MM-DD]

**Cotización verificada:** [moneda y precio] al [fecha/hora y zona horaria] — [fuente enlazada]. Si no hay precio vigente: `N/D` y conclusión provisional.

## Veredicto

**[Strong Buy / Buy / Hold / Sell / Strong Sell / N/D] — valoración [1–5 / N/D] — calidad [alta/media/baja/N/D].**

[En 2–4 frases: acción relevante para el perfil y horizonte de 3–5 años, razón principal, confianza y condición que acota el veredicto. Si faltan datos críticos, señalar que es provisional.]

### En una frase

[Qué ofrece el negocio y qué tiene que salir bien al precio actual.]

## 1. Calidad del negocio

**Hechos:** [modelo de negocio, segmentos y tendencias relevantes; ventaja competitiva observable; resultados o evidencia de retención/distribución; gestión y asignación de capital. Enlazar fuentes.]

**Interpretación:** [durabilidad del moat, retornos de reinversión, previsibilidad, capacidad de adaptación y límites de la ventaja.]

**Confianza:** [alta/media/baja] — [razón o dato que falta].

## 2. Valoración y escenarios

**Precio y base financiera:** [precio verificado, moneda, fecha/fuente; beneficio y FCF normalizados; deuda/caja si son materiales; múltiplos actuales y método. Distinguir GAAP/ajustes y explicar normalizaciones.]

**Interpretación del nivel [1–5]:** [por qué el precio ofrece margen de seguridad, refleja razonablemente las perspectivas o exige buenos/extraordinarios resultados.]

| Escenario | Supuestos clave | Múltiplo de salida | Retorno total a 3 años | Retorno total a 5 años |
|---|---|---:|---:|---:|
| Conservador | [crecimiento, márgenes, inversión/caja, acciones] | [x] | [% / anualizado] | [% / anualizado] |
| Base | [crecimiento, márgenes, inversión/caja, acciones] | [x] | [% / anualizado] | [% / anualizado] |
| Optimista | [crecimiento, márgenes, inversión/caja, acciones] | [x] | [% / anualizado] | [% / anualizado] |

**Cómo leerlos:** [qué hace más sensible el resultado —crecimiento, margen, capex, dilución, múltiplo, moneda— y qué no cubre el modelo. Aclara cómo se tratan dividendos y recompras.]

## 3. Consenso de analistas y precio objetivo — lentes separados

**Rating agregado:** [proveedor y categoría exacta; número de analistas; reparto Buy/Strong Buy/Hold/Sell; método y fecha de actualización; fuente.] [Si no es verificable: `N/D`, no reutilizar datos viejos.]

**Objetivo agregado:** [moneda, objetivo, horizonte, media/mediana u otro método, número/rango de analistas, fuente y actualización].

**Comparación con la cotización:** [porcentaje implícito y clasificación descuento/paridad/prima/N/D] solo si listing, moneda, fechas y horizonte son comparables. Explicar que el objetivo es una expectativa de mercado, no valor intrínseco ni retorno total. No usar rating u objetivo para fijar el nivel fundamental o rescatar una recomendación débil.

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

Antes de cerrar, concilia explícitamente:

- acción con precio verificado y nivel de valoración;
- retorno total anualizado del caso base a 3–5 años y sus supuestos;
- pérdida del caso conservador y riesgos de pérdida permanente;
- pruebas de tesis pendientes y concentración disponible;
- diferencias de horizonte entre la recomendación propia (3–5 años), el consenso/objetivo (a menudo 12 meses) y las narrativas.

**Conclusión de consistencia:** [por qué el retorno base compensa —o no— el downside y la incertidumbre. Si valoración es 3–4, el retorno base es solo medio/alto de un dígito, o hay una prueba material sin resolver, normalmente Hold salvo evidencia específica que justifique Buy. No aplicar un umbral mecánico; explicar el juicio.]

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
- [ ] Hechos e interpretación diferenciados y fuentes trazables.
- [ ] Consenso y precio objetivo no mezclados con valor intrínseco ni con la recomendación propia.
- [ ] Escenarios transparentes; retorno base y downside comparados en la prueba de consistencia.
- [ ] Narrativas general y experta separadas; flag aplicado solo si se cumplen sus condiciones.
- [ ] Acción directa, confianza, condiciones de cambio, concentración conocida y revisión futura explicitadas.
- [ ] El reporte termina con la conclusión y fuentes; cualquier adaptación para publicar se solicita por separado mediante el skill editorial.
