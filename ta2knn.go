package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Hotel struct {
	Id        int
	Clase     int
	Nombre    string
	Direccion string
	Latitud   float64
	Longitud  float64
}

var num = 0
var NLatitud float64
var NLongitud float64
var listaPuntos [30]Hotel
var seleccionado int

func setCoordenadas(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")

	io.WriteString(res,
		`<!DOCTYPE html>
		<html>
		<body>
		
		<p>Click the button to get your coordinates.</p>
		
		<button onclick="getLocation()">Try It</button>
		
		<form action="/coordenadas" method="POST">
			<label for="LatN">Nueva Latitud</label>
			<input type="text" name="LatitudN" id="LatN" placeholder="Nueva Latitud">
			<label for="LatN">Nueva Longitud</label>
			<input type="text" name="LongitudN" id="LongN" placeholder="Nueva Longitud">
			<input type="submit">
		</form>
		
		<p id="demo"></p>
		
		<script>
		var x = document.getElementById("demo");
		
		function getLocation() {
		  if (navigator.geolocation) {
			navigator.geolocation.watchPosition(showPosition);
		  } else { 
			x.innerHTML = "Geolocation is not supported by this browser.";
		  }
		}
			
		function showPosition(position) {
			x.innerHTML="Latitude: " + position.coords.latitude + 
			"<br>Longitude: " + position.coords.longitude;
		}
		</script>
		
		</body>
		</html>`)

	NLat := req.FormValue("LatitudN")
	NLong := req.FormValue("LongitudN")
	NLatitud, _ := strconv.ParseFloat(NLat, 64)
	NLongitud, _ := strconv.ParseFloat(NLong, 64)
	fmt.Println(NLatitud)
	fmt.Println(NLongitud)

}

func listResultado(res http.ResponseWriter, req *http.Request) {

	go Calcular(listaPuntos[num:num+5], NLatitud, NLongitud)

	//seleccionado := <-chPunto

	res.Header().Set("Content-Type", "application/json")

	jsonBytes, _ := json.MarshalIndent(seleccionado, "", " ")

	io.WriteString(res, string(jsonBytes))

}

//endpoint del servicio

func handleRequest() {
	http.HandleFunc("/coordenadas", setCoordenadas)
	http.HandleFunc("/resultado", listResultado)

	//como va a exponer los endpoint, porque puerto se va a escuchar

	log.Fatal(http.ListenAndServe(":9000", nil))

}

func readJSONFromUrl(url string) ([]Hotel, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var hotelsList []Hotel
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	body := bytes.TrimPrefix(respByte, []byte("\xef\xbb\xbf"))

	if err := json.Unmarshal([]byte(body), &hotelsList); err != nil {
		return nil, err
	}

	return hotelsList, nil
}

func Calcular(dato []Hotel, lat float64, long float64) {

	//seleccionado := <-chPunto

	var listaResultado [15]float64
	var listaPunto [15]Hotel
	for i := 0; i < len(dato); i++ {
		x := math.Pow(2, dato[i].Latitud-lat)
		y := math.Pow(2, dato[i].Longitud-long)

		resultado := math.Sqrt(x + y)
		listaResultado[i] = resultado
		listaPunto[i] = dato[i]
	}

	min := listaResultado[0]
	for i := 0; i < len(listaResultado); i++ {
		if listaResultado[i] < min {
			min = listaResultado[i]
			listaPunto[i] = dato[i]
			seleccionado = dato[i].Clase
			//seleccionado <- listaPunto[i]
		}
	}
}

func SaveValues() {
	url := "https://raw.githubusercontent.com/ErnestoVC/friendly-chainsaw/main/FixedData.json"
	hotelsList, err := readJSONFromUrl(url)

	if err != nil {
		panic(err)
	}
	var i int = 0
	for i <= 29 {
		listaPuntos[i].Id = hotelsList[i].Id
		listaPuntos[i].Clase = hotelsList[i].Clase
		listaPuntos[i].Nombre = hotelsList[i].Nombre
		listaPuntos[i].Direccion = hotelsList[i].Direccion
		listaPuntos[i].Latitud = hotelsList[i].Latitud
		listaPuntos[i].Longitud = hotelsList[i].Longitud

		i++

	}
}
func main() {

	go SaveValues()
	handleRequest()

	//chPunto = make(chan Hotel, 1)

	//fmt.Println(data[0])
}
