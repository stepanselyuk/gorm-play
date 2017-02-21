package db_mappings

type field struct {
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

type fieldId struct {
	GeneratorStrategy string
}
