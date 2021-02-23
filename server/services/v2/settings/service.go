package settings

import (
	"net/http"
	"strings"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/errors.go"
	"github.com/lbryio/lbry.go/extras/api"

	"github.com/volatiletech/sqlboiler/boil"
)

// Service is the service struct defined for the comment package for rpc service "moderation.*"
type Service struct{}

// BlockWord takes a list of words to block comments containing these words. These words are added to the existing list
func (s *Service) BlockWord(r *http.Request, args *commentapi.BlockWordArgs, reply *commentapi.BlockWordRespose) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}
	existingWords := strings.Split(settings.MutedWords.String, ",")
	wordsToAdd := strings.Split(args.Words, ",")
	existingWords = append(existingWords, wordsToAdd...)
	settings.MutedWords.SetValid(strings.Join(existingWords, ","))
	err = settings.UpdateG(boil.Infer())
	if err != nil {
		return errors.Err(err)
	}
	reply.WordList = existingWords
	return nil
}

// UnBlockWord takes a list of words to remove from the list of blocked words if they exist.
func (s *Service) UnBlockWord(r *http.Request, args *commentapi.UnBlockWordArgs, reply *commentapi.BlockWordRespose) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	if creatorChannel == nil {
		return api.StatusError{Err: errors.Err("could not find channel %s with channel id %s", args.ChannelName, args.ChannelID)}
	}
	err = lbry.ValidateSignature(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}
	existingWords := strings.Split(settings.MutedWords.String, ",")
	wordsToRemove := strings.Split(args.Words, ",")
	remainingWords := make([]string, 0)
skip:
	for _, word := range existingWords {
		for _, wordToRemove := range wordsToRemove {
			if wordToRemove == word {
				continue skip
			}
		}
		remainingWords = append(remainingWords, word)
	}
	settings.MutedWords.SetValid(strings.Join(remainingWords, ","))
	err = settings.UpdateG(boil.Infer())
	if err != nil {
		return errors.Err(err)
	}
	reply.WordList = remainingWords
	return nil
}

// ListBlockedWords returns a list of all the current blocked words for a channel.
func (s *Service) ListBlockedWords(r *http.Request, args *commentapi.ListBlockedWordsArgs, reply *commentapi.BlockWordRespose) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}
	existingWords := strings.Split(settings.MutedWords.String, ",")
	reply.WordList = existingWords
	return nil
}
