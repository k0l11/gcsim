package ayaka

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 8
		case 1:
			f = 20
		case 2:
			f = 28
		case 3:
			f = 43
		case 4:
			f = 37
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 28, 53 // 28 frames until hitmark, 53 from TCL
	case core.ActionSkill:
		return 56, 56 //should be 82
	case core.ActionBurst:
		return 95, 95 //ok
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 53-28) //charge -> n1, just fill it up to be like original
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 53-28)  //charge -> skill, just fill it up to be like original
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 53-28)  //charge -> burst, just fill it up to be like original
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionDash, 53-28)   //charge -> dash, just fill it up to be like original
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionJump, 53-28)   //charge -> jump, just fill it up to be like original
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 30-28)   //charge -> swap, just fill it up to be 30 frames
}
