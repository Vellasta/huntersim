package hunter

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func (hunter *Hunter) applyEndlessQuiver() {
	if hunter.Talents.EndlessQuiver == 0 {
		return
	}

	spellID := map[int32]int32{
		1: 51517,
		2: 51518,
	}[hunter.Talents.EndlessQuiver]

	hunter.EndlessQuiver = hunter.RegisterSpell(core.SpellConfig{
		SpellCode:    SpellCode_HunterEndlessQuiver,
		ActionID:     core.ActionID{SpellID: spellID},
		SpellSchool:  core.SpellSchoolPhysical,
		DefenseType:  core.DefenseTypeRanged,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | SpellFlagShot,
		CastType:     proto.CastType_CastTypeRanged,
		MissileSpeed: 24,

		CritDamageBonus: hunter.mortalShots(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := hunter.AutoAttacks.Ranged().CalculateWeaponDamage(sim, spell.RangedAttackPower(target, false)) +
				hunter.AmmoDamageBonus

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Endless Quiver Talent",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskEmpty) || !(spell.ProcMask.Matches(core.ProcMaskRangedAuto) || spell.SpellCode == SpellCode_HunterSteadyShot || spell.SpellCode == SpellCode_HunterMultiShot || spell.SpellCode == SpellCode_HunterArcaneShot) {
				return
			}

			EndlessQuiverProcChance := 0.03 * float64(hunter.Talents.EndlessQuiver)
			if sim.Proc(EndlessQuiverProcChance, "Extra Shot") {
				hunter.procEndlessQuiver(sim, result.Target)
			}
		},
	}))
}

func (hunter *Hunter) procEndlessQuiver(sim *core.Simulation, target *core.Unit) {
	hunter.EndlessQuiver.Cast(sim, target)
}
