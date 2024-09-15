package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type Paquete struct {
	Id_paquete  int32
	Seguimiento string
	Tipo        string
	Valor       int32
	Intentos    int32
	Estado      string
}

type Registro_cuentas struct {
	Id_paquete         int32
	Envios_completados bool
	IntentosEntrega    int32
	No_Entregado       bool
	Ganancia           int32
}

// Calcula la ganancia del paquete descontanto por intentos mayores a 1
func Calcular_ganancia(paquete Paquete) int32 {
	estado := paquete.Estado
	tipo_entrega := paquete.Tipo
	valor_paq := paquete.Valor
	intentos := paquete.Intentos
	var valor int32
	valor = 0
	costo_intentos := 100 * (intentos - 1)
	if estado == "Entregado" {
		// Verifico el tipo de envio que se realizo

		if tipo_entrega == "Ostronitas" {
			// Intentos antes de la entrega
			valor = valor_paq - costo_intentos

		} else if tipo_entrega == "Prioritario" {
			valor = int32(math.Round((float64(valor_paq) * 0.3))) + valor_paq
			valor = valor - costo_intentos

		} else {
			//Normal
			valor = valor_paq - costo_intentos
		}
	} else if estado == "No Entregado" {

		if tipo_entrega == "Normal" {
			//No tiene ganancia asociado al paquete
			valor = -costo_intentos

		} else if tipo_entrega == "Ostronitas" {
			//El valor del paquete como ganancia
			valor = valor_paq - costo_intentos

		} else {
			//Prioritario, asegura 30% de ganancia del paquete
			valor = int32(math.Round((float64(valor_paq) * 0.3))) - costo_intentos
		}
	}

	//Aun no hay calculo de costo asociado
	return valor
}

// Actualizo la informacion que se tiene en los registros
// Va cambiando la info de los paquete a medida que logistica le entrega actualizaciones
// Si un paquete no esta en los registros, lo crea
func Actualizar_Cuentas(cuentas *[]Registro_cuentas, paquete Paquete) {
	ganancia := Calcular_ganancia(paquete)
	estado := paquete.Estado

	//Buscar el paquete en el registro financiero para actualizarlo
	for i := range *cuentas {
		if (*cuentas)[i].Id_paquete == paquete.Id_paquete {
			(*cuentas)[i].IntentosEntrega = paquete.Intentos
			(*cuentas)[i].Ganancia = ganancia
			if estado == "Entregado" {
				(*cuentas)[i].Envios_completados = true
				(*cuentas)[i].No_Entregado = false
			} else if estado == "No Entregado" {
				(*cuentas)[i].No_Entregado = true
				(*cuentas)[i].Envios_completados = false
			} else {
				//En camino o en Cetus
				(*cuentas)[i].No_Entregado = false
				(*cuentas)[i].Envios_completados = false
			}
			//finaliza la funcion si encuentra el registro
			return
		}
	}

	//Si el paquete no esta en el registro, se agrega
	nuevoRegistro := Registro_cuentas{
		Id_paquete:         paquete.Id_paquete,
		Envios_completados: estado == "Entregado", //Si es igual se hace true, caso contrario false
		IntentosEntrega:    paquete.Intentos,
		No_Entregado:       estado == "No Entregado", //Si es igual se hace true, caso contrario false
		Ganancia:           ganancia,
	}
	*cuentas = append(*cuentas, nuevoRegistro)
}

// Muestra las ganancias finales de los registros xd
func MostrarGanancia(cuentas []Registro_cuentas) {
	var ganancia_Final int32 = 0
	for _, registro := range cuentas {
		ganancia_Final += registro.Ganancia
	}
	fmt.Printf("Balance final en Créditos: %d\n", ganancia_Final)
}

// Escribe en un archivo "registro_cuentas" los registros de finanzas
func EscribirRegistros(cuentas []Registro_cuentas) {
	// Abrir o crear el archivo
	file, err := os.Create("registro_cuentas.txt")
	if err != nil {
		log.Fatalf("No se pudo crear el archivo: %v", err)
	}
	defer file.Close()

	// Escribir cada cuenta en el archivo
	for _, cuenta := range cuentas {
		_, err := file.WriteString(
			fmt.Sprintf("ID Paquete: %d, Envios Completados: %t, Intentos de Entrega: %d, No Entregado: %t, Ganancia: %d\n",
				cuenta.Id_paquete, cuenta.Envios_completados, cuenta.IntentosEntrega, cuenta.No_Entregado, cuenta.Ganancia))
		if err != nil {
			log.Fatalf("No se pudo escribir en el archivo: %v", err)
		}
	}

	fmt.Println("Registros escritos en registro_cuentas.txt")
}

// Trata de reconectarse a rabbitmq
func connectToRabbitMQ() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ { // Intentar 5 veces con un tiempo de espera
		conn, err = amqp.Dial("amqp://guest:guest@10.35.168.48:5672/")
		if err == nil {
			return conn, nil
		}
		log.Printf("Fallo al conectar a RabbitMQ: %v. Reintentando en 5 segundos...", err)
		time.Sleep(5 * time.Second)
	}
	return nil, err
}

// Funcion main para la conexion con rabbit y la espera de mensajes de logistica
// Si se demoran mas de 60s en recibir actualizaciones, finanzas se detiene y muestra la ganancia total
func main() {

	//Trato de conectar con rabbitmq (se demora en prender)
	conn, err := connectToRabbitMQ()
	if err != nil {
		log.Fatalf("Fallo al conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Fallo al abrir un canal: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"finanzas", // Nombre de la cola
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Fallo al declarar la cola: %v", err)
	}

	//Canal para recibir mensajes de logistica
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Fallo al registrar el consumidor: %v", err)
	}

	var cuentas []Registro_cuentas

	// Temporizador de 60 segundos
	timeout := time.NewTimer(60 * time.Second)

	go func() {
		for d := range msgs {
			var paquete Paquete
			err := json.Unmarshal(d.Body, &paquete)
			if err != nil {
				log.Printf("Fallo al deserializar el paquete: %v", err)
				continue
			} else {
				log.Printf("Info de paquete recibida correctamente")
			}

			// Temporizador cada vez que recibe mensaje
			timeout.Reset(60 * time.Second)

			// Actualizar las cuentas con la información recibida
			Actualizar_Cuentas(&cuentas, paquete)
			log.Printf("Paquete procesado: %+v", paquete)
		}
	}()

	log.Printf("Esperando mensajes..")

	// Espera hasta que el temporizador se muera
	<-timeout.C

	// Muestra el balance final antes de cerrar el programa
	MostrarGanancia(cuentas)
	// Escribe en txt los registros de finanzas
	EscribirRegistros(cuentas)
	log.Printf("Finanzas a determinado que se terminaron de enviar los paquetes")
}
