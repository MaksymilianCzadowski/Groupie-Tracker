package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

type PageData struct {
	Name       string
	Image      string
	FirstAlbum string
}

type ArtistFullData struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    []string            `json:"locations"`
	ConcertDates []string            `json:"concertDates"`
	Relations    map[string][]string `json:"relations"`
}

type MyArtist struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    string              `json:"locations"`
	ConcertDates string              `json:"concertDates"`
	Relations    map[string][]string `json:"relations"`
}

type MyLocation struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}
type LocationData struct {
	Index []MyLocation `json:"index"`
}

type MyConcertDate struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type ConcertDateData struct {
	Index []MyConcertDate `json:"index"`
}

type MyRelationDate struct {
	ID            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}
type RelationData struct {
	Index []MyRelationDate `json:"index"`
}

var Artists []MyArtist
var ArtistData []ArtistFullData
var LocationsData LocationData
var ConcertDatesData ConcertDateData
var RelationsData RelationData

func GetArtistsData() error {

	resp, err := http.Get(baseURL + "/artists")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Artists)
	return nil
}

func GetLocations() error {
	resp, err := http.Get(baseURL + "/locations")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &LocationsData)
	return nil
}

func GetDates() error {
	resp, err := http.Get(baseURL + "/dates")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &ConcertDatesData)
	return nil
}

func GetRelations() error {
	resp, err := http.Get(baseURL + "/relation")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &RelationsData)
	return nil
}

func GetData() {
	GetArtistsData()
	GetLocations()
	GetDates()
	GetRelations()
	var template ArtistFullData
	// var locate MyLocation
	for i := range Artists {
		template.ID = i + 1
		template.Image = Artists[i].Image
		template.Name = Artists[i].Name
		template.Members = Artists[i].Members
		template.CreationDate = Artists[i].CreationDate
		template.FirstAlbum = Artists[i].FirstAlbum
		template.Locations = LocationsData.Index[i].Locations
		template.ConcertDates = ConcertDatesData.Index[i].Dates
		template.Relations = RelationsData.Index[i].DatesLocation

		ArtistData = append(ArtistData, template)
	}
	return
}

func groupieHandler(w http.ResponseWriter, r *http.Request) {

	filter := r.FormValue("Filter")
	search := r.FormValue("search")
	minDate, _ := strconv.Atoi(r.FormValue("minDate"))
	maxDate, _ := strconv.Atoi(r.FormValue("maxDate"))
	groupNbr, _ := strconv.Atoi(r.FormValue("groupNbr"))
	data := ArtistData

	if search != "" && len(data) != 0 && filter == "" {
		data = Search(search)
	} else if filter != "" || (minDate != 0 && maxDate != 0) || groupNbr != 0 && len(data) != 0 {
		data = FilterForType(filter, minDate, maxDate, groupNbr)

	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, data)

}

func FilterForType(filter string, minDate int, maxDate int, groupNbr int) []ArtistFullData {
	var data []ArtistFullData
	// fmt.Println(filter)
	if filter == "group" {
		data = FilterForArtist(false)

	} else if filter == "artist" {
		data = FilterForArtist(true)
	}
	if minDate != 0 && maxDate != 0 {
		data = FilterForDate(minDate, maxDate)
	}
	if groupNbr != 0 {
		data = FilterForNumber(groupNbr)
	}
	return data
}

func FilterForArtist(Group bool) []ArtistFullData {
	var data []ArtistFullData
	var temp ArtistFullData

	for i := range Artists {
		temp.Members = Artists[i].Members
		if len(temp.Members) == 1 && Group {
			data = append(data, getDatabyId(i))
		} else if !(Group) && len(temp.Members) > 1 {
			data = append(data, getDatabyId(i))
		}
	}
	return data
}

func FilterForDate(minDate int, maxDate int) []ArtistFullData {
	var data []ArtistFullData
	var temp ArtistFullData

	for i := range Artists {
		temp.CreationDate = Artists[i].CreationDate
		if temp.CreationDate >= minDate && temp.CreationDate <= maxDate {
			data = append(data, getDatabyId(i))
		}
	}
	return data
}

func FilterForNumber(groupNbr int) []ArtistFullData {
	var data []ArtistFullData
	var temp ArtistFullData

	for i := range Artists {
		temp.Members = Artists[i].Members
		if len(temp.Members) == groupNbr {
			data = append(data, getDatabyId(i))
		}
	}
	return data
}

func getDatabyId(id int) ArtistFullData {
	var data ArtistFullData

	for i := range Artists {
		if i == id {
			// fmt.Println("id de getId =", i)
			data.ID = Artists[i].ID
			data.Image = Artists[i].Image
			data.Name = Artists[i].Name
			data.Members = Artists[i].Members
			data.CreationDate = Artists[i].CreationDate
			data.FirstAlbum = Artists[i].FirstAlbum
			data.Locations = LocationsData.Index[i].Locations
			data.ConcertDates = ConcertDatesData.Index[i].Dates
			data.Relations = RelationsData.Index[i].DatesLocation
			break
		}
	}
	return data
}

func getLocationById(id int) MyLocation {
	var data MyLocation

	for i := range LocationsData.Index {
		if i == id {
			data.ID = id
			data.Locations = LocationsData.Index[i].Locations
			data.Dates = LocationsData.Index[i].Dates
		}
	}
	return data
}

func getDateById(id int) MyConcertDate {
	var data MyConcertDate

	for i := range ConcertDatesData.Index {
		if i == id {
			data.ID = id
			data.Dates = ConcertDatesData.Index[i].Dates
		}
	}
	return data
}

func getRelationById(id int) MyRelationDate {
	var data MyRelationDate

	for i := range RelationsData.Index {
		if i == id {
			data.ID = id
			data.DatesLocation = RelationsData.Index[i].DatesLocation
		}
	}
	return data
}

func Search(search string) []ArtistFullData {
	if search == "" {
		return ArtistData
	}
	var resultSearch []ArtistFullData
	search = strings.ToLower(search)
	reg := regexp.MustCompile(`^` + search)
	for i := range Artists {
		temp := strings.ToLower(Artists[i].Name)
		if reg.Match([]byte(temp)) {
			// fmt.Println("search id = ", i)
			resultSearch = append(resultSearch, getDatabyId(i))
		}
	}

	return resultSearch
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Path[12:])
	data := getDatabyId(id - 1)
	x := template.Must(template.New("").Parse(src))
	t, _ := template.ParseFiles("ArtistPage.html")
	t.Execute(w, data)
	x.Execute(w, data)
}

func main() {

	GetData()
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/", fs)
	http.HandleFunc("/groupie-tracker", groupieHandler)
	http.HandleFunc("/ArtistPage/", artistHandler)
	http.ListenAndServe(":8080", nil)
}

const src = `<script type="text/javascript"> 

function initMap() {
  const map = new google.maps.Map(document.getElementById("map"), {
    zoom: 3,
    center: { lat: 48.8566, lng: 2.3522 },
  });
  const geocoder = new google.maps.Geocoder();
  {{ range $value := .Locations }}  
  geocoder.geocode({ address: {{$value}}}, (results, status) => {
      new google.maps.Marker({
        map,
        position: results[0].geometry.location,
      });
  });
  {{end}}
}
        </script> `
