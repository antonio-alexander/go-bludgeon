// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: employees.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EmployeesClient is the client API for Employees service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmployeesClient interface {
	// employee_create
	EmployeeCreate(ctx context.Context, in *EmployeeCreateRequest, opts ...grpc.CallOption) (*EmployeeCreateResponse, error)
	// employee_read
	EmployeeRead(ctx context.Context, in *EmployeeReadRequest, opts ...grpc.CallOption) (*EmployeeReadResponse, error)
	// employees_read
	EmployeesRead(ctx context.Context, in *EmployeesReadRequest, opts ...grpc.CallOption) (*EmployeesReadResponse, error)
	// employee_update
	EmployeeUpdate(ctx context.Context, in *EmployeeUpdateRequest, opts ...grpc.CallOption) (*EmployeeUpdateResponse, error)
	// employee_delete
	EmployeeDelete(ctx context.Context, in *EmployeeDeleteRequest, opts ...grpc.CallOption) (*EmployeeDeleteResponse, error)
}

type employeesClient struct {
	cc grpc.ClientConnInterface
}

func NewEmployeesClient(cc grpc.ClientConnInterface) EmployeesClient {
	return &employeesClient{cc}
}

func (c *employeesClient) EmployeeCreate(ctx context.Context, in *EmployeeCreateRequest, opts ...grpc.CallOption) (*EmployeeCreateResponse, error) {
	out := new(EmployeeCreateResponse)
	err := c.cc.Invoke(ctx, "/go_bludgeon_employees.Employees/employee_create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *employeesClient) EmployeeRead(ctx context.Context, in *EmployeeReadRequest, opts ...grpc.CallOption) (*EmployeeReadResponse, error) {
	out := new(EmployeeReadResponse)
	err := c.cc.Invoke(ctx, "/go_bludgeon_employees.Employees/employee_read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *employeesClient) EmployeesRead(ctx context.Context, in *EmployeesReadRequest, opts ...grpc.CallOption) (*EmployeesReadResponse, error) {
	out := new(EmployeesReadResponse)
	err := c.cc.Invoke(ctx, "/go_bludgeon_employees.Employees/employees_read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *employeesClient) EmployeeUpdate(ctx context.Context, in *EmployeeUpdateRequest, opts ...grpc.CallOption) (*EmployeeUpdateResponse, error) {
	out := new(EmployeeUpdateResponse)
	err := c.cc.Invoke(ctx, "/go_bludgeon_employees.Employees/employee_update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *employeesClient) EmployeeDelete(ctx context.Context, in *EmployeeDeleteRequest, opts ...grpc.CallOption) (*EmployeeDeleteResponse, error) {
	out := new(EmployeeDeleteResponse)
	err := c.cc.Invoke(ctx, "/go_bludgeon_employees.Employees/employee_delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmployeesServer is the server API for Employees service.
// All implementations must embed UnimplementedEmployeesServer
// for forward compatibility
type EmployeesServer interface {
	// employee_create
	EmployeeCreate(context.Context, *EmployeeCreateRequest) (*EmployeeCreateResponse, error)
	// employee_read
	EmployeeRead(context.Context, *EmployeeReadRequest) (*EmployeeReadResponse, error)
	// employees_read
	EmployeesRead(context.Context, *EmployeesReadRequest) (*EmployeesReadResponse, error)
	// employee_update
	EmployeeUpdate(context.Context, *EmployeeUpdateRequest) (*EmployeeUpdateResponse, error)
	// employee_delete
	EmployeeDelete(context.Context, *EmployeeDeleteRequest) (*EmployeeDeleteResponse, error)
	mustEmbedUnimplementedEmployeesServer()
}

// UnimplementedEmployeesServer must be embedded to have forward compatible implementations.
type UnimplementedEmployeesServer struct {
}

func (UnimplementedEmployeesServer) EmployeeCreate(context.Context, *EmployeeCreateRequest) (*EmployeeCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmployeeCreate not implemented")
}
func (UnimplementedEmployeesServer) EmployeeRead(context.Context, *EmployeeReadRequest) (*EmployeeReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmployeeRead not implemented")
}
func (UnimplementedEmployeesServer) EmployeesRead(context.Context, *EmployeesReadRequest) (*EmployeesReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmployeesRead not implemented")
}
func (UnimplementedEmployeesServer) EmployeeUpdate(context.Context, *EmployeeUpdateRequest) (*EmployeeUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmployeeUpdate not implemented")
}
func (UnimplementedEmployeesServer) EmployeeDelete(context.Context, *EmployeeDeleteRequest) (*EmployeeDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmployeeDelete not implemented")
}
func (UnimplementedEmployeesServer) mustEmbedUnimplementedEmployeesServer() {}

// UnsafeEmployeesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmployeesServer will
// result in compilation errors.
type UnsafeEmployeesServer interface {
	mustEmbedUnimplementedEmployeesServer()
}

func RegisterEmployeesServer(s grpc.ServiceRegistrar, srv EmployeesServer) {
	s.RegisterService(&Employees_ServiceDesc, srv)
}

func _Employees_EmployeeCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmployeeCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmployeesServer).EmployeeCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_bludgeon_employees.Employees/employee_create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmployeesServer).EmployeeCreate(ctx, req.(*EmployeeCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Employees_EmployeeRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmployeeReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmployeesServer).EmployeeRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_bludgeon_employees.Employees/employee_read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmployeesServer).EmployeeRead(ctx, req.(*EmployeeReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Employees_EmployeesRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmployeesReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmployeesServer).EmployeesRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_bludgeon_employees.Employees/employees_read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmployeesServer).EmployeesRead(ctx, req.(*EmployeesReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Employees_EmployeeUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmployeeUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmployeesServer).EmployeeUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_bludgeon_employees.Employees/employee_update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmployeesServer).EmployeeUpdate(ctx, req.(*EmployeeUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Employees_EmployeeDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmployeeDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmployeesServer).EmployeeDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_bludgeon_employees.Employees/employee_delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmployeesServer).EmployeeDelete(ctx, req.(*EmployeeDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Employees_ServiceDesc is the grpc.ServiceDesc for Employees service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Employees_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go_bludgeon_employees.Employees",
	HandlerType: (*EmployeesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "employee_create",
			Handler:    _Employees_EmployeeCreate_Handler,
		},
		{
			MethodName: "employee_read",
			Handler:    _Employees_EmployeeRead_Handler,
		},
		{
			MethodName: "employees_read",
			Handler:    _Employees_EmployeesRead_Handler,
		},
		{
			MethodName: "employee_update",
			Handler:    _Employees_EmployeeUpdate_Handler,
		},
		{
			MethodName: "employee_delete",
			Handler:    _Employees_EmployeeDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "employees.proto",
}
