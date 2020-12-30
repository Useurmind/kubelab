{{- range $pocoDef := .Spec.PocoTypes}}

// {{$pocoDef.Description}}
export interface I{{$pocoDef.PocoName}} {
	{{- range $prop := $pocoDef.Properties}}
	// {{$prop.Description}}
	{{ untitle $prop.Name }}?: {{ $prop.GetType "ts" }};
	{{- end}}
}

{{- end}}