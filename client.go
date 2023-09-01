package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func searchInput() string {
	searchReader := bufio.NewReader(os.Stdin) // create a reader to read from stdin

	fmt.Print("Aeropuerto de origen: ")
	origen, _ := searchReader.ReadString('\n')
	origen = strings.TrimSpace(origen)
	fmt.Print("Aeropuerto de destino: ")
	destino, _ := searchReader.ReadString('\n')
	destino = strings.TrimSpace(destino)
	fmt.Print("Fecha de salida: ")
	fecha, _ := searchReader.ReadString('\n')
	fecha = strings.TrimSpace(fecha) //check format date YYYY-MM-DD
	fmt.Print("Cantitad de Adultos: ")
	cantidad, _ := searchReader.ReadString('\n')
	cantidad = strings.TrimSpace(cantidad)

	stringsToMerge := []string{origen, destino, fecha, cantidad}
	separator := ";"
	response := strings.Join(stringsToMerge, separator)
	return response
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
			search := searchInput()
			fmt.Println("searching for ", search)

			resp, err := http.Get("http://127.0.0.1:5000/api/search") // send a GET request to the server
			if err != nil {                                           // if an error occurred
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

		case "2":
			fmt.Print("Ingrese el ID de la reserva: ")
			idReserva, _ := reader.ReadString('\n')
			idReserva = strings.TrimSpace(idReserva)
			fmt.Println("Buscando reserva n°", idReserva)

			resp, err := http.Get("http://127.0.0.1:5000/api/booking") // send a GET request to the server
			if err != nil {                                            // if an error occurred
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

		case "3":
			fmt.Println("Hasta luego!")
			os.Exit(0) // exit the program
			//break ?????

		default: // if the command is not 1, 2 or 3
			fmt.Println("unknown command") // print an error message
		}
	}
}
