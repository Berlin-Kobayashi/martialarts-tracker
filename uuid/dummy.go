package uuid

var V4Fixture = "b5e57615-0f40-404e-bbe0-6ae81fe8080a"

type Dummy struct {
}

func (g Dummy) Generate() string {
	return V4Fixture
}
