package main

import (
	"bytes"
	"encoding/json"
	"fmt" // import the fmt package
	"log"
	"net/http" // import the http package
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type searchParams struct {
	Origen      string `json:"origen"`
	Destino     string `json:"destino"`
	FechaSalida string `json:"fecha"`
	Adultos     string `json:"adultos"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type flightsOffer struct {
	ID string `json:"id"`
}

type FlighOffers struct {
	Data []struct {
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
	} `json:"data"`
	Dictionaries struct {
		Locations struct {
			BKK struct {
				CityCode    string `json:"cityCode"`
				CountryCode string `json:"countryCode"`
			} `json:"BKK"`
			SYD struct {
				CityCode    string `json:"cityCode"`
				CountryCode string `json:"countryCode"`
			} `json:"SYD"`
		} `json:"locations"`
		Aircraft struct {
			Num747 string `json:"747"`
		} `json:"aircraft"`
		Currencies struct {
			EUR string `json:"EUR"`
		} `json:"currencies"`
		Carriers struct {
			TG string `json:"TG"`
		} `json:"carriers"`
	} `json:"dictionaries"`
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

type PricingResponse struct {
	Data struct {
		Type         string `json:"type"`
		FlightOffers []struct {
			Type                     string `json:"type"`
			ID                       string `json:"id"`
			Source                   string `json:"source"`
			InstantTicketingRequired bool   `json:"instantTicketingRequired"`
			NonHomogeneous           bool   `json:"nonHomogeneous"`
			LastTicketingDate        string `json:"lastTicketingDate"`
			Itineraries              []struct {
				Segments []struct {
					Departure struct {
						IataCode string `json:"iataCode"`
						At       string `json:"at"`
					} `json:"departure,omitempty"`
					Arrival struct {
						IataCode string `json:"iataCode"`
						Terminal string `json:"terminal"`
						At       string `json:"at"`
					} `json:"arrival,omitempty"`
					CarrierCode string `json:"carrierCode"`
					Number      string `json:"number"`
					Aircraft    struct {
						Code string `json:"code"`
					} `json:"aircraft"`
					Operating struct {
						CarrierCode string `json:"carrierCode"`
					} `json:"operating"`
					ID            string `json:"id"`
					NumberOfStops int    `json:"numberOfStops"`
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
				GrandTotal      string `json:"grandTotal"`
				BillingCurrency string `json:"billingCurrency"`
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
					Taxes    []struct {
						Amount string `json:"amount"`
						Code   string `json:"code"`
					} `json:"taxes"`
				} `json:"price"`
				FareDetailsBySegment []struct {
					SegmentID           string `json:"segmentId"`
					Cabin               string `json:"cabin"`
					FareBasis           string `json:"fareBasis"`
					Class               string `json:"class"`
					IncludedCheckedBags struct {
						Quantity int `json:"quantity"`
					} `json:"includedCheckedBags"`
				} `json:"fareDetailsBySegment"`
			} `json:"travelerPricings"`
		} `json:"flightOffers"`
	} `json:"data"`
}

type BookingRequest struct {
	Data struct {
		FlightOffers []struct {
			Type                     string `json:"type"`
			ID                       string `json:"id"`
			Source                   string `json:"source"`
			InstantTicketingRequired bool   `json:"instantTicketingRequired"`
			NonHomogeneous           bool   `json:"nonHomogeneous"`
			LastTicketingDate        string `json:"lastTicketingDate"`
			Itineraries              []struct {
				Segments []struct {
					Departure struct {
						IataCode string `json:"iataCode"`
						At       string `json:"at"`
					} `json:"departure,omitempty"`
					Arrival struct {
						IataCode string `json:"iataCode"`
						Terminal string `json:"terminal"`
						At       string `json:"at"`
					} `json:"arrival,omitempty"`
					CarrierCode string `json:"carrierCode"`
					Number      string `json:"number"`
					Aircraft    struct {
						Code string `json:"code"`
					} `json:"aircraft"`
					Operating struct {
						CarrierCode string `json:"carrierCode"`
					} `json:"operating"`
					ID            string `json:"id"`
					NumberOfStops int    `json:"numberOfStops"`
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
				GrandTotal      string `json:"grandTotal"`
				BillingCurrency string `json:"billingCurrency"`
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
					Taxes    []struct {
						Amount string `json:"amount"`
						Code   string `json:"code"`
					} `json:"taxes"`
				} `json:"price"`
				FareDetailsBySegment []struct {
					SegmentID           string `json:"segmentId"`
					Cabin               string `json:"cabin"`
					FareBasis           string `json:"fareBasis"`
					Class               string `json:"class"`
					IncludedCheckedBags struct {
						Quantity int `json:"quantity"`
					} `json:"includedCheckedBags"`
				} `json:"fareDetailsBySegment"`
			} `json:"travelerPricings"`
		} `json:"flightOffers"`
		Travelers []struct {
			ID          string `json:"id"`
			DateOfBirth string `json:"dateOfBirth"`
			Name        struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
			} `json:"name"`
			Gender  string `json:"gender"`
			Contact struct {
				EmailAddress string `json:"emailAddress"`
				Phones       []struct {
					DeviceType         string `json:"deviceType"`
					CountryCallingCode string `json:"countryCallingCode"`
					Number             string `json:"number"`
				} `json:"phones"`
			} `json:"contact"`
			Documents []struct {
				DocumentType     string `json:"documentType"`
				BirthPlace       string `json:"birthPlace"`
				IssuanceLocation string `json:"issuanceLocation"`
				IssuanceDate     string `json:"issuanceDate"`
				Number           string `json:"number"`
				ExpiryDate       string `json:"expiryDate"`
				IssuanceCountry  string `json:"issuanceCountry"`
				ValidityCountry  string `json:"validityCountry"`
				Nationality      string `json:"nationality"`
				Holder           bool   `json:"holder"`
			} `json:"documents,omitempty"`
		} `json:"travelers"`
		Type string `json:"type"`
	} `json:"data"`
}

type Traveler struct {
	ID          string `json:"id"`
	DateOfBirth string `json:"dateOfBirth"`
	Name        struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"name"`
	Gender  string `json:"gender"`
	Contact struct {
		EmailAddress string `json:"emailAddress"`
		Phones       []struct {
			DeviceType         string `json:"deviceType"`
			CountryCallingCode string `json:"countryCallingCode"`
			Number             string `json:"number"`
		} `json:"phones"`
	} `json:"contact"`
	Documents []struct {
		DocumentType     string `json:"documentType"`
		BirthPlace       string `json:"birthPlace"`
		IssuanceLocation string `json:"issuanceLocation"`
		IssuanceDate     string `json:"issuanceDate"`
		Number           string `json:"number"`
		ExpiryDate       string `json:"expiryDate"`
		IssuanceCountry  string `json:"issuanceCountry"`
		ValidityCountry  string `json:"validityCountry"`
		Nationality      string `json:"nationality"`
		Holder           bool   `json:"holder"`
	} `json:"documents"`
}

type BookingResponse struct {
	Data struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"data"`
}

func getToken() string {

	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	clientID := os.Getenv("CLIENT_ID")
	fmt.Println("client_id:", clientID)
	clientSecret := os.Getenv("SECRET_ID")
	fmt.Println("client_secret:", clientSecret)

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
	var flightSearchResponse FlighOffers
	err = json.NewDecoder(resp.Body).Decode(&flightSearchResponse)
	if err != nil {
		fmt.Println("Error decoding flight search response:", err)
		return
	}

	for i := range flightSearchResponse.Data {
		for j := range flightSearchResponse.Data[i].Itineraries {
			for k := range flightSearchResponse.Data[i].Itineraries[j].Segments {
				if flightSearchResponse.Data[i].Itineraries[j].Segments[k].Operating.CarrierCode == "" {
					// Si 'operating.carrierCode' está vacío, asigna el valor de 'carrierCode' de nivel superior.
					flightSearchResponse.Data[i].Itineraries[j].Segments[k].Operating.CarrierCode = flightSearchResponse.Data[i].Itineraries[j].Segments[k].CarrierCode
				}
			}
		}
	}
	c.IndentedJSON(http.StatusCreated, flightSearchResponse)
}

func priceHandler(c *gin.Context) { // function that handles the request

	var searchPrice FlightPriceRequest

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&searchPrice); err != nil {
		return
	}
	//fmt.Print(searchPrice)

	var accessToken = getToken()
	pricingData, _ := json.Marshal(searchPrice)
	req, err := http.NewRequest("POST", "https://test.api.amadeus.com/v1/shopping/flight-offers/pricing", bytes.NewBuffer(pricingData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var pricingResponse PricingResponse
	err = json.NewDecoder(resp.Body).Decode(&pricingResponse)
	if err != nil {
		fmt.Println("Error decoding flight search response:", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, pricingResponse)

}

func bookingHandler(c *gin.Context) { // function that handles the request

	var bookingRequest BookingRequest
	var accessToken = getToken()

	if err := c.BindJSON(&bookingRequest); err != nil {
		return
	}
	//fmt.Println("Here is booking request: ", bookingRequest)

	bookingData, _ := json.Marshal(bookingRequest)
	req, err := http.NewRequest("POST", "https://test.api.amadeus.com/v1/booking/flight-orders", bytes.NewBuffer(bookingData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Print("booking status", resp.Status)
	//fmt.Print("booking response", resp)
	defer resp.Body.Close()

	var bookingResponse BookingResponse
	err = json.NewDecoder(resp.Body).Decode(&bookingResponse)
	fmt.Print("booking response", bookingResponse)
	if err != nil {
		fmt.Println("Error decoding flight search response:", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, bookingResponse)
}

func main() {

	router := gin.Default()
	router.GET("/api/search", searchHandler)
	router.POST("/api/pricing", priceHandler)
	router.POST("/api/booking", bookingHandler)

	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	server := os.Getenv("SERVER")
	port := os.Getenv("PORT")

	router.Run(server + ":" + port)

	//http.HandleFunc("/api/search", searchHandler)   // handle the request on /api/search
	//http.HandleFunc("/api/booking", bookingHandler) // handle the request on /api/booking
	fmt.Println("Server is listening on : " + port) // print a message
}
