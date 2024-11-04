package consolev2

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"

	v2 "github.com/nekomeowww/factorio-rcon-api/apis/factorioapi/v2"
	"github.com/nekomeowww/factorio-rcon-api/internal/rcon"
	"github.com/nekomeowww/factorio-rcon-api/pkg/apierrors"
	"github.com/nekomeowww/factorio-rcon-api/pkg/utils"
	"github.com/nekomeowww/xo/logger"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NewConsoleServiceParams struct {
	fx.In

	Logger *logger.Logger
	RCON   rcon.RCON
}

type ConsoleService struct {
	v2.UnimplementedConsoleServiceServer

	logger *logger.Logger
	rcon   rcon.RCON
}

func NewConsoleService() func(NewConsoleServiceParams) *ConsoleService {
	return func(params NewConsoleServiceParams) *ConsoleService {
		return &ConsoleService{
			logger: params.Logger,
			rcon:   params.RCON,
		}
	}
}

func (s *ConsoleService) CommandRaw(ctx context.Context, req *v2.CommandRawRequest) (*v2.CommandRawResponse, error) {
	resp, err := s.rcon.Execute(ctx, req.Input)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	return &v2.CommandRawResponse{
		Output: resp,
	}, nil
}

func (s *ConsoleService) CommandMessage(ctx context.Context, req *v2.CommandMessageRequest) (*v2.CommandMessageResponse, error) {
	if req.Message == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("message should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, req.Message)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command message and got response", zap.String("response", resp), zap.String("message", req.Message))

	return &v2.CommandMessageResponse{}, nil
}

func (s *ConsoleService) CommandAlerts(ctx context.Context, req *v2.CommandAlertsRequest) (*v2.CommandAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandAlerts not implemented")
}

func (s *ConsoleService) CommandEnableResearchQueue(ctx context.Context, req *v2.CommandEnableResearchQueueRequest) (*v2.CommandEnableResearchQueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandEnableResearchQueue not implemented")
}

func (s *ConsoleService) CommandMuteProgrammableSpeakerForEveryone(ctx context.Context, req *v2.CommandMuteProgrammableSpeakerForEveryoneRequest) (*v2.CommandMuteProgrammableSpeakerForEveryoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandMuteProgrammableSpeakerForEveryone not implemented")
}

func (s *ConsoleService) CommandUnmuteProgrammableSpeakerForEveryone(ctx context.Context, req *v2.CommandUnmuteProgrammableSpeakerForEveryoneRequest) (*v2.CommandUnmuteProgrammableSpeakerForEveryoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandUnmuteProgrammableSpeakerForEveryone not implemented")
}

func (s *ConsoleService) CommandPermissions(ctx context.Context, req *v2.CommandPermissionsRequest) (*v2.CommandPermissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissions not implemented")
}

func (s *ConsoleService) CommandPermissionsAddPlayer(ctx context.Context, req *v2.CommandPermissionsAddPlayerRequest) (*v2.CommandPermissionsAddPlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsAddPlayer not implemented")
}

func (s *ConsoleService) CommandPermissionsCreateGroup(ctx context.Context, req *v2.CommandPermissionsCreateGroupRequest) (*v2.CommandPermissionsCreateGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsCreateGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsDeleteGroup(ctx context.Context, req *v2.CommandPermissionsDeleteGroupRequest) (*v2.CommandPermissionsDeleteGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsDeleteGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsEditGroup(ctx context.Context, req *v2.CommandPermissionsEditGroupRequest) (*v2.CommandPermissionsEditGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsEditGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsGetPlayerGroup(ctx context.Context, req *v2.CommandPermissionsGetPlayerGroupRequest) (*v2.CommandPermissionsGetPlayerGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsGetPlayerGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsRemovePlayerGroup(ctx context.Context, req *v2.CommandPermissionsRemovePlayerGroupRequest) (*v2.CommandPermissionsRemovePlayerGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsRemovePlayerGroup not implemented")
}

func (s *ConsoleService) CommandPermissionsRenameGroup(ctx context.Context, req *v2.CommandPermissionsRenameGroupRequest) (*v2.CommandPermissionsRenameGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandPermissionsRenameGroup not implemented")
}

func (s *ConsoleService) CommandResetTips(ctx context.Context, req *v2.CommandResetTipsRequest) (*v2.CommandResetTipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandResetTips not implemented")
}

var (
	regexpEvolutionFactor = regexp.MustCompile(`(.*) - Evolution factor: ([0-9.]+)\. \(Time ([0-9.]+)%\) \(Pollution ([0-9.]+)%\) \(Spawner kills ([0-9.]+)%\)`)
)

func parseEvolutionFactorLine(resp string) (*v2.Evolution, error) {
	normalized := strings.TrimSuffix(resp, "\n")

	// output: Nauvis - Evolution factor: 0.0000. (Time 0%) (Pollution 0%) (Spawner kills 0%)
	match := regexpEvolutionFactor.FindStringSubmatch(resp)
	if len(match) != 6 {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse evolution factor: %s due to matches not equals 6", normalized).AsStatus()
	}

	// match[1] is the name
	planetName := match[1]

	// match[2] is the evolution factor
	evolutionFactor, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse evolution factor: %s from %s due to %v", match[1], normalized, err).AsStatus()
	}

	// match[3] is the time
	time, err := strconv.ParseFloat(match[3], 64)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse time: %s from %s due to %v", match[2], normalized, err).AsStatus()
	}

	// match[4] is the pollution
	pollution, err := strconv.ParseFloat(match[4], 64)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse pollution: %s from %s due to %v", match[3], normalized, err).AsStatus()
	}

	// match[5] is the spawner kills
	spawnerKills, err := strconv.ParseFloat(match[5], 64)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse spawner kills: %s from %s due to %v", match[4], normalized, err).AsStatus()
	}

	return &v2.Evolution{
		SurfaceName:     planetName,
		EvolutionFactor: evolutionFactor,
		Time:            time,
		Pollution:       pollution,
		SpawnerKills:    spawnerKills,
	}, nil
}

func (s *ConsoleService) CommandEvolution(ctx context.Context, req *v2.CommandEvolutionRequest) (*v2.CommandEvolutionResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/evolution")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	lines := strings.Split(resp, "\n")

	lines = lo.Map(lines, func(item string, _ int) string {
		return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(item), "\r", ""), "\n", "")
	})
	lines = lo.Filter(lines, func(item string, _ int) bool {
		return item != ""
	})

	evolutions := make([]*v2.Evolution, 0, len(lines))
	for _, line := range lines {
		evolution, err := parseEvolutionFactorLine(line)
		if err != nil {
			return nil, err
		}

		evolutions = append(evolutions, evolution)
	}

	return &v2.CommandEvolutionResponse{
		Evolutions: evolutions,
	}, nil
}

func (s *ConsoleService) CommandEvolutionGet(ctx context.Context, req *v2.CommandEvolutionGetRequest) (*v2.CommandEvolutionGetResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/evolution "+req.SurfaceName)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}
	if strings.Contains(resp, "does not exist") {
		return nil, apierrors.NewErrNotFound().WithDetailf("surface %s does not exist", req.SurfaceName).AsStatus()
	}

	evolution, err := parseEvolutionFactorLine(resp)
	if err != nil {
		return nil, err
	}

	return &v2.CommandEvolutionGetResponse{
		Evolution: evolution,
	}, nil
}

func (s *ConsoleService) CommandSeed(ctx context.Context, req *v2.CommandSeedRequest) (*v2.CommandSeedResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/seed")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	normalized := strings.TrimSuffix(resp, "\n")
	if normalized == "" {
		return &v2.CommandSeedResponse{}, nil
	}

	return &v2.CommandSeedResponse{
		Seed: normalized,
	}, nil
}

func (s *ConsoleService) CommandTime(ctx context.Context, req *v2.CommandTimeRequest) (*v2.CommandTimeResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/time")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	normalized := strings.TrimSuffix(resp, "\n")
	if normalized == "" {
		return &v2.CommandTimeResponse{}, nil
	}

	duration, err := utils.ParseDuration(normalized)
	if err != nil {
		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	return &v2.CommandTimeResponse{
		Time: duration.Seconds(),
	}, nil
}

func (s *ConsoleService) CommandToggleActionLogging(ctx context.Context, req *v2.CommandToggleActionLoggingRequest) (*v2.CommandToggleActionLoggingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandToggleActionLogging not implemented")
}

func (s *ConsoleService) CommandToggleHeavyMode(ctx context.Context, req *v2.CommandToggleHeavyModeRequest) (*v2.CommandToggleHeavyModeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandToggleHeavyMode not implemented")
}

func (s *ConsoleService) CommandUnlockShortcutBar(ctx context.Context, req *v2.CommandUnlockShortcutBarRequest) (*v2.CommandUnlockShortcutBarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandUnlockShortcutBar not implemented")
}

func (s *ConsoleService) CommandUnlockTips(ctx context.Context, req *v2.CommandUnlockTipsRequest) (*v2.CommandUnlockTipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandUnlockTips not implemented")
}

func (s *ConsoleService) CommandVersion(ctx context.Context, req *v2.CommandVersionRequest) (*v2.CommandVersionResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/version")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	return &v2.CommandVersionResponse{
		Version: strings.TrimSuffix(resp, "\n"),
	}, nil
}

func (s *ConsoleService) CommandAdmins(ctx context.Context, req *v2.CommandAdminsRequest) (*v2.CommandAdminsResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/admins")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	admins, err := utils.StringListToPlayers(resp)
	if err != nil {
		return nil, err
	}

	return &v2.CommandAdminsResponse{
		Admins: utils.MapV1PlayersToV2Players(admins),
	}, nil
}

func (s *ConsoleService) CommandBan(ctx context.Context, req *v2.CommandBanRequest) (*v2.CommandBanResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	cmd := "/ban " + req.Username

	resp, err := s.rcon.Execute(ctx, cmd)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command ban and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandBanResponse{}, nil
}

func (s *ConsoleService) CommandBans(ctx context.Context, req *v2.CommandBansRequest) (*v2.CommandBansResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/bans")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	bans, err := utils.StringListToPlayers(resp)
	if err != nil {
		return nil, err
	}

	return &v2.CommandBansResponse{
		Bans: utils.MapV1PlayersToV2Players(bans),
	}, nil
}

func (s *ConsoleService) CommandDemote(ctx context.Context, req *v2.CommandDemoteRequest) (*v2.CommandDemoteResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/demote "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command demote and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandDemoteResponse{}, nil
}

func (s *ConsoleService) CommandIgnore(ctx context.Context, req *v2.CommandIgnoreRequest) (*v2.CommandIgnoreResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/ignore "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command ignore and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandIgnoreResponse{}, nil
}

func (s *ConsoleService) CommandKick(ctx context.Context, req *v2.CommandKickRequest) (*v2.CommandKickResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	cmd := "/kick " + req.Username
	if req.Reason != "" {
		cmd += " " + req.Reason
	}

	resp, err := s.rcon.Execute(ctx, cmd)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command kick and got response", zap.String("response", resp), zap.String("username", req.Username), zap.String("reason", req.Reason))

	return &v2.CommandKickResponse{}, nil
}

func (s *ConsoleService) CommandMute(ctx context.Context, req *v2.CommandMuteRequest) (*v2.CommandMuteResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/mute "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command mute and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandMuteResponse{}, nil
}

func (s *ConsoleService) CommandMutes(ctx context.Context, req *v2.CommandMutesRequest) (*v2.CommandMutesResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/mutes")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	mutes, err := utils.StringListToPlayers(resp)
	if err != nil {
		return nil, err
	}

	return &v2.CommandMutesResponse{
		Mutes: utils.MapV1PlayersToV2Players(mutes),
	}, nil
}

func (s *ConsoleService) CommandPlayers(ctx context.Context, req *v2.CommandPlayersRequest) (*v2.CommandPlayersResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/players")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	lines := strings.Split(resp, "\n")
	lines = lines[1 : len(lines)-1]

	players, err := utils.StringListToPlayers(strings.Join(lines, "\n"))
	if err != nil {
		return nil, err
	}

	return &v2.CommandPlayersResponse{
		Players: utils.MapV1PlayersToV2Players(players),
	}, nil
}

func (s *ConsoleService) CommandPromote(ctx context.Context, req *v2.CommandPromoteRequest) (*v2.CommandPromoteResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/promote "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command promote and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandPromoteResponse{}, nil
}

func (s *ConsoleService) CommandPurge(ctx context.Context, req *v2.CommandPurgeRequest) (*v2.CommandPurgeResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/purge"+" "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command purge and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandPurgeResponse{}, nil
}

func (s *ConsoleService) CommandServerSave(ctx context.Context, req *v2.CommandServerSaveRequest) (*v2.CommandServerSaveResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/server-save")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command server-save and got response", zap.String("response", resp))

	return &v2.CommandServerSaveResponse{}, nil
}

func (s *ConsoleService) CommandUnban(ctx context.Context, req *v2.CommandUnbanRequest) (*v2.CommandUnbanResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/unban "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command unban and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandUnbanResponse{}, nil
}

func (s *ConsoleService) CommandUnignore(ctx context.Context, req *v2.CommandUnignoreRequest) (*v2.CommandUnignoreResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/unignore "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command unignore and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandUnignoreResponse{}, nil
}

func (s *ConsoleService) CommandUnmute(ctx context.Context, req *v2.CommandUnmuteRequest) (*v2.CommandUnmuteResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/unmute "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command unmute and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandUnmuteResponse{}, nil
}

func (s *ConsoleService) CommandWhisper(ctx context.Context, req *v2.CommandWhisperRequest) (*v2.CommandWhisperResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}
	if req.Message == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("message should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/whisper "+req.Username+" "+req.Message)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command whisper and got response", zap.String("response", resp), zap.String("username", req.Username), zap.String("message", req.Message))

	return &v2.CommandWhisperResponse{}, nil
}

func (s *ConsoleService) CommandWhitelistAdd(ctx context.Context, req *v2.CommandWhitelistAddRequest) (*v2.CommandWhitelistAddResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/whitelist add "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command whitelist add and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandWhitelistAddResponse{}, nil
}

func (s *ConsoleService) CommandWhitelistGet(ctx context.Context, req *v2.CommandWhitelistGetRequest) (*v2.CommandWhitelistGetResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/whitelist get")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrInternal().WithDetail(err.Error()).WithError(err).WithCaller().AsStatus()
	}

	resp = strings.TrimPrefix(resp, "Whitelisted players:")
	resp = strings.TrimSpace(resp)

	playerNames := utils.ParseWhitelistedPlayers(resp)

	players := lo.Map(playerNames, func(player string, _ int) *v2.Player {
		return &v2.Player{
			Username: player,
		}
	})

	savePlayers, err := s.CommandPlayers(ctx, &v2.CommandPlayersRequest{})
	if err != nil {
		return nil, err
	}

	mPlayers := lo.SliceToMap(savePlayers.Players, func(player *v2.Player) (string, *v2.Player) {
		return player.Username, player
	})

	for _, player := range players {
		if p, ok := mPlayers[player.Username]; ok {
			player.Online = p.Online
		}
	}

	return &v2.CommandWhitelistGetResponse{
		Whitelist: players,
	}, nil
}

func (s *ConsoleService) CommandWhitelistRemove(ctx context.Context, req *v2.CommandWhitelistRemoveRequest) (*v2.CommandWhitelistRemoveResponse, error) {
	if req.Username == "" {
		return nil, apierrors.NewErrInvalidArgument().WithDetail("username should not be empty").AsStatus()
	}

	resp, err := s.rcon.Execute(ctx, "/whitelist remove "+req.Username)
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command whitelist remove and got response", zap.String("response", resp), zap.String("username", req.Username))

	return &v2.CommandWhitelistRemoveResponse{}, nil
}

func (s *ConsoleService) CommandWhitelistClear(ctx context.Context, req *v2.CommandWhitelistClearRequest) (*v2.CommandWhitelistClearResponse, error) {
	resp, err := s.rcon.Execute(ctx, "/whitelist clear")
	if err != nil {
		if errors.Is(err, rcon.ErrTimeout) {
			return nil, apierrors.NewErrTimeout().WithDetail("RCON connection is not established within deadline threshold").AsStatus()
		}

		return nil, apierrors.NewErrBadRequest().WithDetail(err.Error()).AsStatus()
	}

	s.logger.Info("executed command whitelist clear and got response", zap.String("response", resp))

	return &v2.CommandWhitelistClearResponse{}, nil
}
