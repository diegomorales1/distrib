package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	pb "grpc-server/proto/grpc-server/proto" //Aqui se encuentra la ruta donde se encontrara el archivo generado por el Prococol Buffer

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var registroMutex sync.Mutex

// Esta funcion escribe en un archivo txt todos los paquetes entrantes para tener un registro de cada uno
func escribirRegistro(paquete *pb.Paquete, registro *pb.Registro, entregado bool) {
	registroMutex.Lock()
	defer registroMutex.Unlock()

	file, err := os.OpenFile("Registro_Paquetes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error al abrir el archivo de registro: %v\n", err)
		return
	}
	defer file.Close()

	fechaEntrega := "0"
	if entregado {
		fechaEntrega = time.Now().Format("02-01-2006 15:04")
	}

	linea := fmt.Sprintf("%06d %s %d %s %s %d %s\n",
		paquete.IdPaquete,
		paquete.Tipo,
		paquete.Valor,
		registro.Escolta,
		registro.Destino,
		paquete.Intentos,
		fechaEntrega,
	)

	if _, err := file.WriteString(linea); err != nil {
		fmt.Printf("Error al escribir en el archivo de registro: %v\n", err)
	}
}

type CaravanaServer struct {
	pb.UnimplementedCaravanaServiceServer
}

// Funcion para conectar con Logistica
func conectarLogistica() (pb.LogisticaServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial("dist038:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a logística: %v", err)
	}
	return pb.NewLogisticaServiceClient(conn), conn
}

// En esta funcion se procesan los paquetes para la caravana Ostronitas, con la logica implementada de los 3 intentos como maximo
func procesarPaquetesOstronitas(paquete *pb.Paquete, registro *pb.Registro, clienteLogistica pb.LogisticaServiceClient) {
	intentosMax := 3
	entregado := false

	for paquete.Intentos < int32(intentosMax) {
		time.Sleep(3 * time.Second) // Simular tiempo de entrega

		if rand.Float32() < 0.85 {
			fmt.Printf("Paquete %d entregado con éxito por caravana Ostronitas\n", paquete.IdPaquete)
			paquete.Estado = "Entregado"
			entregado = true
			paquete.Intentos++
			break
		} else {
			fmt.Printf("Paquete %d no pudo ser entregado por caravana Ostronitas.\n", paquete.IdPaquete)
			paquete.Intentos++
		}
	}

	if paquete.Estado != "Entregado" {
		paquete.Estado = "No entregado"
	}

	actualizarEstadoLogistica(clienteLogistica, paquete)
	escribirRegistro(paquete, registro, entregado)
}

// Aqui se procesan los paquetes de Griner, en otras palabras, se simula la entrega de los paquetes que son para las caravanas
// de tipo Prioritario y Normal
func procesarPaquetesGrineer(paquete *pb.Paquete, registro *pb.Registro, clienteLogistica pb.LogisticaServiceClient) {
	maxIntentos := int32(2)
	entregado := false
	valorPaquete := int(paquete.Valor)
	for paquete.Intentos < maxIntentos && valorPaquete > 0 {
		time.Sleep(3 * time.Second)
		if rand.Float32() < 0.85 {
			fmt.Printf("Paquete %d entregado con éxito por caravana %s\n", paquete.IdPaquete, paquete.Tipo)
			paquete.Estado = "Entregado"
			paquete.Intentos++
			entregado = true
			break
		} else {
			fmt.Printf("Paquete %d no pudo ser entregado por caravana %s. Intentando de nuevo...\n", paquete.IdPaquete, paquete.Tipo)
			paquete.Intentos++
			valorPaquete -= 100
		}
	}
	if paquete.Estado != "Entregado" {
		paquete.Estado = "No entregado"
	}
	actualizarEstadoLogistica(clienteLogistica, paquete)
	escribirRegistro(paquete, registro, entregado)
}

// Actualizar estado del paquete en logística
func actualizarEstadoLogistica(clienteLogistica pb.LogisticaServiceClient, paquete *pb.Paquete) {
	_, err := clienteLogistica.ActualizarEstadoPaquete(context.Background(), &pb.ActualizarEstadoRequest{
		IdPaquete: paquete.IdPaquete,
		Estado:    paquete.Estado,
		Intentos:  paquete.Intentos,
	})
	if err != nil {
		log.Printf("Error al actualizar estado del paquete %d: %v", paquete.IdPaquete, err)
	}
}

// Aqui se procesan todos los paquetes enviados por el apartado de logística
func (c *CaravanaServer) ProcesarPaquetes(ctx context.Context, req *pb.ProcesarPaquetesRequest) (*pb.ProcesarPaquetesResponse, error) {
	if len(req.Paquetes) == 0 {
		return nil, fmt.Errorf("no hay paquetes para procesar")
	}

	clienteLogistica, conn := conectarLogistica()
	defer conn.Close()

	for i, paquete := range req.Paquetes {
		registro := req.Registros[i]
		switch paquete.Tipo {
		case "Ostronitas":
			procesarPaquetesOstronitas(paquete, registro, clienteLogistica)
		case "Prioritario", "Normal":
			procesarPaquetesGrineer(paquete, registro, clienteLogistica)
		default:
			fmt.Printf("Tipo de caravana desconocido: %s\n", paquete.Tipo)
		}
	}

	fmt.Println("Caravana completó la entrega de los paquetes.")
	return &pb.ProcesarPaquetesResponse{}, nil
}

// Main donde se prendera el server de caravana.go y el llamado de las funciones para procesar los paquetes
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50052") // Aqui definimos el puerto que utilizara el server de caravana
	if err != nil {
		log.Fatalf("Error al iniciar el servidor de caravana: %v", err)
	}

	grpcServer := grpc.NewServer()
	caravanaServer := &CaravanaServer{}
	pb.RegisterCaravanaServiceServer(grpcServer, caravanaServer)

	log.Printf("Servidor de caravanas escuchando en :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al iniciar el servidor gRPC de caravana: %v", err)
	}
}
