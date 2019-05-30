package taskstreamer

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ttype "github.com/marbar3778/simpleM/x/x/task_streamtypes"
)

type Keeper struct {
	coinKeeper bank.Keeper

	taskStore sdk.StoreKey

	cdc *codec.Codec
}

func NewKeeper(coinKeeper bank.Keeper, taskStore sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		taskStore:  taskStore,
		cdc:        cdc,
	}
}

func (k Keeper) GetTask(ctx sdk.Context, taskTitle string) ttype.Task {
	store := ctx.KVStore(k.taskStore)
	if !store.Has([]byte(taskTitle)) {
		panic("Naa")
	}
	task := store.Get([]byte(taskTitle))
	var taskDetails ttype.Task
	k.cdc.MustUnmarshalBinaryBare(task, &taskDetails)
	return taskDetails
}

func (k Keeper) SetTask(ctx sdk.Context, taskTitle string, taskData ttype.Task) {
	store := ctx.KVStore(k.taskStore)
	store.Set([]byte(taskTitle), k.cdc.MustMarshalBinaryBare(taskData))
}

func (k Keeper) CreateTask(ctx sdk.Context, taskTitle string, taskDescription string, backers []sdk.AccAddress, value sdk.Coins) {
	task := ttype.CreateTask(taskTitle, taskDescription, backers, value)
	k.SetTask(ctx, taskTitle, task)
}

func (k Keeper) BecomeBacker(ctx sdk.Context, taskTitle string, newBacker sdk.AccAddress, addedValue sdk.Coins) {
	task := k.GetTask(ctx, taskTitle)
	task.Value = task.Value.Add(addedValue)
	task.Backers = append(task.Backers, newBacker)
	k.SetTask(ctx, taskTitle, task)
}

func (k Keeper) GetSingleTask(ctx sdk.Context, taskTitle string) {
	k.GetTask(ctx, taskTitle)
}

// func (k Keeper) PayoutTask(ctx sdk.Context, taskTitle string, receiver sdk.AccAddress) {

// }
