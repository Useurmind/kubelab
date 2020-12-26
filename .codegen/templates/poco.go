package {{ .Spec.GoNamespace }}

{{- range $pocoDef := .Spec.PocoTypes}}

// {{$pocoDef.Description}}
type {{$pocoDef.PocoName}} struct {
	{{- range $prop := $pocoDef.Properties}}
	// {{$prop.Description}}
	{{ title $prop.Name }} {{ $prop.GetType "go" }} `json:"{{ untitle $prop.Name }}"`
	{{- end}}
}

{{- end}}