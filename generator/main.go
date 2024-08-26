package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"reflect"
	"strings"
)

type ApiTags map[string]string

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out, imports)

	methods := make([]HandlerInf, 0)
	var structName string
	var paramStructs []StructInfo
	for _, f := range node.Decls {
		funcDecl, ok := f.(*ast.FuncDecl)

		// если declaration это функция
		if ok {
			methods = parseFuncInfo(funcDecl, methods, &structName)
		}

		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if !strings.Contains(currType.Name.Name, "Params") {
				continue
			}

			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				continue
			}

			paramStructs = parseStruct(paramStructs, currStruct, currType)
		}

	}

	// запись функции ServeHttp в выходной файл
	err = serveHttpTemplate.Execute(out, methods)
	if err != nil {
		fmt.Print(err)
	}

	for i, handlerInf := range methods {
		// запись хэндлера для метода структуры в выходной файл
		err = handlerTemplate.Execute(out, handlerInf)
		if err != nil {
			fmt.Print(err)
		}

		// запись функции получения полей структуры из url query в выходной файл
		err = fillParamsGet.Execute(out, paramStructs[i])
	}

}

// пасит теги полей структуры
func parseApiTag(tag reflect.StructTag) ApiTags {
	apiTagsStr := tag.Get("apivalidator")
	var tagValue, tagName string
	tags := make(ApiTags)
	for _, t := range strings.Split(apiTagsStr, ",") {
		tagName = strings.Split(t, "=")[0]
		if strings.Contains(t, "=") {
			tagValue = strings.Split(t, "=")[1]
		} else {
			tagValue = ""
		}

		tags[tagName] = tagValue

	}
	return tags
}

// парсит структуру для получения информации о ней
func parseStruct(paramStructs []StructInfo, currStruct *ast.StructType, currType *ast.TypeSpec) []StructInfo {
	paramStructs = append(paramStructs, StructInfo{currType.Name.Name, make([]Field, 0)})
	currStructFields := &paramStructs[len(paramStructs)-1].Fields
	var fieldName, paramname string
	for _, field := range currStruct.Fields.List {
		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
			parsedTag := parseApiTag(tag)
			val, ok := parsedTag["paramname"]
			fieldName = field.Names[0].Name
			if ok {
				paramname = val
			} else {
				paramname = strings.ToLower(fieldName)
			}
			fieldInf := Field{types.ExprString(field.Type), fieldName, parsedTag, paramname}
			*currStructFields = append(*currStructFields, fieldInf)
		}
	}

	return paramStructs
}

// парсит доку функции в структуру HandlerInf
func parseFuncInfo(funcDecl *ast.FuncDecl, methods []HandlerInf, structName *string) []HandlerInf {
	var str string
	var inf CommentMethodInf
	if funcDecl.Doc != nil {
		for _, com := range funcDecl.Doc.List {
			if strings.Contains(com.Text, "apigen:api") {
				str = com.Text
				str = str[strings.Index(str, "{") : strings.Index(str, "}")+1]
				err := json.Unmarshal([]byte(str), &inf)
				if err != nil {
					log.Println("Error unmarshalling json:", str)
				}
			}
		}
		str = types.ExprString(funcDecl.Recv.List[0].Type)
		if *structName == "" {
			*structName = str
		}

		if *structName != str {
			return methods
		}

		methods = append(methods, HandlerInf{inf, funcDecl.Name.Name, (*structName)[1:]})
	}
	return methods
}
