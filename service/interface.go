package service

type Repository interface {
	Save(data interface{}) error
	Read(id string, result *interface{}) error
}

type UUIDGenerator interface {
	Generate() string
}
