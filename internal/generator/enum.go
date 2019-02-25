package generator

import (
	"strconv"
	"strings"

	"github.com/RettyInc/gqlcodegen/gql"
)

func generateEnum(g *Generator, def *gql.Enum) {
	g.Printf(commentOnTop)
	generateEnumPackageSection(g, def)
	g.Println()
	generateEnumTypeDefSection(g, def)
	g.Println()
	generateEnumConstSection(g, def)
	g.Println()
	generateStringMethod(g, def)
	g.Println()
	generateFromString(g, def)
	g.Println()
	generateImplementGqlType(g, def)
	g.Println()
	generateUnmarshalGraphQL(g, def)
	g.Println()
	generateMarshalJson(g, def)
}

func generateEnumPackageSection(g *Generator, def *gql.Enum) {
	g.Printf("package %s\n", strings.ToLower(def.Name))
	g.Println()
	g.Println("import (")
	g.Println(`"errors"`)
	g.Println(`"strconv"`)
	g.Println(")")
}

func generateEnumTypeDefSection(g *Generator, def *gql.Enum) {
	generateComment(g, def)
	g.Printf("type %s int\n", capitalizeFirst(def.Name))
}

func generateEnumConstSection(g *Generator, def *gql.Enum) {
	g.Println("const(")
	generateEnumConstBody(g, def)
	g.Println(")")
}

func generateEnumConstBody(g *Generator, def *gql.Enum) {
	entries := def.Values
	for i, e := range entries {
		generateComment(g, e)
		g.Printf("%s", capitalizeFirst(e.Name))
		if i == 0 {
			g.Printf(" %s = iota", capitalizeFirst(def.Name))
		}
		g.Println()
	}
}

func generateStringMethod(g *Generator, def *gql.Enum) {
	eName := capitalizeFirst(def.Name)
	typeName := ""
	var index []string
	length := 0
	for _, v := range def.Values {
		typeName += v.Name
		index = append(index, strconv.FormatInt(int64(length), 10))
		length += len(v.Name)
	}
	index = append(index, strconv.FormatInt(int64(length), 10))

	g.Printf("const _%s_Name = \"%s\"\n", eName, typeName)
	g.Printf("var _%s_Index = []int{%s}\n", eName, strings.Join(index, ", "))
	g.Println()
	g.Printf("func (v %s) String() string {\n", eName)
	g.Printf("if v < 0 || v >= %s(len(_%s_Index)-1) {\n", eName, eName)
	g.Printf("return \"%s(\" + strconv.FormatInt(int64(v), 10) + \")\"", eName)
	g.Println()
	g.Println("}")
	g.Printf("return _%s_Name[_%s_Index[v]:_%s_Index[v+1]]", eName, eName, eName)
	g.Println("}")
}

func generateFromString(g *Generator, e *gql.Enum) {
	eName := capitalizeFirst(e.Name)
	g.Printf("func _%sFromString(str string) (%s, error) {\n", eName, eName)
	g.Printf("for i := 0; i < len(_%s_Index) - 1; i++ {\n", eName)
	g.Printf("if v := %s(i); str == v.String() {\n",eName)
	g.Println("return v, nil")
	g.Println("}")
	g.Println("}")
	g.Println(`return -1, errors.New(str + " is not found")`)
	g.Println("}")
}

func generateImplementGqlType(g *Generator, e *gql.Enum) {
	eName := capitalizeFirst(e.Name)
	g.Printf("func (%s) ImplementsGraphQLType(name string) bool {\n", eName)
	g.Printf(`return name == "%s"`, e.Name)
	g.Println()
	g.Println("}")
}

func generateUnmarshalGraphQL(g *Generator, e *gql.Enum) {
	eName := capitalizeFirst(e.Name)
	g.Printf("func (v *%s) UnmarshalGraphQL(input interface{}) error {\n", eName)
	g.Println("switch input := input.(type) {")
	g.Println("case string:")
	g.Printf("value, err := _%sFromString(input)\n", eName)
	g.Println("if err != nil {")
	g.Println("return err")
	g.Println("}")
	g.Println("*v = value")
	g.Println("return nil")
	g.Println("default:")
	g.Println(`return errors.New("wrong type")`)
	g.Println("}")
	g.Println("}")
}

func generateMarshalJson(g *Generator, e *gql.Enum) {
	g.Printf("func (v %s) MarshalJSON() ([]byte, error) {\n", capitalizeFirst(e.Name))
	g.Println("return []byte(`\"`+v.String()+`\"`), nil")
	g.Println("}")
}