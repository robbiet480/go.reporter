package reporter

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/robbiet480/go.reporter"
)

func thingToMap(t *testing.T, thing []byte) map[string]interface{} {
	var outputJSON map[string]interface{}
	err := json.Unmarshal(thing, &outputJSON)
	if err != nil {
		t.Fatal(err)
	}
	return outputJSON
}

func compareOutput(t *testing.T, filepath string) {
	if filepath == "" {
		t.Fatal("No filepath given!")
	}
	fileJSON, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}
	parsedJSON, err := json.Marshal(loadTestFile(t, filepath))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("JSON from go.reporter", string(parsedJSON))
	t.Log("JSON from filesystem", string(fileJSON))
	fileJSONMap := thingToMap(t, []byte(fileJSON))
	parsedJSONMap := thingToMap(t, parsedJSON)
	if reflect.DeepEqual(parsedJSONMap, fileJSONMap) {
		t.Log("Test file JSON matches output JSON of go.reporter")
	} else {
		t.Fatal("Test file JSON does NOT match output JSON of go.reporter")
	}
}

func loadTestFile(t *testing.T, filePath string) (day reporter.Day) {
	backend, err := reporter.NewFilesystemBackend("")
	if err != nil {
		t.Fatal(err)
	}
	fileFromBackend, err := backend.GetReportForPath(filePath)
	if err != nil {
		t.Fatal(err)
	}
	day, err = reporter.DecodeFile(fileFromBackend)
	if err != nil {
		t.Fatal(err)
	}
	return
}

func TestDecodeFileVersionOne(t *testing.T) {
	compareOutput(t, "./testData/2014-01-15-reporter-export.json")
}

func TestDecodeFileVersionTwo(t *testing.T) {
	compareOutput(t, "./testData/2015-10-23-reporter-export.json")
}

func TestAudioPositiveAverageDb(t *testing.T) {
	day := loadTestFile(t, "./testData/2015-10-23-reporter-export.json")
	latestSnapshot := day.GetLatestSnapshot()
	if latestSnapshot.Audio.PositiveAverageDb(true) != 12.32 {
		t.Errorf("Positive Db average does not match expected value! We were expecting 12.32 but got %d", latestSnapshot.Audio.PositiveAverageDb(true))
	}
	if latestSnapshot.Audio.PositiveAverageDb(false) != 12.318460000000002 {
		t.Errorf("Positive Db average does not match expected value! We were expecting 12.32 but got %d", latestSnapshot.Audio.PositiveAverageDb(true))
	}
}
