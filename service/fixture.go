package service

type referencingData struct {
	ID         string
	Data       string
	NestedData nestedData
	References referencingDataReferences
}

type referencingDataReferences struct {
	Single   referencedData
	Multiple []referencedData
}

type nestedData struct {
	Data string
}

type referencedData struct {
	ID   string
	Data string
}

var referencingDataFixture = referencingData{
	ID:         referencingIDFixture,
	Data:       referencingValueFixture,
	NestedData: nestedDataFixture,
	References: referencingDataReferences{
		Single:   referencedDataFixture,
		Multiple: []referencedData{referencedDataFixture},
	},
}

var nestedDataFixture = nestedData{
	Data: nestedValueFixture,
}

var referencedDataFixture = referencedData{
	ID:   referencedIDFixture,
	Data: referencedValueFixture,
}

var referencingIDFixture = "referencingID"
var referencedIDFixture = "referencedID"
var referencedValueFixture = "referencedData"
var referencingValueFixture = "value"
var nestedValueFixture = "nestedValue"

var uuidV4Fixture = "b5e57615-0f40-404e-bbe0-6ae81fe8080a"

var missingIDFixture = "123"
