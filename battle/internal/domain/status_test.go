package domain

import (
	"testing"
)

func TestBattleState(t *testing.T) {
	battle := Battle{
		status: StatusPending,
	}
	stateMachine, _ := NewBattleStateMachine(&battle)

	cases := []struct {
		name     string
		statusTo BattleStatus
		expected BattleStatus
		wantErr  bool
	}{
		{
			name:     "pending to finish",
			statusTo: StatusFinished,
			expected: StatusPending,
			wantErr:  true,
		},
		{
			name:     "pending to inProgress",
			statusTo: StatusProgress,
			expected: StatusProgress,
			wantErr:  false,
		},
		{
			name:     "inProgress to pending",
			statusTo: StatusPending,
			expected: StatusProgress,
			wantErr:  true,
		},
		{
			name:     "inProgress to inProgress",
			statusTo: StatusProgress,
			expected: StatusProgress,
			wantErr:  true,
		},
		{
			name:     "inProgress to finish",
			statusTo: StatusFinished,
			expected: StatusFinished,
			wantErr:  false,
		},
		{
			name:     "finish to pending",
			statusTo: StatusPending,
			expected: StatusFinished,
			wantErr:  true,
		},
		{
			name:     "finish to inProgress",
			statusTo: StatusProgress,
			expected: StatusFinished,
			wantErr:  true,
		},
		{
			name:     "finish to finish",
			statusTo: StatusFinished,
			expected: StatusFinished,
			wantErr:  true,
		},
	}

	var err error
	for _, test := range cases {
		err = stateMachine.Change(test.statusTo)
		if test.wantErr && err == nil {
			t.Errorf("%s: got err == nil, expected err != nil", test.name)
		}
		if !test.wantErr && err != nil {
			t.Errorf("%s: got err == %s, expected err == nil", test.name, err)
		}
		if battle.GetStatus() != test.expected {
			t.Errorf("%s: got status == %s, expected status == %s", test.name, battle.GetStatus(), test.expected)
		}
	}
}
