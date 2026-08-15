package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var conexionBaseDatos *sql.DB

type SolicitudTransferencia struct {
	IDOrigen  int     `json:"id_origen"`
	IDDestino int     `json:"id_destino"`
	Monto     float64 `json:"monto"`
}

type Usuario struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Saldo  float64 `json:"saldo"`
}

func main() {
	var err error
	cadenaConexion := "user=postgres password=sh0ut1nk dbname=transferencias_db sslmode=disable"

	conexionBaseDatos, err = sql.Open("postgres", cadenaConexion)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/api/usuarios", getUsuarios)
	http.HandleFunc("/api/transferir", hacerTransferencia)

	fmt.Println("Servidor corriendo exitosamente en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsuarios(respuestahttp http.ResponseWriter, peticionhttp *http.Request) {
	filasObtenidas, err := conexionBaseDatos.Query("SELECT id, nombre, saldo FROM usuarios")
	if err != nil {
		http.Error(respuestahttp, "Error al consultar usuarios", 500)
		return
	}
	defer filasObtenidas.Close()

	var listaUsuarios []Usuario

	for filasObtenidas.Next() {
		var usuarioActual Usuario
		err := filasObtenidas.Scan(&usuarioActual.ID, &usuarioActual.Nombre, &usuarioActual.Saldo)
		if err != nil {
			continue
		}
		listaUsuarios = append(listaUsuarios, usuarioActual)
	}

	respuestahttp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respuestahttp).Encode(listaUsuarios)
}

func hacerTransferencia(respuestahttp http.ResponseWriter, peticionhttp *http.Request) {
	if peticionhttp.Method != "POST" {
		http.Error(respuestahttp, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var datosSolicitud SolicitudTransferencia
	err := json.NewDecoder(peticionhttp.Body).Decode(&datosSolicitud)
	if err != nil {
		http.Error(respuestahttp, "JSON inválido", 400)
		return
	}

	resultadoResta, err := conexionBaseDatos.Exec(
		"UPDATE usuarios SET saldo = saldo - $1 WHERE id = $2 AND saldo >= $1",
		datosSolicitud.Monto,
		datosSolicitud.IDOrigen,
	)
	if err != nil {
		http.Error(respuestahttp, "Error al restar saldo", 400)
		return
	}

	filasAfectadas, _ := resultadoResta.RowsAffected()
	if filasAfectadas == 0 {
		http.Error(respuestahttp, "Saldo insuficiente u origen inválido", 400)
		return
	}

	_, err = conexionBaseDatos.Exec(
		"UPDATE usuarios SET saldo = saldo + $1 WHERE id = $2",
		datosSolicitud.Monto,
		datosSolicitud.IDDestino,
	)
	if err != nil {
		http.Error(respuestahttp, "Error al sumar saldo", 400)
		return
	}

	_, err = conexionBaseDatos.Exec(
		"INSERT INTO transferencias (id_origen, id_destino, monto) VALUES ($1, $2, $3)",
		datosSolicitud.IDOrigen,
		datosSolicitud.IDDestino,
		datosSolicitud.Monto,
	)
	if err != nil {
		http.Error(respuestahttp, "Error registrando historia", 400)
		return
	}

	respuestahttp.Header().Set("Content-Type", "application/json")
	respuestahttp.WriteHeader(http.StatusOK)
	respuestahttp.Write([]byte(`{"mensaje": "Transferencia exitosa"}`))
}