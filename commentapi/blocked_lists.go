package commentapi

import (
	"net/http"

	"github.com/volatiletech/null"

	"github.com/lbryio/commentron/validator"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	v "github.com/lbryio/ozzo-validation"
)

// SharedBlockedListUpdateArgs use for blockedlist.Update api
type SharedBlockedListUpdateArgs struct {
	Authorization
	SharedBlockedList
}

// Validate validates the data in the update args
func (a SharedBlockedListUpdateArgs) Validate() api.StatusError {
	err := v.ValidateStruct(&a,
		v.Field(&a.ChannelID, validator.ClaimID, v.Required),
		v.Field(&a.ChannelName, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}

	return api.StatusError{}
}

// SharedBlockedList default representation of a shared blocked
type SharedBlockedList struct {
	// A user friendly identifier for the owner/users
	Name *string `json:"name"`
	// The category of block list this is so others search
	Category    *string `json:"category"`
	Description *string `json:"description"`
	// Can members invite others contributors?
	MemberInviteEnabled *bool `json:"member_invite_enabled"`
	// Strikes are number of hours a user should be banned for if
	// part of this blocked list. Strikes 1,2,3 are intended to be
	// progressively higher. Strike 3 is the highest.
	StrikeOne   *uint64 `json:"strike_one"`
	StrikeTwo   *uint64 `json:"strike_two"`
	StrikeThree *uint64 `json:"strike_three"`
	// The number of hours until a sent invite expires.
	InviteExpiration *uint64 `json:"invite_expiration"`
	// Curse jar allows automatic appeals. If they tip the owner of
	// the shared blocked list their appeal is automatically accepted.
	CurseJarAmount *uint64 `json:"curse_jar_amount"`
}

// SharedBlockedListInviteArgs arguments for blocklist.Invite
type SharedBlockedListInviteArgs struct {
	Authorization

	SharedBlockedListID uint64 `json:"blocked_list_id"`
	InviteeChannelName  string `json:"invitee_channel_name"`
	InviteeChannelID    string `json:"invitee_channel_id"`
	Message             string `json:"message"`
}

// SharedBlockedListInviteResponse empty respose for blocklist.Invite
type SharedBlockedListInviteResponse struct {
}

// SharedBlockedListInviteAcceptArgs arguments for blocklist.Accept
type SharedBlockedListInviteAcceptArgs struct {
	Authorization

	SharedBlockedListID uint64 `json:"blocked_list_id"`
	Accepted            bool   `json:"accepted"`
}

// SharedBlockedListInviteAcceptResponse response for blocklist.Accept
type SharedBlockedListInviteAcceptResponse struct {
}

// SharedBlockedListRescindArgs arguments for blocklist.Rescind
type SharedBlockedListRescindArgs struct {
	Authorization

	InvitedChannelName string `json:"invited_channel_name"`
	InvitedChannelID   string `json:"invited_channel_id"`
}

// SharedBlockedListRescindResponse response for blocklist.Rescind
type SharedBlockedListRescindResponse struct {
}

// SharedBlockedListGetArgs arguments for blocklist.Get
type SharedBlockedListGetArgs struct {
	SharedBlockedListID uint64             `json:"blocked_list_id"`
	Status              InviteMemberStatus `json:"status"`
}

// InviteMemberStatus status of invited member
type InviteMemberStatus int

// InviteMemberStatusFrom from a `null.Bool` it provides the functional value
func InviteMemberStatusFrom(v null.Bool) string {
	if v.Valid && v.Bool {
		return "accepted"
	} else if v.Valid && !v.Bool {
		return "rejected"
	}
	return "pending"
}

const (
	// All the defult value for getting all invited members
	All InviteMemberStatus = iota
	// Pending invite has not been accepted or rejected
	Pending
	// Accepted invited member has accepted and their blocked entries merged into list
	Accepted
	// Rejected invited member rejected joining the list and their blocked entries are cleared.
	Rejected
	// None does not return any invited members of the list
	None
)

// SharedBlockedListInvitedMember representation of an InvitedMember
type SharedBlockedListInvitedMember struct {
	InvitedByChannelName string `json:"invited_by_channel_name"`
	InvitedByChannelID   string `json:"invited_by_channel_id"`
	InviterChannelName   string `json:"inviter_channel_name"`
	InviterChannelID     string `json:"inviter_channel_id"`
	Status               string `json:"status"`
	InviteMessage        string `json:"message"`
}

// SharedBlockedListGetResponse response for blocklist.Get
type SharedBlockedListGetResponse struct {
	BlockedList    SharedBlockedList                `json:"shared_blocked_list"`
	InvitedMembers []SharedBlockedListInvitedMember `json:"invited_members"`
}
