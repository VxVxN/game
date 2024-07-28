package quest

import (
	"fmt"

	"github.com/VxVxN/game/pkg/item"
)

type Quest struct {
	Name            string
	Goals           []*Goal
	rewardCallback  func()
	isAwardReceived bool
}

func NewQuest(name string, goals []*Goal, rewardCallback func()) *Quest {
	return &Quest{
		Name:           name,
		Goals:          goals,
		rewardCallback: rewardCallback,
	}
}

func (quest *Quest) UpdateProgress(items []*item.Item) {
	for _, goal := range quest.Goals {
		goal.CheckProgress(items)
	}
	if !quest.IsCompleted() {
		return
	}
	if !quest.isAwardReceived {
		quest.rewardCallback()
		quest.isAwardReceived = true
	}
}

func (quest *Quest) IsCompleted() bool {
	for _, goal := range quest.Goals {
		if !goal.isCompleted {
			return false
		}
	}
	return true
}

func (quest *Quest) GoalsDescription() string {
	var description string
	for _, goal := range quest.Goals {
		if len(goal.NeedItems) != 0 {
			description += "Find: "
		}
		for _, needItem := range goal.NeedItems {
			description += fmt.Sprintf("%s(%d/%d); ", needItem.Type.String(), needItem.NumberHave, needItem.NumberNeed)
		}
	}
	return description
}

type Goal struct {
	NeedItems   []NeedItem
	isCompleted bool
}

type NeedItem struct {
	Type       item.ItemType
	NumberNeed int
	NumberHave int
}

func (goal *Goal) CheckProgress(items []*item.Item) (int, int) {
	var progress int
	for i, needItem := range goal.NeedItems {
		goal.NeedItems[i].NumberHave = 0
		for _, item := range items {
			if item.ItemType == needItem.Type {
				goal.NeedItems[i].NumberHave++
				progress++
				break
			}
		}
	}
	if progress == len(goal.NeedItems) {
		goal.isCompleted = true
	}
	return progress, len(goal.NeedItems)
}
