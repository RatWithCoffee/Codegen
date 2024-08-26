package main

type CommentMethodInf struct {
	Url    string `json:"url"`
	Auth   bool   `json:"auth"`
	Method string `json:"method"`
}

type HandlerInf struct {
	CommentInf CommentMethodInf
	FuncName   string
	StructName string
}

type StructInfo struct {
	StructName string
	Fields     []Field
}

type Field struct {
	FieldType string
	Name      string
	ApigenTag map[string]string
	Paramname string
}
