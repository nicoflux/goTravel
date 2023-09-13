package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
}

type BookingResponse struct {
	Data struct {
		Type string `json:"type"`
		ID   string `json:"id"`
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

type OrderResponse struct {
	Data struct {
		Type      string `json:"type"`
		ID        string `json:"id"`
		Travelers []struct {
			ID          string `json:"id"`
			DateOfBirth string `json:"dateOfBirth"`
			Gender      string `json:"gender"`
			Name        struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
			} `json:"name"`
			Contact struct {
				EmailAddress string `json:"emailAddress"`
				Phones       []struct {
					CountryCallingCode string `json:"countryCallingCode"`
					Number             string `json:"number"`
				} `json:"phones"`
			} `json:"contact,omitempty"`
		} `json:"travelers"`
		FlightOffers []struct {
			ID          string `json:"id"`
			Type        string `json:"type"`
			Source      string `json:"source"`
			Itineraries []struct {
				Duration string `json:"duration"`
				Segments []struct {
					ID       string `json:"id"`
					Duration string `json:"duration"`
					Aircraft struct {
						Code string `json:"code"`
					} `json:"aircraft"`
					CarrierCode string `json:"carrierCode"`
					Operating   struct {
						CarrierCode string `json:"carrierCode"`
					} `json:"operating"`
					Number    string `json:"number"`
					Departure struct {
						At       string `json:"at"`
						Terminal string `json:"terminal"`
						IataCode string `json:"iataCode"`
					} `json:"departure"`
					Arrival struct {
						At       string `json:"at"`
						Terminal string `json:"terminal"`
						IataCode string `json:"iataCode"`
					} `json:"arrival"`
				} `json:"segments"`
			} `json:"itineraries"`
			Price struct {
				Total string `json:"total"`
			} `json:"price"`
		} `json:"flightOffers"`
	} `json:"data"`
}

func searchHandler() {

	////// Searching For Flights //////
	var search searchParams
	fmt.Print("Aeropuerto de origen: ")
	fmt.Scanln(&search.Origen)
	fmt.Print("Aeropuerto de destino: ")
	fmt.Scanln(&search.Destino)
	fmt.Print("Fecha de salida: ")
	fmt.Scanln(&search.FechaSalida)
	fmt.Print("Cantidad de Adultos: ")
	fmt.Scanln(&search.Adultos)

	jsonData := fmt.Sprintf(`{"origen": "%s","destino": "%s","fecha": "%s","adultos": "%s"}`, search.Origen, search.Destino, search.FechaSalida, search.Adultos)

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
	fmt.Println("Se obtuvieron los siguientes resultados:")

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"VUELO", "NÚMERO", "HORA DE SALIDA", "HORA DE LLEGADA", "AVIÓN", "PRECIO TOTAL"})
	for _, dataItem := range flightSearchResponse.Data {
		for _, itinerary := range dataItem.Itineraries {
			for _, segment := range itinerary.Segments {
				id := dataItem.ID
				dTime := segment.Departure.At
				departureTime := dTime[strings.Index(dTime, "T")+1:]
				aTime := segment.Arrival.At
				arrivalTime := dTime[strings.Index(aTime, "T")+1:]
				//carrierCode := segment.CarrierCode
				flightNumber := segment.CarrierCode + segment.Number
				aircraftCode := "A" + segment.Aircraft.Code
				totalPrice := dataItem.Price.Total
				table.Append([]string{id, flightNumber, departureTime, arrivalTime, aircraftCode, totalPrice})
			}
		}
	}
	table.Render()

	// Selection of flight for booking
	var flightID int
	fmt.Print("Seleccione un vuelo (ingrese 0 para realizar nueva búsqueda): ")
	fmt.Scanln(&flightID)
	if flightID == 0 {
		return
	}

	////// Getting final price of flight //////
	flightPriceData := map[string]interface{}{
		"data": map[string]interface{}{
			"type":         "flight-offers-pricing",
			"flightOffers": []interface{}{flightSearchResponse.Data[flightID-1]},
		},
	}
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

	//////  Booking flight //////
	var travelers []Traveler
	adults, _ := strconv.Atoi(search.Adultos)

	for i := 1; i < adults+1; i++ {

		// Create a new traveler
		var traveler Traveler
		traveler.ID = strconv.Itoa(i)
		fmt.Println("Pasajero ", i, ":")
		fmt.Print("Ingrese fecha de nacimiento: ")
		fmt.Scanln(&traveler.DateOfBirth)
		fmt.Print("Ingrese nombre: ")
		fmt.Scanln(&traveler.Name.FirstName)
		fmt.Print("Ingrese apellido: ")
		fmt.Scanln(&traveler.Name.LastName)
		fmt.Print("Ingrese sexo (MALE o FEMALE): ")
		fmt.Scanln(&traveler.Gender)
		fmt.Print("Ingrese correo: ")
		fmt.Scanln(&traveler.Contact.EmailAddress)
		fmt.Print("Ingrese número de teléfono: ")
		var number string
		fmt.Scanln(&number)

		traveler.Contact.Phones = []struct {
			DeviceType         string `json:"deviceType"`
			CountryCallingCode string `json:"countryCallingCode"`
			Number             string `json:"number"`
		}{
			{
				DeviceType:         "MOBILE",
				CountryCallingCode: "56",
				Number:             number,
			},
		}
		travelers = append(travelers, traveler)
	}

	booking := map[string]interface{}{
		"data": map[string]interface{}{
			"type":         "flight-offers-pricing",
			"flightOffers": []interface{}{pricingResponse.Data.FlightOffers[0]},
			"travelers":    travelers,
		},
	}
	bookingData, _ := json.Marshal(booking)
	req, err = http.NewRequest("POST", "http://127.0.0.1:5000/api/booking", bytes.NewBuffer(bookingData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var bookingResponse BookingResponse
	err = json.NewDecoder(resp.Body).Decode(&bookingResponse)
	if err != nil {
		fmt.Println("Error decoding flight booking response:", err)
		return
	}
	fmt.Print("Reserva creada con éxito: ", bookingResponse.Data.ID)
}

func GetBookingHandler() {

	var orderID string
	fmt.Print("Ingrese el ID de la reserva: ")
	fmt.Scanln(&orderID)
	jsonData := fmt.Sprintf(`{"orderID": "%s""}`, orderID)

	var data = strings.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:5000/api/booking", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var orderResponse OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&orderResponse)
	if err != nil {
		fmt.Println("Error decoding flight order response:", err)
		return
	}

	fmt.Println(orderResponse.Data.FlightOffers)

	fmt.Println("Resultado :")
	tableFlight := tablewriter.NewWriter(os.Stdout)
	tableFlight.SetHeader([]string{"NÚMERO", "HORA DE SALIDA", "HORA DE LLEGADA", "AVIÓN", "PRECIO TOTAL"})
	for _, dataItem := range orderResponse.Data.FlightOffers {
		for _, itinerary := range dataItem.Itineraries {
			for _, segment := range itinerary.Segments {
				dTime := segment.Departure.At
				departureTime := dTime[strings.Index(dTime, "T")+1:]
				aTime := segment.Arrival.At
				arrivalTime := dTime[strings.Index(aTime, "T")+1:]
				flightNumber := segment.CarrierCode + segment.Number
				aircraftCode := "A" + segment.Aircraft.Code
				totalPrice := dataItem.Price.Total
				tableFlight.Append([]string{flightNumber, departureTime, arrivalTime, aircraftCode, totalPrice})
			}
		}
	}
	tableFlight.Render()

	fmt.Println("Pasajeros :")
	tableTravelers := tablewriter.NewWriter(os.Stdout)
	tableTravelers.SetHeader([]string{"NOMBRE", "APELLIDO"})
	for _, traveler := range orderResponse.Data.Travelers {
		nombre := traveler.Name.FirstName
		apellido := traveler.Name.LastName
		tableTravelers.Append([]string{nombre, apellido})
	}
	tableTravelers.Render()
}

func main() {
	initText := `Bievenido a goTravel!`
	fmt.Print(initText)
	text := `
1. Realizar búsqueda.
2. Obtener reserva.
3. Salir
Ingrese una opción:`

	for {
		var input string
		fmt.Print(text)
		fmt.Scanln(&input)

		switch input {
		case "1":
			searchHandler()
		case "2":
			GetBookingHandler()
		case "3":
			fmt.Println("Hasta luego!")
			return

		default: // if the command is not 1, 2 or 3
			fmt.Println("Por favor, ingrese 1, 2 o 3") // print an error message
		}
	}
}
