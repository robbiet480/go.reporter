package reporter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AppleEpochTime is an arbitrary epoch date that Apple uses in NSDate on iOS and Mac.
// It is defined as January 1st, 2001 00:00:00 UTC.
var AppleEpochTime = time.Unix(978307200, 0)

// ISO8601 is the standard ISO 8601 timestamp format for Go
const ISO8601 = "2006-01-02T15:04:05-0700"

// DateTime is a special wrapper around time.Time due to complexities around schema differences.
// In version 1 of the schema, timestamps were expressed in seconds since Apple epoch.
// In version 2 of the schema, the app started using standard ISO 8601 timestamps
type DateTime struct{ time.Time }

func (d *DateTime) String() string {
	if SchemaVersion == 1 {
		return strconv.FormatFloat(d.Sub(AppleEpochTime).Seconds(), 'f', -1, 64)
	}
	return d.Format(ISO8601)
}

// MarshalJSON is needed to return either a date string that is ISO 8601 formatted (schema v2) or the number of seconds since Apple epoch (schema v1)
func (d *DateTime) MarshalJSON() ([]byte, error) {
	if SchemaVersion == 1 {
		return json.Marshal(d.Sub(AppleEpochTime).Seconds())
	}
	return json.Marshal(d.Format(ISO8601))
}

// UnmarshalJSON handles deserialization of a timestamp.
// This custom unmarshaling is needed because the input property may be an ISO 8601 timestamp
// or number of seconds since Apple Epoch (January 1st, 2001 00:00:00 UTC)
func (d *DateTime) UnmarshalJSON(data []byte) (err error) {
	var dateTime time.Time
	dateString, rawJSON := "", json.RawMessage{}
	if err = json.Unmarshal(data, &dateString); err == nil {
		dateTime, err = time.Parse(ISO8601, dateString)
		if err != nil {
			return
		}
		SchemaVersion = 2
		d.Time = dateTime
		return
	}
	if err = json.Unmarshal(data, &rawJSON); err == nil {
		var inputDuration time.Duration
		inputDuration, err = time.ParseDuration(string(rawJSON) + "s")
		if err != nil {
			return
		}
		// BUG(robbiet480): For now, this returns older style timestamps in local time according to computer setting
		dateTime = AppleEpochTime.Add(inputDuration).Local()
		SchemaVersion = 1
		d.Time = dateTime
		return
	}
	return
}

// ConnectionType is a struct indicating the network connection of the device at the time of the report.
// The value for the connection attribute cooresponds to the following events:
//
// 0: Device is connected via cellular network
//
// 1: Device is connected via WiFi
//
// 2: Device is not connected
type ConnectionType struct {
	Method      string
	Description string
	Type        int `json:"connection,omitempty"`
}

func (c *ConnectionType) String() string { return c.Method }

// MarshalJSON is needed to return only the Connection type integer
func (c *ConnectionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Type)
}

// UnmarshalJSON provides custom JSON unmarshaling for ConnectionType.
// Given the connection integer, it adds in human readable connection types and descriptions.
func (c *ConnectionType) UnmarshalJSON(data []byte) error {
	var cType int
	if err := json.Unmarshal(data, &cType); err != nil {
		return fmt.Errorf("Connection type should be an int, got %s", data)
	}
	switch cType {
	case 0:
		c.Method = "Cellular"
		c.Description = "Device is connected via cellular network"
	case 1:
		c.Method = "Wi-Fi"
		c.Description = "Device is connected via WiFi"
	case 2:
		c.Method = "Not connected"
		c.Description = "Device is not connected"
	}
	c.Type = cType
	return nil
}

// A ReportImpetus struct indicates how the report was triggered.
// The value for the impetus attribute cooresponds to the following events:
//
// 0: Report button tapped
//
// 1: Report button tapped while Reporter is asleep
//
// 2: Report triggered by notification
//
// 3: Report triggered by setting app to sleep
//
// 4: Report triggered by waking up app
type ReportImpetus struct {
	Description string
	Impetus     int
}

func (r *ReportImpetus) String() string { return r.Description }

// MarshalJSON is needed to return only the report impetus integer
func (r *ReportImpetus) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Impetus)
}

// UnmarshalJSON provides custom JSON unmarshaling for ReportImpetus.
// Given the reportImpetus integer, it adds in a human readable description of the impetus.
func (r *ReportImpetus) UnmarshalJSON(data []byte) error {
	var reportImpetus int
	if err := json.Unmarshal(data, &reportImpetus); err != nil {
		return fmt.Errorf("Connection type should be an int, got %s", data)
	}
	switch reportImpetus {
	case 0:
		r.Description = "Report button tapped"
	case 1:
		r.Description = "Report button tapped while Reporter is asleep"
	case 2:
		r.Description = "Report triggered by notification"
	case 3:
		r.Description = "Report triggered by setting app to sleep"
	case 4:
		r.Description = "Report triggered by waking up app"
	}
	r.Impetus = reportImpetus
	return nil
}

// Photo struct contains the EXIF metadata of a single photo.
// Additionally, the photo struct contains a link to the photo asset within iOS.
// Currently, this information is unused witin the Reporter application and is not of much use outside the iOS system.
type Photo struct {
	ID                string    `json:"uniqueIdentifier,omitempty"`
	Altitude          *float64  `json:"altitude,omitempty"`
	ApertureValue     *float64  `json:"apertureValue,omitempty"`
	AssetURL          string    `json:"assetUrl,omitempty"`
	BrightnessValue   *float64  `json:"brightnessValue,omitempty"`
	DateTime          *DateTime `json:"dateTime,omitempty"`
	Depth             *int      `json:"depth,omitempty"`
	ExposureMode      *int      `json:"exposureMode,omitempty"`
	ExposureProgram   *int      `json:"exposureProgram,omitempty"`
	ExposureTime      *float64  `json:"exposureTime,omitempty"`
	FNumber           *float64  `json:"fNumber,omitempty"`
	Flash             *int      `json:"flash,omitempty"`
	FocalLength       *float64  `json:"focalLength,omitempty"`
	FocalLengthIn35mm *int      `json:"focalLengthIn35mm,omitempty"`
	IsoSpeed          *int      `json:"isoSpeed,omitempty"`
	Latitude          *float64  `json:"latitude,omitempty"`
	LatitudeRef       string    `json:"latitudeRef,omitempty"`
	Longitude         *float64  `json:"longitude,omitempty"`
	LongitudeRef      string    `json:"longitudeRef,omitempty"`
	Make              string    `json:"make,omitempty"`
	MeteringMode      *int      `json:"meteringMode,omitempty"`
	Model             string    `json:"model,omitempty"`
	Orientation       *int      `json:"orientation,omitempty"`
	PixelHeight       *int      `json:"pixelHeight,omitempty"`
	PixelWidth        *int      `json:"pixelWidth,omitempty"`
	ResolutionUnit    *int      `json:"resolutionUnit,omitempty"`
	SceneCaptureType  *int      `json:"sceneCaptureType,omitempty"`
	SensingMode       *int      `json:"sensingMode,omitempty"`
	ShutterSpeed      *float64  `json:"shutterSpeed,omitempty"`
	Software          string    `json:"software,omitempty"`
	WhiteBalance      *int      `json:"whiteBalance,omitempty"`
}

// PhotoSet is a struct with a single array of photos written to the snapshot if the user has taken photos between reports.
type PhotoSet struct {
	ID     string  `json:"uniqueIdentifier,omitempty"`
	Photos []Photo `json:"photos,omitempty"`
}

// Altitude is a struct containing detailed altitude information at the time of the report.
type Altitude struct {
	ID                      string   `json:"uniqueIdentifier,omitempty"`
	AdjustedPressure        *float64 `json:"adjustedPressure,omitempty"`
	FloorsAscended          *int     `json:"floorsAscended,omitempty"`
	FloorsDescended         *int     `json:"floorsDescended,omitempty"`
	GPSAltitudeFromLocation *float64 `json:"gpsAltitudeFromLocation,omitempty"`
	GPSRawAltitude          *float64 `json:"gpsRawAltitude,omitempty"`
	Pressure                *float64 `json:"pressure,omitempty"`
}

// Audio is measured decibels, which is "a logarithmic unit used to express the ratio between two values of a physical quantity, often power or intensity."
// Because it is easier to define a reference sound at the upper limit (where the microphone is overloaded and "clips"), decibels are often expressed as negative values.
// This is true for the iPhone, so the values that are delivered in this property are the raw output from the iOS CoreAudio API, reflecting the average and peak volume recorded over a single second.
// The lower the number, the quieter the noise. The closer the number is to zero (where the audio would clip), the louder the ambient noise.
type Audio struct {
	ID      string   `json:"uniqueIdentifier,omitempty"`
	Average *float64 `json:"avg,omitempty"`
	Peak    *float64 `json:"peak,omitempty"`
}

// PositiveAverageDb does the same calculation the app does to show a positive Db average value instead of the standard negative Db.
// According to the author, Nick Felton (as per https://gist.github.com/dbreunig/9315705#gistcomment-1350866):
// Here's the conversion we are using:
// The raw value from Apple (-160 dB to 0 dB) is in the JSON output.
// Simply adding to shift the scale makes for nonsense “dBA” values (e.g. 105 dB in a quiet room).
// We very roughly approximated our display value so it seemed reasonable in this way:
// (x + 65) * 2 where x is the raw value Apple gives us, again, -160 dB to 0 dB.
// You can still use the raw values from Apple (in JSON) and apply any correction or calibration as they see to be appropriate.
func (a *Audio) PositiveAverageDb(rounded bool) float64 {
	value := (float64(*a.Average) + float64(65)) * 2
	if rounded {
		return roundPlus(value, 2)
	}
	return value
}

// PositivePeakDb does the same calculation the app does to show a positive Db peak value instead of the standard negative Db.
func (a *Audio) PositivePeakDb(rounded bool) float64 {
	value := (float64(*a.Peak) + float64(65)) * 2
	if rounded {
		return roundPlus(value, 2)
	}
	return value
}

// A Region is a struct containing a parsed CLPlacemark Region
type Region struct {
	Latitude   float64 `json:"-"`
	Longitude  float64 `json:"-"`
	Radius     float64 `json:"-"`
	Identifier string  `json:"-"`
}

func (r *Region) String() string { return r.Identifier }

// MarshalJSON is needed to return only the Region identifier
func (r *Region) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Identifier)
}

// UnmarshalJSON provides custom JSON unmarshaling for region.
// It splits up the string format of CLPlacemark.region into a Golang usable format
func (r *Region) UnmarshalJSON(b []byte) (err error) {
	var placemark string
	if err = json.Unmarshal(b, &placemark); err == nil {
		replacer := strings.NewReplacer("<", "", ">", "", ",", " ", "+", "")
		cleanedString := replacer.Replace(placemark)
		splitFields := strings.Fields(cleanedString)
		lat, err := strconv.ParseFloat(splitFields[0], 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(splitFields[1], 64)
		if err != nil {
			return err
		}
		rad, err := strconv.ParseFloat(splitFields[3], 64)
		if err != nil {
			return err
		}
		r.Identifier = placemark
		r.Latitude = lat
		r.Longitude = lon
		r.Radius = rad
	}
	return err
}

// Placemark struct is the result of reverse geocoding the latitude and longitude deribed from iOS's location services.
// It will often get addresses wrong, but will usually be accurate with ZIP, county, neighborhood, city, and state attributes.
type Placemark struct {
	ID                    string  `json:"uniqueIdentifier,omitempty"`
	SubAdministrativeArea string  `json:"subAdministrativeArea,omitempty"`
	SubLocality           string  `json:"subLocality,omitempty"`
	SubThoroughfare       string  `json:"subThoroughfare,omitempty"`
	Thoroughfare          string  `json:"thoroughfare,omitempty"`
	AdministrativeArea    string  `json:"administrativeArea,omitempty"`
	PostalCode            string  `json:"postalCode,omitempty"`
	Region                *Region `json:"region,omitempty"`
	Country               string  `json:"country,omitempty"`
	Locality              string  `json:"locality,omitempty"`
	Name                  string  `json:"name,omitempty"`
}

// A Location struct is essentially a CoreLocation CLLocation (https://developer.apple.com/library/ios/documentation/CoreLocation/Reference/CLLocation_Class/CLLocation/CLLocation.html#//apple_ref/doc/uid/TP40007126) object, with a CLPlacemark embedded (https://developer.apple.com/library/ios/documentation/CoreLocation/Reference/CLPlacemark_class/Reference/Reference.html#//apple_ref/doc/uid/TP40009574).
// Refer to the linked documentation for each class for details on their properties.
type Location struct {
	ID                 string     `json:"uniqueIdentifier,omitempty"`
	Speed              *int       `json:"speed,omitempty"`
	Placemark          *Placemark `json:"placemark,omitempty"`
	Timestamp          *DateTime  `json:"timestamp,omitempty"`
	Longitude          *float64   `json:"longitude,omitempty"`
	Latitude           *float64   `json:"latitude,omitempty"`
	VerticalAccuracy   *float64   `json:"verticalAccuracy,omitempty"`
	Course             *int       `json:"course,omitempty"`
	Altitude           *float64   `json:"altitude,omitempty"`
	HorizontalAccuracy *float64   `json:"horizontalAccuracy,omitempty"`
}

// The Weather struct is perhaps the most self-explanitory of the data captured.
// struct keys are descriptive, detailing the metric and the units used.
type Weather struct {
	ID                        string   `json:"uniqueIdentifier,omitempty"`
	RelativeHumidity          string   `json:"relativeHumidity,omitempty"`
	VisibilityKilometers      *float64 `json:"visibilityKM,omitempty"`
	TemperatureCelsius        *float64 `json:"tempC,omitempty"`
	PrecipitationTodayInches  *float64 `json:"precipTodayIn,omitempty"`
	WindKilometersPerHour     *float64 `json:"windKPH,omitempty"`
	WindDegrees               *int     `json:"windDegrees,omitempty"`
	Latitude                  *float64 `json:"latitude,omitempty"`
	StationID                 string   `json:"stationID,omitempty"`
	VisibilityMiles           *float64 `json:"visibilityMi,omitempty"`
	PressureInches            *float64 `json:"pressureIn,omitempty"`
	PressureMillibars         *float64 `json:"pressureMb,omitempty"`
	FeelsLikeFarenheit        *float64 `json:"feelslikeF,omitempty"`
	Longitude                 *float64 `json:"longitude,omitempty"`
	FeelsLikeCelsius          *float64 `json:"feelslikeC,omitempty"`
	TemperatureFarenheit      *float64 `json:"tempF,omitempty"`
	PrecipitationTodayMetric  *float64 `json:"precipTodayMetric,omitempty"`
	WindGustKilometersPerHour *float64 `json:"windGustKPH,omitempty"`
	WindDirection             string   `json:"windDirection,omitempty"`
	DewPoint                  *float64 `json:"dewpointC,omitempty"`
	UVIndex                   *float64 `json:"uv,omitempty"`
	WeatherDescription        string   `json:"weather,omitempty"`
	WindGustMilesPerHour      *float64 `json:"windGustMPH,omitempty"`
	WindMilesPerHour          *float64 `json:"windMPH,omitempty"`
}

// Token is an individual common repsonses, either words or phrases
type Token struct {
	ID   string `json:"uniqueIdentifier,omitempty"`
	Text string `json:"text,omitempty"`
}

type token Token

func (t *Token) String() string { return t.Text }

// MarshalJSON is needed to return either a Token object with uniqueIdentifier (schema v2) or a single text element (schema v1)
func (t *Token) MarshalJSON() ([]byte, error) {
	if SchemaVersion == 1 {
		return json.Marshal(t.Text)
	}
	return json.Marshal(*t)
}

// UnmarshalJSON provides custom JSON unmarshaling for Token.
// In version 1 of the schema, tokens were expressed as arrays of strings.
// In version 2 of the schema, the app started expressing tokens as arrays of objects containing uniqueIdentifier and text
func (t *Token) UnmarshalJSON(b []byte) (err error) {
	j, n := token{}, ""
	if err = json.Unmarshal(b, &j); err == nil {
		*t = Token(j)
		SchemaVersion = 2
		return
	}
	if err = json.Unmarshal(b, &n); err == nil {
		t.Text = n
		SchemaVersion = 1
	}
	return
}

// A LocationResponse describes the users present location at the time of the report.
// The locationResponse includes the current location data from the iOS location services API
// as well as a foursquareVenueID, which is provided by the FourSquare Venues Platform API.
type LocationResponse struct {
	ID                string    `json:"uniqueIdentifier,omitempty"`
	Text              string    `json:"text,omitempty"`
	Location          *Location `json:"location,omitempty"`
	FoursquareVenueID string    `json:"foursquareVenueId,omitempty"`
}

// TextResponse contains free form, user generated text
type TextResponse struct {
	ID   string `json:"uniqueIdentifier,omitempty"`
	Text string `json:"text,omitempty"`
}

// Response is a struct containing any information entered by the user in Reporter survey questions.
// Each question answered is captured as a single struct within the array, containing the questionPrompt and the user input or selected responses.
// If a question is not answered, it will not be written to the array.
type Response struct {
	ID              string            `json:"uniqueIdentifier,omitempty"`
	Tokens          []*Token          `json:"tokens,omitempty"`
	AnsweredOptions []string          `json:"answeredOptions,omitempty"`
	Location        *LocationResponse `json:"locationResponse,omitempty"`
	QuestionPrompt  string            `json:"questionPrompt,omitempty"`
	NumericResponse string            `json:"numericResponse,omitempty"`
	TextResponses   []*TextResponse   `json:"textResponses,omitempty"` // v2
	TextResponse    string            `json:"textResponse,omitempty"`  // v1
}

// A Snapshot is single report for the day
type Snapshot struct {
	ID                string          `json:"uniqueIdentifier,omitempty"`  //
	Steps             *int            `json:"steps,omitempty"`             // The steps property provides a single numerical value reflecting the number of steps taken between the last report filed and the current report. It is only captured if the user is using an iPhone 5S or above, which features the M7 motion coprocessor.
	Responses         []*Response     `json:"responses,omitempty"`         //
	Battery           *float64        `json:"battery,omitempty"`           // The battery key refers to a double numerical value, between 0 and 1, reflecting the power stored in the iPhone's battery at the time of report.
	SectionIdentifier string          `json:"sectionIdentifier,omitempty"` // A convenience variable used by the application when displaying reports in a UITableView.
	Audio             *Audio          `json:"audio,omitempty"`             //
	Background        *int            `json:"background,omitempty"`        // A state variable indicating the report was captured in the background. We are not captuing reports in the background. Therefore, this attribute is not in use.
	Date              *DateTime       `json:"date,omitempty"`              //
	Day               *DateTime       `json:"day,omitempty"`               //
	Location          *Location       `json:"location,omitempty"`          //
	PhotoSet          *PhotoSet       `json:"photoSet,omitempty"`          //
	Weather           *Weather        `json:"weather,omitempty"`           //
	Connection        *ConnectionType `json:"connection,omitempty"`        // The connection attribute indicates the current network connection of the device.
	Altitude          *Altitude       `json:"altitude,omitempty"`          //
	ReportImpetus     *ReportImpetus  `json:"reportImpetus,omitempty"`     // The attribute reportImpetus indicates how the report was triggered.
	Draft             *int            `json:"draft,omitempty"`             // A state variable indicating the report is being edited. If it is, it won't be saved. Therefore, this will always be 0.
	DwellStatus       *int            `json:"dwellStatus,omitempty"`       // Debug variable. Not in use.
	Sync              *int            `json:"sync,omitempty"`              // This is a state variable to ensure each report is saved to Dropbox. It will always be 0 because once it is 1 (or true) the app will not attempt to write it to Dropbox.
}
