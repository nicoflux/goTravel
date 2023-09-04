package main

import (
	"bytes"
	"encoding/json"
	"fmt"      // import the fmt package
	"net/http" // import the http package

	"github.com/gin-gonic/gin"
)

type searchParams struct {
	Origen      string `json:"origen"`
	Destino     string `json:"destino"`
	FechaSalida string `json:"fecha"`
	Adultos     string `json:"adultos"`
}

// Read implements io.Reader.
func (searchParams) Read(p []byte) (n int, err error) {
	panic("unimplemented")
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type flightsOffer struct {
	ID string `json:"id"`
}

type FlightSearchResponse struct {
	Data []struct {
		ID          string `json:"id"`
		Itineraries []struct {
			Segments []struct {
				Departure struct {
					DerpartureTime string `json:"at"`
				} `json:"departure"`
				Arrival struct {
					ArrivalTime string `json:"at"`
				} `json:"arrival"`
				CarrierCode string `json:"carrierCode"`
				Number      string `json:"number"`
				Aircraft    struct {
					Code string `json:"code"`
				} `json:"aircraft"`
			} `json:"segments"`
		} `json:"itineraries"`
		Price struct {
			Total string `json:"total"`
		} `json:"price,omitempty"`
	} `json:"data"`
}

type FlightPriceRequest struct {
	Data struct {
		Type         string `json:"type"`
		FlightOffers []struct {
			Type                     string `json:"type"`
			ID                       string `json:"id"`
			Source                   string `json:"source"`
			InstantTicketingRequired bool   `json:"instantTicketingRequired"`
			NonHomogeneous           bool   `json:"nonHomogeneous"`
			OneWay                   bool   `json:"oneWay"`
			LastTicketingDate        string `json:"lastTicketingDate"`
			NumberOfBookableSeats    int    `json:"numberOfBookableSeats"`
			Itineraries              []struct {
				Duration string `json:"duration"`
				Segments []struct {
					Departure struct {
						IataCode string `json:"iataCode"`
						Terminal string `json:"terminal"`
						At       string `json:"at"`
					} `json:"departure"`
					Arrival struct {
						IataCode string `json:"iataCode"`
						Terminal string `json:"terminal"`
						At       string `json:"at"`
					} `json:"arrival"`
					CarrierCode string `json:"carrierCode"`
					Number      string `json:"number"`
					Aircraft    struct {
						Code string `json:"code"`
					} `json:"aircraft"`
					Operating struct {
						CarrierCode string `json:"carrierCode"`
					} `json:"operating"`
					Duration        string `json:"duration"`
					ID              string `json:"id"`
					NumberOfStops   int    `json:"numberOfStops"`
					BlacklistedInEU bool   `json:"blacklistedInEU"`
				} `json:"segments"`
			} `json:"itineraries"`
			Price struct {
				Currency string `json:"currency"`
				Total    string `json:"total"`
				Base     string `json:"base"`
				Fees     []struct {
					Amount string `json:"amount"`
					Type   string `json:"type"`
				} `json:"fees"`
				GrandTotal string `json:"grandTotal"`
			} `json:"price"`
			PricingOptions struct {
				FareType                []string `json:"fareType"`
				IncludedCheckedBagsOnly bool     `json:"includedCheckedBagsOnly"`
			} `json:"pricingOptions"`
			ValidatingAirlineCodes []string `json:"validatingAirlineCodes"`
			TravelerPricings       []struct {
				TravelerID   string `json:"travelerId"`
				FareOption   string `json:"fareOption"`
				TravelerType string `json:"travelerType"`
				Price        struct {
					Currency string `json:"currency"`
					Total    string `json:"total"`
					Base     string `json:"base"`
				} `json:"price"`
				FareDetailsBySegment []struct {
					SegmentID           string `json:"segmentId"`
					Cabin               string `json:"cabin"`
					FareBasis           string `json:"fareBasis"`
					Class               string `json:"class"`
					IncludedCheckedBags struct {
						Weight     int    `json:"weight"`
						WeightUnit string `json:"weightUnit"`
					} `json:"includedCheckedBags"`
				} `json:"fareDetailsBySegment"`
			} `json:"travelerPricings"`
		} `json:"flightOffers"`
	} `json:"data"`
}

type booking struct {
	NAME string `json:"name"`
}

var flights = flightsOffer{
	ID: "123",
}

func getToken() string {
	clientID := "FAMrbww1l8rybRMtw8tiHgLcuWVBB7Z9"
	clientSecret := "iGwczbmWUgBjN2Ny"

	// Set the token endpoint URL
	tokenURL := "https://test.api.amadeus.com/v1/security/oauth2/token"

	// Create a POST request payload
	tokenRequestData := bytes.NewBufferString(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", clientID, clientSecret))

	// Make the POST request to obtain the token
	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", tokenRequestData)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "null"
	}
	defer resp.Body.Close()

	// Decode the response JSON
	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "null"
	}

	// Store the access token
	accessToken := tokenResponse.AccessToken
	fmt.Println("Access Token:", accessToken)
	return accessToken
}

func searchHandler(c *gin.Context) { // function that handles the request
	var search searchParams

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&search); err != nil {
		return
	}

	var accessToken = getToken()
	//fmt.Println("origen:", search.Origen)
	//fmt.Println("destino:", search.Destino)
	//fmt.Println("fecha:", search.FechaSalida)
	//fmt.Println("adultos:", search.Adultos)

	url := fmt.Sprintf("https://test.api.amadeus.com/v2/shopping/flight-offers?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&adults=%v&includedAirlineCodes=LA,JA,H2&nonStop=true&currencyCode=CLP&travelClass=ECONOMY", search.Origen, search.Destino, search.FechaSalida, search.Adultos)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Decode the flight search response JSON (customize FlightSearchResponse structure)
	var flightSearchResponse FlightSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&flightSearchResponse)
	if err != nil {
		fmt.Println("Error decoding flight search response:", err)
		return
	}
	c.IndentedJSON(http.StatusCreated, flightSearchResponse)
}

func priceHandler(c *gin.Context) { // function that handles the request
	//var accessToken = getToken()
}

func bookingHandler(c *gin.Context) { // function that handles the request
	/* if r.Method != http.MethodGet { // if the request method is not GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // return an error
		return
	} */

	var newBooking booking

	if err := c.BindJSON(&newBooking); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, newBooking)
}

func main() {

	router := gin.Default()
	router.GET("/api/search", searchHandler)
	router.GET("/api/pricing", priceHandler)
	//router.POST("/api/booking", bookingHandler)

	router.Run("127.0.0.1:5000")

	//http.HandleFunc("/api/search", searchHandler)   // handle the request on /api/search
	//http.HandleFunc("/api/booking", bookingHandler) // handle the request on /api/booking
	fmt.Println("Server is listening on : 5000") // print a message
	//http.ListenAndServe("127.0.0.1:5000", nil)      // listen on port 5000
}
