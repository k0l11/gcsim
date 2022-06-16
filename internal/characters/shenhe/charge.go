package shenhe

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 49

func init() {
	chargeFrames = frames.InitAbilSlice(49)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupPole,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charged[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), chargeHitmark, chargeHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
