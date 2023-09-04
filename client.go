package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type searchParams struct {
	Origen      string `json:"origen"`
	Destino     string `json:"destino"`
	FechaSalida string `json:"fecha"`
	Adultos     string `json:"adultos"`
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

func searchHandler() {
	var search searchParams
	searchReader := bufio.NewReader(os.Stdin) // create a reader to read from stdin

	fmt.Print("Aeropuerto de origen: ")
	search.Origen, _ = searchReader.ReadString('\n')
	search.Origen = strings.TrimSpace(search.Origen)
	fmt.Print("Aeropuerto de destino: ")
	search.Destino, _ = searchReader.ReadString('\n')
	search.Destino = strings.TrimSpace(search.Destino)
	fmt.Print("Fecha de salida: ")
	search.FechaSalida, _ = searchReader.ReadString('\n')
	search.FechaSalida = strings.TrimSpace(search.FechaSalida) //check format date YYYY-MM-DD
	fmt.Print("Cantitad de Adultos: ")
	search.Adultos, _ = searchReader.ReadString('\n')
	search.Adultos = strings.TrimSpace(search.Adultos)

	jsonData := fmt.Sprintf(`{"origen": "%s","destino": "%s","fecha": "%s","adultos": "%s"}`, search.Origen, search.Destino, search.FechaSalida, search.Adultos)

	// Create a request with the JSON data
	var data = strings.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:5000/api/search", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var flightSearchResponse FlightSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&flightSearchResponse)
	if err != nil {
		fmt.Println("Error decoding flight search response:", err)
		return
	}

	// Process and print flight search results as needed
	fmt.Println("Flight Search Results:")

	for _, flight := range flightSearchResponse.Data {
		fmt.Println("Flight ID:", flight.ID)
		for _, itinerary := range flight.Itineraries {
			for _, segment := range itinerary.Segments {
				fmt.Println("Departure Time:", segment.Departure.DerpartureTime)
				fmt.Println("Arrival Time:", segment.Arrival.ArrivalTime)
				fmt.Println("Flight Number:", segment.CarrierCode+segment.Number)
				fmt.Println("Aircraft Code:", segment.Aircraft.Code)
			}
		}
		fmt.Println("Price:", flight.Price.Total)
	}

	var flightID int
	fmt.Print("Ingrese una opción: ")
	_, err = fmt.Scanf("%d", &flightID)

	var pricingData FlightPriceRequest

	pricingData.Data.FlightOffers[0].Itineraries[0].Segments[0].CarrierCode = flightSearchResponse.Data[flightID].Itineraries[0].Segments[0].CarrierCode
	pricingData.Data.FlightOffers[0].Itineraries[0].Segments[0].Number = flightSearchResponse.Data[flightID].Itineraries[0].Segments[0].Number

	fmt.Println("pricingData:", pricingData.Data.FlightOffers[0].Itineraries[0].Segments[0].CarrierCode, pricingData.Data.FlightOffers[0].Itineraries[0].Segments[0].Number)

	/* 	var data = strings.NewReader(pricingData)
	   	req, err = http.NewRequest("POST", "http://127.0.0.1:5000/api/searchpricing", bytes.NewBuffer(pricingData))
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	req.Header.Set("Content-Type", "application/json")
	   	resp, err = client.Do(req)
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	defer resp.Body.Close()
	   	bodyText, err := io.ReadAll(resp.Body)
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	fmt.Printf("%s\n", bodyText) */

}

// no funciona
func GetBookingHandler() {
	reader := bufio.NewReader(os.Stdin) // create a reader to read from stdin

	fmt.Print("Ingrese el ID de la reserva: ")
	idReserva, _ := reader.ReadString('\n')
	idReserva = strings.TrimSpace(idReserva)
	fmt.Println("Buscando reserva n°", idReserva)

	reservaData := strings.NewReader(idReserva)
	resp, err := http.NewRequest("GET", "http://127.0.0.1:5000/api/booking", reservaData) // send a GET request to the server
	if err != nil {                                                                       // if an error occurred
		fmt.Println("Error:", err) // print the error
		return                     // exit the program
	}
	defer resp.Body.Close() // close the response body when the function returns

	body, err := io.ReadAll(resp.Body) // read the response body
	if err != nil {                    // if an error occurred
		fmt.Println("Error reading response body:", err) // print the error
		return                                           // exit the program
	}

	fmt.Println("Response:", string(body)) // print the response body

}

func main() {
	initText := `Bievenido a goTravel!
1. Realizar búsqueda.
2. Obtener reserva.
3. Salir
Ingrese una opción:`

	for { // infinite loop
		reader := bufio.NewReader(os.Stdin) // create a reader to read from stdin
		fmt.Print(initText)                 // print a message
		input, _ := reader.ReadString('\n') // read from stdin until a newline character is found
		input = strings.TrimSpace(input)    // remove the trailing newline character

		switch input {
		case "1":
			searchHandler()
		case "2":
			GetBookingHandler()
		case "3":
			fmt.Println("Hasta luego!")
			os.Exit(0) // exit the program
			//break ?????

		default: // if the command is not 1, 2 or 3
			fmt.Println("unknown command") // print an error message
		}
	}
}
