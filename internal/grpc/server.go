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
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	pb.RegisterPersonaServiceServer(s, &PersonaServer{})

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

	p := &types.Persona{
		Name:    req.Persona.Name,
		Topic:   req.Persona.Topic,
		Prompt:  req.Persona.Prompt,
		Context: req.Persona.Context,
		RAG:     req.Persona.Rag,
	}

	var err error
	if s.service != nil {
		err = s.service.CreatePersona(p)
	} else {
		// Fallback to legacy global service
		err = persona.CreatePersona(p)
	}

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.CreatePersonaResponse{
		Persona: &pb.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// GetPersona retrieves a persona by ID
func (s *PersonaServer) GetPersona(ctx context.Context, req *pb.GetPersonaRequest) (*pb.GetPersonaResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "persona ID is required")
	}

	var p types.Persona
	var err error

	if s.service != nil {
		p, err = s.service.GetPersona(req.Id)
	} else {
		// Fallback to legacy global service
		p, err = persona.GetPersona(req.Id)
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &pb.GetPersonaResponse{
		Persona: &pb.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// ListPersonas returns all personas
func (s *PersonaServer) ListPersonas(ctx context.Context, req *pb.ListPersonasRequest) (*pb.ListPersonasResponse, error) {
	var personas []types.Persona
	var err error

	if s.service != nil {
		personas, err = s.service.ListPersonas()
	} else {
		// Fallback to legacy global service
		personas = persona.ListPersonas()
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list personas: %v", err)
	}

	var protoPersonas []*pb.Persona
	for _, p := range personas {
		protoPersonas = append(protoPersonas, &pb.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		})
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

	p := types.Persona{
		ID:      req.Id,
		Name:    req.Persona.Name,
		Topic:   req.Persona.Topic,
		Prompt:  req.Persona.Prompt,
		Context: req.Persona.Context,
		RAG:     req.Persona.Rag,
	}

	var err error
	if s.service != nil {
		err = s.service.UpdatePersona(req.Id, p)
	} else {
		// Fallback to legacy global service
		err = persona.UpdatePersona(req.Id, p)
	}

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.UpdatePersonaResponse{
		Persona: &pb.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// DeletePersona removes a persona by ID
func (s *PersonaServer) DeletePersona(ctx context.Context, req *pb.DeletePersonaRequest) (*pb.DeletePersonaResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "persona ID is required")
	}

	var err error
	if s.service != nil {
		err = s.service.DeletePersona(req.Id)
	} else {
		// Fallback to legacy global service
		err = persona.DeletePersona(req.Id)
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &pb.DeletePersonaResponse{}, nil
}

// CreateIdentity creates a new identity
func (s *PersonaServer) CreateIdentity(ctx context.Context, req *pb.CreateIdentityRequest) (*pb.CreateIdentityResponse, error) {
	if req.Identity == nil {
		return nil, status.Errorf(codes.InvalidArgument, "identity is required")
	}

	// Convert pb.Identity to types.Identity
	identity, err := pbToTypesIdentity(req.Identity)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if s.service != nil {
		err = s.service.CreateIdentity(identity)
	} else {
		return nil, status.Errorf(codes.Unimplemented, "identity service not available")
	}

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateIdentityResponse{
		Identity: typesToPbIdentity(identity),
	}, nil
}

func pbToTypesIdentity(pbId *pb.Identity) (*types.Identity, error) {
	if pbId == nil {
		return nil, nil
	}
	id := &types.Identity{
		ID:          pbId.Id,
		PersonaID:   pbId.PersonaId,
		Name:        pbId.Name,
		Description: pbId.Description,
		Background:  pbId.Background,
		IsActive:    pbId.IsActive,
		Tags:        pbId.Tags,
	}
	if pbId.CreatedAt != nil {
		id.CreatedAt = pbId.CreatedAt.AsTime()
	}
	if pbId.UpdatedAt != nil {
		id.UpdatedAt = pbId.UpdatedAt.AsTime()
	}
	if pbId.GetRichAttributes() != nil {
		id.RichAttributes = pbToTypesRichAttributes(pbId.GetRichAttributes())
	}
	return id, nil
}

func typesToPbIdentity(id *types.Identity) *pb.Identity {
	if id == nil {
		return nil
	}
	pbId := &pb.Identity{
		Id:          id.ID,
		PersonaId:   id.PersonaID,
		Name:        id.Name,
		Description: id.Description,
		Background:  id.Background,
		IsActive:    id.IsActive,
		Tags:        id.Tags,
	}
	pbId.CreatedAt = timestamppb.New(id.CreatedAt)
	pbId.UpdatedAt = timestamppb.New(id.UpdatedAt)
	if id.RichAttributes != nil {
		pbId.RichAttributes = typesToPbRichAttributes(id.RichAttributes)
	}
	return pbId
}

// RichAttributes conversion helpers
func pbToTypesRichAttributes(pbAttr *pb.RichAttributes) *types.RichAttributes {
	if pbAttr == nil {
		return nil
	}
	return &types.RichAttributes{
		Demographics:         pbToTypesDemographics(pbAttr.GetDemographics()),
		Psychographics:       pbToTypesPsychographics(pbAttr.GetPsychographics()),
		LifeHistory:          pbToTypesLifeHistory(pbAttr.GetLifeHistory()),
		CulturalReligious:    pbToTypesCulturalReligious(pbAttr.GetCulturalReligious()),
		PoliticalSocial:      pbToTypesPoliticalSocial(pbAttr.GetPoliticalSocial()),
		Health:               pbToTypesHealth(pbAttr.GetHealth()),
		Preferences:          pbToTypesPreferences(pbAttr.GetPreferences()),
		BehavioralTendencies: pbToTypesBehavioralTendencies(pbAttr.GetBehavioralTendencies()),
		CurrentContext:       pbToTypesCurrentContext(pbAttr.GetCurrentContext()),
		Custom:               pbAttr.GetCustom(),
	}
}

func typesToPbRichAttributes(attr *types.RichAttributes) *pb.RichAttributes {
	if attr == nil {
		return nil
	}
	return &pb.RichAttributes{
		Demographics:         typesToPbDemographics(attr.Demographics),
		Psychographics:       typesToPbPsychographics(attr.Psychographics),
		LifeHistory:          typesToPbLifeHistory(attr.LifeHistory),
		CulturalReligious:    typesToPbCulturalReligious(attr.CulturalReligious),
		PoliticalSocial:      typesToPbPoliticalSocial(attr.PoliticalSocial),
		Health:               typesToPbHealth(attr.Health),
		Preferences:          typesToPbPreferences(attr.Preferences),
		BehavioralTendencies: typesToPbBehavioralTendencies(attr.BehavioralTendencies),
		CurrentContext:       typesToPbCurrentContext(attr.CurrentContext),
		Custom:               attr.Custom,
	}
}

// Demographics conversion
func pbToTypesDemographics(pbDemo *pb.Demographics) *types.Demographics {
	if pbDemo == nil {
		return nil
	}
	return &types.Demographics{
		Age:                 int(pbDemo.GetAge()),
		Gender:              pbDemo.GetGender(),
		Ethnicity:           pbDemo.GetEthnicity(),
		Nationality:         pbDemo.GetNationality(),
		Education:           pbDemo.GetEducation(),
		Occupation:          pbDemo.GetOccupation(),
		SocioeconomicStatus: pbDemo.GetSocioeconomicStatus(),
		Location:            pbToTypesLocation(pbDemo.GetLocation()),
		Languages:           pbDemo.GetLanguages(),
		MaritalStatus:       pbDemo.GetMaritalStatus(),
		Children:            int(pbDemo.GetChildren()),
	}
}

func typesToPbDemographics(demo *types.Demographics) *pb.Demographics {
	if demo == nil {
		return nil
	}
	return &pb.Demographics{
		Age:                 int32(demo.Age),
		Gender:              demo.Gender,
		Ethnicity:           demo.Ethnicity,
		Nationality:         demo.Nationality,
		Education:           demo.Education,
		Occupation:          demo.Occupation,
		SocioeconomicStatus: demo.SocioeconomicStatus,
		Location:            typesToPbLocation(demo.Location),
		Languages:           demo.Languages,
		MaritalStatus:       demo.MaritalStatus,
		Children:            int32(demo.Children),
	}
}

// Location conversion
func pbToTypesLocation(pbLoc *pb.Location) *types.Location {
	if pbLoc == nil {
		return nil
	}
	return &types.Location{
		Country:    pbLoc.GetCountry(),
		Region:     pbLoc.GetRegion(),
		City:       pbLoc.GetCity(),
		UrbanRural: pbLoc.GetUrbanRural(),
		Timezone:   pbLoc.GetTimezone(),
	}
}

func typesToPbLocation(loc *types.Location) *pb.Location {
	if loc == nil {
		return nil
	}
	return &pb.Location{
		Country:    loc.Country,
		Region:     loc.Region,
		City:       loc.City,
		UrbanRural: loc.UrbanRural,
		Timezone:   loc.Timezone,
	}
}

// Psychographics conversion
func pbToTypesPsychographics(pbPsycho *pb.Psychographics) *types.Psychographics {
	if pbPsycho == nil {
		return nil
	}
	return &types.Psychographics{
		Personality:      pbToTypesPersonality(pbPsycho.GetPersonality()),
		Values:           pbPsycho.GetValues(),
		CoreBeliefs:      pbPsycho.GetCoreBeliefs(),
		CognitiveStyle:   pbPsycho.GetCognitiveStyle(),
		LearningStyle:    pbPsycho.GetLearningStyle(),
		RiskTolerance:    pbPsycho.GetRiskTolerance(),
		OpennessToChange: pbPsycho.GetOpennessToChange(),
	}
}

func typesToPbPsychographics(psycho *types.Psychographics) *pb.Psychographics {
	if psycho == nil {
		return nil
	}
	return &pb.Psychographics{
		Personality:      typesToPbPersonality(psycho.Personality),
		Values:           psycho.Values,
		CoreBeliefs:      psycho.CoreBeliefs,
		CognitiveStyle:   psycho.CognitiveStyle,
		LearningStyle:    psycho.LearningStyle,
		RiskTolerance:    psycho.RiskTolerance,
		OpennessToChange: psycho.OpennessToChange,
	}
}

// Personality conversion
func pbToTypesPersonality(pbPers *pb.Personality) *types.Personality {
	if pbPers == nil {
		return nil
	}
	return &types.Personality{
		Openness:          pbPers.GetOpenness(),
		Conscientiousness: pbPers.GetConscientiousness(),
		Extraversion:      pbPers.GetExtraversion(),
		Agreeableness:     pbPers.GetAgreeableness(),
		Neuroticism:       pbPers.GetNeuroticism(),
	}
}

func typesToPbPersonality(pers *types.Personality) *pb.Personality {
	if pers == nil {
		return nil
	}
	return &pb.Personality{
		Openness:          pers.Openness,
		Conscientiousness: pers.Conscientiousness,
		Extraversion:      pers.Extraversion,
		Agreeableness:     pers.Agreeableness,
		Neuroticism:       pers.Neuroticism,
	}
}

// LifeHistory conversion
func pbToTypesLifeHistory(pbLife *pb.LifeHistory) *types.LifeHistory {
	if pbLife == nil {
		return nil
	}
	return &types.LifeHistory{
		ChildhoodTraumas: pbLife.GetChildhoodTraumas(),
		AdultTraumas:     pbLife.GetAdultTraumas(),
		MajorEvents:      pbToTypesLifeEvents(pbLife.GetMajorEvents()),
		EducationHistory: pbToTypesEducationHistory(pbLife.GetEducationHistory()),
		CareerHistory:    pbToTypesCareerHistory(pbLife.GetCareerHistory()),
	}
}

func typesToPbLifeHistory(life *types.LifeHistory) *pb.LifeHistory {
	if life == nil {
		return nil
	}
	return &pb.LifeHistory{
		ChildhoodTraumas: life.ChildhoodTraumas,
		AdultTraumas:     life.AdultTraumas,
		MajorEvents:      typesToPbLifeEvents(life.MajorEvents),
		EducationHistory: typesToPbEducationHistory(life.EducationHistory),
		CareerHistory:    typesToPbCareerHistory(life.CareerHistory),
	}
}

// LifeEvent conversion
func pbToTypesLifeEvents(pbEvents []*pb.LifeEvent) []*types.LifeEvent {
	if pbEvents == nil {
		return nil
	}
	events := make([]*types.LifeEvent, len(pbEvents))
	for i, pbEvent := range pbEvents {
		events[i] = &types.LifeEvent{
			Type:        pbEvent.GetType(),
			Description: pbEvent.GetDescription(),
			Age:         int(pbEvent.GetAge()),
			Impact:      pbEvent.GetImpact(),
		}
		if pbEvent.GetDate() != nil {
			events[i].Date = pbEvent.GetDate().AsTime()
		}
	}
	return events
}

func typesToPbLifeEvents(events []*types.LifeEvent) []*pb.LifeEvent {
	if events == nil {
		return nil
	}
	pbEvents := make([]*pb.LifeEvent, len(events))
	for i, event := range events {
		pbEvents[i] = &pb.LifeEvent{
			Type:        event.Type,
			Description: event.Description,
			Age:         int32(event.Age),
			Impact:      event.Impact,
		}
		if !event.Date.IsZero() {
			pbEvents[i].Date = timestamppb.New(event.Date)
		}
	}
	return pbEvents
}

// Education conversion
func pbToTypesEducationHistory(pbEdu []*pb.Education) []*types.Education {
	if pbEdu == nil {
		return nil
	}
	edu := make([]*types.Education, len(pbEdu))
	for i, pbE := range pbEdu {
		edu[i] = &types.Education{
			Level:       pbE.GetLevel(),
			Field:       pbE.GetField(),
			Institution: pbE.GetInstitution(),
			Performance: pbE.GetPerformance(),
		}
		if pbE.GetGraduation() != nil {
			edu[i].Graduation = pbE.GetGraduation().AsTime()
		}
	}
	return edu
}

func typesToPbEducationHistory(edu []*types.Education) []*pb.Education {
	if edu == nil {
		return nil
	}
	pbEdu := make([]*pb.Education, len(edu))
	for i, e := range edu {
		pbEdu[i] = &pb.Education{
			Level:       e.Level,
			Field:       e.Field,
			Institution: e.Institution,
			Performance: e.Performance,
		}
		if !e.Graduation.IsZero() {
			pbEdu[i].Graduation = timestamppb.New(e.Graduation)
		}
	}
	return pbEdu
}

// Career conversion
func pbToTypesCareerHistory(pbCareer []*pb.Career) []*types.Career {
	if pbCareer == nil {
		return nil
	}
	career := make([]*types.Career, len(pbCareer))
	for i, pbC := range pbCareer {
		career[i] = &types.Career{
			Title:     pbC.GetTitle(),
			Industry:  pbC.GetIndustry(),
			Company:   pbC.GetCompany(),
			IsCurrent: pbC.GetIsCurrent(),
			Salary:    pbC.GetSalary(),
		}
		if pbC.GetStartDate() != nil {
			career[i].StartDate = pbC.GetStartDate().AsTime()
		}
		if pbC.GetEndDate() != nil {
			career[i].EndDate = pbC.GetEndDate().AsTime()
		}
	}
	return career
}

func typesToPbCareerHistory(career []*types.Career) []*pb.Career {
	if career == nil {
		return nil
	}
	pbCareer := make([]*pb.Career, len(career))
	for i, c := range career {
		pbCareer[i] = &pb.Career{
			Title:     c.Title,
			Industry:  c.Industry,
			Company:   c.Company,
			IsCurrent: c.IsCurrent,
			Salary:    c.Salary,
		}
		if !c.StartDate.IsZero() {
			pbCareer[i].StartDate = timestamppb.New(c.StartDate)
		}
		if !c.EndDate.IsZero() {
			pbCareer[i].EndDate = timestamppb.New(c.EndDate)
		}
	}
	return pbCareer
}

// CulturalReligious conversion
func pbToTypesCulturalReligious(pbCult *pb.CulturalReligious) *types.CulturalReligious {
	if pbCult == nil {
		return nil
	}
	return &types.CulturalReligious{
		Religion:            pbCult.GetReligion(),
		Spirituality:        pbCult.GetSpirituality(),
		CulturalBackground:  pbCult.GetCulturalBackground(),
		Traditions:          pbCult.GetTraditions(),
		Holidays:            pbCult.GetHolidays(),
		DietaryRestrictions: pbCult.GetDietaryRestrictions(),
	}
}

func typesToPbCulturalReligious(cult *types.CulturalReligious) *pb.CulturalReligious {
	if cult == nil {
		return nil
	}
	return &pb.CulturalReligious{
		Religion:            cult.Religion,
		Spirituality:        cult.Spirituality,
		CulturalBackground:  cult.CulturalBackground,
		Traditions:          cult.Traditions,
		Holidays:            cult.Holidays,
		DietaryRestrictions: cult.DietaryRestrictions,
	}
}

// PoliticalSocial conversion
func pbToTypesPoliticalSocial(pbPol *pb.PoliticalSocial) *types.PoliticalSocial {
	if pbPol == nil {
		return nil
	}
	return &types.PoliticalSocial{
		PoliticalLeaning: pbPol.GetPoliticalLeaning(),
		Activism:         pbPol.GetActivism(),
		SocialGroups:     pbPol.GetSocialGroups(),
		Causes:           pbPol.GetCauses(),
		VotingHistory:    pbPol.GetVotingHistory(),
		MediaConsumption: pbPol.GetMediaConsumption(),
	}
}

func typesToPbPoliticalSocial(pol *types.PoliticalSocial) *pb.PoliticalSocial {
	if pol == nil {
		return nil
	}
	return &pb.PoliticalSocial{
		PoliticalLeaning: pol.PoliticalLeaning,
		Activism:         pol.Activism,
		SocialGroups:     pol.SocialGroups,
		Causes:           pol.Causes,
		VotingHistory:    pol.VotingHistory,
		MediaConsumption: pol.MediaConsumption,
	}
}

// Health conversion
func pbToTypesHealth(pbHealth *pb.Health) *types.Health {
	if pbHealth == nil {
		return nil
	}
	return &types.Health{
		PhysicalHealth:    pbHealth.GetPhysicalHealth(),
		MentalHealth:      pbHealth.GetMentalHealth(),
		Disabilities:      pbHealth.GetDisabilities(),
		ChronicConditions: pbHealth.GetChronicConditions(),
		Addictions:        pbHealth.GetAddictions(),
		Medications:       pbHealth.GetMedications(),
	}
}

func typesToPbHealth(health *types.Health) *pb.Health {
	if health == nil {
		return nil
	}
	return &pb.Health{
		PhysicalHealth:    health.PhysicalHealth,
		MentalHealth:      health.MentalHealth,
		Disabilities:      health.Disabilities,
		ChronicConditions: health.ChronicConditions,
		Addictions:        health.Addictions,
		Medications:       health.Medications,
	}
}

// Preferences conversion
func pbToTypesPreferences(pbPref *pb.Preferences) *types.Preferences {
	if pbPref == nil {
		return nil
	}
	return &types.Preferences{
		Hobbies:        pbPref.GetHobbies(),
		Interests:      pbPref.GetInterests(),
		FavoriteFoods:  pbPref.GetFavoriteFoods(),
		FavoriteMusic:  pbPref.GetFavoriteMusic(),
		FavoriteMovies: pbPref.GetFavoriteMovies(),
		FavoriteBooks:  pbPref.GetFavoriteBooks(),
		TechnologyUse:  pbPref.GetTechnologyUse(),
		TravelStyle:    pbPref.GetTravelStyle(),
	}
}

func typesToPbPreferences(pref *types.Preferences) *pb.Preferences {
	if pref == nil {
		return nil
	}
	return &pb.Preferences{
		Hobbies:        pref.Hobbies,
		Interests:      pref.Interests,
		FavoriteFoods:  pref.FavoriteFoods,
		FavoriteMusic:  pref.FavoriteMusic,
		FavoriteMovies: pref.FavoriteMovies,
		FavoriteBooks:  pref.FavoriteBooks,
		TechnologyUse:  pref.TechnologyUse,
		TravelStyle:    pref.TravelStyle,
	}
}

// BehavioralTendencies conversion
func pbToTypesBehavioralTendencies(pbBehav *pb.BehavioralTendencies) *types.BehavioralTendencies {
	if pbBehav == nil {
		return nil
	}
	return &types.BehavioralTendencies{
		DecisionMaking:     pbBehav.GetDecisionMaking(),
		ConflictResolution: pbBehav.GetConflictResolution(),
		CommunicationStyle: pbBehav.GetCommunicationStyle(),
		LeadershipStyle:    pbBehav.GetLeadershipStyle(),
		CopingMechanisms:   pbBehav.GetCopingMechanisms(),
		StressResponse:     pbBehav.GetStressResponse(),
	}
}

func typesToPbBehavioralTendencies(behav *types.BehavioralTendencies) *pb.BehavioralTendencies {
	if behav == nil {
		return nil
	}
	return &pb.BehavioralTendencies{
		DecisionMaking:     behav.DecisionMaking,
		ConflictResolution: behav.ConflictResolution,
		CommunicationStyle: behav.CommunicationStyle,
		LeadershipStyle:    behav.LeadershipStyle,
		CopingMechanisms:   behav.CopingMechanisms,
		StressResponse:     behav.StressResponse,
	}
}

// CurrentContext conversion
func pbToTypesCurrentContext(pbContext *pb.CurrentContext) *types.CurrentContext {
	if pbContext == nil {
		return nil
	}
	return &types.CurrentContext{
		Mood:         pbContext.GetMood(),
		StressLevel:  pbContext.GetStressLevel(),
		CurrentGoals: pbContext.GetCurrentGoals(),
		RecentEvents: pbContext.GetRecentEvents(),
		LifeStage:    pbContext.GetLifeStage(),
	}
}

func typesToPbCurrentContext(context *types.CurrentContext) *pb.CurrentContext {
	if context == nil {
		return nil
	}
	return &pb.CurrentContext{
		Mood:         context.Mood,
		StressLevel:  context.StressLevel,
		CurrentGoals: context.CurrentGoals,
		RecentEvents: context.RecentEvents,
		LifeStage:    context.LifeStage,
	}
}
