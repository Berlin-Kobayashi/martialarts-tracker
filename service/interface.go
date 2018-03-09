package service

type Repository interface {
	Save(collectionName string, data interface{}) error
	Read(collectionName string, id string, result *interface{}) error
}

type IDGenerator interface {
	Generate() string
}
