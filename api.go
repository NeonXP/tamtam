/*
 * TamTam Bot API
 */
package tamtam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Api struct {
	key string
	url *url.URL
}

// New TamTam Api object
func New(key string) *Api {
	u, _ := url.Parse("https://botapi.tamtam.chat/")
	return &Api{
		key: key,
		url: u,
	}
}

// region Misc methods

func (a *Api) GetMe() (*UserWithPhoto, error) {
	result := new(UserWithPhoto)
	values := url.Values{}
	body, err := a.request(http.MethodGet, "me", values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) GetUploadURL(uploadType UploadType) (*UploadEndpoint, error) {
	result := new(UploadEndpoint)
	values := url.Values{}
	values.Set("type", string(uploadType))
	body, err := a.request(http.MethodPost, "uploads", values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) GetHandler(updates chan interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		u := new(Update)
		b, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(b, u)
		switch u.UpdateType {
		case UpdateTypeMessageCallback:
			upd := new(UpdateMessageCallback)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeMessageCreated:
			upd := new(UpdateMessageCreated)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeMessageRemoved:
			upd := new(UpdateMessageRemoved)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeMessageEdited:
			upd := new(UpdateMessageEdited)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeMessageRestored:
			upd := new(UpdateMessageRestored)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeBotAdded:
			upd := new(UpdateBotAdded)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeBotRemoved:
			upd := new(UpdateBotRemoved)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeUserAdded:
			upd := new(UpdateUserAdded)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeUserRemoved:
			upd := new(UpdateUserRemoved)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeBotStarted:
			upd := new(UpdateBotStarted)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		case UpdateTypeChatTitleChanged:
			upd := new(UpdateChatTitleChanged)
			_ = json.Unmarshal(b, upd)
			updates <- *upd
		}
	}
}

// endregion

// region Chat methods
func (a *Api) GetChats(count, marker int64) (*ChatList, error) {
	result := new(ChatList)
	values := url.Values{}
	if count > 0 {
		values.Set("count", strconv.Itoa(int(count)))
	}
	if marker > 0 {
		values.Set("marker", strconv.Itoa(int(marker)))
	}
	body, err := a.request(http.MethodGet, "chats", values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) GetChat(chatID int64) (*Chat, error) {
	result := new(Chat)
	values := url.Values{}
	body, err := a.request(http.MethodGet, fmt.Sprintf("chats/%d", chatID), values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) GetChatMembership(chatID int64) (*ChatMember, error) {
	result := new(ChatMember)
	values := url.Values{}
	body, err := a.request(http.MethodGet, fmt.Sprintf("chats/%d/members/me", chatID), values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) GetChatMembers(chatID, count, marker int64) (*ChatMembersList, error) {
	result := new(ChatMembersList)
	values := url.Values{}
	if count > 0 {
		values.Set("count", strconv.Itoa(int(count)))
	}
	if marker > 0 {
		values.Set("marker", strconv.Itoa(int(marker)))
	}
	body, err := a.request(http.MethodGet, fmt.Sprintf("chats/%d/members", chatID), values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) LeaveChat(chatID int64) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	body, err := a.request(http.MethodDelete, fmt.Sprintf("chats/%d/members/me", chatID), values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) EditChat(chatID int64, update *ChatPatch) (*Chat, error) {
	result := new(Chat)
	values := url.Values{}
	body, err := a.request(http.MethodPatch, fmt.Sprintf("chats/%d", chatID), values, update)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) AddMember(chatID int64, users UserIdsList) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	body, err := a.request(http.MethodPost, fmt.Sprintf("chats/%d/members", chatID), values, users)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) RemoveMember(chatID int64, userID int64) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	values.Set("user_id", strconv.Itoa(int(userID)))
	body, err := a.request(http.MethodDelete, fmt.Sprintf("chats/%d/members", chatID), values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) SendAction(chatID int64, action SenderAction) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	body, err := a.request(http.MethodPost, fmt.Sprintf("chats/%d/actions", chatID), values, ActionRequestBody{Action: action})
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

// endregion

// region Message methods

func (a *Api) GetMessages(chatID int64, messageIDs []string, from int64, to int64, count int64) (*MessageList, error) {
	result := new(MessageList)
	values := url.Values{}
	if chatID > 0 {
		values.Set("chat_id", strconv.Itoa(int(chatID)))
	}
	if len(messageIDs) > 0 {
		for _, mid := range messageIDs {
			values.Add("message_ids", mid)
		}
	}
	if from > 0 {
		values.Set("from", strconv.Itoa(int(from)))
	}
	if to > 0 {
		values.Set("count", strconv.Itoa(int(to)))
	}
	if count > 0 {
		values.Set("count", strconv.Itoa(int(count)))
	}
	body, err := a.request(http.MethodGet, "messages", values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) SendMessage(chatID int64, userID int64, message *NewMessageBody) (*Message, error) {
	result := new(Message)
	values := url.Values{}
	if chatID > 0 {
		values.Set("chat_id", strconv.Itoa(int(chatID)))
	}
	if userID > 0 {
		values.Set("user_id", strconv.Itoa(int(userID)))
	}
	body, err := a.request(http.MethodPost, "messages", values, message)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) EditMessage(messageID int64, message *NewMessageBody) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	values.Set("message_id", strconv.Itoa(int(messageID)))
	body, err := a.request(http.MethodPut, "messages", values, message)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) DeleteMessage(messageID int64) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	values.Set("message_id", strconv.Itoa(int(messageID)))
	body, err := a.request(http.MethodDelete, "messages", values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) AnswerOnCallback(callbackID int64, callback *CallbackAnswer) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	values.Set("callback_id", strconv.Itoa(int(callbackID)))
	body, err := a.request(http.MethodPost, "answers", values, callback)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

// endregion

// region Subscriptions

func (a *Api) GetSubscriptions() (*GetSubscriptionsResult, error) {
	result := new(GetSubscriptionsResult)
	values := url.Values{}
	body, err := a.request(http.MethodGet, "subscriptions", values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) Subscribe(subscription *SubscriptionRequestBody) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	body, err := a.request(http.MethodPost, "subscriptions", values, subscription)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) Unsubscribe(subscriptionURL string) (*SimpleQueryResult, error) {
	result := new(SimpleQueryResult)
	values := url.Values{}
	values.Set("url", subscriptionURL)
	body, err := a.request(http.MethodDelete, "subscriptions", values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *Api) GetUpdates(limit int, timeout int, marker int64, types []string) (*UpdateList, error) {
	result := new(UpdateList)
	values := url.Values{}
	if limit > 0 {
		values.Set("limit", strconv.Itoa(limit))
	}
	if timeout > 0 {
		values.Set("timeout", strconv.Itoa(timeout))
	}
	if marker > 0 {
		values.Set("marker", strconv.Itoa(int(marker)))
	}
	if len(types) > 0 {
		for _, t := range types {
			values.Add("types", t)
		}
	}
	body, err := a.request(http.MethodDelete, "updates", values, nil)
	if err != nil {
		return result, err
	}
	defer body.Close()
	return result, json.NewDecoder(body).Decode(result)
}

// endregion

// region Internal

func (a *Api) request(method, path string, query url.Values, body interface{}) (io.ReadCloser, error) {
	c := http.DefaultClient
	u := *a.url
	u.Path = path
	query.Set("access_token", a.key)
	u.RawQuery = query.Encode()
	j, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(j))
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if resp.StatusCode != http.StatusOK {
		errObj := new(Error)
		err = json.NewDecoder(resp.Body).Decode(errObj)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("code=%s message=%s error=%s", errObj.Code, errObj.Message, errObj.Error)
	}
	return resp.Body, err
}

// endregion