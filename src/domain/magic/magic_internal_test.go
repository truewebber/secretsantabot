package magic

import (
	"reflect"
	"testing"

	"github.com/truewebber/secretsantabot/domain/chat"
)

func Test_makePairs(t *testing.T) {
	t.Parallel()

	type args struct {
		participants []chat.Person
	}

	tests := []struct {
		name string
		args args
		want []chat.GiverReceiverPair
	}{
		{
			name: "minimum amount of participants",
			args: args{
				participants: []chat.Person{
					{
						TelegramUserID: 1,
					},
					{
						TelegramUserID: 2,
					},
				},
			},
			want: []chat.GiverReceiverPair{
				{
					Giver:    chat.Person{TelegramUserID: 1},
					Receiver: chat.Person{TelegramUserID: 2},
				},
				{
					Giver:    chat.Person{TelegramUserID: 2},
					Receiver: chat.Person{TelegramUserID: 1},
				},
			},
		},
		{
			name: "odd amount of participants",
			args: args{
				participants: []chat.Person{
					{
						TelegramUserID: 1,
					},
					{
						TelegramUserID: 2,
					},
					{
						TelegramUserID: 3,
					},
				},
			},
			want: []chat.GiverReceiverPair{
				{
					Giver:    chat.Person{TelegramUserID: 1},
					Receiver: chat.Person{TelegramUserID: 2},
				},
				{
					Giver:    chat.Person{TelegramUserID: 2},
					Receiver: chat.Person{TelegramUserID: 3},
				},
				{
					Giver:    chat.Person{TelegramUserID: 3},
					Receiver: chat.Person{TelegramUserID: 1},
				},
			},
		},
		{
			name: "even amount of participants but not minimum",
			args: args{
				participants: []chat.Person{
					{
						TelegramUserID: 1,
					},
					{
						TelegramUserID: 2,
					},
					{
						TelegramUserID: 3,
					},
					{
						TelegramUserID: 4,
					},
				},
			},
			want: []chat.GiverReceiverPair{
				{
					Giver:    chat.Person{TelegramUserID: 1},
					Receiver: chat.Person{TelegramUserID: 2},
				},
				{
					Giver:    chat.Person{TelegramUserID: 2},
					Receiver: chat.Person{TelegramUserID: 3},
				},
				{
					Giver:    chat.Person{TelegramUserID: 3},
					Receiver: chat.Person{TelegramUserID: 4},
				},
				{
					Giver:    chat.Person{TelegramUserID: 4},
					Receiver: chat.Person{TelegramUserID: 1},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := makePairs(tt.args.participants); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makePairs() = %v, want %v", got, tt.want)
			}
		})
	}
}
