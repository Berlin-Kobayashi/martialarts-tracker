package service

type indexedData struct {
	ID                string
	Data              string
	NestedData        nestedData
	NestedIndexedData nestedIndexedData
}

type nestedData struct {
	Data                    string
	DeeplyNestedIndexedData deeplyNestedIndexedData
}

type nestedIndexedData struct {
	ID                      string
	Data                    string
	DeeplyNestedIndexedData deeplyNestedIndexedData
}

type deeplyNestedIndexedData struct {
	ID   string
	Data string
}

var indexedDataFixture = indexedData{
	ID:                idFixture,
	Data:              dataValueFixture,
	NestedData:        nestedDataFixture,
	NestedIndexedData: nestedIndexedDataFixture,
}

var nestedDataFixture = nestedData{
	Data:                    nestedDataValueFixture,
	DeeplyNestedIndexedData: deeplyNestedIndexedDataFixture,
}

var nestedIndexedDataFixture = nestedIndexedData{
	ID:                      nestedIDFixture,
	Data:                    nestedIndexedDataFixtureValue,
	DeeplyNestedIndexedData: deeplyNestedIndexedDataFixture,
}

var deeplyNestedIndexedDataFixture = deeplyNestedIndexedData{
	ID:   deeplyNestedIDFixture,
	Data: "myDeeplyNestedData",
}

var idFixture = "myID"
var nestedIDFixture = "myNestedID"
var deeplyNestedIDFixture = "deeplyNestedID"

var dataValueFixture = "myData"
var nestedDataValueFixture = "myNestedData"
var nestedIndexedDataFixtureValue = "myNestedIndexedData"

var uuidV4Fixture = "b5e57615-0f40-404e-bbe0-6ae81fe8080a"
