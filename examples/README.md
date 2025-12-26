# ATemplateB Examples

This directory contains example templates demonstrating various ATemplateB features and use cases.

## Quick Examples

### 1. Simple Nginx Reverse Proxy

**Files**: [nginx.tmpl](nginx.tmpl) + [config.yaml](config.yaml)

The simplest example showing basic templating features:
- Value substitution
- Conditionals (`if` statements)
- Loops (`range`)
- Whitespace control

**Generate**:
```bash
ATemplateB --template=examples/nginx.tmpl --values=examples/config.yaml --output=nginx.conf
```

**Features demonstrated**:
- `.Values.<key>` access
- `{{ range .upstreams }}` iteration
- `{{ if .ssl_enabled }}` conditionals
- `{{-` and `-}}` whitespace trimming

**Expected output**: See [nginx.conf](nginx.conf)

---

### 2. YAML File Inclusion

**File**: [landing.tmpl.yaml](landing.tmpl.yaml)

Demonstrates the `getYaml` custom function for dynamic file inclusion.

**Generate**:
```bash
ATemplateB --template=examples/landing.tmpl.yaml --output=landing.yaml
```

**Features demonstrated**:
- `getYaml` function to read external YAML files
- Sprig `indent` function for formatting
- Sprig `trim` function for cleanup
- Template without values file

**Use cases**:
- Composing multiple YAML configurations
- Including shared snippets
- Building complex configs from modular pieces

---

### 3. Multi-Template Composition

**Files**: [sca-nginx.tmpl.conf](sca-nginx.tmpl.conf) + [pdm.tmpl](pdm.tmpl) + [proc.tmpl](proc.tmpl) + [stream.tmpl](stream.tmpl) + [remoteplay.tmpl](remoteplay.tmpl)

Advanced example showing template composition with multiple files.

**Features demonstrated**:
- `{{ template "name" }}` for including other templates
- `{{ define "name" }}` for creating reusable template blocks
- Passing context to sub-templates with `dict`
- Complex nested data structures
- Multiple range iterations
- Parent context access with `$.Values`
- Default values with `| default`

**How it works**:

1. **Main template** ([sca-nginx.tmpl.conf](sca-nginx.tmpl.conf)):
   - Defines the outer nginx structure
   - Iterates over endpoint arrays
   - Calls sub-templates for each endpoint type

2. **Sub-templates** (pdm.tmpl, proc.tmpl, etc.):
   - Define reusable server block configurations
   - Accept parameters via `dict` function
   - Generate specialized nginx configurations

**Example sub-template call**:
```go
{{- range .Values.endpoints.procs }}
    {{ template "proc" (dict "item" . "global" $) }}
{{- end }}
```

This passes:
- `item`: Current processor endpoint from the array
- `global`: Root template context (accessible as `.global` in sub-template)

**Expected values structure**:
```yaml
hostname: example.com
node_role: primary
max_body_size: 100M
endpoints:
  pdms:
    - site: syd
      port: 411
    - site: mel
      port: 412
  procs:
    - site: syd
      port: 421
  streams:
    - site: syd
      port: 431
  remoteplays:
    - site: syd
      port: 441
```

---

## Template Naming Conventions

Templates can use two naming patterns:

### Pattern 1: Basic `.tmpl`
```
template.tmpl → generates any output extension
```

### Pattern 2: Extension Hint `.tmpl.<ext>`
```
template.tmpl.yaml → helps IDE with YAML syntax
template.tmpl.conf → helps IDE with conf syntax
template.tmpl.json → helps IDE with JSON syntax
```

The extension after `.tmpl` doesn't affect functionality - it only helps your IDE provide better syntax highlighting and autocomplete.

---

## Common Template Patterns

### Accessing Nested Values

```yaml
# values.yaml
app:
  database:
    host: localhost
    port: 5432
```

```go
# template
host: {{ .Values.app.database.host }}
port: {{ .Values.app.database.port }}
```

### Iterating with Index

```yaml
# values.yaml
servers:
  - web-01
  - web-02
  - web-03
```

```go
# template
{{- range $index, $server := .Values.servers }}
server_{{ $index }}: {{ $server }}
{{- end }}
```

**Output**:
```
server_0: web-01
server_1: web-02
server_2: web-03
```

### Accessing Parent Context in Loops

```yaml
# values.yaml
domain: example.com
servers:
  - name: web-01
  - name: web-02
```

```go
# template
{{- range .Values.servers }}
server_name: {{ .name }}.{{ $.Values.domain }}
{{- end }}
```

**Output**:
```
server_name: web-01.example.com
server_name: web-02.example.com
```

Notice `$.Values` (with `$`) accesses the root context, not the current loop item.

### Default Values

```go
# Provide fallback if value is missing
port: {{ .Values.port | default "8080" }}
timeout: {{ .Values.timeout | default 30 }}
cache_enabled: {{ .Values.cache_enabled | default "false" }}
```

### Required Values

```go
# Fail with error message if value is missing
database_url: {{ .Values.database_url | required "database_url must be set" }}
api_key: {{ .Values.api_key | required "api_key is required" }}
```

### Conditional Sections

```go
{{- if .Values.features.caching }}
# Caching enabled
cache_size: {{ .Values.cache_size | default "100M" }}
cache_ttl: {{ .Values.cache_ttl | default "3600" }}
{{- else }}
# Caching disabled
{{- end }}
```

### String Manipulation

```go
# Uppercase
APP_NAME={{ .Values.app_name | upper }}

# Lowercase
username={{ .Values.username | lower }}

# Quote
password="{{ .Values.password | quote }}"

# Indent (useful for YAML)
config: |
{{ .Values.config_block | indent 2 }}
```

---

## Testing Your Templates

### 1. Start Simple
Begin with a minimal template and values file, then add complexity incrementally.

### 2. Use Meaningful Test Data
Use realistic values that represent your actual use case.

### 3. Check Whitespace
Template whitespace can be tricky. Use `{{-` and `-}}` to control it:

```go
# Bad: creates blank lines
{{ range .items }}
{{ . }}
{{ end }}

# Good: clean output
{{- range .items }}
{{ . }}
{{- end }}
```

### 4. Validate Output
Always validate generated configs before deploying:

```bash
# For nginx
nginx -t -c nginx.conf

# For YAML
yamllint output.yaml

# For JSON
jq . output.json
```

---

## Creating Your Own Examples

### Step 1: Define Your Use Case
What configuration file do you need to generate?

### Step 2: Identify Variables
What values should be configurable?

### Step 3: Create Values File
```yaml
# my-values.yaml
variable1: value1
variable2: value2
```

### Step 4: Create Template
```go
# my-template.tmpl
setting1: {{ .Values.variable1 }}
setting2: {{ .Values.variable2 }}
```

### Step 5: Generate and Test
```bash
ATemplateB --template=my-template.tmpl --values=my-values.yaml --output=output.conf
```

---

## Additional Resources

- **Main README**: [../README.md](../README.md)
- **Sprig Functions**: https://masterminds.github.io/sprig/
- **Go Templates**: https://pkg.go.dev/text/template
- **YAML Spec**: https://yaml.org/spec/

---

## Need Help?

If these examples don't cover your use case:
1. Check the main [README](../README.md) for additional documentation
2. Review the [Sprig function reference](https://masterminds.github.io/sprig/)
3. Open an issue on GitHub with your use case
