package discgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

// TODO expose all events except those replaced by state events.

type EventReady struct {
	V               int                 `json:"v"`
	User            *ModelUser          `json:"user"`
	PrivateChannels []*ModelChannel     `json:"private_channels"`
	Guilds          []*eventGuildCreate `json:"guilds"`
	SessionID       string              `json:"session_id"`
	Trace           []string            `json:"_trace"`
}

type eventResumed struct {
	Trace []string `json:"_trace"`
}

type eventChannelCreate struct {
	ModelChannel
}

type eventChannelUpdate struct {
	ModelChannel
}

type eventChannelDelete struct {
	ModelChannel
}

type eventGuildCreate struct {
	ModelGuild
	Large       bool                `json:"large"`
	Unavailable bool                `json:"unavailable"`
	MemberCount int                 `json:"member_count"`
	VoiceStates []*ModelVoiceState  `json:"voice_states"` // without guild_id key
	Members     []*ModelGuildMember `json:"members"`
	Channels    []*ModelChannel     `json:"channels"`
	Presences   []*ModelPresence    `json:"presences"`
}

type ModelPresence struct {
	User   ModelUser  `json:"user"`
	Game   *ModelGame `json:"game"`
	Status string     `json:"status"`
}

type eventGuildUpdate struct {
	ModelGuild
}

type eventGuildDelete struct {
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

type eventGuildEmojisUpdate struct {
	GuildID string             `json:"guild_id"`
	Emojis  []*ModelGuildEmoji `json:"emojis"`
}

type EventGuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

type eventGuildMemberAdd struct {
	ModelGuildMember
	GuildID string `json:"guild_id"`
}

type eventGuildMemberRemove struct {
	User    ModelUser `json:"user"`
	GuildID string    `json:"guild_id"`
}

type eventGuildMemberUpdate struct {
	GuildID string    `json:"guild_id"`
	Roles   []string  `json:"roles"`
	User    ModelUser `json:"user"`
	Nick    string    `json:"nick"`
}

type eventGuildMembersChunk struct {
	GuildID string              `json:"guild_id"`
	Members []*ModelGuildMember `json:"members"`
}

type eventGuildRoleCreate struct {
	GuildID string    `json:"guild_id"`
	Role    ModelRole `json:"role"`
}

type eventGuildRoleUpdate struct {
	GuildID string    `json:"guild_id"`
	Role    ModelRole `json:"role"`
}

type eventGuildRoleDelete struct {
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
	// TODO why is there even a user here?
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
	ModelGameTypeGame = iota
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

type EventHandlerError struct {
	Err       error
	Event     interface{}
	EventName string
}

func (e *EventHandlerError) Error() string {
	eventJSON, err := json.MarshalIndent(e.Event, "", "    ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v handler error: %v\nevent: %v", e.EventName, e.Err, eventJSON)
}

var errUnknownEvent = errors.New("unknown event")

func getEventStruct(eventType string) (interface{}, error) {
	switch eventType {
	case "READY":
		return new(EventReady), nil
	case "RESUMED":
		return new(eventResumed), nil
	case "CHANNEL_CREATE":
		return new(eventChannelCreate), nil
	case "CHANNEL_UPDATE":
		return new(eventChannelUpdate), nil
	case "CHANNEL_DELETE":
		return new(eventChannelDelete), nil
	case "GUILD_CREATE":
		return new(eventGuildCreate), nil
	case "GUILD_UPDATE":
		return new(eventGuildUpdate), nil
	case "GUILD_DELETE":
		return new(eventGuildDelete), nil
	case "GUILD_BAN_ADD":
		return new(EventGuildBanAdd), nil
	case "GUILD_BAN_REMOVE":
		return new(EventGuildBanRemove), nil
	case "GUILD_EMOJIS_UPDATE":
		return new(eventGuildEmojisUpdate), nil
	case "GUILD_INTEGRATIONS_UPDATE":
		return new(EventGuildIntegrationsUpdate), nil
	case "GUILD_MEMBER_ADD":
		return new(eventGuildMemberAdd), nil
	case "GUILD_MEMBER_REMOVE":
		return new(eventGuildMemberRemove), nil
	case "GUILD_MEMBER_UPDATE":
		return new(eventGuildMemberUpdate), nil
	case "GUILD_MEMBERS_CHUNK":
		return new(eventGuildMembersChunk), nil
	case "GUILD_ROLE_CREATE":
		return new(eventGuildRoleCreate), nil
	case "GUILD_ROLE_UPDATE":
		return new(eventGuildRoleUpdate), nil
	case "GUILD_ROLE_DELETE":
		return new(eventGuildRoleDelete), nil
	case "MESSAGE_CREATE":
		return new(EventMessageCreate), nil
	case "MESSAGE_UPDATE":
		return new(EventMessageUpdate), nil
	case "MESSAGE_DELETE":
		return new(EventMessageDelete), nil
	case "MESSAGE_DELETE_BULK":
		return new(EventMessageDeleteBulk), nil
	case "MESSAGE_REACTION_ADD":
		return new(EventMessageReactionAdd), nil
	case "MESSAGE_REACTION_REMOVE":
		return new(EventMessageReactionRemove), nil
	case "MESSAGE_REACTION_REMOVE_ALL":
		return new(EventMessageReactionRemoveAll), nil
	case "PRESENCE_UPDATE":
		return new(EventPresenceUpdate), nil
	case "TYPING_START":
		return new(EventTypingStart), nil
	case "USER_UPDATE":
		return new(EventUserUpdate), nil
	case "VOICE_STATE_UPDATE":
		return new(EventVoiceStateUpdate), nil
	case "VOICE_SERVER_UPDATE":
		return new(eventVoiceServerUpdate), nil
	}
	return nil, errUnknownEvent
}
