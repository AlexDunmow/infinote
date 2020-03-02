package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/volatiletech/null"
)

// MarshalNullString to handle custom type
func MarshalNullString(ns null.String) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalString(ns.String)
}

// UnmarshalNullString to handle custom type
func UnmarshalNullString(v interface{}) (null.String, error) {
	if v == nil {
		return null.String{Valid: false}, nil
	}
	s, err := graphql.UnmarshalString(v)
	return null.StringFrom(s), err
}

// MarshalNullInt to handle custom type
func MarshalNullInt(ns null.Int) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalInt(ns.Int)
}

// UnmarshalNullInt to handle custom type
func UnmarshalNullInt(v interface{}) (null.Int, error) {
	if v == nil {
		return null.Int{Valid: false}, nil
	}
	s, err := graphql.UnmarshalInt(v)
	return null.IntFrom(s), err
}
