package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
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
					} `json:"arrival"`
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
					Duration      string `json:"duration"`
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
			PaymentCardRequired bool `json:"paymentCardRequired"`
		} `json:"flightOffers"`
	} `json:"data"`
}

type FlightOffer struct {
	Itineraries []interface{}
	Price       map[string]interface{}
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

	client := &http.Client{}
	var data = strings.NewReader(`{"data":{"type":"flight-offers-pricing","flightOffers":[{"type":"flight-offer","id":"1","source":"GDS","instantTicketingRequired":false,"nonHomogeneous":false,"oneWay":false,"lastTicketingDate":"2023-09-06","numberOfBookableSeats":7,"itineraries":[{"duration":"PT2H30M","segments":[{"departure":{"iataCode":"ARI","terminal":"","at":"2023-12-02T00:57:00"},"arrival":{"iataCode":"SCL","at":"2023-12-02T03:27:00"},"carrierCode":"LA","number":"197","aircraft":{"code":"320"},"operating":{"carrierCode":"LA"},"duration":"PT2H30M","id":"1","numberOfStops":0,"blacklistedInEU":false}]}],"price":{"currency":"CLP","total":"52312.00","base":"34973.00","fees":[{"amount":"0.00","type":"SUPPLIER"},{"amount":"0.00","type":"TICKETING"}],"grandTotal":"52312.00"},"pricingOptions":{"fareType":["PUBLISHED"],"includedCheckedBagsOnly":false},"validatingAirlineCodes":["LA"],"travelerPricings":[{"travelerId":"1","fareOption":"STANDARD","travelerType":"ADULT","price":{"currency":"CLP","total":"52312.00","base":"34973.00"},"fareDetailsBySegment":[{"segmentId":"1","cabin":"ECONOMY","fareBasis":"G00QP5ZJ","class":"G","includedCheckedBags":{"weight":0,"weightUnit":""}}]}]}]}}
}`)
	req, err := http.NewRequest("POST", "https://test.api.amadeus.com/v1/shopping/flight-offers/pricing", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
	fmt.Printf("%s\n", bodyText)

}
