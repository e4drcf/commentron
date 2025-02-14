package moderation

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func amI(_ *http.Request, args *commentapi.AmIArgs, reply *commentapi.AmIResponse) error {
	channel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(channel.ClaimID, args.Signature, args.SigningTS, channel.Name)
	if err != nil {
		return errors.Err(err)
	}
	moderations, err := channel.ModChannelDelegatedModerators(qm.Load(model.DelegatedModeratorRels.CreatorChannel)).All(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	approvedChannels := make(map[string]string)
	for _, moderation := range moderations {
		if moderation.R != nil && moderation.R.CreatorChannel != nil {
			reply.Type = "Channel"
			approvedChannels[moderation.R.CreatorChannel.Name] = moderation.R.CreatorChannel.ClaimID
		}
	}
	reply.ChannelName = args.ChannelName
	reply.ChannelID = args.ChannelID
	reply.AuthorizedChannels = approvedChannels

	moderator, err := model.Moderators(model.ModeratorWhere.ModChannelID.EQ(null.StringFrom(args.ChannelID)), model.ModeratorWhere.ModLevel.EQ(1)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if moderator != nil {
		reply.Type = "Global"
	}
	return nil
}
