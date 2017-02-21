package db_mappings

import (
	"github.com/metaleap/go-xsd/types"
	doctrine "github.com/stepanselyuk/doctrine-mappings-xsd-go/doctrine-project.org/schemas/orm/doctrine-mapping.xsd_go"
)

// universal relations interface for doctrine.ToneToOne and doctrine.TmanyToOne to avoid the logic copying
type relation interface {
	GetField() xsdt.Nmtoken
	GetJoinColumns() *doctrine.TjoinColumns
	GetTargetEntity() xsdt.String
}

type one2one struct {
	*doctrine.ToneToOne
}

func (r one2one) GetField() xsdt.Nmtoken                 { return r.Field }
func (r one2one) GetJoinColumns() *doctrine.TjoinColumns { return r.JoinColumns }
func (r one2one) GetTargetEntity() xsdt.String           { return r.TargetEntity }

type many2one struct {
	*doctrine.TmanyToOne
}

func (r many2one) GetField() xsdt.Nmtoken                 { return r.Field }
func (r many2one) GetJoinColumns() *doctrine.TjoinColumns { return r.JoinColumns }
func (r many2one) GetTargetEntity() xsdt.String           { return r.TargetEntity }
