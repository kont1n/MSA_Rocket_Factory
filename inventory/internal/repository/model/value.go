package model

// Value для metadata
type Value struct {
	StringValue  string  `bson:"string_value,omitempty"`
	Int64Value   int64   `bson:"int64_value,omitempty"`
	Float64Value float64 `bson:"float64_value,omitempty"`
	BoolValue    bool    `bson:"bool_value,omitempty"`
}
