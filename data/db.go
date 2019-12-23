package data

// Entity ...
type Entity interface {
	TableName() string
}

// Strategy ...
type Strategy interface {
	Create(entity Entity) error
}
