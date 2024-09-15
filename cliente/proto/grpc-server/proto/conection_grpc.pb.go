// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: conection.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	LogisticaService_EnviarOrden_FullMethodName             = "/conection.LogisticaService/EnviarOrden"
	LogisticaService_ConsultarEstado_FullMethodName         = "/conection.LogisticaService/ConsultarEstado"
	LogisticaService_ObtenerPaquetes_FullMethodName         = "/conection.LogisticaService/ObtenerPaquetes"
	LogisticaService_ActualizarEstadoPaquete_FullMethodName = "/conection.LogisticaService/ActualizarEstadoPaquete"
)

// LogisticaServiceClient is the client API for LogisticaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogisticaServiceClient interface {
	EnviarOrden(ctx context.Context, in *Orden, opts ...grpc.CallOption) (*RespuestaOrden, error)
	ConsultarEstado(ctx context.Context, in *ConsultaEstadoRequest, opts ...grpc.CallOption) (*EstadoPaquete, error)
	ObtenerPaquetes(ctx context.Context, in *ObtenerPaquetesRequest, opts ...grpc.CallOption) (*ObtenerPaquetesResponse, error)
	ActualizarEstadoPaquete(ctx context.Context, in *ActualizarEstadoRequest, opts ...grpc.CallOption) (*ActualizarEstadoResponse, error)
}

type logisticaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLogisticaServiceClient(cc grpc.ClientConnInterface) LogisticaServiceClient {
	return &logisticaServiceClient{cc}
}

func (c *logisticaServiceClient) EnviarOrden(ctx context.Context, in *Orden, opts ...grpc.CallOption) (*RespuestaOrden, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RespuestaOrden)
	err := c.cc.Invoke(ctx, LogisticaService_EnviarOrden_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticaServiceClient) ConsultarEstado(ctx context.Context, in *ConsultaEstadoRequest, opts ...grpc.CallOption) (*EstadoPaquete, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EstadoPaquete)
	err := c.cc.Invoke(ctx, LogisticaService_ConsultarEstado_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticaServiceClient) ObtenerPaquetes(ctx context.Context, in *ObtenerPaquetesRequest, opts ...grpc.CallOption) (*ObtenerPaquetesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ObtenerPaquetesResponse)
	err := c.cc.Invoke(ctx, LogisticaService_ObtenerPaquetes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticaServiceClient) ActualizarEstadoPaquete(ctx context.Context, in *ActualizarEstadoRequest, opts ...grpc.CallOption) (*ActualizarEstadoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ActualizarEstadoResponse)
	err := c.cc.Invoke(ctx, LogisticaService_ActualizarEstadoPaquete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogisticaServiceServer is the server API for LogisticaService service.
// All implementations must embed UnimplementedLogisticaServiceServer
// for forward compatibility.
type LogisticaServiceServer interface {
	EnviarOrden(context.Context, *Orden) (*RespuestaOrden, error)
	ConsultarEstado(context.Context, *ConsultaEstadoRequest) (*EstadoPaquete, error)
	ObtenerPaquetes(context.Context, *ObtenerPaquetesRequest) (*ObtenerPaquetesResponse, error)
	ActualizarEstadoPaquete(context.Context, *ActualizarEstadoRequest) (*ActualizarEstadoResponse, error)
	mustEmbedUnimplementedLogisticaServiceServer()
}

// UnimplementedLogisticaServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLogisticaServiceServer struct{}

func (UnimplementedLogisticaServiceServer) EnviarOrden(context.Context, *Orden) (*RespuestaOrden, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnviarOrden not implemented")
}
func (UnimplementedLogisticaServiceServer) ConsultarEstado(context.Context, *ConsultaEstadoRequest) (*EstadoPaquete, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConsultarEstado not implemented")
}
func (UnimplementedLogisticaServiceServer) ObtenerPaquetes(context.Context, *ObtenerPaquetesRequest) (*ObtenerPaquetesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ObtenerPaquetes not implemented")
}
func (UnimplementedLogisticaServiceServer) ActualizarEstadoPaquete(context.Context, *ActualizarEstadoRequest) (*ActualizarEstadoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActualizarEstadoPaquete not implemented")
}
func (UnimplementedLogisticaServiceServer) mustEmbedUnimplementedLogisticaServiceServer() {}
func (UnimplementedLogisticaServiceServer) testEmbeddedByValue()                          {}

// UnsafeLogisticaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogisticaServiceServer will
// result in compilation errors.
type UnsafeLogisticaServiceServer interface {
	mustEmbedUnimplementedLogisticaServiceServer()
}

func RegisterLogisticaServiceServer(s grpc.ServiceRegistrar, srv LogisticaServiceServer) {
	// If the following call pancis, it indicates UnimplementedLogisticaServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LogisticaService_ServiceDesc, srv)
}

func _LogisticaService_EnviarOrden_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Orden)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaServiceServer).EnviarOrden(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticaService_EnviarOrden_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaServiceServer).EnviarOrden(ctx, req.(*Orden))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticaService_ConsultarEstado_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConsultaEstadoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaServiceServer).ConsultarEstado(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticaService_ConsultarEstado_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaServiceServer).ConsultarEstado(ctx, req.(*ConsultaEstadoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticaService_ObtenerPaquetes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ObtenerPaquetesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaServiceServer).ObtenerPaquetes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticaService_ObtenerPaquetes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaServiceServer).ObtenerPaquetes(ctx, req.(*ObtenerPaquetesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticaService_ActualizarEstadoPaquete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActualizarEstadoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaServiceServer).ActualizarEstadoPaquete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticaService_ActualizarEstadoPaquete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaServiceServer).ActualizarEstadoPaquete(ctx, req.(*ActualizarEstadoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LogisticaService_ServiceDesc is the grpc.ServiceDesc for LogisticaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogisticaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "conection.LogisticaService",
	HandlerType: (*LogisticaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EnviarOrden",
			Handler:    _LogisticaService_EnviarOrden_Handler,
		},
		{
			MethodName: "ConsultarEstado",
			Handler:    _LogisticaService_ConsultarEstado_Handler,
		},
		{
			MethodName: "ObtenerPaquetes",
			Handler:    _LogisticaService_ObtenerPaquetes_Handler,
		},
		{
			MethodName: "ActualizarEstadoPaquete",
			Handler:    _LogisticaService_ActualizarEstadoPaquete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "conection.proto",
}

const (
	CaravanaService_ProcesarPaquetes_FullMethodName = "/conection.CaravanaService/ProcesarPaquetes"
)

// CaravanaServiceClient is the client API for CaravanaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CaravanaServiceClient interface {
	ProcesarPaquetes(ctx context.Context, in *ProcesarPaquetesRequest, opts ...grpc.CallOption) (*ProcesarPaquetesResponse, error)
}

type caravanaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCaravanaServiceClient(cc grpc.ClientConnInterface) CaravanaServiceClient {
	return &caravanaServiceClient{cc}
}

func (c *caravanaServiceClient) ProcesarPaquetes(ctx context.Context, in *ProcesarPaquetesRequest, opts ...grpc.CallOption) (*ProcesarPaquetesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProcesarPaquetesResponse)
	err := c.cc.Invoke(ctx, CaravanaService_ProcesarPaquetes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CaravanaServiceServer is the server API for CaravanaService service.
// All implementations must embed UnimplementedCaravanaServiceServer
// for forward compatibility.
type CaravanaServiceServer interface {
	ProcesarPaquetes(context.Context, *ProcesarPaquetesRequest) (*ProcesarPaquetesResponse, error)
	mustEmbedUnimplementedCaravanaServiceServer()
}

// UnimplementedCaravanaServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCaravanaServiceServer struct{}

func (UnimplementedCaravanaServiceServer) ProcesarPaquetes(context.Context, *ProcesarPaquetesRequest) (*ProcesarPaquetesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcesarPaquetes not implemented")
}
func (UnimplementedCaravanaServiceServer) mustEmbedUnimplementedCaravanaServiceServer() {}
func (UnimplementedCaravanaServiceServer) testEmbeddedByValue()                         {}

// UnsafeCaravanaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CaravanaServiceServer will
// result in compilation errors.
type UnsafeCaravanaServiceServer interface {
	mustEmbedUnimplementedCaravanaServiceServer()
}

func RegisterCaravanaServiceServer(s grpc.ServiceRegistrar, srv CaravanaServiceServer) {
	// If the following call pancis, it indicates UnimplementedCaravanaServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CaravanaService_ServiceDesc, srv)
}

func _CaravanaService_ProcesarPaquetes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcesarPaquetesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaravanaServiceServer).ProcesarPaquetes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaravanaService_ProcesarPaquetes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaravanaServiceServer).ProcesarPaquetes(ctx, req.(*ProcesarPaquetesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CaravanaService_ServiceDesc is the grpc.ServiceDesc for CaravanaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CaravanaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "conection.CaravanaService",
	HandlerType: (*CaravanaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcesarPaquetes",
			Handler:    _CaravanaService_ProcesarPaquetes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "conection.proto",
}
