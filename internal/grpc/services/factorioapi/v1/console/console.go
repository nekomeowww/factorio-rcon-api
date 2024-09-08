package consolev1

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorcon/rcon"
	v1 "github.com/nekomeowww/factorio-rcon-api/apis/factorioapi/v1"
	"github.com/nekomeowww/factorio-rcon-api/pkg/apierrors"
	"github.com/nekomeowww/xo/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NewConsoleServiceParams struct {
	fx.In

	Logger *logger.Logger
	RCON   *rcon.Conn
}

type ConsoleService struct {
	v1.UnimplementedConsoleServiceServer

	logger *logger.Logger
	rcon   *rcon.Conn
}

func NewConsoleService() func(NewConsoleServiceParams) *ConsoleService {
	return func(params NewConsoleServiceParams) *ConsoleService {
		return &ConsoleService{
			logger: params.Logger,
			rcon:   params.RCON,
		}
	}
}

func (s *ConsoleService) CommandRaw(ctx context.Context, req *v1.CommandRawRequest) (*v1.CommandRawResponse, error) {
	resp, err := s.rcon.Execute(req.Input)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	return &v1.CommandRawResponse{
		Output: resp,
	}, nil
}

func (s *ConsoleService) CommandMessage(ctx context.Context, req *v1.CommandMessageRequest) (*v1.CommandMessageResponse, error) {
	if req.Message == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("message should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(req.Message)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command message and got response", zap.String("response", resp), zap.String("message", req.Message))

	return &v1.CommandMessageResponse{}, nil
}

func (s *ConsoleService) CommandAlerts(ctx context.Context, req *v1.CommandAlertsRequest) (*v1.CommandAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandAlerts not implemented")
}

func (s *ConsoleService) CommandEnableResearchQueue(ctx context.Context, req *v1.CommandEnableResearchQueueRequest) (*v1.CommandEnableResearchQueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandEnableResearchQueue not implemented")
}

func (s *ConsoleService) CommandMuteProgrammableSpeakerForEveryone(ctx context.Context, req *v1.CommandMuteProgrammableSpeakerForEveryoneRequest) (*v1.CommandMuteProgrammableSpeakerForEveryoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandMuteProgrammableSpeakerForEveryone not implemented")
}

func (s *ConsoleService) CommandUnmuteProgrammableSpeakerForEveryone(ctx context.Context, req *v1.CommandUnmuteProgrammableSpeakerForEveryoneRequest) (*v1.CommandUnmuteProgrammableSpeakerForEveryoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandUnmuteProgrammableSpeakerForEveryone not implemented")
}

func (s *ConsoleService) CommandPermissions(ctx context.Context, req *v1.CommandPermissionsRequest) (*v1.CommandPermissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissions not implemented")
}

func (s *ConsoleService) CommandPermissionsAddPlayer(ctx context.Context, req *v1.CommandPermissionsAddPlayerRequest) (*v1.CommandPermissionsAddPlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsAddPlayer not implemented")
}

func (s *ConsoleService) CommandPermissionsCreateGroup(ctx context.Context, req *v1.CommandPermissionsCreateGroupRequest) (*v1.CommandPermissionsCreateGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsCreateGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsDeleteGroup(ctx context.Context, req *v1.CommandPermissionsDeleteGroupRequest) (*v1.CommandPermissionsDeleteGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsDeleteGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsEditGroup(ctx context.Context, req *v1.CommandPermissionsEditGroupRequest) (*v1.CommandPermissionsEditGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsEditGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsGetPlayerGroup(ctx context.Context, req *v1.CommandPermissionsGetPlayerGroupRequest) (*v1.CommandPermissionsGetPlayerGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsGetPlayerGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsRemovePlayerGroup(ctx context.Context, req *v1.CommandPermissionsRemovePlayerGroupRequest) (*v1.CommandPermissionsRemovePlayerGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsRemovePlayerGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsRenameGroup(ctx context.Context, req *v1.CommandPermissionsRenameGroupRequest) (*v1.CommandPermissionsRenameGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsRenameGroup not implemented")
}

func (s *ConsoleService) CommandResetTips(ctx context.Context, req *v1.CommandResetTipsRequest) (*v1.CommandResetTipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandResetTips not implemented")
}

var (
	regexpEvolutionFactor = regexp.MustCompile(`Evolution factor: ([0-9.]+)\. \(Time ([0-9.]+)%\) \(Pollution ([0-9.]+)%\) \(Spawner kills ([0-9.]+)%\)`)
)

func (s *ConsoleService) CommandEvolution(ctx context.Context, req *v1.CommandEvolutionRequest) (*v1.CommandEvolutionResponse, error) {
	resp, err := s.rcon.Execute("/evolution")
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	normalized := strings.TrimSuffix(resp, "\n")

	// output: Evolution factor: 0.0000. (Time 0%) (Pollution 0%) (Spawner kills 0%)
	match := regexpEvolutionFactor.FindStringSubmatch(resp)
	if len(match) != 5 {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse evolution factor: %s due to matches not equals 5", normalized).AsStatus()
	}

	// match[1] is the evolution factor
	evolutionFactor, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse evolution factor: %s from %s due to %v", match[1], normalized, err).AsStatus()
	}

	// match[2] is the time
	time, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse time: %s from %s due to %v", match[2], normalized, err).AsStatus()
	}

	// match[3] is the pollution
	pollution, err := strconv.ParseFloat(match[3], 64)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse pollution: %s from %s due to %v", match[3], normalized, err).AsStatus()
	}

	// match[4] is the spawner kills
	spawnerKills, err := strconv.ParseFloat(match[4], 64)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse spawner kills: %s from %s due to %v", match[4], normalized, err).AsStatus()
	}

	return &v1.CommandEvolutionResponse{
		EvolutionFactor: evolutionFactor,
		Time:            time,
		Pollution:       pollution,
		SpawnerKills:    spawnerKills,
	}, nil
}

func (s *ConsoleService) CommandSeed(ctx context.Context, req *v1.CommandSeedRequest) (*v1.CommandSeedResponse, error) {
	resp, err := s.rcon.Execute("/seed")
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	normalized := strings.TrimSuffix(resp, "\n")
	if normalized == "" {
		return &v1.CommandSeedResponse{}, nil
	}

	return &v1.CommandSeedResponse{
		Seed: normalized,
	}, nil
}

func (s *ConsoleService) CommandTime(ctx context.Context, req *v1.CommandTimeRequest) (*v1.CommandTimeResponse, error) {
	resp, err := s.rcon.Execute("/time")
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	normalized := strings.TrimSuffix(resp, "\n")
	if normalized == "" {
		return &v1.CommandTimeResponse{}, nil
	}

	return &v1.CommandTimeResponse{}, nil
}

func (s *ConsoleService) CommandToggleActionLogging(ctx context.Context, req *v1.CommandToggleActionLoggingRequest) (*v1.CommandToggleActionLoggingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandToggleActionLogging not implemented")
}

func (s *ConsoleService) CommandToggleHeavyMode(ctx context.Context, req *v1.CommandToggleHeavyModeRequest) (*v1.CommandToggleHeavyModeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandToggleHeavyMode not implemented")
}

func (s *ConsoleService) CommandUnlockShortcutBar(ctx context.Context, req *v1.CommandUnlockShortcutBarRequest) (*v1.CommandUnlockShortcutBarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandUnlockShortcutBar not implemented")
}

func (s *ConsoleService) CommandUnlockTips(ctx context.Context, req *v1.CommandUnlockTipsRequest) (*v1.CommandUnlockTipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandUnlockTips not implemented")
}

func (s *ConsoleService) CommandVersion(ctx context.Context, req *v1.CommandVersionRequest) (*v1.CommandVersionResponse, error) {
	resp, err := s.rcon.Execute("/version")
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	return &v1.CommandVersionResponse{
		Version: strings.TrimSuffix(resp, "\n"),
	}, nil
}

func (s *ConsoleService) CommandAdmins(ctx context.Context, req *v1.CommandAdminsRequest) (*v1.CommandAdminsResponse, error) {
	resp, err := s.rcon.Execute("/admins")
	if err != nil {
		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	admins, err := StringListToPlayers(resp)
	if err != nil {
		return nil, err
	}

	return &v1.CommandAdminsResponse{
		Admins: admins,
	}, nil
}

func (s *ConsoleService) CommandBan(ctx context.Context, req *v1.CommandBanRequest) (*v1.CommandBanResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	cmd := "/ban " + req.Username

	resp, err := s.rcon.Execute(cmd)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command ban and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandBanResponse{}, nil
}

func (s *ConsoleService) CommandBans(ctx context.Context, req *v1.CommandBansRequest) (*v1.CommandBansResponse, error) {
	resp, err := s.rcon.Execute("/bans")
	if err != nil {
		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	bans, err := StringListToPlayers(resp)
	if err != nil {
		return nil, err
	}

	return &v1.CommandBansResponse{
		Bans: bans,
	}, nil
}

func (s *ConsoleService) CommandDemote(ctx context.Context, req *v1.CommandDemoteRequest) (*v1.CommandDemoteResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/demote " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command demote and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandDemoteResponse{}, nil
}

func (s *ConsoleService) CommandIgnore(ctx context.Context, req *v1.CommandIgnoreRequest) (*v1.CommandIgnoreResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/ignore " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command ignore and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandIgnoreResponse{}, nil
}

func (s *ConsoleService) CommandKick(ctx context.Context, req *v1.CommandKickRequest) (*v1.CommandKickResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	cmd := "/kick " + req.Username
	if req.Reason != "" {
		cmd += " " + req.Reason
	}

	resp, err := s.rcon.Execute(cmd)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command kick and got response", zap.String("response", resp), zap.String("username", req.Username), zap.String("reason", req.Reason))

	return &v1.CommandKickResponse{}, nil
}

func (s *ConsoleService) CommandMute(ctx context.Context, req *v1.CommandMuteRequest) (*v1.CommandMuteResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/mute " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command mute and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandMuteResponse{}, nil
}

func (s *ConsoleService) CommandMutes(ctx context.Context, req *v1.CommandMutesRequest) (*v1.CommandMutesResponse, error) {
	resp, err := s.rcon.Execute("/mutes")
	if err != nil {
		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	mutes, err := StringListToPlayers(resp)
	if err != nil {
		return nil, err
	}

	return &v1.CommandMutesResponse{
		Mutes: mutes,
	}, nil
}

func (s *ConsoleService) CommandPlayers(ctx context.Context, req *v1.CommandPlayersRequest) (*v1.CommandPlayersResponse, error) {
	resp, err := s.rcon.Execute("/players")
	if err != nil {
		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	players, err := StringListToPlayers(resp)
	if err != nil {
		return nil, err
	}

	return &v1.CommandPlayersResponse{
		Players: players,
	}, nil
}

func (s *ConsoleService) CommandPromote(ctx context.Context, req *v1.CommandPromoteRequest) (*v1.CommandPromoteResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/promote " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command promote and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandPromoteResponse{}, nil
}

func (s *ConsoleService) CommandPurge(ctx context.Context, req *v1.CommandPurgeRequest) (*v1.CommandPurgeResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/purge" + " " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command purge and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandPurgeResponse{}, nil
}

func (s *ConsoleService) CommandServerSave(ctx context.Context, req *v1.CommandServerSaveRequest) (*v1.CommandServerSaveResponse, error) {
	resp, err := s.rcon.Execute("/server-save")
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command server-save and got response", zap.String("response", resp))

	return &v1.CommandServerSaveResponse{}, nil
}

func (s *ConsoleService) CommandUnban(ctx context.Context, req *v1.CommandUnbanRequest) (*v1.CommandUnbanResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/unban " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command unban and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandUnbanResponse{}, nil
}

func (s *ConsoleService) CommandUnignore(ctx context.Context, req *v1.CommandUnignoreRequest) (*v1.CommandUnignoreResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/unignore " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command unignore and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandUnignoreResponse{}, nil
}

func (s *ConsoleService) CommandUnmute(ctx context.Context, req *v1.CommandUnmuteRequest) (*v1.CommandUnmuteResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/unmute " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command unmute and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandUnmuteResponse{}, nil
}

func (s *ConsoleService) CommandWhisper(ctx context.Context, req *v1.CommandWhisperRequest) (*v1.CommandWhisperResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}
	if req.Message == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("message should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/whisper " + req.Username + " " + req.Message)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command whisper and got response", zap.String("response", resp), zap.String("username", req.Username), zap.String("message", req.Message))

	return &v1.CommandWhisperResponse{}, nil
}

func (s *ConsoleService) CommandWhitelistAdd(ctx context.Context, req *v1.CommandWhitelistAddRequest) (*v1.CommandWhitelistAddResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/whitelist add " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command whitelist add and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandWhitelistAddResponse{}, nil
}

func (s *ConsoleService) CommandWhitelistGet(ctx context.Context, req *v1.CommandWhitelistGetRequest) (*v1.CommandWhitelistGetResponse, error) {
	resp, err := s.rcon.Execute("/whitelist get")
	if err != nil {
		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	whitelist, err := StringListToPlayers(resp)
	if err != nil {
		return nil, err
	}

	return &v1.CommandWhitelistGetResponse{
		Whitelist: whitelist,
	}, nil
}

func (s *ConsoleService) CommandWhitelistRemove(ctx context.Context, req *v1.CommandWhitelistRemoveRequest) (*v1.CommandWhitelistRemoveResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute("/whitelist remove " + req.Username)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command whitelist remove and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v1.CommandWhitelistRemoveResponse{}, nil
}

func (s *ConsoleService) CommandWhitelistClear(ctx context.Context, req *v1.CommandWhitelistClearRequest) (*v1.CommandWhitelistClearResponse, error) {
	resp, err := s.rcon.Execute("/whitelist clear")
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command whitelist clear and got response", zap.String("response", resp))

	return &v1.CommandWhitelistClearResponse{}, nil
}
