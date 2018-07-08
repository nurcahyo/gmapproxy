package place

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
)

var (
	ENDPOINT = "https://maps.googleapis.com/maps/api/place"
	TAG      = reflect.TypeOf(Place{}).PkgPath()
)

type Place struct {
	ID        []byte    `gorm:"type:binary(16);primary"`
	City      string    `gorm:"unique_index:city_country_types"`
	Country   string    `gorm:"unique_index:city_country_types"`
	Types     string    `gorm:"unique_index:city_country_types"`
	Results   string    `gorm:"type:json;null"`
	CreatedAt time.Time `gorm:"type:datetime;" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;"  sql:"DEFAULT:current_timestamp;UPDATE:current_timestamp"`
}

func (p *Place) TableName() string {
	return "places"
}

func (p *Place) BeforeCreate(scope *gorm.Scope) (err error) {
	uuid, err := uuid.NewV4().MarshalBinary()
	scope.SetColumn("ID", uuid)
	return err
}

func (p *Place) GetResultsStruct() []*Response_RESULTS {
	var results []*Response_RESULTS
	json.Unmarshal([]byte(p.Results), &results)
	return results
}

var db *gorm.DB

func init() {
	db, _ = gorm.Open("mysql", "homestead:secret@tcp(localhost)/gmapproxy?charset=utf8&parseTime=true")

	// register model
	db.AutoMigrate(&Place{})
}

func FindNearbyPlaceByCityAndLatLong(in *Request) (out *Response, err error) {

	place := &Place{}
	er := db.Where("country = ? and city = ? and types = ?", in.GetCountry(), in.GetCity(), in.GetTypes()).Find(place).Error

	if er == nil {
		fmt.Println(place.Results)
		if time.Since(place.CreatedAt).Hours() < 720 {
			results := place.GetResultsStruct()
			out = &Response{
				Results: results,
			}
			return out, nil
		}
		db.Delete(place)
	}

	fmt.Println(TAG, er)
	// Search From place api
	fmt.Println(TAG, "Getting places from place api")
	types := strings.Split(in.GetTypes(), ",")
	results := make([]*Response_RESULTS, 0)

	for _, t := range types {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/nearbysearch/json", ENDPOINT), nil)
		q := req.URL.Query()
		q.Add("location", in.GetLatlong())
		q.Add("radius", fmt.Sprint(in.GetRadius()))
		q.Add("type", t)
		q.Add("sensor", "true")
		q.Add("key", in.GetKey())
		req.URL.RawQuery = q.Encode()
		res, e := http.DefaultClient.Do(req)
		if e != nil {
			err = e
			fmt.Println(e.Error())
			continue
		}
		if res.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(res.Body)
			jsonBody := fmt.Sprintf("Error when fetching places from place api. Reason: ( %s )", string(body))
			return nil, errors.New(jsonBody)
		}

		placeResponse := Response{}
		body, _ := ioutil.ReadAll(res.Body)
		jsonBody := string(body)
		json.Unmarshal(body, &placeResponse)
		if placeResponse.GetStatus() != "OK" {

		}
		if len(placeResponse.GetResults()) < 1 {
			return nil, errors.New(jsonBody)
		}
		defer res.Body.Close()
		for _, result := range placeResponse.GetResults() {
			results = append(results, result)
		}
	}
	rawResults, _ := json.Marshal(results)
	place = &Place{
		City:    in.GetCity(),
		Country: in.GetCountry(),
		Types:   in.GetTypes(),
		Results: string(rawResults),
	}
	err = db.Create(place).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	out = &Response{}
	out.Results = results
	return out, err
}
