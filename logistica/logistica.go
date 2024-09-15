package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	pb "grpc-server/proto/grpc-server/proto" //Aqui se encuentra la ruta donde se encontrara el archivo generado por el Prococol Buffer

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Aqui se genera un timer general para manejar el envio entre paquetes, para simular el recibimiento y el re intento de entregar los paquetes
var timers = map[string]*time.Timer{}
var mu sync.Mutex

type Server struct {
	pb.UnimplementedLogisticaServiceServer
	mu                   sync.Mutex
	ordenes              []Registro
	paquetesOstronitas   []Paquete
	paquetesPrioritarios []Paquete
	paquetesNormales     []Paquete
	listaPaquetes        []Paquete
	siguienteID          int32
	caravanas            map[string]*Caravana
}

type Registro struct {
	Timestamp   time.Time
	Id_paquete  int32
	Tipo        string
	Nombre      string
	Valor       int
	Escolta     string
	Destino     string
	Seguimiento string
}

type Paquete struct {
	Id_paquete  int32
	Seguimiento string
	Tipo        string
	Valor       int32
	Intentos    int32
	Estado      string
}

type Caravana struct {
	Estado  string
	Cliente pb.CaravanaServiceClient // Cliente gRPC para enviar paquetes a la caravana
}

// Aqui se genera la conexión a RabbitMQ y envío de mensajes
func sendToQueue(body string) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir un canal: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"finanzas", // nombre de la cola de RabbitMQ
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar una cola: %s", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("Error al enviar mensaje a la cola: %s", err)
	}
	log.Printf("Mensaje enviado a RabbitMQ: %s", body)
}

// En esta funcion se escribe en archivo Registro_Paquete.txt para tener un historial de cada paquete entrante
func (s *Server) escribirRegistroArchivo(registro Registro) error {

	// En un principio se creara dicho archivo en la misma altura que el codigo de logistica.go
	file, err := os.OpenFile("Registro_Paquete.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("no se pudo abrir o crear el archivo: %v", err)
	}
	defer file.Close()

	timestamp := registro.Timestamp.Format("02-01-2006 15:04")

	line := fmt.Sprintf("%s %06d %s %s %d %s %s %s\n",
		timestamp,
		registro.Id_paquete,
		registro.Tipo,
		registro.Nombre,
		registro.Valor,
		registro.Escolta,
		registro.Destino,
		registro.Seguimiento,
	)

	_, err = file.WriteString(line)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	return nil
}

// En esta funcion se escribe en Registro_Memoria.txt de todo lo que esta pasando en cada entrega o entregas fallidas de paquetes, para tener
// un historial de cada paquete
func (s *Server) escribirRegistroMemoria(paquete Paquete, estado string) error {
	file, err := os.OpenFile("Registro_Memoria.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("no se pudo abrir o crear el archivo: %v", err)
	}
	defer file.Close()

	// Formato de fecha estandar que vamos a utilizar
	timestamp := time.Now().Format("02-01-2006 15:04")

	line := fmt.Sprintf("%06d %s %s %s %d %s\n",
		paquete.Id_paquete,
		estado,       // Estado del paquete ("En Cetus", "Entregado", etc.)
		paquete.Tipo, // ID de la caravana
		paquete.Seguimiento,
		paquete.Intentos,
		timestamp,
	)

	_, err = file.WriteString(line)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	return nil
}

// Aqui se inicializan las caravanas con estado y cliente gRPC con el mismo puerto de conexion
func (s *Server) InitCaravanas() {

	s.caravanas = map[string]*Caravana{
		"Ostronitas":  {Estado: "Disponible", Cliente: s.connectCaravana("caravana_container:50052")},
		"Prioritario": {Estado: "Disponible", Cliente: s.connectCaravana("caravana_container:50052")},
		"Normal":      {Estado: "Disponible", Cliente: s.connectCaravana("caravana_container:50052")},
	}
}

// En esta funcion se establece la conexión gRPC a la caravana correspondiente para cada paquete
func (s *Server) connectCaravana(address string) pb.CaravanaServiceClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a la caravana en %s: %v", address, err)
	}
	return pb.NewCaravanaServiceClient(conn)
}

// Aqui se generan los numeros de seguimiento para cada paquete entrante
// Aquel numero de seguimiento se genera de manera aleatoria para entregarselo al cliente
func generarNumeroSeguimiento() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(1000000)+1)
}

// En esta funcion se actualiza el estado del paquete, por ejemplo en los casos que va en camino, entregado o no entregado
func (s *Server) ActualizarEstadoPaquete(ctx context.Context, req *pb.ActualizarEstadoRequest) (*pb.ActualizarEstadoResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.listaPaquetes {
		if s.listaPaquetes[i].Id_paquete == req.IdPaquete {
			s.listaPaquetes[i].Estado = req.Estado
			s.listaPaquetes[i].Intentos = req.Intentos
			fmt.Printf("Estado del paquete %d actualizado a: %s\n", req.IdPaquete, req.Estado)
			err := s.escribirRegistroMemoria(s.listaPaquetes[i], req.Estado)
			if err != nil {
				log.Printf("Error al escribir en el archivo de memoria: %v", err)
			}

			return &pb.ActualizarEstadoResponse{Mensaje: "Estado actualizado correctamente"}, nil
		}
	}

	return nil, fmt.Errorf("Paquete con ID %d no encontrado", req.IdPaquete)
}

// En esta funcion se envian las ordenes a las caravanas correspondientes para cada caravana
func (s *Server) EnviarOrden(ctx context.Context, orden *pb.Orden) (*pb.RespuestaOrden, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// ID incremental
	s.siguienteID++
	idPaquete := s.siguienteID

	seguimiento := generarNumeroSeguimiento()

	registro := Registro{
		Timestamp:   time.Now(),
		Id_paquete:  idPaquete,
		Tipo:        orden.Caravana,
		Nombre:      orden.Recurso,
		Valor:       int(orden.Valor),
		Escolta:     orden.Escolta,
		Destino:     orden.Destino,
		Seguimiento: seguimiento,
	}
	s.ordenes = append(s.ordenes, registro)

	paquete := Paquete{
		Id_paquete:  idPaquete,
		Seguimiento: seguimiento,
		Tipo:        orden.Caravana,
		Valor:       orden.Valor,
		Intentos:    0,
		Estado:      "En Cetus",
	}

	// Agregar el paquete a la lista global de todos los paquetes
	s.listaPaquetes = append(s.listaPaquetes, paquete)

	// Llamado de la funcion para escribir en el registro como archivo txt
	err := s.escribirRegistroArchivo(registro)
	if err != nil {
		log.Printf("Error al escribir en el archivo de registro: %v", err)
	}

	// Llamado de la funcion para escribir en el la memoria como archivo txt
	err = s.escribirRegistroMemoria(paquete, "En Cetus")
	if err != nil {
		log.Printf("Error al escribir en el archivo de memoria: %v", err)
	}

	// Aqui se asigna el paquete a la cola correspondiente de cada caravana
	switch orden.Caravana {
	case "Ostronitas":
		s.paquetesOstronitas = append(s.paquetesOstronitas, paquete)
	case "Prioritario":
		s.paquetesPrioritarios = append(s.paquetesPrioritarios, paquete)
	case "Normal":
		s.paquetesNormales = append(s.paquetesNormales, paquete)
	}

	enviarActualizacionFinanzas(paquete)

	fmt.Printf("Orden recibida y registrada: %+v\n", registro)

	// Verificar si la caravana correspondiente está disponible y si hay al menos dos paquetes para el envio de la caravana
	s.checkCaravanas()

	return &pb.RespuestaOrden{
		Seguimiento: seguimiento,
		Mensaje:     "Orden procesada con éxito",
	}, nil
}

// En esta funcion se verifica el estado de la caravana (Si es que esta disponible), ademas de ver el tema
// de los dos paquetes para procesarlos en las caravanas correspondientes.
// Otro tema que ve esta funcion es cuando un paquete se queda solo, esperando un lapso de 30 segundos para enviarlo
// sin esperar por un segundo paquete
func (s *Server) checkCaravanas() {
	for tipo, caravana := range s.caravanas {
		if caravana.Estado == "Disponible" {
			var paquetes []Paquete

			switch tipo {
			case "Ostronitas":
				paquetes = s.paquetesOstronitas
			case "Prioritario":
				paquetes = s.paquetesPrioritarios
			case "Normal":
				paquetes = s.paquetesNormales
			}

			if len(paquetes) >= 2 {
				// Marcar la caravana como ocupada
				caravana.Estado = "Ocupada"

				pbPaquetes := convertirPaquetesAPB(paquetes[:2])

				go s.enviarPaquetes(tipo, pbPaquetes)

				// Eliminar los paquetes enviados de la cola correspondiente por caravana
				switch tipo {
				case "Ostronitas":
					s.paquetesOstronitas = s.paquetesOstronitas[2:]
				case "Prioritario":
					s.paquetesPrioritarios = s.paquetesPrioritarios[2:]
				case "Normal":
					s.paquetesNormales = s.paquetesNormales[2:]
				}

				// Resetear el temporizador si ya existía y desbloquear el lock para la hebra entrante
				mu.Lock()
				if timers[tipo] != nil {
					timers[tipo].Stop()
				}
				delete(timers, tipo)
				mu.Unlock()

			} else if len(paquetes) == 1 {
				mu.Lock()
				if timers[tipo] == nil {
					timers[tipo] = time.AfterFunc(30*time.Second, func() {
						mu.Lock()
						defer mu.Unlock()
						var actualPaquetes []Paquete
						switch tipo {
						case "Ostronitas":
							actualPaquetes = s.paquetesOstronitas
						case "Prioritario":
							actualPaquetes = s.paquetesPrioritarios
						case "Normal":
							actualPaquetes = s.paquetesNormales
						}

						if len(actualPaquetes) == 1 {
							// Enviar el paquete solitario
							caravana.Estado = "Ocupada"
							pbPaquetes := convertirPaquetesAPB(actualPaquetes[:1])
							go s.enviarPaquetes(tipo, pbPaquetes)

							// Eliminar el paquete de la cola cuando se toma por una caravana
							switch tipo {
							case "Ostronitas":
								s.paquetesOstronitas = s.paquetesOstronitas[1:]
							case "Prioritario":
								s.paquetesPrioritarios = s.paquetesPrioritarios[1:]
							case "Normal":
								s.paquetesNormales = s.paquetesNormales[1:]
							}
							fmt.Printf("Caravana %s procesó un solo paquete debido a timeout\n", tipo)
						}

						// Remover el temporizador después de procesar
						delete(timers, tipo)
					})
				} else {
					timers[tipo].Reset(30 * time.Second)
				}
				mu.Unlock()
			}
		}
	}
}

// Convertir un slice de Paquete a un slice de *pb.Paquete para trabajar con el struct del gRPC
func convertirPaquetesAPB(paquetes []Paquete) []*pb.Paquete {
	var pbPaquetes []*pb.Paquete
	for _, p := range paquetes {
		pbPaquetes = append(pbPaquetes, &pb.Paquete{
			IdPaquete:   p.Id_paquete,
			Seguimiento: p.Seguimiento,
			Tipo:        p.Tipo,
			Valor:       p.Valor,
			Intentos:    p.Intentos,
			Estado:      p.Estado,
		})
	}
	return pbPaquetes
}

// Es una funcion simple para consultar el estado de un paquete, devolviendo los parametros solicitados
func (s *Server) ConsultarEstado(ctx context.Context, req *pb.ConsultaEstadoRequest) (*pb.EstadoPaquete, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Buscar el paquete en la lista global por ID
	for _, paquete := range s.listaPaquetes {
		if paquete.Seguimiento == req.Seguimiento {
			fmt.Printf("Consulta de estado para número de seguimiento: %s\n", req.Seguimiento)
			return &pb.EstadoPaquete{
				IdPaquete: paquete.Id_paquete,
				Estado:    paquete.Estado,
				Tipo:      paquete.Tipo,
				Valor:     paquete.Valor,
				Intentos:  paquete.Intentos,
			}, nil
		}
	}

	return nil, fmt.Errorf("Paquete con número de seguimiento:: %s no encontrado", req.Seguimiento)
}

// Envia las actualizacion de paquetes a finanzas para que este ultimo los pueda procesar y dar resultados
func enviarActualizacionFinanzas(paquete Paquete) {
	// Convertir el paquete a JSON
	paqueteJson, err := json.Marshal(paquete)
	if err != nil {
		log.Fatalf("Error al convertir el paquete a JSON para finanzas: %v", err)
	}

	// Envia el paquete a la cola de Rabbit para finanzas
	sendToQueue(string(paqueteJson))
	log.Printf("Actualizacion enviada a finanzas para el paquete ID: %d con estado: %s", paquete.Id_paquete, paquete.Estado)
}

// Enviar paquetes a la caravana para que los procese y comience la simulacion de la entrega de los paquetes
func (s *Server) enviarPaquetes(tipo string, paquetes []*pb.Paquete) {
	fmt.Printf("Caravana %s procesando paquetes: %+v\n", tipo, paquetes)
	s.mu.Lock()
	for _, p := range paquetes {
		for i := range s.listaPaquetes {
			if s.listaPaquetes[i].Id_paquete == p.IdPaquete {
				s.listaPaquetes[i].Estado = "En camino"
			}
		}
	}
	s.mu.Unlock()

	registros := s.buscarRegistros(paquetes)

	caravana := s.caravanas[tipo]
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Enviar los paquetes a la caravana junto con los registros (escolta y destino)
	_, err := caravana.Cliente.ProcesarPaquetes(ctx, &pb.ProcesarPaquetesRequest{Paquetes: paquetes, Registros: registros})
	if err != nil {
		fmt.Printf("Error al enviar paquetes a la caravana %s: %v\n", tipo, err)

		// Si hubo un error, regresar los paquetes a la cola
		s.mu.Lock()
		paquetes := convertirPaquetesDePB(paquetes)
		switch tipo {
		case "Prioritario":
			s.paquetesPrioritarios = append(paquetes, s.paquetesPrioritarios...)
		case "Ostronitas":
			s.paquetesOstronitas = append(paquetes, s.paquetesOstronitas...)
		case "Normal":
			s.paquetesNormales = append(paquetes, s.paquetesNormales...)
		}
		s.mu.Unlock()
	} else {
		fmt.Printf("Caravana %s completó la entrega de los paquetes.\n", tipo)
	}

	// Marcar la caravana como disponible de nuevo para volver a recibir paquetes
	s.mu.Lock()

	s.caravanas[tipo].Estado = "Disponible"

	for i := range s.listaPaquetes {
		if s.listaPaquetes[i].Estado == "No entregado" || s.listaPaquetes[i].Estado == "Entregado" {
			enviarActualizacionFinanzas(s.listaPaquetes[i])
		}
	}

	s.mu.Unlock()

	// Volver a verificar las caravanas después de completar la entrega
	s.checkCaravanas()
}

// Aqui se busca entre los paquetes si corresponde al mismo ID del paquete entrante para saber cuales quitar y enviar
func (s *Server) buscarRegistros(paquetes []*pb.Paquete) []*pb.Registro {
	var registros []*pb.Registro
	for _, paquete := range paquetes {
		for _, registro := range s.ordenes {
			if registro.Id_paquete == paquete.IdPaquete {
				registros = append(registros, &pb.Registro{
					Escolta:   registro.Escolta,
					Destino:   registro.Destino,
					IdPaquete: registro.Id_paquete,
				})
			}
		}
	}
	return registros
}

// Convertir un slice de *pb.Paquete a un slice de Paquete para enviarlocomo struct en modo gRPC
func convertirPaquetesDePB(pbPaquetes []*pb.Paquete) []Paquete {
	var paquetes []Paquete
	for _, p := range pbPaquetes {
		paquetes = append(paquetes, Paquete{
			Id_paquete:  p.IdPaquete,
			Seguimiento: p.Seguimiento,
			Tipo:        p.Tipo,
			Valor:       p.Valor,
			Intentos:    p.Intentos,
			Estado:      p.Estado,
		})
	}
	return paquetes
}

// Funcion principal para escuchar en el puerto 50051 esperando ordenes de los clientes y procesar
// aquellos paquetes para enviarselos a las caravanas
func main() {
	server := &Server{siguienteID: 0}
	server.InitCaravanas()

	lis, err := net.Listen("tcp", "10.35.168.48:50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLogisticaServiceServer(grpcServer, server)

	log.Printf("Servidor logistica escuchando en :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al iniciar el servidor gRPC: %v", err)
	}
}
