package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "grpc-server/proto/grpc-server/proto" //Aqui se encuentra la ruta donde se encontrara el archivo generado por el Prococol Buffer

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Funcion para leer el el input.txt, aquel input debe de estar al mismo nivel que el codigo de cliente.go
func leerArchivo(filename string) ([]*pb.Orden, error) {
	var ordenes []*pb.Orden
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error al abrir archivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linea := strings.TrimSpace(scanner.Text())
		campos := strings.Split(linea, ",")
		if len(campos) != 5 {
			return nil, fmt.Errorf("formato de archivo incorrecto, se esperan 5 campos por línea")
		}

		valorInt, err := strconv.Atoi(strings.TrimSpace(campos[4]))
		if err != nil {
			log.Fatalf("Error al convertir el valor a entero: %v", err)
		}

		orden := &pb.Orden{
			Caravana: campos[0],
			Escolta:  campos[1],
			Destino:  campos[2],
			Recurso:  campos[3],
			Valor:    int32(valorInt),
		}
		ordenes = append(ordenes, orden)
	}
	return ordenes, nil
}

// Funcion para enviar las ordenes cada 5 segundos desde la lectura del input.txt previo
func enviarOrdenes(client pb.LogisticaServiceClient, ordenes []*pb.Orden) ([]string, error) {
	var seguimientos []string
	for _, orden := range ordenes {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		resp, err := client.EnviarOrden(ctx, orden)
		if err != nil {
			return nil, fmt.Errorf("error al enviar la orden: %v", err)
		}

		fmt.Printf("Orden enviada: Número de seguimiento %s\n", resp.Seguimiento)
		seguimientos = append(seguimientos, resp.Seguimiento)

		// Aqui se esperan 5 segundos entre lecturas para enviar las ordenes
		time.Sleep(5 * time.Second)
	}
	return seguimientos, nil
}

// Funcion para consultar el estado de un paquete con el numero de seguimiento formado por logistica.go
// Aquella consulta se hacen cada 5 segundos
func consultarSeguimiento(client pb.LogisticaServiceClient, seguimientos []string) {
	for _, seguimiento := range seguimientos {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		req := &pb.ConsultaEstadoRequest{Seguimiento: seguimiento}
		estado, err := client.ConsultarEstado(ctx, req)
		if err != nil {
			fmt.Printf("Error al consultar el estado para seguimiento %s: %v\n", seguimiento, err)
			continue
		}
		fmt.Printf("Estado del paquete con seguimiento %s: %+v\n", seguimiento, estado)

		// Esperar 5 segundos entre cada consulta del seguimiento del paquete
		time.Sleep(5 * time.Second)
	}
}

// Funcion principal para simular las acciones del cliente, donde se lle un archivo con todas las ordenes y sus caracteristicas, para luego
// guardar los numeros de seguimiento y hacer consulta del estado se sus paquetes
func main() {

	//Time sleep para que se conecten bien todos los servicios antes de leer el input
	time.Sleep(15 * time.Second)

	conn, err := grpc.Dial("dist038:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor logística: %v", err)
	}
	defer conn.Close()

	client := pb.NewLogisticaServiceClient(conn)

	// Aqui se hace el llamdo del input.txt para procesar las ordenes
	ordenes, err := leerArchivo("input.txt")
	if err != nil {
		log.Fatalf("Error al leer órdenes: %v", err)
	}

	// Enviar órdenes cada 5 segundos y obtener los números de seguimientos generados por logistica.go
	seguimientos, err := enviarOrdenes(client, ordenes)
	if err != nil {
		log.Fatalf("Error al enviar órdenes: %v", err)
	}
	consultarSeguimiento(client, seguimientos)
}
