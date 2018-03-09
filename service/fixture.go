package service

type indexedData struct {
	ID                string
	Data              string
	NestedData        nestedData
	NestedIndexedData nestedIndexedData
	MappedIndexedData map[string]deeplyNestedIndexedData
	SlicedIndexedData []deeplyNestedIndexedData
	MappedData        map[string]string
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
	MappedIndexedData: map[string]deeplyNestedIndexedData{
		mapIndexFixture: deeplyNestedIndexedDataFixture,
	},
	SlicedIndexedData: []deeplyNestedIndexedData{deeplyNestedIndexedDataFixture},
	MappedData:        map[string]string{mapIndexFixture: dataValueFixture},
}

var nestedDataFixture = nestedData{
	Data:                    nestedDataValueFixture,
	DeeplyNestedIndexedData: deeplyNestedIndexedDataFixture,
}

var nestedIndexedDataFixture = nestedIndexedData{
	ID:                      nestedIDFixture,
	Data:                    nestedIndexedDataValueFixture,
	DeeplyNestedIndexedData: deeplyNestedIndexedDataFixture,
}

var deeplyNestedIndexedDataFixture = deeplyNestedIndexedData{
	ID:   deeplyNestedIDFixture,
	Data: deeplyNestedDataValueFixture,
}

var idFixture = "myID"
var nestedIDFixture = "myNestedID"
var deeplyNestedIDFixture = "deeplyNestedID"
var deeplyNestedDataValueFixture = "myDeeplyNestedData"

var dataValueFixture = "myData"
var nestedDataValueFixture = "myNestedData"
var nestedIndexedDataValueFixture = "myNestedIndexedData"

var mapIndexFixture = "myMapID"

var uuidV4Fixture = "b5e57615-0f40-404e-bbe0-6ae81fe8080a"
