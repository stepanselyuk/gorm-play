package db_mappings

import (
	"encoding/xml"
	"fmt"
	"github.com/metaleap/go-util-fs"
	atom "github.com/metaleap/go-xsd-pkg/kbcafe.com/rss/atom.xsd.xml_go"
	doctrine "github.com/stepanselyuk/doctrine-mappings-xsd-go/doctrine-project.org/schemas/orm/doctrine-mapping.xsd_go"
	"log"
	"path/filepath"
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

	//v := doc.Entities[0].Fields[4].Options.Options[0].XsdGoPkgCDATA
	v := doc.Entities[0].Table

	fmt.Printf("Product: %+v %p\n", v, v)
}

type AtomEntryDoc struct {
	XMLName xml.Name `xml:"entry"`
	atom.TentryType
}

type AtomFeedDoc struct {
	XMLName xml.Name `xml:"feed"`
	atom.TfeedType
}

func AtomMapping() {

	filename, _ := filepath.Abs("db_mappings/xml/samplefeed.xml")
	log.Printf("Loading %s", filename)

	doc, dataOrig := &AtomFeedDoc{}, ufs.ReadBinaryFile(filename, true)

	err := xml.Unmarshal(dataOrig, doc)

	if err != nil {
		panic(err)
	}

	//v := doc.Ids[0].Lang
	v := doc.Titles[0].XsdGoPkgCDATA

	fmt.Printf("Atom: %+v %p\n", v, v)
}
