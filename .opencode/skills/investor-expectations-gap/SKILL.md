name: investor-expectations-gap
description: Unified protocol for high-multiple (P/E > 20x) justification and contrarian-thesis review. Reuses existing specialists; does not create a new agent.

# Protocolo de brecha de expectativas (expectations gap)

## Propósito

Esta skill define un protocolo unificado que el **investor-evaluator** ejecuta cuando hay que juzgar si un precio elevado o una divergencia con el consenso pueden estar justificados por los fundamentos. No crea un subagente nuevo; reutiliza los especialistas existentes:

- `investor-valuation`: preguntas 1, 3, 4, 5 y escenarios.
- `investor-business-quality`: pregunta 2 (ROIC, reinversión, TAM).
- `investor-risk-killer`: umbral cuantitativo de ruptura (pregunta 5).
- `investor-target-consensus`: lente de expectativas del mercado (rating y precio objetivo).


El resultado se integra en el veredicto final y se registra en el snapshot de evaluación bajo **"valuation assumptions"** y **"killer switch"**.

## Cuándo activarse

El evaluador carga esta skill cuando se cumpla cualquiera de los dos disparadores:

### Modo A — Múltiplo alto (P/E > 20x)
El P/E verificado (trailing, forward o normalizado) supera 20 veces utilidades. Hay que responder si el múltiplo está justificado por el crecimiento futuro.

### Modo B — Tesis contraria (divergencia con el consenso)
Existe una divergencia verificable entre el precio, las expectativas del consenso y los fundamentos de la empresa (por ejemplo, el consenso es muy pesimista o muy optimista frente a una lectura distinta de la calidad del negocio, y la valoración ofrece margen o exige demasiado). La divergencia con el consenso **fuerza más análisis, pero no cambia la recomendación por sí sola**.

## Pasos del protocolo

El evaluador ejecuta los pasos en este orden e integra las respuestas. Este control se añade a los tres escenarios obligatorios de cada evaluación; no los sustituye.

### Paso 1 — investor-valuation responde las preguntas de control
Responde explícitamente, antes del veredicto de valoración:

1. **¿Qué crecimiento del beneficio por acción necesito durante 5 años para ganar un retorno que me compense?** Usa la calculadora local y el requisito de crecimiento inverso (reverse valuation). Indica el retorno si el múltiplo de salida se mantiene en 20x.
2. *(delegada al siguiente paso)* ¿Ese crecimiento está respaldado por ROIC, reinversión y mercado disponible, o es extrapolación del pasado?
3. **¿Qué retorno obtengo si el crecimiento es la mitad y el múltiplo baja a 18-20x?** Calcula el escenario con crecimiento anual a la mitad del caso base y salida a 18-20x. Si el retorno es negativo o cercano a cero, el margen de error es fino.
4. **¿Es el P/E sobre beneficio normalizado o sobre uno en pico?** Declara trailing, forward o normalizado, y si el año de base refleja márgenes cíclicos en máximo, ajustes contables o SBC no reconocida.
5. **¿Cuál es mi señal de que la tesis se rompe?** Define el crecimiento mínimo por debajo del cual el retorno base se anula, y la condición de múltiplo de salida o deterioro del negocio que invalida la tesis (umbral cuantitativo del siguiente paso).

Usa la calculadora local con entradas verificadas para preguntas 1 y 3; nunca calcules mentalmente. Marca `N/D` lo no verificable.

Para la pregunta 3, añade a la sensibilidad del mismo archivo de entrada el crecimiento anual a la mitad del caso base y múltiplos de salida de 18x y 20x. Informa ambos resultados para 3 y 5 años; no ejecutes una corrida separada. Si falta un input verificado, marca ese resultado `N/D`.

### Paso 2 — investor-business-quality valida el crecimiento requerido
Responde si el crecimiento del EPS necesario para justificar el precio está respaldado por ROIC, tasa de reinversión y TAM, o si es solo extrapolación del pasado. No da conclusiones de precio ni de valoración.

### Paso 3 — investor-risk-killer define el umbral de ruptura
Define el crecimiento mínimo del EPS por debajo del cual el retorno base se anula, y la condición de múltiplo de salida o deterioro del negocio que invalida la tesis.

### Paso 4 — investor-target-consensus aporta el lente de expectativas del mercado
- `investor-target-consensus`: rating agregado, desglose, precio objetivo, diferencia implícita vs. precio verificado, limitaciones. Es una expectativa de mercado, no valor intrínseca.
- Factores macro/industriales relevantes: si hay evidencia verificable sobre industria, geopolítica, regulación o cadena de suministro que afecte la tesis, el evaluador la cita directamente; no hay un subagente dedicado a este dominio.

### Paso 5 — El evaluador integra y registra
- Combina las respuestas en un veredicto único basado en el perfil del inversor (horizonte 3–5 años, maximizar retorno total, tolerar volatilidad).
- Registra en el snapshot de evaluación bajo **"valuation assumptions"** y **"killer switch"**: preguntas 1–5, respuestas de cada especialista, escenario de mitad de crecimiento + salida a 18-20x, y umbral de ruptura.
- Si el caso base ofrece solo un retorno de un dígito medio/bajo, la lectura interna del precio es razonable-a-cara, o un test de tesis a corto plazo está sin resolver, usa `Hold` a menos que exista evidencia específica que justifique `Buy`.

## Reglas de integración

- Un múltiplo bajo no prueba que el precio esté por debajo del valor razonable; un múltiplo alto no prueba que esté caro. Siempre hay que validar el crecimiento requerido contra ROIC, reinversión y TAM.
- La divergencia con el consenso es un motivo para investigar, no una señal automática de compra ni de venta.
- No uses el precio objetivo de consenso para cambiar la lectura fundamental del precio. Ni un `Strong Buy` de analistas rescata un caso base poco atractivo.
- Separa hechos de interpretación. No inventes precios, múltiplos, fechas de ganancias ni catalizadores. Si un dato no se puede verificar, usa `N/D`.
- La recomendación final es la síntesis del evaluador para el perfil del inversor, no la calificación de los analistas.

## Registro en el snapshot

Cuando esta skill se active, el snapshot del ticker debe incluir:
- Modo activado (A: P/E > 20x, o B: tesis contraria/divergencia).
- Respuestas a las preguntas 1–5 con fuente y especialista de origen.
- Escenario de mitad de crecimiento y salida a 18-20x (retorno y margen de error).
- Umbral cuantitativo de ruptura (killer switch).
- Cómo se integró en el veredicto final.
