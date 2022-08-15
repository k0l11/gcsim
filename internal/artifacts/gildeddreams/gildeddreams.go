package gildeddreams

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.GildedDreams, NewSet)
}

type Set struct {
	atkCount int
	emCount  int
	Index    int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

// 2-Piece Bonus: Elemental Mastery +80.
// 4-Piece Bonus: Within 8s of triggering an Elemental Reaction, the character equipping this will obtain buffs based on the Elemental
// Type of the other party members. ATK is increased by 14% for each party member whose Elemental Type is the same as the equipping
// character, and Elemental Mastery is increased by 50 for every party member with a different Elemental Type. Each of the aforementioned
// buffs will count up to 3 characters. This effect can be triggered once every 8s. The character who equips this can still trigger its
// effects when not on the field.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("gd-2pc", -1),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		add := func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if c.Player.Active() != char.Index {
				return false
			}
			if atk.Info.ActorIndex != char.Index {
				return false
			}

			s.emCount = 0
			s.atkCount = 0

			for _, this := range c.Player.Chars() {
				if char.Index == this.Index {
					continue
				}
				if this.Base.Element != char.Base.Element {
					s.emCount++
				} else {
					s.atkCount++
				}
			}

			if s.emCount > 3 {
				s.emCount = 3
			}
			if s.atkCount > 3 {
				s.atkCount = 3
			}

			m := make([]float64, attributes.EndStatType)
			m[attributes.EM] = 50 * float64(s.emCount)
			m[attributes.ATKP] = 0.14 * float64(s.atkCount)
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase("gd-4pc", 8*60),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
			c.Log.NewEvent("gilde ddreams proc'd", glog.LogArtifactEvent, char.Index).
				Write("em", s.emCount).
				Write("atk", s.atkCount)
			return false
		}

		for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
			c.Events.Subscribe(i, add, fmt.Sprintf("gd-4pc-%v", char.Base.Key.String()))
		}
	}

	return &s, nil
}
