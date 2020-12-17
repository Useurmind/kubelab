package {{ .Spec.GoNamespace }}

{{- range $pocoDef := .Spec.PocoTypes}}

type {{$pocoDef.PocoName}} struct {
	{{- range $prop := $pocoDef.Properties}}
	{{ title $prop.Name }} {{ $prop.GetType "go" }} `json:"{{ untitle $prop.Name }}"`
	{{- end}}
}

{{- end}}