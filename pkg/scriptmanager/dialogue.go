package scriptmanager

type PieceDialogue struct {
	Replicas         []string
	currentReplica   int
	Answers          []Answer
	activeAnswer     int
	needAnswer       bool
	IsEndDialogue    bool
	CanStartDialogue bool
}

type Answer struct {
	Text              string
	Action            func()
	NextPieceDialogue *PieceDialogue
}

func (dialogue *PieceDialogue) NextReplica() {
	dialogue.currentReplica++
	if dialogue.currentReplica >= len(dialogue.Replicas)-1 {
		dialogue.currentReplica = len(dialogue.Replicas) - 1
		if len(dialogue.Answers) == 0 {
			dialogue.activeAnswer = 0
			dialogue.currentReplica = 0
			dialogue.IsEndDialogue = true
			return
		}
		dialogue.needAnswer = true
	}
}

func (dialogue *PieceDialogue) NextAnswer() {
	dialogue.activeAnswer++
	if dialogue.activeAnswer >= len(dialogue.Answers) {
		dialogue.activeAnswer = 0
	}
}

func (dialogue *PieceDialogue) BeforeAnswer() {
	dialogue.activeAnswer--
	if dialogue.activeAnswer < 0 {
		dialogue.activeAnswer = len(dialogue.Answers) - 1
	}
}

func (dialogue *PieceDialogue) NeedAnswer() bool {
	return dialogue.needAnswer
}

func (dialogue *PieceDialogue) CurrentReplica() string {
	return dialogue.Replicas[dialogue.currentReplica]
}

func (dialogue *PieceDialogue) IsActiveAnswer(index int) bool {
	return dialogue.activeAnswer == index
}

func (dialogue *PieceDialogue) DoAnswer() {
	if len(dialogue.Answers) == 0 {
		return
	}
	if dialogue.Answers[dialogue.activeAnswer].Action != nil {
		dialogue.Answers[dialogue.activeAnswer].Action()
	}
	if dialogue.Answers[dialogue.activeAnswer].NextPieceDialogue == nil {
		dialogue.IsEndDialogue = true
	}
}

func (dialogue *PieceDialogue) NextPieceDialogue() *PieceDialogue {
	if len(dialogue.Answers) == 0 {
		dialogue.IsEndDialogue = true
		return dialogue
	}
	nextPieceDialogue := dialogue.Answers[dialogue.activeAnswer].NextPieceDialogue
	if nextPieceDialogue == nil {
		dialogue.IsEndDialogue = true
		return dialogue
	}
	dialogue.currentReplica = 0
	dialogue.activeAnswer = 0
	if len(nextPieceDialogue.Replicas) == 1 {
		nextPieceDialogue.needAnswer = true
	}
	return nextPieceDialogue
}
