package messages

type SwitchBranchMsg struct {
	Err          error
	TargetBranch string
	NextAction   string
	SourceBranch string
}

type MessageClearMsg struct {
	MessageID int
}
