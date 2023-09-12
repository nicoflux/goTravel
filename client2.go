package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type searchParams struct {
	Origen      string `json:"origen"`
	Destino     string `json:"destino"`
	FechaSalida string `json:"fecha"`
	Adultos     string `json:"adultos"`
}

type FlightOffers struct {
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
	fmt.Print("Cantidad de Adultos: ")
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
	var flightSearchResponse FlightOffers
	err = json.NewDecoder(resp.Body).Decode(&flightSearchResponse)
	if err != nil {
		fmt.Println("Error decoding flight search response:", err)
		return
	}

	// Process and print flight search results as needed
	fmt.Println("Se obtuvieron los siguientes resultados:")

	table := tablewriter.NewWriter(os.Stdout)

	// Definir las cabeceras de la tabla
	table.SetHeader([]string{"VUELO", "NÚMERO", "HORA DE SALIDA", "HORA DE LLEGADA", "AVIÓN", "PRECIO TOTAL"})

	// Recorrer los datos y agregar filas a la tabla
	for _, dataItem := range flightSearchResponse.Data {
		for _, itinerary := range dataItem.Itineraries {
			for _, segment := range itinerary.Segments {
				// Obtener los valores de cada campo
				id := dataItem.ID
				departureTime := segment.Departure.At
				arrivalTime := segment.Arrival.At
				//carrierCode := segment.CarrierCode
				flightNumber := segment.CarrierCode + segment.Number
				aircraftCode := segment.Aircraft.Code
				totalPrice := dataItem.Price.Total

				// Agregar una fila a la tabla
				table.Append([]string{id, flightNumber, departureTime, arrivalTime, aircraftCode, totalPrice})
			}
		}
	}

	// Renderizar la tabla
	table.Render()

	var flightID int
	fmt.Print("Seleccione un vuelo: ")
	_, err = fmt.Scanf("%d", &flightID)

	flightPriceData := map[string]interface{}{
		"data": map[string]interface{}{
			"type":         "flight-offers-pricing",
			"flightOffers": []interface{}{flightSearchResponse.Data[flightID]},
		},
	}
	fmt.Println(flightPriceData)
	pricingData, _ := json.Marshal(flightPriceData)

	req, err = http.NewRequest("POST", "http://127.0.0.1:5000/api/pricing", bytes.NewBuffer(pricingData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var pricingResponse PricingResponse
	err = json.NewDecoder(resp.Body).Decode(&pricingResponse)
	if err != nil {
		fmt.Println("Error decoding flight pricing response:", err)
		return
	}

	fmt.Println("El precio total final es de: ", pricingResponse.Data.FlightOffers[0].Price.Total)

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
	initText := `Bievenido a goTravel!`
	fmt.Print(initText)
	text := `
1. Realizar búsqueda.
2. Obtener reserva.
3. Salir
Ingrese una opción:`

	for { // infinite loop
		reader := bufio.NewReader(os.Stdin) // create a reader to read from stdin
		fmt.Print(text)                     // print a message
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
