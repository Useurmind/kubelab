{{- range $pocoDef := .Spec.PocoTypes}}

export interface I{{$pocoDef.PocoName}} {
	{{- range $prop := $pocoDef.Properties}}
	{{ untitle $prop.Name }}: {{ $prop.GetType "ts" }};
	{{- end}}
}

{{- end}}