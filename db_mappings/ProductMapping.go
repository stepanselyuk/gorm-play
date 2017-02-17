package db_mappings

import (
	"encoding/xml"
	doctrine "github.com/stepanselyuk/doctrine-mappings-xsd-go/doctrine-project.org/schemas/orm/doctrine-mapping.xsd_go"
	"log"
	"path/filepath"
	//"io/ioutil"
	"github.com/metaleap/go-util-fs"
)

type ProductMapping struct {
	XMLName xml.Name `xml:"doctrine-mapping"`
	doctrine.TxsdDoctrineMapping
}

func GetProductMapping() {

	filename, _ := filepath.Abs("db_mappings/xml/Products.orm.xml")
	log.Printf("Loading %s", filename)

	doc, dataOrig := &ProductMapping{}, ufs.ReadBinaryFile(filename, true)

	err := xml.Unmarshal(dataOrig, doc)

	if err != nil {
		panic(err)
	}

	log.Print(doc.Entities[0].Table.String())
}
