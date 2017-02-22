package db_mappings

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gedex/inflector"
	"github.com/serenize/snaker"
)

type TemplateField struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type TemplateParams struct {
	Name            string
	Fields          []*TemplateField
	NeedTimePackage bool
	TableName       string
}

var hasMany = make(map[string][]string)

func GenerateModel(modelDetails *modelDetails) *TemplateParams {

	var needTimePackage bool

	templateFields := []*TemplateField{}

	for _, field := range modelDetails.normalFields {

		fieldType := gormDataType(field)

		if fieldType == "time.Time" || fieldType == "*time.Time" {

			needTimePackage = true

			if field.Nullable {
				fieldType = "*time.Time"
			} else {
				fieldType = "time.Time"
			}
		}

		templateFields = append(templateFields, &TemplateField{
			Name:    gormColumnName(field.Name),
			Type:    fieldType,
			Tag:     genJSON(field.Name, field.Default, modelDetails.primaryKeys),
			Comment: field.Comment,
		})
	}

	for _, rel := range modelDetails.relations {

		relTypeDesc := strings.Replace(rel.TypeDescription, "{x}", modelDetails.getModelSimpleName(), 1)
		relTypeDesc = strings.Replace(relTypeDesc, "{y}", rel.ModelName, 1)

		switch rel.Type {
		case "BelongsTo", "HasOne":

			templateFields = append(templateFields, &TemplateField{
				Name:    rel.ModelName,
				Type:    "*" + rel.ModelName,
				Tag:     "gorm:\"ForeignKey:" + rel.ThisColumnName + ";AssociationForeignKey:" + rel.ReferencedColumnName + "\"",
				Comment: rel.Type + " \"" + rel.ModelName + "\" with relation " + relTypeDesc + ".",
			})

			// adding holder field if it was not added before (id columns with association)
			// FIXME need to extract logic to a reused func
			if _, ok := modelDetails.normalFields[rel.ThisColumnName]; !ok {

				templateFields = append(templateFields, &TemplateField{
					Name:    gormColumnName(rel.ThisColumnName),
					Type:    gormDataType(&field{Type: rel.ReferencedColumnType}),
					Tag:     genJSON(rel.ThisColumnName, "", nil),
					Comment: "Value holder field for relation with \"" + rel.ModelName + "\".",
				})
			}

		case "HasMany":

			templateFields = append(templateFields, &TemplateField{
				Name:    inflector.Pluralize(rel.ModelName),
				Type:    "[]*" + rel.ModelName,
				Tag:     "gorm:\"ForeignKey:" + rel.ReferencedColumnName + ";AssociationForeignKey:" + rel.ThisColumnName + "\"",
				Comment: rel.Type + " \"" + rel.ModelName + "\" with relation " + relTypeDesc + ".",
			})

			// adding holder field if it was not added before (id columns with association)
			// FIXME need to extract logic to a reused func
			if _, ok := modelDetails.normalFields[rel.ThisColumnName]; !ok {

				templateFields = append(templateFields, &TemplateField{
					Name:    gormColumnName(rel.ThisColumnName),
					Type:    gormDataType(&field{Type: rel.ReferencedColumnType}),
					Tag:     genJSON(rel.ThisColumnName, "", nil),
					Comment: "Value holder field for relation with \"" + rel.ModelName + "\".",
				})
			}
		}
	}

	params := &TemplateParams{
		Name:            modelDetails.getModelSimpleName(),
		Fields:          templateFields,
		NeedTimePackage: needTimePackage,
		TableName:       modelDetails.getTableName(),
	}

	return params
}

func AddHasMany(params *TemplateParams) {
	if _, ok := hasMany[params.Name]; ok {
		for _, infColName := range hasMany[params.Name] {
			params.Fields = append(params.Fields, &TemplateField{
				Name:    gormColumnName(infColName),
				Type:    "[]*" + gormTableName(infColName),
				Tag:     genJSON(strings.ToLower(infColName), "", nil),
				Comment: "This line is infered from other tables.",
			})
		}
	}
}

func SaveModel(modelName string, params *TemplateParams, outPath string) error {

	body, err := Asset("db_mappings/_templates/model.go.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("").Parse(string(body))
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, params); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	modelFile := filepath.Join(outPath, snaker.CamelToSnake(modelName)+".go")

	if err := ioutil.WriteFile(modelFile, src, 0644); err != nil {
		return err
	}

	return nil
}

// Infer belongs_to Relation from column's name
func inferORM(s string, tables []string) (bool, string) {
	s = strings.ToLower(s)
	ss := strings.Split(s, "_")

	newSS := []string{}
	var containsID bool = false
	for _, word := range ss {
		if word == "id" {
			containsID = true
			continue
		}

		newSS = append(newSS, word)
	}

	if containsID == false || len(newSS) == 0 {
		return false, ""
	}

	infColName := strings.Join(newSS, "_")

	// Check the table is existed or not
	tableName := snaker.CamelToSnake(infColName)
	tableName = inflector.Pluralize(tableName)

	exist := false
	for _, table := range tables {
		if table == tableName {
			exist = true
		}
	}

	if !exist {
		return false, ""
	}

	return true, infColName
}

// Generate json
func genJSON(columnName, columnDefault string, primaryKeys map[string]*fieldId) (json string) {
	json = "json:\"" + columnName + "\""

	// FIXME can be other strategies
	// see http://docs.doctrine-project.org/projects/doctrine-orm/en/latest/reference/basic-mapping.html#identifier-generation-strategies
	if pk, ok := primaryKeys[columnName]; ok {

		var p string

		if pk.GeneratorStrategy == "IDENTITY" {
			p = "gorm:\"primary_key;AUTO_INCREMENT\" "
		} else {
			p = "gorm:\"primary_key\" "
		}

		json = p + json
	}

	if columnDefault != "" && !strings.Contains(columnDefault, "nextval") {
		d := " sql:\"DEFAULT:" + columnDefault + "\""
		json += d
	}

	return
}

// Singlarlize table name and upper initial character
func gormTableName(s string) string {
	var tableName string

	tableName = strings.ToLower(s)
	tableName = inflector.Singularize(tableName)
	tableName = snaker.SnakeToCamel(tableName)

	return strings.Title(tableName)
}

// Ex: facebook_uid → FacebookUID
func gormColumnName(s string) string {
	s = strings.ToLower(s)
	ss := strings.Split(s, "_")

	for i, word := range ss {
		if word == "id" || word == "uid" || word == "url" {
			word = strings.ToUpper(word)
		}

		ss[i] = strings.Title(word)
	}

	return strings.Join(ss, "")
}

// See http://docs.doctrine-project.org/projects/doctrine-orm/en/latest/reference/basic-mapping.html#doctrine-mapping-types
// FIXME maybe better to use general types instead of concrete types
func gormDataType(f *field) string {
	switch {

	case f.Type == "string" || f.Type == "text" || f.Type == "decimal" || f.Type == "guid":
		return "string"

	case f.Type == "integer" && f.Unsigned:
		return "uint"

	case f.Type == "smallint" && f.Unsigned:
		return "uint16"

	case f.Type == "bigint" && f.Unsigned:
		return "uint64"

	case f.Type == "integer" && !f.Unsigned:
		return "int"

	case f.Type == "smallint" && !f.Unsigned:
		return "int16"

	case f.Type == "bigint" && !f.Unsigned:
		return "int64"

	case f.Type == "boolean":
		return "bool"

	case f.Type == "date" || f.Type == "time" || f.Type == "datetime":
		return "*time.Time"

	case f.Type == "datetimetz":
		return "time.Time"

	case f.Type == "float":
		return "float64"

	// Questions:
	// object
	// array
	// simple_array
	// json_array
	// blob --> []byte ?

	default:
		return f.Type
	}
}
