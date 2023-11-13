package mongodb

import "testing"

func TestGetConnectionUriHost(t *testing.T) {
	mcd := MongoConnectionDetails{
		Host:     "localhost",
		Port:     27017,
		Database: "test",
	}
	uri, err := mcd.GetConnectionURI()
	if err != nil {
		t.Errorf("No error should have been returned when creating uri. Got %+v", err)
	}
	expectedUri := "mongodb://localhost:27017/test"
	if uri != expectedUri {
		t.Errorf("Generated uri should have been %s, was %s", expectedUri, uri)
	}
}

func TestGetConnectionUriUrl(t *testing.T) {
	mcd := MongoConnectionDetails{
		Url:      "mongodb://localhost:27017/test",
		Host:     "not_used",
		Port:     11111,
		Database: "not_used",
	}
	uri, err := mcd.GetConnectionURI()
	if err != nil {
		t.Errorf("No error should have been returned when creating uri. Got %+v", err)
	}
	expectedUri := "mongodb://localhost:27017/test"
	if uri != expectedUri {
		t.Errorf("Generated uri should have been %s, was %s", expectedUri, uri)
	}
}

func TestGetConnectionUriMissingHost(t *testing.T) {
	mcdMissingHost := MongoConnectionDetails{
		Port:     27017,
		Database: "test",
	}
	uri, err := mcdMissingHost.GetConnectionURI()
	if err == nil {
		t.Error("An error should have been returned when creating uri")
	}
	if err.Error() != "missing mongo host" {
		t.Errorf("Error should have been 'missing mongo host', was %s", err.Error())
	}
	if uri != "" {
		t.Errorf("Uri should have been unset, was %s", uri)
	}
}

func TestGetConnectionUriMissingPort(t *testing.T) {
	mcdMissingPort := MongoConnectionDetails{
		Host:     "localhost",
		Database: "test",
	}
	uri, err := mcdMissingPort.GetConnectionURI()
	if err == nil {
		t.Error("An error should have been returned when creating uri")
	}
	if err.Error() != "missing mongo port" {
		t.Errorf("Error should have been 'missing mongo port', was %s", err.Error())
	}
	if uri != "" {
		t.Errorf("Uri should have been unset, was %s", uri)
	}
}

func TestGetConnectionUriMissingDatabase(t *testing.T) {
	mcdMissingDatabase := MongoConnectionDetails{
		Host: "localhost",
		Port: 27017,
	}
	uri, err := mcdMissingDatabase.GetConnectionURI()
	if err == nil {
		t.Error("An error should have been returned when creating uri")
	}
	if err.Error() != "missing mongo database name" {
		t.Errorf("Error should have been 'missing mongo database name', was %s", err.Error())
	}
	if uri != "" {
		t.Errorf("Uri should have been unset, was %s", uri)
	}
}
