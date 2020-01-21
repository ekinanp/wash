package rql

type ASTNode interface {
	Marshal() interface{}
	Unmarshal(interface{}) error
}
