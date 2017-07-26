package discgo

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
)

type eventReady struct {
	V               int             `json:"v"`
	User            *ModelUser      `json:"user"`
	PrivateChannels []*ModelChannel `json:"private_channels"`
	Guilds          []*ModelGuild   `json:"guilds"`
	SessionID       string          `json:"session_id"`
	Trace           []string        `json:"_trace"`
}

type eventResumed struct {
	Trace []string `json:"_trace"`
}

type EventChannelCreate struct {
	ModelChannel
}

type EventChannelUpdate struct {
	ModelChannel
}

type EventChannelDelete struct {
	ModelChannel
}

type EventGuildCreate struct {
	ModelGuild
}

type EventGuildUpdate struct {
	ModelGuild
}

type EventGuildDelete struct {
	ID          string `json:"id"`
	Unavailable bool   `json:"unavailable"`
}

type EventGuildBanAdd struct {
	ModelUser
	GuildID string `json:"guild_id"`
}

type EventGuildBanRemove struct {
	ModelUser
	GuildID string `json:"guild_id"`
}

type EventGuildEmojisUpdate struct {
	GuildID string             `json:"guild_id"`
	Emojis  []*ModelGuildEmoji `json:"emojis"`
}

type EventGuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

type EventGuildMemberAdd struct {
	GuildID string `json:"guild_id"`
}

type EventGuildMemberRemove struct {
	User    *ModelUser `json:"user"`
	GuildID string     `json:"guild_id"`
}

type EventGuildMemberUpdate struct {
	GuildID string    `json:"guild_id"`
	Roles   []string  `json:"roles"`
	User    ModelUser `json:"user"`
	Nick    string    `json:"nick"`
}

type EventGuildMembersChunk struct {
	GuildID string              `json:"guild_id"`
	Members []*ModelGuildMember `json:"members"`
}

type EventGuildRoleCreate struct {
	GuildID string    `json:"guild_id"`
	Role    ModelRole `json:"role"`
}

type EventGuildRoleUpdate struct {
	GuildID string    `json:"guild_id"`
	Role    ModelRole `json:"role"`
}

type EventGuildRoleDelete struct {
	GuildID string    `json:"guild_id"`
	Role    ModelRole `json:"role"`
}

type EventMessageCreate struct {
	ModelMessage
}

// May not be full message.
type EventMessageUpdate struct {
	ModelMessage
}

type EventMessageDelete struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
}

type EventMessageDeleteBulk struct {
	IDs       []string `json:"ids"`
	ChannelID string   `json:"channel_id"`
}

type EventMessageReactionAdd struct {
	UserID    string          `json:"user_id"`
	ChannelID string          `json:"channel_id"`
	MessageID string          `json:"message_id"`
	Emoji     ModelGuildEmoji `json:"emoji"`
}

type EventMessageReactionRemove struct {
	UserID    string           `json:"user_id"`
	ChannelID string           `json:"channel_id"`
	MessageID string           `json:"message_id"`
	Emoji     ModelGuildMember `json:"emoji"`
}

type EventMessageReactionRemoveAll struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
}

type EventPresenceUpdate struct {
	User    ModelUser  `json:"user"`
	Roles   []string   `json:"roles"`
	Game    *ModelGame `json:"game"`
	GuildID string     `json:"guild_id"`
	Status  string     `json:"status"`
}

const (
	StatusIdle    = "idle"
	StatusDND     = "dnd"
	StatusOnline  = "online"
	StatusOffline = "offline"
)

type ModelGame struct {
	Name string  `json:"name"`
	Type *int    `json:"type"`
	URL  *string `json:"url"`
}

const (
	// Yes this is actually what Discord calls it.
	ModelGameTypeGame      = iota
	ModelGameTypeStreaming
)

type EventTypingStart struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	Timestamp int    `json:"timestamp"`
}

type EventUserUpdate struct {
	ModelUser
}

type EventVoiceStateUpdate struct {
	ModelVoiceState
}

type eventVoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

type eventMux map[string]interface{}

func newEventMux() eventMux {
	return make(eventMux)
}

func (em eventMux) Register(fn interface{}) {
	switch fn.(type) {
	case func(ctx context.Context, conn *Conn, e *eventReady):
		em["READY"] = fn
	case func(ctx context.Context, conn *Conn, e *eventResumed):
		em["RESUMED"] = fn
	case func(ctx context.Context, conn *Conn, e *EventChannelCreate):
		em["CHANNEL_CREATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventChannelUpdate):
		em["CHANNEL_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventChannelDelete):
		em["CHANNEL_DELETE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildCreate):
		em["GUILD_CREATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildUpdate):
		em["GUILD_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildDelete):
		em["GUILD_DELETE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildBanAdd):
		em["GUILD_BAN_ADD"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildBanRemove):
		em["GUILD_BAN_REMOVE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildEmojisUpdate):
		em["GUILD_EMOJIS_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildIntegrationsUpdate):
		em["GUILD_INTEGRATIONS_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildMemberAdd):
		em["GUILD_MEMBER_ADD"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildMemberRemove):
		em["GUILD_MEMBER_REMOVE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildMemberUpdate):
		em["GUILD_MEMBER_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildMembersChunk):
		em["GUILD_MEMBERS_CHUNK"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildRoleCreate):
		em["GUILD_ROLE_CREATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildRoleUpdate):
		em["GUILD_ROLE_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventGuildRoleDelete):
		em["GUILD_ROLE_DELETE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventMessageCreate):
		em["MESSAGE_CREATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventMessageUpdate):
		em["MESSAGE_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventMessageDelete):
		em["MESSAGE_DELETE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventMessageDeleteBulk):
		em["MESSAGE_DELETE_BULK"] = fn
	case func(ctx context.Context, conn *Conn, e *EventMessageReactionAdd):
		em["MESSAGE_REACTION_ADD"] = fn
	case func(ctx context.Context, conn *Conn, e *EventMessageReactionRemove):
		em["MESSAGE_REACTION_REMOVE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventMessageReactionRemoveAll):
		em["MESSAGE_REACTION_REMOVE_ALL"] = fn
	case func(ctx context.Context, conn *Conn, e *EventPresenceUpdate):
		em["PRESENCE_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventTypingStart):
		em["TYPING_START"] = fn
	case func(ctx context.Context, conn *Conn, e *EventUserUpdate):
		em["USER_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *EventVoiceStateUpdate):
		em["VOICE_STATE_UPDATE"] = fn
	case func(ctx context.Context, conn *Conn, e *eventVoiceServerUpdate):
		em["VOICE_SERVER_UPDATE"] = fn
	default:
		panic("unknown event handler signature")
	}
}

func (em eventMux) route(ctx context.Context, conn *Conn, p *receivedPayload, sync bool) error {
	h, ok := em[p.Type]
	if !ok {
		// Discord better not be sending unknown events.
		return nil
	}
	e := reflect.New(reflect.TypeOf(h).In(2).Elem())
	err := json.Unmarshal(p.Data, e.Interface())
	if err != nil {
		return err
	}
	args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(conn), e}
	fn := func() {
		reflect.ValueOf(h).Call(args)
	}
	if sync {
		fn()
	} else {
		conn.runWorker(fn)
	}
	return nil
}

type State struct {
	sync.RWMutex

	user       *ModelUser
	dmChannels map[string]*ModelChannel

	guildMap   map[string]*ModelGuild
	channelMap map[string]*ModelChannel
}

func (s *State) registerEventHandlers(em *eventMux) {
	em.Register(s.ready)
	em.Register(s.createChannel)
}

func (s *State) ready(ctx context.Context, conn *Conn, e *eventReady) {
	s.Lock()
	defer s.Unlock()
	s.user = e.User
	for _, c := range e.PrivateChannels {
		s.dmChannels[c.ID] = c
	}
}

func (s *State) insertChannel(c *ModelChannel) {
	s.Lock()
	defer s.Unlock()
	if c.Type == ModelChannelTypeDM || c.Type == ModelChannelTypeGroupDM {
		s.dmChannels[c.ID] = c
	} else {
		g, ok := s.guildMap[*c.GuildID]
		if !ok {
			// TODO on panics, print out the associated event.
			panic("a channel created for an unknown guild")
		}
		*g.Channels = append(*g.Channels, c)
	}
}

func (s *State) createChannel(ctx context.Context, conn *Conn, e *EventChannelCreate) {
	s.insertChannel(&e.ModelChannel)
}

func (s *State) updateChannel(ctx context.Context, conn *Conn, e *EventChannelUpdate) {
	s.insertChannel(&e.ModelChannel)
}

func (s *State) deleteChannel(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
	if e.Type == ModelChannelTypeDM || e.Type == ModelChannelTypeGroupDM {
		s.dmChannels = append(s.dmChannels, c)
	} else {
		g, ok := s.guildMap[*c.GuildID]
		if !ok {
			// TODO on panics, print out the associated event.
			panic("a channel created for an unknown guild")
		}
		*g.Channels = append(*g.Channels, c)
	}
}

func (s *State) createGuild(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updateGuild(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) deleteGuild(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) addGuildBan(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updateGuildEmojis(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) addGuildMember(ctx context.Context, conn *Conn, e *EventGuildMemberAdd) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) removeGuildMember(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updateGuildMember(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) chunkGuildMembers(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) createGuildRole(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updateGuildRole(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) deleteGuildRole(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) createMessage(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updateMessage(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) deleteMessage(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) bulkDeleteMessages(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) addMessageReaction(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) removeMessageReaction(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updatePresence(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) startTyping(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updateUser(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updateVoiceState(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}

func (s *State) updateVoiceServer(ctx context.Context, conn *Conn, e *EventChannelDelete) {
	s.Lock()
	defer s.Unlock()
}
