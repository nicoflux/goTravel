package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

/* type PricingResponse struct {
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
} */

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

/* type FlightPrice struct {
	Data struct {
		Type         string        `json:"type"`
		FlightOffers []interface{} `json:"flightOffers"`
	} `json:"data"`
} */

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

type BookingRequest struct {
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
		Travelers []Traveler
		Remarks   struct {
			General []struct {
				SubType string `json:"subType"`
				Text    string `json:"text"`
			} `json:"general"`
		} `json:"remarks"`
		TicketingAgreement struct {
			Option string `json:"option"`
			Delay  string `json:"delay"`
		} `json:"ticketingAgreement"`
		Contacts []struct {
			AddresseeName struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
			} `json:"addresseeName"`
			CompanyName string `json:"companyName"`
			Purpose     string `json:"purpose"`
			Phones      []struct {
				DeviceType         string `json:"deviceType"`
				CountryCallingCode string `json:"countryCallingCode"`
				Number             string `json:"number"`
			} `json:"phones"`
			EmailAddress string `json:"emailAddress"`
			Address      struct {
				Lines       []string `json:"lines"`
				PostalCode  string   `json:"postalCode"`
				CityName    string   `json:"cityName"`
				CountryCode string   `json:"countryCode"`
			} `json:"address"`
		} `json:"contacts"`
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

	url := fmt.Sprintf("https://test.api.amadeus.com/v2/shopping/flight-offers?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&adults=%v&includedAirlineCodes=H2,LA,HA&nonStop=true&currencyCode=CLP&travelClass=ECONOMY", origin, destination, departureDate, adults)
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
	totalPrice := flightSearchResponse.Data[0].Price.Total
	fmt.Printf("Total Price: %s\n", totalPrice)

	flightOfferJSON, err := json.Marshal(flightSearchResponse.Data[0])
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}
	flightOffersJSON := `{"data":{"type": "flight-offers-pricing","flightOffers": [`
	flightOffersJSON += string(flightOfferJSON) // Supongamos que flightSearchResponse.Data[0] es un JSON válido
	flightOffersJSON += `]}}`

	flightData := map[string]interface{}{
		"data": map[string]interface{}{
			"type":         "flight-offers-pricing",
			"flightOffers": []interface{}{flightSearchResponse.Data[0]},
		},
	}

	pricingData, _ := json.Marshal(flightData)
	req2, err2 := http.NewRequest("POST", "https://test.api.amadeus.com/v1/shopping/flight-offers/pricing", bytes.NewBuffer(pricingData))

	if err != nil {
		panic(err2)
	}
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	req2.Header.Set("Content-Type", "application/json")

	client2 := &http.Client{}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer resp.Body.Close()

	var pricingResponse PricingResponse
	err = json.NewDecoder(resp2.Body).Decode(&pricingResponse)
	if err != nil {
		fmt.Println("Error decoding flight search response:", err)
		return
	}
	fmt.Println("Status Code Pricing request: ", resp2.StatusCode)
	/* 	bodyText, err3 := io.ReadAll(resp2.Body)
	   	if err3 != nil {
	   		log.Fatal(err3)
	   	}
	   	fmt.Println(resp.StatusCode)
	   	fmt.Printf("%s\n", bodyText) */

	var bookingRequest BookingRequest
	var traveler Traveler
	bookingRequest.Data.FlightOffers = pricingResponse.Data.FlightOffers
	traveler.DateOfBirth = "1998-03-10"
	traveler.Name.LastName = "Gonzalez"
	traveler.Name.FirstName = "Jorge"
	traveler.Gender = "male"
	traveler.Contact.EmailAddress = "jorge.gonzalez833@telefonica.es"
	//append(traveler.Contact.Phones, ) := "123456789" //country code to 56

	bookingRequest.Data.TicketingAgreement.Delay = "3D"

	// Set valid option
	bookingRequest.Data.TicketingAgreement.Option = "DELAY_TO_CANCEL"
	fmt.Println("-------------------------------------------------------")
	fmt.Println(traveler)
	fmt.Println("-------------------------------------------------------")
	bookingRequest.Data.Travelers = append(bookingRequest.Data.Travelers, traveler)
	fmt.Println(bookingRequest.Data.Travelers)
	fmt.Println("-------------------------------------------------------")

	/* 	bookingRequest.Data.Travelers[0].DateOfBirth = birthDate
	   	bookingRequest.Data.Travelers[0].Name.FirstName = firstName
	   	bookingRequest.Data.Travelers[0].Name.LastName = lastName
	   	bookingRequest.Data.Travelers[0].Gender = gender
	   	bookingRequest.Data.Travelers[0].Contact.EmailAddress = email
	   	bookingRequest.Data.Travelers[0].Contact.Phones[0].DeviceType = "MOBILE"
	   	bookingRequest.Data.Travelers[0].Contact.Phones[0].CountryCallingCode = "56"
	   	bookingRequest.Data.Travelers[0].Contact.Phones[0].Number = phoneNumber */

	bookingData, _ := json.Marshal(bookingRequest)
	req3, err3 := http.NewRequest("POST", "https://test.api.amadeus.com/v1/booking/flight-orders", bytes.NewBuffer(bookingData))
	if err != nil {
		log.Fatal(err3)
	}
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", "Bearer "+accessToken)
	resp3, err3 := client.Do(req3)
	if err3 != nil {
		log.Fatal(err3)
	}
	defer resp3.Body.Close()
	bodyText3, err := io.ReadAll(resp3.Body)
	if err3 != nil {
		log.Fatal(err3)
	}
	fmt.Printf("%s\n", bodyText3)

}
