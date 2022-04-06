package hello

import (
	"context"
	"fmt"
	"log"

	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/types"
)

const (
	_greettingSeed = "hello"
)

type HelloWorldScheme struct {
	Counter uint32
}

type HelloWorld struct {
	Payer         types.Account
	PayerPubkey   common.PublicKey
	ProgramPubkey common.PublicKey
	GreetedPubKey common.PublicKey
	conn          *client.Client
}

func GetHelloWorld(payer types.Account, programIdB58 string) *HelloWorld {
	conn := GetConnection()

	ctx := context.Background()

	info, err := conn.GetAccountInfo(ctx, programIdB58)
	if err != nil {
		log.Fatalf("fail to get program info, err: %v", err)
	}

	if !info.Executable {
		log.Fatalf("this program is unexcutable")
	}

	fmt.Println("Using program ", programIdB58)

	payerPubkey := payer.PublicKey
	programPubkey := common.PublicKeyFromString(programIdB58)

	greetedPubkey := common.CreateWithSeed(payerPubkey, _greettingSeed, programPubkey)
	accountInfo, _ := conn.GetAccountInfo(ctx, greetedPubkey.String())
	if accountInfo.Owner == "" {
		blockhash, _ := conn.GetLatestBlockhash(ctx)

		msg := types.NewMessage(
			types.NewMessageParam{
				FeePayer: payerPubkey,
				Instructions: []types.Instruction{
					sysprog.CreateAccountWithSeed(
						sysprog.CreateAccountWithSeedParam{
							From:  payerPubkey,
							Base:  payerPubkey,
							Seed:  _greettingSeed,
							New:   greetedPubkey,
							Owner: programPubkey,
						},
					),
				},
				RecentBlockhash: blockhash.Blockhash,
			},
		)

		tx, err := types.NewTransaction(
			types.NewTransactionParam{
				Message: msg,
				Signers: []types.Account{payer},
			},
		)
		if err != nil {
			log.Fatalf("fail to gen tx, err: %v", err)
		}

		_, err = conn.SendTransaction(ctx, tx)
		if err != nil {
			log.Fatalf("fail to send tx, err: %v", err)
		}

		fmt.Println("Creating account ", greetedPubkey.String())
	}

	return &HelloWorld{
		Payer:         payer,
		PayerPubkey:   payerPubkey,
		ProgramPubkey: programPubkey,
		GreetedPubKey: common.CreateWithSeed(payerPubkey, _greettingSeed, programPubkey),
		conn:          conn,
	}
}

func (h *HelloWorld) SayHello() {
	fmt.Println("Saying hello to ", h.GreetedPubKey.String())

	ctx := context.Background()
	blockhash, _ := h.conn.GetLatestBlockhash(ctx)

	msg := types.NewMessage(
		types.NewMessageParam{
			FeePayer: h.PayerPubkey,
			Instructions: []types.Instruction{
				{
					ProgramID: h.ProgramPubkey,
					Accounts: []types.AccountMeta{
						{
							PubKey:     h.GreetedPubKey,
							IsSigner:   false,
							IsWritable: true,
						},
					},
				},
			},
			RecentBlockhash: blockhash.Blockhash,
		},
	)

	tx, err := types.NewTransaction(
		types.NewTransactionParam{
			Message: msg,
			Signers: []types.Account{h.Payer},
		},
	)
	if err != nil {
		log.Fatalf("fail to gen tx, err: %v", err)
	}

	_, err = h.conn.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatalf("fail to send tx, err: %v", err)
	}
}

func (h *HelloWorld) ReportGreetings() {
	ctx := context.Background()
	accountInfo, err := h.conn.GetAccountInfo(ctx, h.GreetedPubKey.String())
	if err != nil {
		log.Fatalf("fail to get account info, err: %v", err)
	}

	var greeting HelloWorldScheme
	if err = borsh.Deserialize(&greeting, accountInfo.Data); err != nil {
		log.Fatalf("fail to deserialize account info data, err: %v", err)
	}

	fmt.Println(h.GreetedPubKey.String(), " has been greeted ", greeting.Counter, " time(s)")
}
