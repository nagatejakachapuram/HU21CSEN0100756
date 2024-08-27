package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	// Define the registration details
	data := map[string]string{
		"companyName": "Afformed",            // Replace with your company name
		"ownerName":   "Kachapuram Nagateja", // Replace with your name
		"rollNo":      "HU21CSEN0100756",     // Replace with your roll number
		"ownerEmail":  "nkachapu@gitam.in",   // Replace with your email
		"accessCode":  "RHFsxX",              // Replace with the access code provided to you
	}

	// Convert the data to JSON
	jsonData, _ := json.Marshal(data)

	// Send the POST request
	resp, err := http.Post("http://20.244.56.144/test/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Print the response status
	log.Println("Response Status:", resp.Status)
}
