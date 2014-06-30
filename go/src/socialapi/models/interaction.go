package models

import (
	"errors"
	"fmt"
	"socialapi/request"
	"time"

	"github.com/koding/bongo"
)

type Interaction struct {
	// unique identifier of the Interaction
	Id int64 `json:"id,string"`

	// Id of the interacted message
	MessageId int64 `json:"messageId,string"      sql:"NOT NULL"`

	// Id of the actor
	AccountId int64 `json:"accountId,string"      sql:"NOT NULL"`

	// holds troll, unsafe, etc
	MetaBits MetaBits `json:"metaBits"`

	// Type of the interaction
	TypeConstant string `json:"typeConstant"      sql:"NOT NULL;TYPE:VARCHAR(100);"`

	// Creation of the interaction
	CreatedAt time.Time `json:"createdAt"         sql:"NOT NULL"`
}

var AllowedInteractions = map[string]struct{}{
	"like":     {},
	"upvote":   {},
	"downvote": {},
}

const (
	Interaction_TYPE_LIKE     = "like"
	Interaction_TYPE_UPVOTE   = "upvote"
	Interaction_TYPE_DONWVOTE = "downvote"
)

func (i *Interaction) BeforeCreate() error {
	return i.MarkIfExempt()
}

func (i *Interaction) BeforeUpdate() error {
	return i.MarkIfExempt()
}

func (i *Interaction) AfterCreate() {
	bongo.B.AfterCreate(i)
}

func (i *Interaction) AfterUpdate() {
	bongo.B.AfterUpdate(i)
}

func (i Interaction) AfterDelete() {
	bongo.B.AfterDelete(i)
}

func (i Interaction) GetId() int64 {
	return i.Id
}

func (i Interaction) TableName() string {
	return "api.interaction"
}

func NewInteraction() *Interaction {
	return &Interaction{}
}

func (i *Interaction) One(q *bongo.Query) error {
	return bongo.B.One(i, i, q)
}

func (i *Interaction) ById(id int64) error {
	return bongo.B.ById(i, id)
}

func (i *Interaction) Create() error {
	return bongo.B.Create(i)
}

func (i *Interaction) Update() error {
	return bongo.B.Update(i)
}

func (i *Interaction) CreateRaw() error {
	insertSql := "INSERT INTO " +
		i.TableName() +
		` ("message_id","account_id","type_constant","created_at") VALUES ($1,$2,$3,$4) ` +
		"RETURNING ID"

	return bongo.B.DB.CommonDB().
		QueryRow(insertSql, i.MessageId, i.AccountId, i.TypeConstant, i.CreatedAt).
		Scan(&i.Id)
}

func (i *Interaction) MarkIfExempt() error {
	isExempt, err := i.isExempt()
	if err != nil {
		return err
	}

	if isExempt {
		i.MetaBits.Mark(Troll)
	}

	return nil
}

func (i *Interaction) isExempt() (bool, error) {
	if i.MetaBits.Is(Troll) {
		return true, nil
	}

	accountId, err := i.getAccountId()
	if err != nil {
		return false, err
	}

	account, err := ResetAccountCache(accountId)
	if err != nil {
		return false, err
	}

	if account == nil {
		return false, fmt.Errorf("account is nil, accountId:%d", i.AccountId)
	}

	if account.IsTroll {
		return true, nil
	}

	return false, nil
}

func (i *Interaction) getAccountId() (int64, error) {
	if i.AccountId != 0 {
		return i.AccountId, nil
	}

	if i.Id == 0 {
		return 0, fmt.Errorf("couldnt find accountId from content %+v", i)
	}

	ii := NewInteraction()
	if err := ii.ById(i.Id); err != nil {
		return 0, err
	}

	return ii.AccountId, nil
}

func (i *Interaction) Some(data interface{}, q *bongo.Query) error {
	return bongo.B.Some(i, data, q)
}

func (i *Interaction) Delete() error {
	selector := map[string]interface{}{
		"message_id": i.MessageId,
		"account_id": i.AccountId,
	}

	if err := i.One(bongo.NewQS(selector)); err != nil {
		return err
	}

	if err := bongo.B.Delete(i); err != nil {
		return err
	}

	return nil
}

func (i *Interaction) List(query *request.Query) ([]int64, error) {
	var interactions []int64

	if i.MessageId == 0 {
		return interactions, errors.New("Message is not set")
	}

	return i.FetchInteractorIds(query)
}

func (i *Interaction) FetchInteractorIds(query *request.Query) ([]int64, error) {
	interactorIds := make([]int64, 0)
	q := &bongo.Query{
		Selector: map[string]interface{}{
			"message_id":    i.MessageId,
			"type_constant": query.Type,
		},
		Pagination: *bongo.NewPagination(query.Limit, query.Skip),
		Pluck:      "account_id",
		Sort: map[string]string{
			"created_at": "desc",
		},
	}

	q.AddScope(RemoveTrollContent(i, query.ShowExempt))

	if err := i.Some(&interactorIds, q); err != nil {
		// TODO log this error
		return make([]int64, 0), nil
	}

	return interactorIds, nil
}

func (c *Interaction) Count(q *request.Query) (int, error) {
	if c.MessageId == 0 {
		return 0, errors.New("messageId is not set")
	}

	if q.Type == "" {
		return 0, errors.New("query type is not set")
	}

	query := &bongo.Query{
		Selector: map[string]interface{}{
			"message_id":    c.MessageId,
			"type_constant": q.Type,
		},
	}

	query.AddScope(RemoveTrollContent(
		c, q.ShowExempt,
	))

	return c.CountWithQuery(query)
}

func (c *Interaction) CountWithQuery(q *bongo.Query) (int, error) {
	return bongo.B.CountWithQuery(c, q)
}

func (c *Interaction) FetchAll(interactionType string) ([]Interaction, error) {
	var interactions []Interaction

	if c.MessageId == 0 {
		return interactions, errors.New("ChannelId is not set")
	}

	selector := map[string]interface{}{
		"message_id":    c.MessageId,
		"type_constant": interactionType,
	}

	err := c.Some(&interactions, bongo.NewQS(selector))
	if err != nil {
		return interactions, err
	}

	return interactions, nil
}

func (i *Interaction) IsInteracted(accountId int64) (bool, error) {
	if i.MessageId == 0 {
		return false, errors.New("Message Id is not set")
	}

	selector := map[string]interface{}{
		"message_id": i.MessageId,
		"account_id": accountId,
	}

	// do not set
	err := NewInteraction().One(bongo.NewQS(selector))
	if err == nil {
		return true, nil
	}

	if err == bongo.RecordNotFound {
		return false, nil
	}

	return false, err
}

func (i *Interaction) FetchInteractorCount() (int, error) {
	return bongo.B.Count(i, "message_id = ?", i.MessageId)
}
