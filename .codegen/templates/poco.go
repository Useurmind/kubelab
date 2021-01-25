package {{ .Spec.GoNamespace }}

{{- range $enumDef := .Spec.StringEnums}}

const (
	{{- range $key, $value := $enumDef.Values}}
	{{ $enumDef.Name }}{{ $key }} = "{{ $value }}"
	{{- end}}
)

{{- end }}

{{- range $pocoDef := .Spec.PocoTypes}}

// {{$pocoDef.Description}}
type {{$pocoDef.PocoName}} struct {
	{{- range $prop := $pocoDef.Properties}}
	// {{$prop.Description}}
	{{ title $prop.Name }} {{ $prop.GetType "go" }} `json:"{{ untitle $prop.Name }},omitempty"`
	{{- end}}
}

{{- end}}