package db_mappings

import (
	"encoding/xml"
	"fmt"
	ufs "github.com/metaleap/go-util-fs"
	doctrine "github.com/stepanselyuk/doctrine-mappings-xsd-go/doctrine-project.org/schemas/orm/doctrine-mapping.xsd_go"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Field struct {
	Name      string
	Type      string
	Nullable  bool
	Unique    bool
	Length    uint64
	Scale     int64
	Precision int64
	// options
	Default  string
	Unsigned bool
	Fixed    bool
	Comment  string
}

type FieldId struct {
	GeneratorStrategy string
}

func GenerateGormModels() {

	dir, _ := filepath.Abs("db_mappings/xml/")

	ufs.WalkFilesIn(dir, func(fullPath string) (keepWalking bool) {
		keepWalking = true
		if strings.HasSuffix(fullPath, ".xml") {
			generateModelForXml(fullPath)
		}
		return
	})
}

type DoctrineMapping struct {
	XMLName xml.Name `xml:"doctrine-mapping"`
	doctrine.TxsdDoctrineMapping
}

func generateModelForXml(filePath string) {

	doc, dataOrig := &DoctrineMapping{}, ufs.ReadBinaryFile(filePath, true)
	err := xml.Unmarshal(dataOrig, doc)

	if err != nil {
		panic(err)
	}

	entity := doc.Entities[0]

	v := entity.Table
	defer fmt.Println(filePath)
	defer fmt.Printf("Table: %+v %p\n", v, v)

	pkeys := map[string]*FieldId{}
	fields := []*Field{}

	for _, id := range entity.Ids {

		pkeys[id.Column.String()] = &FieldId{
			GeneratorStrategy: id.Generator.Strategy.String(),
		}

		field := &doctrine.Tfield{}

		field.Name.Set(id.Name.String())
		field.Column.Set(id.Column.String())
		field.Type.Set(id.Type.String())
		field.Length.Set(id.Length.String())
		field.Options = id.Options

		//fmt.Printf("Id field: %+v %p\n", field, field)

		entity.Fields = append([]*doctrine.Tfield{field}, entity.Fields...)
	}

	for _, f := range entity.Fields {

		optionDefaultValue := ""
		optionUnsigned := false
		optionFixed := false
		optionComment := ""

		if f.Options != nil {
			for _, o := range f.Options.Options {
				switch o.Name {
				case "default":
					optionDefaultValue = string(o.XsdGoPkgCDATA)
				case "unsigned":
					optionUnsigned, _ = strconv.ParseBool(o.XsdGoPkgCDATA)
				case "fixed":
					optionFixed, _ = strconv.ParseBool(o.XsdGoPkgCDATA)
				case "comment":
					optionComment = string(o.XsdGoPkgCDATA)
				}
			}
		}

		fieldLength, _ := strconv.ParseUint(f.Length.String(), 10, 64)

		fields = append(fields, &Field{
			Name:      f.Column.String(),
			Type:      f.Type.String(),
			Nullable:  f.Nullable.B(),
			Unique:    f.Unique.B(),
			Length:    fieldLength,
			Scale:     f.Scale.N(),
			Precision: f.Precision.N(),
			// options
			Default:  optionDefaultValue,
			Unsigned: optionUnsigned,
			Fixed:    optionFixed,
			Comment:  optionComment,
		})
	}

	fmt.Println(pkeys)

	for _, f1 := range fields {
		fmt.Println(f1.Name)
	}

	modelParams := map[string]*TemplateParams{}
	tables := []string{}

	modelParams[entity.Table.String()] = GenerateModel(entity.Table.String(), pkeys, fields, tables)

	outDir, _ := filepath.Abs("models/generated/")

	for table, param := range modelParams {

		fmt.Println("Add relation for Table name: " + table)

		AddHasMany(param)

		if err := SaveModel(table, param, outDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
