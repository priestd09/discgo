package discgo

import (
	"time"

	"encoding/json"

	"path"

	"net/url"
	"strconv"

	"bytes"
	"io"
	"mime/multipart"
)

// Channel represents a channel in Discord.
// If Recipient is set this is a DM Channel otherwise this is a Guild Channel.
// Guild Channel represents an isolated set of users and messages within a Guild.
// DM Channel represent a one-to-one conversation between two users, outside of the
// scope of guilds.
type Channel struct {
	ID                   string       `json:"id"`
	GuildID              string       `json:"guild_id"`
	Name                 string       `json:"name"`
	Type                 string       `json:"type"`
	Position             int          `json:"position"`
	IsPrivate            bool         `json:"is_private"`
	PermissionOverwrites []*Overwrite `json:"permission_overwrites"`
	Topic                string       `json:"topic"`
	Recipient            *User        `json:"recipient"`
	LastMessageID        string       `json:"last_message_id"`
	Bitrate              int          `json:"bitrate"`
	UserLimit            int          `json:"user_limit"`
}

// Message represents a message sent in a channel within Discord.
// The author object follows the structure of the user object, but
// is only a valid user in the case where the message is generated
// by a user or bot user. If the message is generated by a webhook,
// the author object corresponds to the webhook's id, username, and avatar.
// You can tell if a message is generated by a webhook by checking for the
// webhook_id on the message object.
type Message struct {
	ID              string        `json:"id"`
	ChannelID       string        `json:"channel_id"`
	Author          *User         `json:"author"`
	Content         string        `json:"content"`
	Timestamp       *time.Time    `json:"timestamp"`
	EditedTimestamp *time.Time    `json:"edited_timestamp"`
	TTS             bool          `json:"tts"`
	MentionEveryone bool          `json:"mention_everyone"`
	Mentions        []*User       `json:"mentions"`
	MentionRoles    []string      `json:"mention_roles"`
	Attachments     []*Attachment `json:"attachments"`
	Embeds          []*Embed      `json:"embeds"`
	Reactions       []*Reaction   `json:"reactions"`
	Nonce           string        `json:"nonce"`
	Pinned          bool          `json:"pinned"`
	WebhookID       string        `json:"webhook_id"`
}

type Reaction struct {
	Count int
	Me    bool
	Emoji *ReactionEmoji
}

type ReactionEmoji struct {
	ID   *string // nullable
	Name string
}

type Overwrite struct {
	ID    string
	Type  string
	Allow int
	Deny  int
}

type Embed struct {
	Title       string          `json:"title,omitempty"`
	Type        string          `json:"type,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Timestamp   *time.Time      `json:"timestamp,omitempty"`
	Color       int             `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []*EmbedField   `json:"fields,omitempty"`
}

type EmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Attachment struct {
	ID       string `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	Size     int    `json:"size,omitempty"`
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// TODO create mention method on User, Channel, Role and Custom ReactionEmoji structs
// https://discordapp.com/developers/docs/resources/channel#message-formatting

func (c *Client) GetChannel(cID string) (ch *Channel, err error) {
	endpoint := path.Join("channels", cID)
	req := c.newRequest("GET", endpoint, nil)
	body, err := c.do(req, endpoint, 0)
	if err != nil {
		return nil, err
	}
	return ch, json.Unmarshal(body, &ch)
}

type ParamsModifyChannel struct {
	Name      string `json:"name,omitempty"`
	Position  int    `json:"position,omitempty"`
	Topic     string `json:"topic,omitempty"`
	Bitrate   int    `json:"bitrate,omitempty"`
	UserLimit int    `json:"user_limit,omitempty"`
}

func (c *Client) ModifyChannel(cID string, params *ParamsModifyChannel) error {
	endpoint := path.Join("channels", cID)
	req := c.newRequestJSON("PATCH", endpoint, params)
	_, err := c.do(req, endpoint, 0)
	return err
}

func (c *Client) DeleteChannel(cID string) (ch *Channel, err error) {
	endpoint := path.Join("channels", cID)
	req := c.newRequest("DELETE", endpoint, nil)
	body, err := c.do(req, endpoint, 0)
	if err != nil {
		return nil, err
	}
	return ch, json.Unmarshal(body, &ch)
}

type ParamsGetMessages struct {
	AroundID string
	BeforeID string
	AfterID  string
	Limit    int
}

func (params *ParamsGetMessages) rawQuery() string {
	v := make(url.Values)
	if params.AroundID != "" {
		v.Set("around", params.AroundID)
	}
	if params.BeforeID != "" {
		v.Set("before", params.BeforeID)
	}
	if params.AfterID != "" {
		v.Set("after", params.AfterID)
	}
	if params.Limit > 0 {
		v.Set("limit", strconv.Itoa(params.Limit))
	}
	return v.Encode()
}

func (c *Client) GetMessages(cID string, params *ParamsGetMessages) (msgs []*Message, err error) {
	endpoint := path.Join("channels", cID, "messages")
	req := c.newRequest("GET", endpoint, nil)
	if params != nil {
		req.URL.RawQuery = params.rawQuery()
	}
	body, err := c.do(req, endpoint, 0)
	if err != nil {
		return nil, err
	}
	return msgs, json.Unmarshal(body, &msgs)
}

func (c *Client) GetMessage(cID, mID string) (m *Message, err error) {
	endpoint := path.Join("channels", cID, "messages", mID)
	req := c.newRequest("GET", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "messages", "*")
	body, err := c.do(req, rateLimitPath, 0)
	if err != nil {
		return nil, err
	}
	return m, json.Unmarshal(body, &m)
}

type ParamsCreateMessage struct {
	Content string `json:"content,omitempty"`
	Nonce   string `json:"nonce,omitempty"`
	TTS     bool   `json:"tts,omitempty"`
	File    *File  `json:"-"`
	Embed   *Embed `json:"embed,omitempty"`
}

type File struct {
	Name    string
	Content io.Reader
}

func (c *Client) CreateMessage(cID string, params *ParamsCreateMessage) (m *Message, err error) {
	reqBody := &bytes.Buffer{}
	reqBodyWriter := multipart.NewWriter(reqBody)

	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	w, err := reqBodyWriter.CreateFormField("payload_json")
	if err != nil {
		return nil, err
	}
	_, err = w.Write(payloadJSON)
	if err != nil {
		return nil, err
	}

	if params.File != nil {
		w, err := reqBodyWriter.CreateFormFile("file", params.File.Name)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(w, params.File.Content)
		if err != nil {
			return nil, err
		}
	}

	err = reqBodyWriter.Close()
	if err != nil {
		return nil, err
	}

	endpoint := path.Join("channels", cID, "messages")
	req := c.newRequest("POST", endpoint, reqBody)
	req.Header.Set("Content-Type", reqBodyWriter.FormDataContentType())

	body, err := c.do(req, endpoint, 0)
	if err != nil {
		return nil, err
	}
	return m, json.Unmarshal(body, &m)
}

func (c *Client) CreateReaction(cID, mID, emoji string) error {
	endpoint := path.Join("channels", cID, "messages", mID, "reactions", emoji, "@me")
	req := c.newRequest("PUT", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "messages", "*", "reactions", "*", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

// uID = "@me" for your own reaction
func (c *Client) DeleteReaction(cID, mID, emoji, uID string) error {
	endpoint := path.Join("channels", cID, "messages", mID, "reactions", emoji, uID)
	req := c.newRequest("DELETE", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "messages", "*", "reactions", "*", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

func (c *Client) GetReactions(cID, mID, emoji string) (users []*User, err error) {
	endpoint := path.Join("channels", cID, "messages", mID, "reactions", emoji)
	req := c.newRequest("GET", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "messages", "*", "reactions", "*")
	body, err := c.do(req, rateLimitPath, 0)
	if err != nil {
		return nil, err
	}
	return users, json.Unmarshal(body, &users)
}

func (c *Client) DeleteReactions(cID, mID string) error {
	endpoint := path.Join("channels", cID, "messages", mID, "reactions")
	req := c.newRequest("DELETE", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "messages", "*", "reactions")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

type ParamsEditMessage struct {
	Content string `json:"content,omitempty"`
	Embed   *Embed `json:"embed,omitempty"`
}

func (c *Client) EditMessage(cID, mID string, params *ParamsEditMessage) (m *Message, err error) {
	endpoint := path.Join("channels", cID, "messages", mID)
	req := c.newRequestJSON("PATCH", endpoint, params)
	rateLimitPath := path.Join("channels", cID, "messages", "*")
	body, err := c.do(req, rateLimitPath, 0)
	if err != nil {
		return nil, err
	}
	return m, json.Unmarshal(body, &m)
}

func (c *Client) DeleteMessage(cID, mID string) error {
	endpoint := path.Join("channels", cID, "messages", mID)
	req := c.newRequest("DELETE", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "messages", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

func (c *Client) BulkDeleteMessages(cID string, mIDs []string) error {
	endpoint := path.Join("channels", cID, "messages", "bulk-delete")
	req := c.newRequestJSON("POST", endpoint, mIDs)
	_, err := c.do(req, endpoint, 0)
	return err
}

type ParamsEditPermissions struct {
	Allow int    `json:"allow"`
	Deny  int    `json:"deny"`
	Type  string `json:"type"`
}

func (c *Client) EditPermissions(cID, overwriteID string, params *ParamsEditPermissions) error {
	endpoint := path.Join("channels", cID, "permissions", overwriteID)
	req := c.newRequestJSON("PUT", endpoint, params)
	rateLimitPath := path.Join("channels", cID, "permissions", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

func (c *Client) GetInvites(cID string) (invites []*Invite, err error) {
	endpoint := path.Join("channels", cID, "invites")
	req := c.newRequest("GET", endpoint, nil)
	body, err := c.do(req, endpoint, 0)
	if err != nil {
		return nil, err
	}
	return invites, json.Unmarshal(body, &invites)
}

type ParamsCreateInvite struct {
	MaxAge    int `json:"max_age,omitempty"`
	MaxUses   int `json:"max_uses,omitempty"`
	Temporary bool `json:"temporary,omitempty"`
	Unique    bool `json:"unique,omitempty"`
}

func (c *Client) CreateInvite(cID string, params *ParamsCreateInvite) (invite *Invite, err error) {
	endpoint := path.Join("channels", cID, "invites")
	req := c.newRequestJSON("POST", cID, params)
	body, err := c.do(req, endpoint, 0)
	if err != nil {
		return nil, err
	}
	return invite, json.Unmarshal(body, &invite)
}

func (c *Client) DeletePermission(cID, overwriteID string) error {
	endpoint := path.Join("channels", cID, "overwrites", overwriteID)
	req := c.newRequest("DELETE", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "overwrites", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

func (c *Client) TriggerTypingIndicator(cID string) error {
	endpoint := path.Join("channels", cID, "typing")
	req := c.newRequest("POST", endpoint, nil)
	_, err := c.do(req, endpoint, 0)
	return err
}

func (c *Client) GetPinnedMessages(cID string) (msgs []*Message, err error) {
	endpoint := path.Join("channels", cID, "pins")
	req := c.newRequest("GET", endpoint, nil)
	body, err := c.do(req, endpoint, 0)
	if err != nil {
		return nil, err
	}
	return msgs, json.Unmarshal(body, &msgs)
}

func (c *Client) PinMessage(cID, mID string) error {
	endpoint := path.Join("channels", cID, "pins", mID)
	req := c.newRequest("PUT", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "pins", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

func (c *Client) DeletePinnedMessage(cID, mID string) error {
	endpoint := path.Join("channels", cID, "pins", mID)
	req := c.newRequest("DELETE", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "pins", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

func (c *Client) AddRecipient(cID, uID string) error {
	endpoint := path.Join("channels", cID, "recipients", uID)
	req := c.newRequest("PUT", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "recipients", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}

func (c *Client) RemoveRecipient(cID, uID string) error {
	endpoint := path.Join("channels", cID, "recipients", uID)
	req := c.newRequest("DELETE", endpoint, nil)
	rateLimitPath := path.Join("channels", cID, "recipients", "*")
	_, err := c.do(req, rateLimitPath, 0)
	return err
}