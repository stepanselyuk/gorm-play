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

// get absolute path of directory with xml files
func getXmlFilesDirPath() string {
	dir, _ := filepath.Abs("db_mappings/xml/")
	return dir
}

func GenerateGormModels() {

	dir := getXmlFilesDirPath()

	ufs.WalkFilesIn(dir, func(fullPath string) (keepWalking bool) {
		keepWalking = true
		if strings.HasSuffix(fullPath, ".orm.xml") {
			generateModelForXml(fullPath)
		}
		return
	})
}

type DoctrineMapping struct {
	XMLName xml.Name `xml:"doctrine-mapping"`
	doctrine.TxsdDoctrineMapping
}

func getDocForFilePath(filePath string) *DoctrineMapping {

	doc, dataOrig := &DoctrineMapping{}, ufs.ReadBinaryFile(filePath, true)
	err := xml.Unmarshal(dataOrig, doc)

	if err != nil {
		panic(err)
	}

	return doc
}

func getDocForModelName(modelName string) *DoctrineMapping {
	dir := getXmlFilesDirPath()
	filePath := dir + "/" + modelName + ".orm.xml"
	return getDocForFilePath(filePath)
}

func generateModelForXml(filePath string) {

	entity := getDocForFilePath(filePath).Entities[0]

	//v := entity.Table
	defer fmt.Println(filePath)
	//defer fmt.Printf("Table: %+v %p\n", v, v)

	pkeys := map[string]*FieldId{}
	fields := []*Field{}

	for _, id := range entity.Ids {

		var pkey *FieldId

		if id.AssociationKey.B() {

			relations := []relation{}
			for _, v := range entity.OneToOnes {
				relations = append(relations, one2one{v})
			}
			for _, v := range entity.ManyToOnes {
				relations = append(relations, many2one{v})
			}

			// try to find in one2one relations
			fillAssociatedIdThroughRelations(relations, id)

			// case of we cannot fill data
			if id.Column.String() == "" {
				continue
			}
		}

		pkey = &FieldId{
			GeneratorStrategy: "NONE",
		}

		if id.Generator != nil {
			pkey.GeneratorStrategy = id.Generator.Strategy.String()
		}

		pkeys[id.Column.String()] = pkey

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

	modelParams := &TemplateParams{}
	tables := []string{}

	entityName := entity.Name.String()[strings.LastIndex(entity.Name.String(), "\\")+1:]

	modelParams = GenerateModel(entityName, entity.Table.String(), pkeys, fields, tables)

	outDir, _ := filepath.Abs("models/generated/")

	//fmt.Println("Add relation for Table name: " + table)
	//AddHasMany(param)

	if err := SaveModel(entity.Table.String(), modelParams, outDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// get field by its name from specified slice with fields
func getFieldByName(fields []*doctrine.Tfield, name string) *doctrine.Tfield {
	for _, field := range fields {
		if field.Column.String() == name {
			return field
		}
	}
	return nil
}

// get ID-node by its name from specified slice with IDs
func getIdByName(ids []*doctrine.Tid, name string) *doctrine.Tid {
	for _, id := range ids {
		if id.Column.String() == name {
			return id
		}
	}
	return nil
}

// see http://docs.doctrine-project.org/projects/doctrine-orm/en/latest/tutorials/composite-primary-keys.html
func fillAssociatedIdThroughRelations(rels []relation, id *doctrine.Tid) {

	if rels == nil {
		return
	}

	for _, rel := range rels {

		if rel.GetField().String() == id.Name.String() && rel.GetJoinColumns() != nil && rel.GetJoinColumns().JoinColumns != nil {

			joinColumn := rel.GetJoinColumns().JoinColumns[0]
			id.Column.Set(joinColumn.Name.String())

			// loading doc of related model to get type of id-column, cause both should have
			// the same type to be linked
			relDoc := getDocForModelName(rel.GetTargetEntity().String())

			relId := getIdByName(relDoc.Entities[0].Ids, joinColumn.ReferencedColumnName.String())
			if relId != nil {

				id.Type.Set(relId.Type.String())
				id.Length.Set(relId.Length.String())

			} else {

				relField := getFieldByName(relDoc.Entities[0].Fields, joinColumn.ReferencedColumnName.String())

				if relField == nil {
					panic(fmt.Sprintf("Cannot find neither related ID or normal field for column '%s' in linked document of model '%s'", id.Name.String(), rel.GetTargetEntity().String()))
				}

				id.Type.Set(relField.Type.String())
				id.Length.Set(relField.Length.String())
			}

		}
	}
}
