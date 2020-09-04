package commentapi

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc"
)

// DefaultPort the default port that is used for commentron client
const DefaultPort = 5900

type Client struct {
	conn    jsonrpc.RPCClient
	address string
}

func NewClient(address string) *Client {
	d := Client{}

	if address == "" {
		address = "http://localhost:" + strconv.Itoa(DefaultPort)
	}

	d.conn = jsonrpc.NewClient(address)
	d.address = address

	return &d
}

///////////////////////
//  REACTION SERVICE //
///////////////////////

func (d *Client) ReactionList(args ReactionListArgs) (*ReactionListResponse, error) {
	structs.DefaultTagName = "json"
	response := new(ReactionListResponse)
	return response, d.call(response, "reaction.List", structs.Map(args))
}

func (d *Client) ReactionReact(args ReactArgs) (*ReactResponse, error) {
	structs.DefaultTagName = "json"
	response := new(ReactResponse)
	return response, d.call(response, "reaction.React", structs.Map(args))
}

//////////////////////
//  COMMENT SERVICE //
//////////////////////

func (d *Client) CommentList(args ListArgs) (*ListResponse, error) {
	structs.DefaultTagName = "json"
	response := new(ListResponse)
	return response, d.call(response, "comment.List", structs.Map(args))
}

func (d *Client) CommentByID(args ByIDArgs) (*ByIDResponse, error) {
	structs.DefaultTagName = "json"
	response := new(ByIDResponse)
	return response, d.call(response, "comment.ByID", structs.Map(args))
}

func (d *Client) CommentAbandon(args AbandonArgs) (*AbandonResponse, error) {
	structs.DefaultTagName = "json"
	response := new(AbandonResponse)
	return response, d.call(response, "comment.Abandon", structs.Map(args))
}

func (d *Client) CommentCreate(args CreateArgs) (*CreateResponse, error) {
	structs.DefaultTagName = "json"
	response := new(CreateResponse)
	return response, d.call(response, "comment.Create", structs.Map(args))
}

func (d *Client) CommentEdit(args EditArgs) (*EditResponse, error) {
	structs.DefaultTagName = "json"
	response := new(EditResponse)
	return response, d.call(response, "comment.Edit", structs.Map(args))
}

func (d *Client) GetChannelForComment(args ChannelArgs) (*ChannelResponse, error) {
	structs.DefaultTagName = "json"
	response := new(ChannelResponse)
	return response, d.call(response, "comment.Abandon", structs.Map(args))
}

////////////////
//  INTERNALS //
////////////////

func decode(data interface{}, targetStruct interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   targetStruct,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	err = decoder.Decode(data)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	return nil
}

func debugParams(params map[string]interface{}) string {
	var s []string
	for k, v := range params {
		r := reflect.ValueOf(v)
		if r.Kind() == reflect.Ptr {
			if r.IsNil() {
				continue
			}
			v = r.Elem().Interface()
		}
		s = append(s, fmt.Sprintf("%s=%+v", k, v))
	}
	sort.Strings(s)
	return strings.Join(s, " ")
}

func (d *Client) callNoDecode(command string, params map[string]interface{}) (interface{}, error) {
	log.Debugln("jsonrpc: " + command + " " + debugParams(params))
	r, err := d.conn.Call(command, params)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if r.Error != nil {
		return nil, errors.Err("Error in daemon: " + r.Error.Message)
	}

	return r.Result, nil
}

func (d *Client) call(response interface{}, command string, params map[string]interface{}) error {
	result, err := d.callNoDecode(command, params)
	if err != nil {
		return err
	}
	return decode(result, response)
}

func (d *Client) SetRPCTimeout(timeout time.Duration) {
	d.conn = jsonrpc.NewClientWithOpts(d.address, &jsonrpc.RPCClientOpts{
		HTTPClient: &http.Client{Timeout: timeout},
	})
}
