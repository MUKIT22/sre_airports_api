package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)


type Airport struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	IATA    string `json:"iata"`
	ImageURL string `json:"image_url"`
}

type AirportV2 struct {
	Airport
	RunwayLength int `json:"runway_length"`
}

// Mock data for airports in Bangladesh
var airports = []Airport{
	{"Hazrat Shahjalal International Airport", "Dhaka", "DAC", "https://storage.googleapis.com/bd-airport-data/dac.jpg"},
	{"Shah Amanat International Airport", "Chittagong", "CGP", "https://storage.googleapis.com/bd-airport-data/cgp.jpg"},
	{"Osmani International Airport", "Sylhet", "ZYL", "https://storage.googleapis.com/bd-airport-data/zyl.jpg"},
}

// Mock data for airports in Bangladesh (with runway length for V2)
var airportsV2 = []AirportV2{
	{Airport{"Hazrat Shahjalal International Airport", "Dhaka", "DAC", "https://storage.googleapis.com/bd-airport-data/dac.jpg"}, 3200},
	{Airport{"Shah Amanat International Airport", "Chittagong", "CGP", "https://storage.googleapis.com/bd-airport-data/cgp.jpg"}, 2900},
	{Airport{"Osmani International Airport", "Sylhet", "ZYL", "https://storage.googleapis.com/bd-airport-data/zyl.jpg"}, 2500},
}

// HomePage handler
func HomePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: OK"))
}

// Airports handler for the first endpoint
func Airports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airports)
}

// AirportsV2 handler for the second version endpoint
func AirportsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airportsV2)
}

// ##############################
// ## TODO: Edit this function ##
// ##############################

func UpdateAirportImage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string `json:"name"`
		ImageData []byte `json:"image_data"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find the airport by name
	var airport *Airport
	for i := range airports {
		if airports[i].Name == req.Name {
			airport = &airports[i]
			break
		}
	}
	if airport == nil {
		http.Error(w, "Airport not found", http.StatusNotFound)
		return
	}

	// Initialize AWS session with IAM user credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
		return
	}

	// Generate presigned URL
	svc := s3.New(sess)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("my-unique-bucket-name"),
		Key:    aws.String(fmt.Sprintf("%s.jpg", req.Name)),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		http.Error(w, "Failed to generate presigned URL", http.StatusInternalServerError)
		return
	}

	// Upload image using presigned URL
	req, err = http.NewRequest("PUT", urlStr, bytes.NewReader(req.ImageData))
	if err != nil {
		http.Error(w, "Failed to create upload request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "image/jpeg")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	// Update the airport's image URL
	airport.ImageURL = fmt.Sprintf("https://my-unique-bucket-name.s3.amazonaws.com/%s.jpg", req.Name)

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airport)
}


func main() {
	// Setup routes
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/airports", Airports)
	http.HandleFunc("/airports_v2", AirportsV2)

	// TODO: complete the UpdateAirportImage handler function
	http.HandleFunc("/update_airport_image", UpdateAirportImage)

	// Start the server
	http.ListenAndServe(":8080", nil)
}
