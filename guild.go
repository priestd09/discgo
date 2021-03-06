package discgo

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

type ModelGuild struct {
	ID                              string             `json:"id"`
	Name                            string             `json:"name"`
	Icon                            string             `json:"icon"`
	Splash                          string             `json:"splash"`
	OwnerID                         string             `json:"owner_id"`
	Region                          string             `json:"region"`
	AFKChannelID                    string             `json:"afk_channel_id"`
	AFKTimeout                      int                `json:"afk_timeout"`
	EmbedEnabled                    bool               `json:"embed_enabled"`
	EmbedChannelID                  string             `json:"embed_channel_id"`
	VerificationLevel               int                `json:"verification_level"`
	DefaultMessageNotificationLevel int                `json:"default_message_notifications"`
	Roles                           []*ModelRole       `json:"roles"`
	Emojis                          []*ModelGuildEmoji `json:"emojis"`
	Features                        []string           `json:"features"` // not sure if this is right, DiscordGo doesn't have anything
	MFALevel                        int                `json:"mfa_level"`
	JoinedAt                        time.Time          `json:"joined_at"`
}

const (
	LevelMessageNotificationAllMessages = iota
	LevelMessageNotificationOnlyMentions
)

const (
	LevelExplicitContentFilterDisabled = iota
	LevelExplicitContentFilterMembersWithoutRoles
	LevelExplicitContentFilterAllMembers
)

const (
	LevelMFANone = iota
	LevelMFAElevated
)

const (
	LevelVerificationNone = iota
	LevelVerificationLow
	LevelVerificationMedium
	LevelVerificationHigh
	LevelVerificationVeryHigh
)

type ModelGuildEmbed struct {
	Enabled   bool   `json:"enabled,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
}

type ModelGuildMember struct {
	User     *ModelUser `json:"user"`
	Nick     *string    `json:"nick"`
	Roles    []string   `json:"roles"`
	JoinedAt time.Time  `json:"joined_at"`
	// TODO this two are never updated?
	Deaf bool `json:"deaf"`
	Mute bool `json:"mute"`
}

type ModelIntegration struct {
	ID                string                   `json:"id"`
	Name              string                   `json:"name"`
	Type              string                   `json:"type"`
	Enabled           bool                     `json:"enabled"`
	Syncing           bool                     `json:"syncing"`
	RoleID            string                   `json:"role_id"`
	ExpireBehaviour   int                      `json:"expire_behaviour"`
	ExpireGracePeriod int                      `json:"expire_grace_period"`
	User              *ModelUser               `json:"user"`
	Account           *ModelIntegrationAccount `json:"account"`
	SyncedAt          time.Time                `json:"synced_at"`
}

type ModelIntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ModelGuildEmoji struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	RequireColons bool     `json:"require_colons"`
	Managed       bool     `json:"managed"`
}

type EndpointGuilds struct {
	*endpoint
}

func (c *RESTClient) Guilds() EndpointGuilds {
	e2 := c.rootEndpoint().appendMajor("guilds")
	return EndpointGuilds{e2}
}

type ParamsGuildsCreate struct {
	Name                        string                      `json:"name,omitempty"`
	Region                      string                      `json:"region,omitempty"`
	Icon                        string                      `json:"icon,omitempty"`
	VerificationLevel           int                         `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                         `json:"default_message_notifications,omitempty"`
	Roles                       []*ModelRole                `json:"roles,omitempty"`
	Channels                    []*ParamsGuildChannelCreate `json:"channels,omitempty"`
}

// TODO not sure about the naming on this?
type ParamsGuildChannelCreate struct {
	Name                 string                      `json:"name"`
	Type                 string                      `json:"type,omitempty"`
	Bitrate              int                         `json:"bitrate,omitempty"`
	UserLimit            int                         `json:"user_limit,omitempty"`
	PermissionOverwrites []*ModelPermissionOverwrite `json:"permission_overwrites,omitempty"`
}

// TODO Docs for this are not clear on what the Channels field should be, and the link for that field is broken.
func (e EndpointGuilds) Create(ctx context.Context, params *ParamsGuildsCreate) (g *ModelGuild, err error) {
	return g, e.doMethod(ctx, "POST", params, &g)
}

// TODO all endpoints need to be value receivers and returns!!!
type EndpointGuild struct {
	*endpoint
}

func (c *RESTClient) Guild(gID string) EndpointGuild {
	e2 := c.Guilds().appendMajor(gID)
	return EndpointGuild{e2}
}

func (e EndpointGuild) Get(ctx context.Context, gID string) (g *ModelGuild, err error) {
	return g, e.doMethod(ctx, "GET", nil, &g)
}

type ParamsGuildModify struct {
	Name                        string `json:"name,omitempty"`
	Region                      string `json:"region,omitempty"`
	VerificationLevel           int    `json:"verification_level,omitempty"`
	DefaultMessageNotifications int    `json:"default_message_notifications,omitempty"`
	AFKChannelID                string `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int    `json:"afk_tiemout,omitempty"`
	Icon                        string `json:"icon,omitempty"`
	OwnerID                     string `json:"owner_id,omitempty"`
	Splash                      string `json:"splash,omitempty"`
}

func (e EndpointGuild) Modify(ctx context.Context, params *ParamsGuildModify) (g *ModelGuild, err error) {
	return g, e.doMethod(ctx, "PATCH", params, &g)
}

func (e EndpointGuild) Delete(ctx context.Context) (g *ModelGuild, err error) {
	return g, e.doMethod(ctx, "DELETE", nil, &g)
}

type EndpointGuildMe struct {
	*endpoint
}

func (e EndpointGuild) Me() EndpointGuildMe {
	e2 := e.Member("@me").appendMajor("nick")
	return EndpointGuildMe{e2}
}

// TODO not sure if necessary, there is a EndpointGuildMember.Modify, not sure if it takes @me. @abalabahaha#9421 on discord said that perhaps its different because modifying my nick and managing others is a different permission but that does not really make sense to me. But whatever...
// TODO i also don't really like the name, doesn't fit in
func (e EndpointGuildMe) ModifyNick(ctx context.Context, nick string) (newNick string, err error) {
	nickStruct := struct {
		Nick string `json:"nick"`
	}{
		Nick: nick,
	}
	return nickStruct.Nick, e.doMethod(ctx, "PATCH", nickStruct, &nickStruct)
}

type EndpointGuildChannels struct {
	*endpoint
}

func (e EndpointGuild) Channels() EndpointGuildChannels {
	e2 := e.appendMajor("channels")
	return EndpointGuildChannels{e2}
}

func (e EndpointGuildChannels) Get(ctx context.Context) (channels []*ModelChannel, err error) {
	return channels, e.doMethod(ctx, "GET", nil, &channels)
}

func (e EndpointGuildChannels) Create(ctx context.Context, params *ParamsGuildChannelCreate) (ch *ModelChannel, err error) {
	return ch, e.doMethod(ctx, "POST", params, &ch)
}

type ParamsGuildChannelsModifyPositions struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

func (e EndpointGuildChannels) ModifyPositions(ctx context.Context, params []*ParamsGuildChannelsModifyPositions) (channels *ModelChannel, err error) {
	return channels, e.doMethod(ctx, "PATCH", params, &channels)
}

type EndpointGuildMembers struct {
	*endpoint
}

func (g EndpointGuild) Members() EndpointGuildMembers {
	e2 := g.appendMajor("members")
	return EndpointGuildMembers{e2}
}

// TODO necessary???
type ParamsGuildMembersGet struct {
	Limit   int
	AfterID string
}

func (params *ParamsGuildMembersGet) rawQuery() string {
	v := make(url.Values)
	if params.AfterID != "" {
		v.Set("after", params.AfterID)
	}
	if params.Limit > 0 {
		v.Set("limit", strconv.Itoa(params.Limit))
	}
	return v.Encode()
}

func (e EndpointGuildMembers) Get(ctx context.Context, params *ParamsGuildMembersGet) (guildMembers []*ModelGuildMember, err error) {
	req := e.newRequest(ctx, "GET", nil)
	if params != nil {
		req.URL.RawQuery = params.rawQuery()
	}
	return guildMembers, e.do(req, &guildMembers)
}

type EndpointGuildMember struct {
	*endpoint
}

func (e EndpointGuild) Member(uID string) EndpointGuildMember {
	e2 := e.Members().appendMinor(uID)
	return EndpointGuildMember{e2}
}

type ParamsGuildMemberAdd struct {
	AccessToken string       `json:"access_token"`
	Nick        string       `json:"nick,omitempty"`
	Roles       []*ModelRole `json:"roles,omitempty"`
	Mute        bool         `json:"mute,omitempty"`
	Deaf        bool         `json:"deaf,omitempty"`
}

func (e EndpointGuildMember) Add(ctx context.Context, params *ParamsGuildMemberAdd) (gm *ModelGuildMember, err error) {
	return gm, e.doMethod(ctx, "PUT", params, &gm)
}

func (e EndpointGuildMember) Get(ctx context.Context) (gm *ModelGuildMember, err error) {
	return gm, e.doMethod(ctx, "GET", nil, &gm)
}

// TODO rename this and all other params to postfix params
type ParamsGuildMemberModify struct {
	Nick      string       `json:"nick,omitempty"`
	Roles     []*ModelRole `json:"roles,omitempty"`
	Mute      *bool        `json:"mute,omitempty"` // pointer so that you can set false
	Deaf      *bool        `json:"deaf,omitempty"` // pointer so that you can set false
	ChannelID string       `json:"channel_id,omitempty"`
}

func (e EndpointGuildMember) Modify(ctx context.Context, params *ParamsGuildMemberModify) error {
	return e.doMethod(ctx, "PATCH", params, nil)
}

type EndpointGuildMemberRole struct {
	*endpoint
}

func (e EndpointGuildMember) Role(roleID string) EndpointGuildMemberRole {
	e2 := e.appendMajor("roles").appendMinor(roleID)
	return EndpointGuildMemberRole{e2}
}

func (e EndpointGuildMemberRole) Add(ctx context.Context) error {
	return e.doMethod(ctx, "PUT", nil, nil)
}

func (e EndpointGuildMemberRole) Remove(ctx context.Context) error {
	return e.doMethod(ctx, "DELETE", nil, nil)
}

func (e EndpointGuildMember) Remove(ctx context.Context) error {
	return e.doMethod(ctx, "DELETE", nil, nil)
}

type EndpointGuildBans struct {
	*endpoint
}

func (e EndpointGuild) Bans() EndpointGuildBans {
	e2 := e.appendMajor("bans")
	return EndpointGuildBans{e2}
}

func (e EndpointGuildBans) Get(ctx context.Context) (users []*ModelUser, err error) {
	return users, e.doMethod(ctx, "GET", nil, &users)
}

type EndpointGuildBan struct {
	*endpoint
}

func (e EndpointGuild) Ban(uID string) EndpointGuildBan {
	e2 := e.Bans().appendMinor(uID)
	return EndpointGuildBan{e2}
}

type ParamsGuildBanCreate struct {
	DeleteMessageDays int `json:"delete-message-days"`
}

func (e EndpointGuildBan) Create(ctx context.Context, params *ParamsGuildBanCreate) error {
	return e.doMethod(ctx, "PUT", params, nil)
}
func (e EndpointGuildBan) Remove(ctx context.Context) error {
	return e.doMethod(ctx, "DELETE", nil, nil)
}

// TODO maybe guildrole instead?
type EndpointRoles struct {
	*endpoint
}

func (e EndpointGuild) Roles() EndpointRoles {
	e2 := e.appendMajor("roles")
	return EndpointRoles{e2}
}

func (e EndpointRoles) Get(ctx context.Context) (roles []*ModelRole, err error) {
	return roles, e.doMethod(ctx, "GET", nil, &roles)
}

type ParamsRoleCreate struct {
	Name string `json:"name,omitempty"`
	// TODO should be null?
	Permissions int  `json:"permissions,omitempty"`
	Color       int  `json:"color,omitempty"`
	Hoist       bool `json:"hoist,omitempty"`
	Mentionable bool `json:"mentionable,omitempty"`
}

func (e EndpointRoles) Create(ctx context.Context, params *ParamsRoleCreate) (r *ModelRole, err error) {
	return r, e.doMethod(ctx, "GET", params, &r)
}

type ParamsRolesModifyPositions struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

func (e EndpointRoles) ModifyPositions(ctx context.Context, params *ParamsRolesModifyPositions) (roles []*ModelRole, err error) {
	return roles, e.doMethod(ctx, "PATCH", params, &roles)
}

type EndpointRole struct {
	*endpoint
}

func (e EndpointGuild) Role(roleID string) EndpointRole {
	e2 := e.Roles().appendMinor(roleID)
	return EndpointRole{e2}
}

// TODO nulls
type ParamsRoleModify struct {
	Name        string `json:"name,omitempty"`
	Permissions int    `json:"permissions,omitempty"`
	Color       int    `json:"color,omitempty"`
	Hoist       bool   `json:"hoist,omitempty"`
	Mentionable bool   `json:"mentionable,omitempty"`
}

func (e EndpointRole) Modify(ctx context.Context, params *ParamsRoleModify) (r *ModelRole, err error) {
	return r, e.doMethod(ctx, "PATCH", params, &r)
}

func (e EndpointRole) Delete(ctx context.Context) error {
	return e.doMethod(ctx, "DELETE", nil, nil)
}

// TODO i don't like the api because prune is a verb :(. same as bulk-delete ****
type EndpointPrune struct {
	*endpoint
}

func (e EndpointGuild) Prune() EndpointPrune {
	e2 := e.appendMajor("prune")
	return EndpointPrune{e2}
}

func (e EndpointPrune) GetCount(ctx context.Context, days int) (count int, err error) {
	req := e.newRequest(ctx, "GET", nil)
	if days > 0 {
		v := url.Values{}
		v.Set("days", strconv.Itoa(days))
		req.URL.RawQuery = v.Encode()
	}
	countStruct := struct {
		// TODO should I stick with discord's naming?
		Count int `json:"pruned"`
	}{}
	return countStruct.Count, e.do(req, &countStruct)
}

func (e EndpointPrune) Begin(ctx context.Context, days int) (pruned int, err error) {
	req := e.newRequest(ctx, "POST", nil)
	if days > 0 {
		v := url.Values{}
		v.Set("days", strconv.Itoa(days))
		req.URL.RawQuery = v.Encode()
	}
	prunedStruct := struct {
		Pruned int `json:"pruned"`
	}{}
	return prunedStruct.Pruned, e.do(req, &prunedStruct)
}

func (e EndpointGuild) VoiceRegions() EndpointVoiceRegions {
	e2 := e.appendMajor("regions")
	return EndpointVoiceRegions{e2}
}

func (e EndpointGuild) Invites() EndpointInvites {
	e2 := e.appendMajor("invites")
	return EndpointInvites{e2}
}

type EndpointIntegrations struct {
	*endpoint
}

func (e EndpointGuild) Integrations() EndpointIntegrations {
	e2 := e.appendMajor("integrations")
	return EndpointIntegrations{e2}
}

func (e EndpointIntegrations) Get(ctx context.Context) (integrations []*ModelIntegration, err error) {
	return integrations, e.doMethod(ctx, "GET", nil, &integrations)
}

type ParamsIntegrationsCreate struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (e EndpointIntegrations) Create(ctx context.Context, params *ParamsIntegrationsCreate) error {
	return e.doMethod(ctx, "POST", params, nil)
}

type EndpointIntegration struct {
	*endpoint
}

func (e EndpointGuild) Integration(integrationID string) EndpointIntegration {
	e2 := e.Integrations().appendMinor(integrationID)
	return EndpointIntegration{e2}
}

type ParamsIntegrationModify struct {
	// TODO impossible to not send or send 0 value :(
	ExpireBehaviour   int  `json:"expire_behaviour,omitempty"`
	ExpireGracePeriod int  `json:"expire_grace_period,omitempty"`
	EnableEmoticons   bool `json:"enable_emoticons,omitempty"`
}

func (e EndpointIntegration) Modify(ctx context.Context, params *ParamsIntegrationModify) error {
	return e.doMethod(ctx, "PATCH", params, nil)
}

func (e EndpointIntegration) Delete(ctx context.Context, gID, integrationID string) error {
	return e.doMethod(ctx, "DELETE", nil, nil)
}

func (e EndpointIntegration) Sync(ctx context.Context) error {
	e2 := e.appendMajor("sync")
	return e2.doMethod(ctx, "POST", nil, nil)
}

type EndpointGuildEmbed struct {
	*endpoint
}

func (e EndpointGuild) Embed() EndpointGuildEmbed {
	e2 := e.appendMajor("embed")
	return EndpointGuildEmbed{e2}
}

func (e EndpointGuildEmbed) Get(ctx context.Context) (ge *ModelGuildEmbed, err error) {
	return ge, e.doMethod(ctx, "GET", nil, &ge)
}

func (e EndpointGuildEmbed) Modify(ctx context.Context, ge *ModelGuildEmbed) (newGE *ModelGuildEmbed, err error) {
	return newGE, e.doMethod(ctx, "PATCH", ge, &newGE)
}
