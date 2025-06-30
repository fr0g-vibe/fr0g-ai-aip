package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc/pb"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// PersonaServer implements the gRPC PersonaService
type PersonaServer struct {
	pb.UnimplementedPersonaServiceServer
	service *persona.Service
	config  *config.Config
}

// StartGRPCServer starts a real gRPC server using protobuf
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	
	// Create a default service for the standalone server
	memStorage := storage.NewMemoryStorage()
	service := persona.NewService(memStorage)
	personaServer := &PersonaServer{
		service: service,
	}
	
	pb.RegisterPersonaServiceServer(s, personaServer)

	fmt.Printf("gRPC server listening on port %s\n", port)
	fmt.Println("Using real gRPC with protobuf")

	return s.Serve(lis)
}

// NewPersonaServer creates a new gRPC persona server
func NewPersonaServer(cfg *config.Config, service *persona.Service) *PersonaServer {
	return &PersonaServer{
		service: service,
		config:  cfg,
	}
}

// StartGRPCServerWithConfig starts a gRPC server with full configuration
func StartGRPCServerWithConfig(cfg *config.Config, service *persona.Service) error {
	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Configure gRPC server options
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(cfg.GRPC.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.GRPC.MaxSendMsgSize),
	}

	s := grpc.NewServer(opts...)

	// Register the persona service
	personaServer := NewPersonaServer(cfg, service)
	pb.RegisterPersonaServiceServer(s, personaServer)

	fmt.Printf("gRPC server listening on port %s\n", cfg.GRPC.Port)
	fmt.Println("Using real gRPC with protobuf")

	return s.Serve(lis)
}

// CreatePersona creates a new persona
func (s *PersonaServer) CreatePersona(ctx context.Context, req *pb.CreatePersonaRequest) (*pb.CreatePersonaResponse, error) {
	if req.Persona == nil {
		return nil, status.Errorf(codes.InvalidArgument, "persona is required")
	}

	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	// Convert proto to internal type
	p := types.ProtoToPersona(req.Persona)
	
	err := s.service.CreatePersona(p)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create persona: %v", err)
	}

	return &pb.CreatePersonaResponse{
		Persona: types.PersonaToProto(p),
	}, nil
}

// GetPersona retrieves a persona by ID
func (s *PersonaServer) GetPersona(ctx context.Context, req *pb.GetPersonaRequest) (*pb.GetPersonaResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "persona ID is required")
	}

	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	p, err := s.service.GetPersona(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "persona not found: %v", err)
	}

	return &pb.GetPersonaResponse{
		Persona: types.PersonaToProto(&p),
	}, nil
}

// ListPersonas returns all personas
func (s *PersonaServer) ListPersonas(ctx context.Context, req *pb.ListPersonasRequest) (*pb.ListPersonasResponse, error) {
	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	personas, err := s.service.ListPersonas()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list personas: %v", err)
	}

	var protoPersonas []*pb.Persona
	for _, p := range personas {
		protoPersonas = append(protoPersonas, types.PersonaToProto(&p))
	}

	return &pb.ListPersonasResponse{
		Personas: protoPersonas,
	}, nil
}

// UpdatePersona updates an existing persona
func (s *PersonaServer) UpdatePersona(ctx context.Context, req *pb.UpdatePersonaRequest) (*pb.UpdatePersonaResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "persona ID is required")
	}

	if req.Persona == nil {
		return nil, status.Errorf(codes.InvalidArgument, "persona is required")
	}

	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	// Convert proto to internal type
	p := types.ProtoToPersona(req.Persona)
	p.Id = req.Id

	err := s.service.UpdatePersona(req.Id, *p)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update persona: %v", err)
	}

	return &pb.UpdatePersonaResponse{
		Persona: types.PersonaToProto(p),
	}, nil
}

// DeletePersona removes a persona by ID
func (s *PersonaServer) DeletePersona(ctx context.Context, req *pb.DeletePersonaRequest) (*pb.DeletePersonaResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "persona ID is required")
	}

	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	err := s.service.DeletePersona(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to delete persona: %v", err)
	}

	return &pb.DeletePersonaResponse{}, nil
}

// CreateIdentity creates a new identity
func (s *PersonaServer) CreateIdentity(ctx context.Context, req *pb.CreateIdentityRequest) (*pb.CreateIdentityResponse, error) {
	if req.Identity == nil {
		return nil, status.Errorf(codes.InvalidArgument, "identity is required")
	}

	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	// Convert proto to internal type
	identity := types.ProtoToIdentity(req.Identity)
	
	err := s.service.CreateIdentity(identity)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create identity: %v", err)
	}

	return &pb.CreateIdentityResponse{
		Identity: types.IdentityToProto(identity),
	}, nil
}

// GetIdentity retrieves an identity by ID
func (s *PersonaServer) GetIdentity(ctx context.Context, req *pb.GetIdentityRequest) (*pb.GetIdentityResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "identity ID is required")
	}

	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	identity, err := s.service.GetIdentity(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "identity not found: %v", err)
	}

	return &pb.GetIdentityResponse{
		Identity: types.IdentityToProto(&identity),
	}, nil
}

// ListIdentities returns identities with optional filtering
func (s *PersonaServer) ListIdentities(ctx context.Context, req *pb.ListIdentitiesRequest) (*pb.ListIdentitiesResponse, error) {
	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	var filter *types.IdentityFilter
	if req.Filter != nil {
		filter = &types.IdentityFilter{
			PersonaID: req.Filter.PersonaId,
			Tags:      req.Filter.Tags,
			Search:    req.Filter.Search,
		}
		isActive := req.Filter.IsActive
		filter.IsActive = &isActive
	}

	identities, err := s.service.ListIdentities(filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list identities: %v", err)
	}

	var pbIdentities []*pb.Identity
	for _, identity := range identities {
		pbIdentities = append(pbIdentities, types.IdentityToProto(&identity))
	}

	return &pb.ListIdentitiesResponse{
		Identities: pbIdentities,
	}, nil
}

// UpdateIdentity updates an existing identity
func (s *PersonaServer) UpdateIdentity(ctx context.Context, req *pb.UpdateIdentityRequest) (*pb.UpdateIdentityResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "identity ID is required")
	}
	if req.Identity == nil {
		return nil, status.Errorf(codes.InvalidArgument, "identity is required")
	}

	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	// Convert proto to internal type
	identity := types.ProtoToIdentity(req.Identity)
	identity.Id = req.Id

	err := s.service.UpdateIdentity(req.Id, *identity)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update identity: %v", err)
	}

	return &pb.UpdateIdentityResponse{
		Identity: types.IdentityToProto(identity),
	}, nil
}

// DeleteIdentity removes an identity by ID
func (s *PersonaServer) DeleteIdentity(ctx context.Context, req *pb.DeleteIdentityRequest) (*pb.DeleteIdentityResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "identity ID is required")
	}

	if s.service == nil {
		return nil, status.Errorf(codes.Internal, "persona service not available")
	}

	err := s.service.DeleteIdentity(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to delete identity: %v", err)
	}

	return &pb.DeleteIdentityResponse{}, nil
}

