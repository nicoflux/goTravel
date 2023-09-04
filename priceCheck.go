package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type PricingResponse struct {
	FlightOffers []struct {
		Price struct {
			Total string `json:"total"`
		} `json:"price"`

		TravelerPricings []struct {
			Price struct {
				Total string `json:"total"`
			} `json:"price"`
		} `json:"travelerPricings"`
	} `json:"flightOffers"`
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

type FlightData struct {
	Data []struct {
		PricingDate []struct {
			ID          int `json:"id"`
			Itineraries []struct {
				Segments []struct {
					Departure struct {
						DepartureTime string `json:"at"`
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
				Total float64 `json:"total"`
			} `json:"price"`
		} `json:"pricingDate"`
	} `json:"data"`
}

/*type FlightSearchResponse struct {
	Meta struct {
		Count int `json:"count"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"meta"`
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
				Duration        string `json:"duration"`
				ID              string `json:"id"`
				NumberOfStops   int    `json:"numberOfStops"`
				BlacklistedInEU bool   `json:"blacklistedInEU"`
				Arrival0        struct {
					IataCode string `json:"iataCode"`
					At       string `json:"at"`
				} `json:"arrival,omitempty"`
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
			Bkk struct {
				CityCode    string `json:"cityCode"`
				CountryCode string `json:"countryCode"`
			} `json:"BKK"`
			Mnl struct {
				CityCode    string `json:"cityCode"`
				CountryCode string `json:"countryCode"`
			} `json:"MNL"`
			Syd struct {
				CityCode    string `json:"cityCode"`
				CountryCode string `json:"countryCode"`
			} `json:"SYD"`
		} `json:"locations"`
		Aircraft struct {
			Num320 string `json:"320"`
			Num321 string `json:"321"`
			Num333 string `json:"333"`
		} `json:"aircraft"`
		Currencies struct {
			Eur string `json:"EUR"`
		} `json:"currencies"`
		Carriers struct {
			Pr string `json:"PR"`
		} `json:"carriers"`
	} `json:"dictionaries"`
}*/

func main() {
	// Set your Amadeus API credentials
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
		return
	}
	defer resp.Body.Close()

	// Decode the response JSON
	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Store the access token
	accessToken := tokenResponse.AccessToken
	fmt.Println("Access Token:", accessToken)

	// Define your flight search parameters
	origin := "ARI"
	destination := "SCL"
	departureDate := "2023-12-02"
	adults := 1

	url := fmt.Sprintf("https://test.api.amadeus.com/v2/shopping/flight-offers?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&adults=%v&includedAirlineCodes=H2,LA,JA,H2&nonStop=true&currencyCode=CLP&travelClass=ECONOMY", origin, destination, departureDate, adults)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Close = true
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var responseDatamap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseDatamap); err != nil {
		fmt.Println("Error al decodificar la respuesta JSON del GET:", err)
		return
	}
	//fmt.Println(responseDatamap)
	/*var test_data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&test_data); err != nil {
		fmt.Println("Error al decodificar la respuesta JSON:", err)
		return
	}
	for key, value := range test_data {
		fmt.Printf("Clave: %s, Valor: %v\n", key, value)
	}*/

	// Decode the flight search response JSON (customize FlightSearchResponse structure)
	/*var flightSearchResponse FlightSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&flightSearchResponse)
	if err != nil {
		fmt.Println("Error decoding flight search response:", err)
		return
	}
	*/
	// Process and print flight search results as needed
	//fmt.Println("Flight Search Results:")

	/*for _, flight := range flightSearchResponse.Data {
		fmt.Println("Flight ID:", flight.ID)
		for _, itinerary := range flight.Itineraries {
			for _, segment := range itinerary.Segments {
				fmt.Println("Departure Time:", segment.Departure.DerpartureTime)
				fmt.Println("Arrival Time:", segment.Arrival.ArrivalTime)
				fmt.Println("Carrier Code:", segment.CarrierCode)
				fmt.Println("Flight Number:", segment.CarrierCode+segment.Number)
				fmt.Println("Aircraft Code:", segment.Aircraft.Code)
			}
		}
		fmt.Println("Price:", flight.Price.Total)
	}*/

	/* 	var flightPriceRequest FlightPriceRequest
	   	flightPriceRequest.Data.FlightOffers[0].Itineraries[0].Segments[0].Departure.At = "2023-12-15T01:02:00"
	   	flightPriceRequest.Data.FlightOffers[0].Itineraries[0].Segments[0].Arrival.At = "2023-12-15T03:32:00"
	   	flightPriceRequest.Data.FlightOffers[0].Itineraries[0].Segments[0].CarrierCode = "LA"
	   	flightPriceRequest.Data.FlightOffers[0].Itineraries[0].Segments[0].Number = "197"
	   	flightPriceRequest.Data.FlightOffers[0].Itineraries[0].Segments[0].Aircraft.Code = "320" */

	//Create a POST request for pricing
	pricingURL := "https://test.api.amadeus.com/v2/shopping/flight-offers/pricing"
	//pricingData := flightSearchResponse.Data[0]
	//fmt.Println("princingDate:", flightSearchResponse)
	//fmt.Println("Pricing Request:", pricingData)
	//var flightdata FlightData // Ajusta el tipo de estructura según la respuesta JSON

	/*if err := json.NewDecoder(resp.Body).Decode(&flightdata); err != nil {
		fmt.Println("Error al decodificar la respuesta JSON del GET:", err)
		return
	}*/
	pricingRequestData, _ := json.Marshal(responseDatamap)
	req2, err := http.NewRequest("POST", pricingURL, bytes.NewBuffer(pricingRequestData))
	if err != nil {
		panic(err)
	}
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	req2.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp2, err2 := client.Do(req)
	if err != nil {
		panic(err2)
	}
	defer resp.Body.Close()

	fmt.Println("status code:", resp2.StatusCode)
	//fmt.Println("resp2", resp2)

	var responsePricemap map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&responsePricemap); err != nil {
		fmt.Println("Error al decodificar la respuesta JSON del GET:", err)
		return
	}
	fmt.Println(responsePricemap)

	// Decode the pricing response JSON
	/*var pricingResponse PricingResponse
	err = json.NewDecoder(resp2.Body).Decode(&pricingResponse)
	if err != nil {
		fmt.Println("Error decoding pricing response:", err)
		return
	}

	fmt.Println("Pricing Response:", pricingResponse)
	for _, flight := range pricingResponse.FlightOffers {
		fmt.Println("Price: ", flight.Price.Total)
	}*/
	var response struct {
		Data []struct {
			ID    int `json:"id"`
			Price struct {
				GrandTotal float64 `json:"grandTotal"`
			} `json:"price"`
		} `json:"data"`
	}

	// Decodificar el JSON en la estructura
	if err := json.Unmarshal([]byte(responsePricemap), &response); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return
	}

	// Acceder a los valores específicos
	for _, item := range response.Data {
		fmt.Printf("ID: %d, GrandTotal: %.2f\n", item.ID, item.Price.GrandTotal)
	}
}
