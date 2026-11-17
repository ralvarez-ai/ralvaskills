---
status: active
project: ralvaskills
date: 2026-11-17
related: [designs/2026-05-20-technical-specification.md, plans/2026-10-02-mcp-server-plan.md, guides/registry-flow.md]
---

# Permisos por skill — declarar, aplicar al instalar y reconciliar al actualizar

Una skill declara en su frontmatter los permisos de herramienta que necesita (`allowed-tools`), Claude Code ya los respeta al activarla, y `rsk` los traduce a los clientes que no leen ese campo, los reconcilia al actualizar y retira lo que concedió al desinstalar. Hoy eso es un comando manual (`rsk tools allow`) que el usuario tiene que recordar por proyecto, por skill y por cliente.

Estado: fases 1-4 implementadas (2026-11-17); fase 5 pendiente, en plan aparte.
Origen: conversación sobre `notify-on-finish` (2026-11-17): la skill funciona, pero cada lanzamiento pide aprobación hasta que alguien añade la regla a mano.

---

## 1. Contexto y problema

`rsk tools list|allow|deny|remove` (`cmd/rsk/claude.go:79`) ya escribe `.claude/settings.json` con `ReadPermissions`/`WritePermissions` (`internal/tool/claude.go:189,208`). Lo que no existe:

1. **La skill no declara el permiso que necesita.** El agente lo descubre cuando la primera llamada pide aprobación, y el usuario lo arregla a mano. Con 67 skills y 10 clientes, eso no escala.
2. **Nada se aplica al instalar** en los clientes que no leen el frontmatter. `rsk install` enlaza, escribe `rsk.mod`/`rsk.lock` y regenera los pines; los permisos quedan fuera.
3. **Nada se reconcilia al actualizar.** Si `notify-on-finish@0.2.0` añade una entrada a `allowed-tools`, el usuario no se entera; si la retira, queda una concesión huérfana.
4. **Solo Claude.** `rsk tools --for` responde "other tools are not supported yet", y OpenCode usa otro formato (`permissions: [{action, resource, effect}]`, gana la última regla).
5. **La instalación desde `rsk mcp`** (`install_skills`) tampoco puede dejar el proyecto en un estado usable sin un paso manual.

**Hallazgo que reduce el alcance:** `allowed-tools` ya es un campo del frontmatter de Agent Skills y Claude Code lo aplica como pre-aprobación del turno, con `${CLAUDE_SKILL_DIR}` sustituido también dentro de las reglas `Bash(...)`. Para Claude no hay que escribir configuración: el permiso viaja en el tarball de la skill. Lo que falta es el resto de clientes, la memoria de lo concedido y la reconciliación.

**Verificado el 2026-11-17** con `claude -p --permission-mode manual --permission-prompts none --allowedTools Skill` sobre `notify-on-finish`:

| Cuerpo de la skill | Comando que ejecutó el agente | Resultado |
|---|---|---|
| `scripts/notify-on-finish …` (relativo) | `.claude/skills/notify-on-finish/scripts/notify-on-finish …` | **denegado**: la regla expande a ruta absoluta y no casa |
| `${CLAUDE_SKILL_DIR}/scripts/notify-on-finish …` | `/home/…/.claude/skills/notify-on-finish/scripts/notify-on-finish …` | **allowed**, id devuelto, `state=done code=0 duration=12s` |

El control negativo (`touch /tmp/…`) fue denegado en ambas ejecuciones: la concesión es estrecha, no un bypass.

**Verificado en OpenCode v2.1.4** (`opencode run --format json`): la skill se carga con el tool `skill`, la shell ejecuta el script en una sola llamada y devuelve el id (`state=done code=0 duration=8s`), sin prompt. OpenCode **no lee `allowed-tools`** —su frontmatter admite `name`, `description`, `metadata` y `opencode/autoinvoke`— así que corre por su propia política de permisos. La fase 2 (traducción y escritura) solo hace falta cuando el proyecto declara `{action:"shell", resource:"*", effect:"ask"}` o `deny`; con la política por defecto no hay nada que escribir.

## 2. Objetivo / no-objetivo

**Objetivo.** Un contrato declarativo por skill (`allowed-tools`), aplicado al instalar en los clientes que lo necesitan, reconciliado al actualizar y retirado al desinstalar, sin conceder nada en silencio.

**No-objetivo.** Ejecutar herramientas por la skill, conceder permisos a nivel de sistema (sudo, sandbox del SO), ni sustituir el sistema de permisos del cliente: `rsk` solo escribe la regla que el cliente ya entiende.

## 3. Decisión: declarar en `allowed-tools`, recordar lo que `rsk` escribió

| Pieza | Dónde | Para qué |
|---|---|---|
| **Declaración** | `SKILL.md` frontmatter, campo `allowed-tools` | Qué pide la skill, versionado con ella y publicado en el tarball |
| **Concesión del cliente** | Claude: el propio frontmatter · OpenCode: `opencode.json` | Lo que el cliente lee |
| **Registro** | `rsk.lock` (proyecto) y ledger de máquina (global) | Qué escribió `rsk`, a nombre de quién y en qué versión — solo para los clientes donde `rsk` escribe |

Claude Code no necesita la tercera pieza: la concesión es del skill y dura el turno. El registro existe para OpenCode —donde `rsk` sí escribe— y para poder revocar y detectar drift sin tocar reglas del usuario.

## 4. Esquema declarativo

`allowed-tools` es el campo estándar. Formato recomendado, lista YAML (el string separado por espacios es ambiguo cuando la regla lleva espacios dentro del paréntesis):

```yaml
---
name: notify-on-finish
version: 0.1.0
allowed-tools:
  - Bash(${CLAUDE_SKILL_DIR}/scripts/notify-on-finish *)
  - Bash(notify-on-finish *)
description: ...
---
```

| Aspecto | Regla |
|---|---|
| Formas aceptadas | Lista YAML o string separado por comas; `rsk` normaliza a lista |
| Sintaxis de cada entrada | `Tool(specifier)` del cliente propietario (`Bash(…)`, `Read(…)`, `Write(…)`, `WebFetch(…)`, `mcp__server__tool`) |
| `${CLAUDE_SKILL_DIR}` | Se sustituye en la regla **y** en el cuerpo; en OpenCode se resuelve a la ruta real del skill al instalar |
| Coherencia cuerpo ↔ regla | El cuerpo debe invocar el script con `${CLAUDE_SKILL_DIR}`; con una ruta relativa la regla no casa y la llamada se deniega (verificado, §1) |
| Vacío o ausente | Sin cambios en la configuración; una skill sin permisos sigue funcionando con prompts |
| Prohibido | `Bash(*)`, `Bash(**)`, comodín inicial suelto (`Bash(*foo *)`), entradas sin paréntesis |
| Límite | 5 entradas por skill; 200 caracteres por entrada |
| `disallowed-tools` | Fuera de alcance: una skill no restringe el entorno del usuario |

**Mapeo por cliente:**

| Declarado | Claude Code | OpenCode V2 |
|---|---|---|
| `Bash(<cmd> *)` | se usa tal cual (lo aplica el cliente) | `{action:"shell", resource:"<cmd> *", effect:"allow"}` |
| `Read(<glob>)` | se usa tal cual | `{action:"read", resource:"<glob>", effect:"allow"}` |
| `Write(<glob>)` | se usa tal cual | `{action:"edit", resource:"<glob>", effect:"allow"}` |
| `WebFetch(domain:<host>)` | se usa tal cual | `{action:"webfetch", resource:"<host>", effect:"allow"}` |
| `mcp__<server>__<tool>` | se usa tal cual | *sin equivalente, se omite con aviso* |

OpenCode aplica **la última regla que casa**, así que las concesiones se insertan al final del array y nunca se reordenan las del usuario. Los otros ocho adaptadores exponen `ErrUnsupported` y no escriben nada.

## 5. Dónde se aplica (matriz de eventos)

| Evento | Entrada de código | Comportamiento |
|---|---|---|
| `rsk install` (proyecto) | `cmd/rsk/install.go:219`, `install.Apply` (`internal/install/apply.go:45`) | Claude: nada que escribir. OpenCode: muestra los permisos en el preview (`printInstallPreview`, `:285`) y los escribe tras confirmar |
| `rsk install` (global) | `cmd/rsk/install.go:183` | Igual, contra `opencode.json` global y ledger de máquina |
| `rsk new` / `rsk init` | `cmd/rsk/new.go`, `init.go` | Mismo camino que install para las skills que enlaza |
| `rsk update` | `cmd/rsk/update.go:166` (local), `:360` (registry) | **Reconcilia**: diff entre lo declarado por la versión nueva y lo concedido; pide permiso solo si cambió |
| `rsk pin` | `cmd/rsk/pin.go` | No toca permisos (fijar no instala) |
| `rsk sync` | `cmd/rsk/sync.go` | Repara drift: reaplica lo registrado si falta, avisa si sobra |
| `rsk uninstall` / `destroy` | `cmd/rsk/uninstall.go`, `destroy.go` | Revoca solo las reglas registradas a nombre de esa skill |
| `rsk mcp` → `install_skills` | `internal/mcp` | Devuelve los permisos en el resultado; exige interacción y nunca los escribe sin aprobación |
| `rsk status` | `cmd/rsk/status.go` | Columna de drift: declarado y no concedido, concedido y ya no declarado |

Puntos de extensión concretos:

- `internal/skill/skill.go:5` — `AllowedTools []string` en el struct `Skill`.
- `internal/skill/registry.go:117` — el parser de frontmatter (hoy solo `version:` y `description:`, línea a línea) gana `allowed-tools` en sus dos formas y la validación.
- `internal/tool/tool.go:18` — `GrantPermissions(projectDir string, rules []Rule) ([]string, error)`, `RevokePermissions(projectDir string, rules []string) error` y `SupportsPermissions() bool`.
- `internal/tool/claude.go` — `SupportsPermissions() == false`: el cliente ya lo aplica. `ReadPermissions`/`WritePermissions` quedan para `rsk tools`.
- `internal/tool/opencode.go:74` — traducción y merge al final del array `permissions`, conservando el resto del archivo.
- `internal/install/apply.go:33` — `Result.Permissions []Grant` para que el comando informe y el ledger se escriba en el mismo paso que el lock.

## 6. Algoritmo de reconciliación (instalar y actualizar)

```
declarado := allowed-tools(versión nueva)
concedido := ledger.skills[<skill>].granted      # reglas que escribió rsk
preexistente := configActual − concedido         # reglas del usuario

nuevas    := declarado − concedido − preexistente
sobrantes := concedido − declarado
iguales   := declarado ∩ concedido
```

- **nuevas** → se muestran con la entrada y el cliente; se escriben si el usuario acepta (o con `--apply-permissions`).
- **sobrantes** → no se borran solas: se ofrecen (`--revoke-removed`) y, si nadie responde, se quedan y `rsk status` las marca.
- **iguales** → no-op silencioso; idempotente, sin reescrituras.
- **preexistente** → se registra como "ya estaba" y no se toca nunca, ni al desinstalar.

En Claude el diff solo se **informa** (el cliente aplica la versión que tiene enlazada, así que un `update` ya cambia el permiso efectivo); en OpenCode es donde se escribe.

**Conteo de referencias.** `ledger.rules[<regla>] = [skills...]`. Una regla compartida por dos skills se escribe una vez y solo se revoca cuando se queda sin dueños.

**Al instalar** (`install.Apply`): mismo algoritmo con `concedido` vacío — todo lo declarado son `nuevas`, y solo en clientes donde `rsk` escribe.

**Al actualizar** (`rsk update`): el diff se hace contra **la versión del lock**, no contra la instalada en disco, para que un lock desincronizado no se interprete como "el usuario retiró el permiso". Si el lock no tiene entrada (proyecto anterior a este cambio), se asume vacío y todo lo declarado es `nuevas`; una sola vez, y se avisa.

### Formato del registro

Recomendación: `rsk.lock` gana una tabla por skill, porque el lock **ya** es el registro de estado resuelto y el diff natural es "versión vieja → versión nueva":

```toml
[[skills]]
  name = "notify-on-finish"
  version = "0.1.0"
  source = "registry"
  granted = ["shell:*notify-on-finish*"]

[rule-owners]
  "shell:*notify-on-finish*" = ["notify-on-finish"]
```

Alternativa descartada: fichero nuevo `.rsk/permissions.toml` — un artefacto más que sincronizar y un `schema.Version` que subir en `internal/manifest/mod.go:27`. Se reconsidera si el lock crece demasiado.

Instalación global: no hay lock de máquina; el registro va a `~/.config/rsk/permissions.json` (junto a `config.json`), con la misma forma y clave por skill.

## 7. Seguridad (no negociable)

Instalar una skill del registry no puede convertirse en una vía para escribir `Bash(*)` en el allowlist de nadie.

- **Nunca en silencio.** Sin TTY: se imprimen y se omiten. Con TTY: prompt por defecto *no*, con `a` para aceptar todas las de esa skill.
- **`--apply-permissions`** existe para automatizar, es explícito y se documenta como "confía en esta skill".
- **Validación en dos sitios:** al publicar (`cmd/generate-registry` rechaza un `allowed-tools` inválido) y al instalar (nunca confiar en el índice).
- Prohibido `Bash(*)`, `Bash(**)`, comodín inicial suelto y entradas sin paréntesis.
- `deny` del cliente siempre gana; `rsk` nunca escribe en `deny`.
- El merge de OpenCode no reordena ni elimina reglas existentes; el de Claude conserva el resto del JSON.
- Las reglas de usuario nunca se revocan; solo lo registrado a nombre de una skill.
- En Claude la concesión es **de turno**; no deja rastro permanente. Si la organización activa `allowManagedPermissionRulesOnly`, no aplica y `rsk status` lo dice.
- `rsk mcp`'s `install_skills` mantiene `_meta["anthropic/requiresUserInteraction"] = true` y añade los permisos al resultado, sin aplicarlos.
- Windows: normalizar `\` a `/` antes de comparar, igual que hace el matcher del cliente.

## 8. Fases

| # | Fase | Salida |
|---|---|---|
| 0 | **Documental y declarativa** | `notify-on-finish` lleva su `allowed-tools`; su `SKILL.md` documenta el equivalente manual para clientes sin el campo. Sin código |
| 1 | **Lectura y validación** | `allowed-tools` parseado (lista y string) y validado; `rsk catalog` y `rsk install --dry-run` lo muestran; no se escribe nada |
| 2 | **Escritura al instalar (OpenCode)** | Traducción con `${CLAUDE_SKILL_DIR}` resuelto; ledger en `rsk.lock`; `--apply-permissions`/`--skip-permissions`; `rsk tools --for opencode` |
| 3 | **Reconciliación al actualizar** | `rsk update` informa del diff en Claude y escribe en OpenCode; no-op cuando no cambió nada |
| 4 | **Retirada y drift** | `uninstall`/`destroy` revocan; `rsk sync` repara; `rsk status` marca drift |
| 5 | **Superficie pública** | El índice del registry lleva el resumen de permisos; `rsk mcp` los devuelve; `rsk catalog` los muestra antes de instalar |

Fase 0 hecha. Las fases 1-2 son el núcleo; 3 es la que se pidió; 4 y 5 pueden ir después.

## 9. Tests

- **Unitarios:** parser (`allowed-tools` como lista, como string con espacios dentro del paréntesis, ausente, inválido); validación (regla trivial, comodín inicial, más de 5 entradas); traducción por cliente; merge de OpenCode conservando el orden y el resto del fichero; refcount de reglas compartidas; idempotencia (aplicar dos veces = un cambio).
- **Golden:** `opencode.json` resultante y `.claude/settings.json` intacto (Claude no debe escribirse).
- **Integración** con `HOME` y proyecto temporales: instalar → regla escrita solo en OpenCode; actualizar a una versión que añade una entrada → se ofrece; a una que la retira → se ofrece revocar; desinstalar → desaparece solo lo nuestro.
- **No-TTY:** el comando sale 0 sin escribir y lo dice en stderr.
- **Grant de Claude:** `claude -p --permission-mode manual --permission-prompts none --allowedTools Skill` invoca la skill, corre el script sin prompt y sigue denegando un comando no declarado (control negativo).
- **Regresión:** una skill sin `allowed-tools` no cambia ningún byte de configuración.

## 10. Riesgos y decisiones abiertas

| Riesgo | Mitigación |
|---|---|
| Supply chain: una skill del registry concede de más | Consentimiento obligatorio, validación al publicar, límite de 5 entradas, revisión visible en el preview |
| El permiso de Claude dura un turno y no la sesión | Aceptado: se lanza el trabajo desatado en el turno en que la skill está activa; si hace falta persistente, el usuario tiene `rsk tools allow` |
| `opencode.jsonc` pierde comentarios al reescribir | Ya pasa con `instructions`; documentar y, si molesta, preservar comentarios en un cambio aparte |
| Scope equivocado en Claude (proyecto vs `settings.local.json`) | Claude no se escribe; `rsk tools` sigue siendo la vía manual y respeta el scope actual |
| Lock desincronizado provoca revocaciones falsas | El diff se hace contra el lock y las revocaciones nunca son automáticas |
| El registro crece y el lock deja de ser legible | Alternativa `permissions.toml` ya evaluada |

Decisiones resueltas: declaración en `allowed-tools` del `SKILL.md` (no bloque nuevo ni sidecar); registro en `rsk.lock`; la fase 5 se abre como plan aparte cuando 1-4 estén cerradas.

## 11. Estado de implementación (2026-11-17)

Lo que ya está en el árbol:

| Pieza | Dónde |
|---|---|
| Parser de `allowed-tools` (lista YAML y string) + validación | `internal/skill/registry.go` (`AllowedToolsPrefix`, `ValidateAllowedTools`, `NormalizeAllowedTools`, `SplitRule`) |
| `Skill.AllowedTools` | `internal/skill/skill.go` |
| Registro en el lock | `internal/manifest/lock.go` (`permissions`, `granted`, `FindLockEntry`) |
| Regla neutra y traducción por cliente | `internal/tool/permissions.go` (`Rule`, `Translate`, `PermissionWriter`, `WriterFor`) |
| Escritura en OpenCode (append, última regla gana, sin reordenar) | `internal/tool/opencode.go` (`GrantPermissions`, `RevokePermissions`) |
| Reconciliación, consentimiento y registro | `internal/permission/reconcile.go` |
| Retirada al desinstalar | `internal/permission/revoke.go` + `cmd/rsk/uninstall.go` |
| Cableado CLI y flags | `cmd/rsk/permissions.go`, `install.go`, `update_lock.go`, `setup.go` |

Verificado con el binario en un proyecto temporal: instalar escribe la regla solo en OpenCode y la registra; reinstalar no duplica ni reescribe; una versión que añade una entrada escribe solo esa; una que la retira **avisa una vez** y no borra; desinstalar retira exactamente lo escrito y deja `{}`.

Pendiente, fuera de esta entrega (backlog `qk4m2v`, `v7pd1a`):

- `rsk tools --for opencode` (hoy sigue siendo solo Claude).
- Un `rsk update` real de extremo a extremo: la reconciliación comparte `reconcileLock` con el install y está cubierta por tests, pero el camino `update` no se ejecutó contra una versión nueva publicada.
- Fase 5 (índice del registry y `rsk mcp`).

## 12. Criterios de salida

- [x] `notify-on-finish` declara su `allowed-tools` y lanza su script sin prompt en Claude Code.
- [x] `rsk install` en OpenCode escribe la regla al final, conserva las existentes y es idempotente al repetir.
- [x] En Claude, `rsk install`/`rsk update` no escriben `.claude/settings.json` y lo informan.
- [x] Subir una skill a una versión con una entrada nueva en `allowed-tools` hace que `rsk update` la ofrezca una sola vez.
- [x] Una versión que retira una entrada la ofrece revocar y no la borra sin respuesta.
- [x] `rsk uninstall` deja `opencode.json` como estaba antes de instalar la skill (salvo reglas del usuario).
- [x] Una skill sin `allowed-tools` no produce ningún cambio en la configuración.
- [ ] Sin TTY, ningún comando escribe permisos.
