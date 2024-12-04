package authz

import (
	"encoding/json"
	"testing"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmod "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
)

func TestMessageAuthzSerialization(t *testing.T, msg sdk.Msg, module module.AppModuleBasic) {
	someDate := time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
	const (
		mockGranter string = "cosmos1abc"
		mockGrantee string = "cosmos1xyz"
	)

	var (
		mockMsgGrant  authz.MsgGrant
		mockMsgRevoke authz.MsgRevoke
		mockMsgExec   authz.MsgExec
	)

	// mutates mockMsg
	testSerDeser := func(msg proto.Message, mockMsg proto.Message) {
		encCfg := moduletestutil.MakeTestEncodingConfig(authzmod.AppModuleBasic{}, module)
		msgGrantBytes := json.RawMessage(sdk.MustSortJSON(encCfg.Codec.MustMarshalJSON(msg)))
		err := encCfg.Codec.UnmarshalJSON(msgGrantBytes, mockMsg)
		require.NoError(t, err)
	}

	// Authz: Grant Msg
	typeURL := sdk.MsgTypeURL(msg)
	expiryTime := someDate.Add(time.Hour)
	grant, err := authz.NewGrant(someDate, authz.NewGenericAuthorization(typeURL), &expiryTime)
	require.NoError(t, err)

	msgGrant := authz.MsgGrant{Granter: mockGranter, Grantee: mockGrantee, Grant: grant}
	testSerDeser(&msgGrant, &mockMsgGrant)

	// Authz: Revoke Msg
	msgRevoke := authz.MsgRevoke{Granter: mockGranter, Grantee: mockGrantee, MsgTypeUrl: typeURL}
	testSerDeser(&msgRevoke, &mockMsgRevoke)

	// Authz: Exec Msg
	msgAny, err := cdctypes.NewAnyWithValue(msg)
	require.NoError(t, err)
	msgExec := authz.MsgExec{Grantee: mockGrantee, Msgs: []*cdctypes.Any{msgAny}}
	testSerDeser(&msgExec, &mockMsgExec)
	require.Equal(t, msgExec.Msgs[0].Value, mockMsgExec.Msgs[0].Value)
}
