package db_mappings

import (
	"encoding/xml"
	ufs "github.com/metaleap/go-util-fs"
	doctrine "github.com/stepanselyuk/doctrine-mappings-xsd-go/doctrine-project.org/schemas/orm/doctrine-mapping.xsd_go"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// model name -> data name -> some value
var data map[string]*modelDetails
var logger *log.Logger

type modelDetails struct {
	entity       *doctrine.Tentity
	primaryKeys  map[string]*fieldId
	normalFields map[string]*field
	relations    map[string]*modelRelation
}

// initialize struct
func (d *modelDetails) init() {

	d.primaryKeys = make(map[string]*fieldId)
	d.normalFields = make(map[string]*field)
	d.relations = make(map[string]*modelRelation)
}

// get table name of entity
func (d *modelDetails) getTableName() string {
	return d.entity.Table.String()
}

// get simple name of entity
func (d *modelDetails) getModelSimpleName() string {
	return getEntitySimpleName(d.entity)
}

// init data map
func initDataMap() {
	if data == nil {
		data = make(map[string]*modelDetails)
	}
}

// init logger
func initLogger() {

	if logger == nil {
		flag := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
		logger = log.New(os.Stdout, "", flag)
		logger.Println("Logger has been set")
	}
}

// get absolute path of directory with xml files
func getXmlFilesDirPath() string {
	dir, _ := filepath.Abs("db_mappings/xml/")
	return dir
}

func loadParseAllXmlModels() {

	dir := getXmlFilesDirPath()

	ufs.WalkFilesIn(dir, func(fullPath string) (keepWalking bool) {
		keepWalking = true
		if strings.HasSuffix(fullPath, ".orm.xml") {

			entity := getDocForFilePath(fullPath).Entities[0]
			entityName := getEntitySimpleName(entity)

			data[entityName] = &modelDetails{
				entity: entity,
			}

			// initialize internal fields
			data[entityName].init()

			logger.Println("Loaded file " + fullPath)
		}
		return
	})
}

// get simple name for entity: AppBundle\Entity\ModelName -> ModelName
func getEntitySimpleName(entity *doctrine.Tentity) string {
	return entity.Name.String()[strings.LastIndex(entity.Name.String(), "\\")+1:]
}

// load and parse xml for specified filepath
func getDocForFilePath(filePath string) *DoctrineMapping {

	doc, dataOrig := &DoctrineMapping{}, ufs.ReadBinaryFile(filePath, true)
	err := xml.Unmarshal(dataOrig, doc)

	if err != nil {
		panic(err)
	}

	return doc
}

// load and parse xml for specified model name
func getDocForModelName(modelName string) *DoctrineMapping {

	dir := getXmlFilesDirPath()
	filePath := dir + "/" + modelName + ".orm.xml"

	return getDocForFilePath(filePath)
}

func GenerateGormModels() {

	initLogger()
	initDataMap()
	loadParseAllXmlModels()

	// handle all models before saving, cause we can affect each other
	for _, modelDetails := range data {
		processModelDetails(modelDetails)
	}

	// save models to files
	for _, modelDetails := range data {
		templateParams := generateModelForEntity(modelDetails)
		if err := SaveModel(templateParams.Name, templateParams, getOutputDirPath()); err != nil {
			logger.Fatal(err)
		}
	}
}

// ------------------------------------------------------
// ------------------------------------------------------

type DoctrineMapping struct {
	XMLName xml.Name `xml:"doctrine-mapping"`
	doctrine.TxsdDoctrineMapping
}

// get absolute path of directory with xml files
func getOutputDirPath() string {
	dir, _ := filepath.Abs("models/generated/")
	return dir
}

func processModelDetails(modelDetails *modelDetails) {

	handlePrimaryKeys(modelDetails)
	handleNormalFields(modelDetails)
	handleRelations(modelDetails)
}

func generateModelForEntity(modelDetails *modelDetails) *TemplateParams {
	return GenerateModel(modelDetails)
}

func handlePrimaryKeys(modelDetails *modelDetails) {

	for _, id := range modelDetails.entity.Ids {

		if id.AssociationKey.B() {

			relations := []relation{}
			for _, v := range modelDetails.entity.OneToOnes {
				relations = append(relations, one2one{v})
			}
			for _, v := range modelDetails.entity.ManyToOnes {
				relations = append(relations, many2one{v})
			}

			// try to find in one2one relations
			fillAssociatedPrimaryKeyThroughDefinedRelations(relations, id)

			// case of we cannot fill data
			if id.Column.String() == "" {
				logger.Printf("Cannot fill details for primary key '%s' of model '%s'.\n", id.Column.String(), modelDetails.getModelSimpleName())
				continue
			}
		}

		modelDetails.primaryKeys[id.Column.String()] = &fieldId{
			GeneratorStrategy: "NONE",
		}

		if id.Generator != nil {
			modelDetails.primaryKeys[id.Column.String()].GeneratorStrategy = id.Generator.Strategy.String()
		}

		// creating new field from primary key

		fieldVar := &doctrine.Tfield{}

		fieldVar.Name.Set(id.Name.String())
		fieldVar.Column.Set(id.Column.String())
		fieldVar.Type.Set(id.Type.String())
		fieldVar.Length.Set(id.Length.String())
		fieldVar.Options = id.Options

		//fmt.Printf("Id field: %+v %p\n", fieldVar, fieldVar)

		// FIXME maybe need to find other way to not spoil in entity data
		modelDetails.entity.Fields = append([]*doctrine.Tfield{fieldVar}, modelDetails.entity.Fields...)
	}
}

// see http://docs.doctrine-project.org/projects/doctrine-orm/en/latest/tutorials/composite-primary-keys.html
func fillAssociatedPrimaryKeyThroughDefinedRelations(rels []relation, id *doctrine.Tid) {

	if rels == nil {
		return
	}

	for _, rel := range rels {

		if rel.GetField().String() == id.Name.String() && rel.GetJoinColumns() != nil && rel.GetJoinColumns().JoinColumns != nil {

			joinColumn := rel.GetJoinColumns().JoinColumns[0]
			id.Column.Set(joinColumn.Name.String())

			// loading doc of related model to get type of id-column, cause both should have the same type to be linked
			relEntity := data[rel.GetTargetEntity().String()].entity

			if relId := getIdByName(relEntity.Ids, joinColumn.ReferencedColumnName.String()); relId != nil {

				id.Type.Set(relId.Type.String())
				id.Length.Set(relId.Length.String())

			} else if relField := getFieldByName(relEntity.Fields, joinColumn.ReferencedColumnName.String()); relField != nil {

				id.Type.Set(relField.Type.String())
				id.Length.Set(relField.Length.String())

			} else {
				logger.Panicf("Cannot find neither related primary key or normal field for column '%s' in linked document of model '%s'.\n", id.Name.String(), rel.GetTargetEntity().String())
			}

		}
	}
}

func handleNormalFields(modelDetails *modelDetails) {

	for _, f := range modelDetails.entity.Fields {

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

		modelDetails.normalFields[f.Column.String()] = &field{
			Name:     f.Column.String(),
			Type:     f.Type.String(),
			Nullable: f.Nullable.B(),
			Unique:   f.Unique.B(),
			// only for "string"
			Length: fieldLength,
			// only for "decimal"
			Scale: f.Scale.N(),
			// only for "decimal"
			Precision: f.Precision.N(),
			// --------------------------
			// options
			Default:  optionDefaultValue,
			Unsigned: optionUnsigned,
			// for "string"
			Fixed:   optionFixed,
			Comment: optionComment,
		}
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

type modelRelation struct {
	Type                 string
	TypeDescription      string
	ModelName            string
	ThisColumnName       string
	ReferencedColumnName string
	ReferencedColumnType string
}

func handleRelations(modelDetails *modelDetails) {

	for _, rel := range modelDetails.entity.OneToOnes {

		thisColumnName := rel.JoinColumns.JoinColumns[0].Name.String()
		referencedColumnName := rel.JoinColumns.JoinColumns[0].ReferencedColumnName.String()

		modelDetails.relations[thisColumnName+":"+rel.TargetEntity.String()] = &modelRelation{
			Type:                 "BelongsTo",
			TypeDescription:      "One({x})-To-One({y})",
			ModelName:            rel.TargetEntity.String(),
			ThisColumnName:       thisColumnName,
			ReferencedColumnName: referencedColumnName,
			ReferencedColumnType: getReferencedColumnType(rel.TargetEntity.String(), referencedColumnName),
		}

		// we can infer as opposite has-one relation from belongs-to/one-to-one relation
		modelDetailsTarget := data[rel.TargetEntity.String()]

		// many models get refer to the same column on target entity so we are using complex key
		if _, ok := modelDetailsTarget.relations[referencedColumnName+":"+modelDetails.getModelSimpleName()]; !ok {

			modelDetailsTarget.relations[referencedColumnName+":"+modelDetails.getModelSimpleName()] = &modelRelation{
				Type:                 "HasOne",
				TypeDescription:      "One({x})-Has-One({y})",
				ModelName:            modelDetails.getModelSimpleName(),
				ThisColumnName:       referencedColumnName,
				ReferencedColumnName: thisColumnName,
				ReferencedColumnType: getReferencedColumnType(modelDetails.getModelSimpleName(), thisColumnName),
			}
		}
	}

	for _, rel := range modelDetails.entity.ManyToOnes {

		thisColumnName := rel.JoinColumns.JoinColumns[0].Name.String()
		referencedColumnName := rel.JoinColumns.JoinColumns[0].ReferencedColumnName.String()

		modelDetails.relations[thisColumnName+":"+rel.TargetEntity.String()] = &modelRelation{
			Type:                 "BelongsTo",
			TypeDescription:      "Many({x})-To-One({y})",
			ModelName:            rel.TargetEntity.String(),
			ThisColumnName:       thisColumnName,
			ReferencedColumnName: referencedColumnName,
			ReferencedColumnType: getReferencedColumnType(rel.TargetEntity.String(), referencedColumnName),
		}

		// we can infer as opposite has-many relation from belongs-to/many-to-one relation
		modelDetailsTarget := data[rel.TargetEntity.String()]

		// many models get refer to the same column on target entity so we are using complex key
		if _, ok := modelDetailsTarget.relations[referencedColumnName+":"+modelDetails.getModelSimpleName()]; !ok {

			modelDetailsTarget.relations[referencedColumnName+":"+modelDetails.getModelSimpleName()] = &modelRelation{
				Type:                 "HasMany",
				TypeDescription:      "One({x})-Has-Many({y})",
				ModelName:            modelDetails.getModelSimpleName(),
				ThisColumnName:       referencedColumnName,
				ReferencedColumnName: thisColumnName,
				ReferencedColumnType: getReferencedColumnType(modelDetails.getModelSimpleName(), thisColumnName),
			}
		}
	}

	// FIXME need to implement many-to-many relation
}

func getReferencedColumnType(modelName string, columnName string) string {

	entity := data[modelName].entity

	for _, id := range entity.Ids {

		if id.AssociationKey.B() {
			continue
		}

		if id.Column.String() == columnName {
			return id.Type.String()
		}
	}

	for _, f := range entity.Fields {

		if f.Column.String() == columnName {
			return f.Type.String()
		}
	}

	return "string"
}
