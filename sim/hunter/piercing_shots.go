package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (hunter *Hunter) applyPiercingShots() {
	if hunter.Talents.PiercingShots == 0 {
		return
	}

	spellID := map[int32]int32{
		1: 51512,
		2: 51513,
	}[hunter.Talents.PiercingShots]

	hunter.PiercingShots = hunter.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_HunterPiercingShots,
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 0,
		BonusCoefficient: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Piercing Shots",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim) //Resets the tick counter with Apply vs ApplyorRefresh
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Piercing Shots Talent",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskEmpty) || !(spell.SpellCode == SpellCode_HunterSteadyShot || spell.SpellCode == SpellCode_HunterMultiShot || spell.SpellCode == SpellCode_HunterAimedShot) {
				return
			}

			if result.Outcome.Matches(core.OutcomeCrit) {
				hunter.procPiercingShots(sim, result.Target, result)
			}
		},
	}))
}

func (hunter *Hunter) procPiercingShots(sim *core.Simulation, target *core.Unit, result *core.SpellResult) {
	dot := hunter.PiercingShots.Dot(target)

	newDamage := result.Damage * 0.15 * float64(hunter.Talents.PiercingShots) // 30% of casted shot damage

	dot.SnapshotBaseDamage = newDamage / 4.0 // spread over 4 ticks of the dot
	dot.SnapshotAttackerMultiplier = 1

	hunter.PiercingShots.Cast(sim, target)
}
