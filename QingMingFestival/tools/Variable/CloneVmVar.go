package Variable

import (
	"QingMingFestival/types/vmTypes"
)

var CloneVmMessageChan chan vmTypes.ResultCloneVmMessage

func init() {
	CloneVmMessageChan = make(chan vmTypes.ResultCloneVmMessage)
}
