package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

// UserData defines the structure of the request payload
type UserData struct {
    UserID           int   `json:"user_id"`
    Favorites        []int `json:"favorites"`
    ViewedProperties []int `json:"viewed_properties"`
}

// Recommendation defines the structure of each recommended property in the response
type Recommendation struct {
    Location  string  `json:"location"`
    Bedrooms  float64 `json:"Bedrooms"`
    Bathrooms float64 `json:"Bathrooms"`
    Toilets   float64 `json:"Toilets"`
    Price     float64 `json:"price"`
}

func main() {
    // Create an instance of UserData with sample data
    userData := UserData{
        UserID:           1,
        Favorites:        []int{15, 18, 200, 25, 26},
        ViewedProperties: []int{56, 60, 91},
    }

    // Convert the UserData instance to JSON
    requestBody, err := json.Marshal(userData)
    if err != nil {
        fmt.Printf("Error encoding JSON: %v\n", err)
        return
    }

    // Send the POST request to the FastAPI endpoint
    response, err := http.Post("http://127.0.0.1:8000/recommend", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        fmt.Printf("Error making POST request: %v\n", err)
        return
    }
    defer response.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Printf("Error reading response body: %v\n", err)
        return
    }

    // Check if the response status is OK
    if response.StatusCode != http.StatusOK {
        fmt.Printf("Error: received status code %d\nResponse: %s\n", response.StatusCode, body)
        return
    }

    // Parse the JSON response
    var recommendations []Recommendation
    err = json.Unmarshal(body, &recommendations)
    if err != nil {
        fmt.Printf("Error decoding JSON response: %v\n", err)
        return
    }

    // Print the recommendations
    fmt.Println("Recommended Properties:")
    for _, rec := range recommendations {
        fmt.Printf("Location: %s, Bedrooms: %.0f, Bathrooms: %.0f, Toilets: %.0f, Price: %.2f\n",
            rec.Location, rec.Bedrooms, rec.Bathrooms, rec.Toilets, rec.Price)
    }
}
