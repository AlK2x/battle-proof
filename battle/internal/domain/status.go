package domain

import (
	"errors"
	"fmt"
)

type BattleStatus string

const (
	StatusPending  BattleStatus = "pending"
	StatusProgress BattleStatus = "in_progress"
	StatusFinished BattleStatus = "finished"
)

func NewBattleStateMachine(battle *Battle) (*BattleStateMachine, error) {
	sm := &BattleStateMachine{
		battle: battle,
	}
	pendingState := &PendingState{
		stateMachine: sm,
	}
	inProgressState := &InProgressState{
		stateMachine: sm,
	}
	finishedState := &FinishedState{
		stateMachine: sm,
	}

	var state BattleState
	switch battle.Status {
	case StatusPending:
		state = pendingState
	case StatusProgress:
		state = inProgressState
	case StatusFinished:
		state = finishedState
	default:
		return nil, fmt.Errorf("unknown battle status %s", battle.Status)
	}

	sm.pendingState = pendingState
	sm.inProgressState = inProgressState
	sm.finishedState = finishedState

	sm.state = state

	return sm, nil
}

type BattleStateMachine struct {
	pendingState    BattleState
	inProgressState BattleState
	finishedState   BattleState

	state BattleState

	battle *Battle
}

func (b *BattleStateMachine) Change(to BattleStatus) error {
	return b.state.Change(to)
}

func (b *BattleStateMachine) setState(next BattleState, status BattleStatus) {
	b.state = next
	b.battle.Status = status
}

type BattleState interface {
	Change(to BattleStatus) error
	CanTransition(to BattleStatus) bool
}

type PendingState struct {
	stateMachine *BattleStateMachine
}

func (s *PendingState) Change(to BattleStatus) error {
	if !s.CanTransition(to) {
		return fmt.Errorf("cannot change battle status from %s to %s", StatusPending, to)
	}

	s.stateMachine.setState(s.stateMachine.inProgressState, to)
	return nil
}

func (s *PendingState) CanTransition(to BattleStatus) bool {
	return to == StatusProgress
}

type InProgressState struct {
	stateMachine *BattleStateMachine
}

func (s *InProgressState) Change(to BattleStatus) error {
	if !s.CanTransition(to) {
		return fmt.Errorf("cannot change battle status from %s to %s", StatusProgress, to)
	}

	s.stateMachine.setState(s.stateMachine.finishedState, to)
	return nil
}

func (s *InProgressState) CanTransition(to BattleStatus) bool {
	return to == StatusFinished
}

type FinishedState struct {
	stateMachine *BattleStateMachine
}

func (s *FinishedState) Change(to BattleStatus) error {
	return errors.New("cannot change status of finished battle")
}

func (s *FinishedState) CanTransition(to BattleStatus) bool {
	return false
}
