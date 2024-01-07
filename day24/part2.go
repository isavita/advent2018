package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Group struct {
	Units        int
	HitPoints    int
	AttackDamage int
	AttackType   string
	Initiative   int
	Immunities   []string
	Weaknesses   []string

	Attacker *Group
	Target   *Group
}

func (g Group) EffectivePower() int {
	return g.Units * g.AttackDamage
}

func (g Group) DamageDealt(e *Group) int {
	if contrainsString(e.Immunities, g.AttackType) {
		return 0
	}
	if contrainsString(e.Weaknesses, g.AttackType) {
		return g.EffectivePower() * 2
	}
	return g.EffectivePower()
}

func contrainsString(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func (g Group) String() string {
	out := fmt.Sprintf("%d units each with %d hit points", g.Units, g.HitPoints)
	if len(g.Immunities) > 0 || len(g.Weaknesses) > 0 {
		out += " ("
		if len(g.Immunities) > 0 {
			out += "immune to " + strings.Join(g.Immunities, " and ")
			if len(g.Weaknesses) > 0 {
				out += "; "
			}
		}
		if len(g.Weaknesses) > 0 {
			out += "weak to " + strings.Join(g.Weaknesses, " and ")
		}
		out += ")"
	}
	out += fmt.Sprintf(" with an attack that does %d %s damage at initiative %d", g.AttackDamage, g.AttackType, g.Initiative)
	return out
}

type Initiative []*Group

func (in Initiative) Len() int {
	return len(in)
}

func (in Initiative) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

func (in Initiative) Less(i, j int) bool {
	return in[i].Initiative > in[j].Initiative
}

func (in Initiative) Attack() {
	sort.Sort(in)

	for _, group := range in {
		if group.Units > 0 && group.Target != nil && group.Target.Units > 0 {
			group.Target.Units -= group.DamageDealt(group.Target) / group.Target.HitPoints
		}
		if group.Target != nil {
			group.Target.Attacker = nil
			group.Target = nil
		}
	}
}

func (in *Initiative) Clean() {
	c := (*in)[:0]
	for _, g := range *in {
		if g.Units > 0 {
			c = append(c, g)
		}
	}
	sort.Sort(c)
	*in = c
}

type Army []*Group

const (
	ArmyImmuneSystem = iota + 1
	ArmyInfection
	ArmyCount
)

var StringArmies = map[string]int{
	"Immune System": ArmyImmuneSystem,
	"Infection":     ArmyInfection,
}

func (a Army) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Army) Len() int {
	return len(a)
}

func (a Army) Less(i, j int) bool {
	if a[i].EffectivePower() > a[j].EffectivePower() {
		return true
	}
	return a[i].EffectivePower() == a[j].EffectivePower() && a[i].Initiative > a[j].Initiative
}

func (a Army) Alive() bool {
	for _, g := range a {
		if g.Units > 0 {
			return true
		}
	}
	return false
}

func (a Army) Boost(amount int) {
	for _, g := range a {
		g.AttackDamage += amount
	}
}

type Battlefield map[int]Army

var (
	armyName         = regexp.MustCompile(`^(.*):$`)
	groupImmunities  = regexp.MustCompile(`immune to (.*?)[;)]`)
	groupWeaknesses  = regexp.MustCompile(`weak to (.*?)[;)]`)
	groupDescription = regexp.MustCompile(`^(\d+) units each with (\d+) hit points.*with an attack that does (\d+) (\w+) damage at initiative (\d+)$`)
)

const (
	descriptionCount = iota + 1
	descriptionHitPoints
	descriptionDamage
	descriptionDamageType
	descriptionInitiative
)

func PrepareForBattle(input []string) (Battlefield, Initiative) {
	var initiative Initiative
	battle := make(Battlefield)
	var currentArmy int

	for _, line := range input {
		if armyName.MatchString(line) {
			if id, ok := StringArmies[armyName.FindStringSubmatch(line)[1]]; ok {
				currentArmy = id
			} else {
				panic(fmt.Errorf("unknown army: %s", armyName.FindStringSubmatch(line)[1]))
			}
		} else {
			if currentArmy <= 0 || currentArmy >= ArmyCount {
				panic(fmt.Errorf("tried to assign group to invalid army: %d", currentArmy))
			}
			description := groupDescription.FindStringSubmatch(line)
			if len(description) == 0 {
				continue
			}

			group := &Group{
				Units:        intOrPanic(description[descriptionCount]),
				HitPoints:    intOrPanic(description[descriptionHitPoints]),
				AttackDamage: intOrPanic(description[descriptionDamage]),
				AttackType:   description[descriptionDamageType],
				Initiative:   intOrPanic(description[descriptionInitiative]),
			}

			immunities := groupImmunities.FindStringSubmatch(line)
			if len(immunities) > 0 {
				group.Immunities = strings.Split(immunities[1], ", ")
			}

			weaknesses := groupWeaknesses.FindStringSubmatch(line)
			if len(weaknesses) > 0 {
				group.Weaknesses = strings.Split(weaknesses[1], ", ")
			}

			battle[currentArmy] = append(battle[currentArmy], group)
			initiative = append(initiative, group)
		}
	}

	return battle, initiative
}

func (b Battlefield) FindTargets() {
	for army, groups := range b {
		sort.Sort(groups)
		for _, group := range groups {
			for enemyArmy, enemyGroups := range b {
				if army == enemyArmy || group.Units <= 0 {
					continue
				}

				var mostDamage int
				var targetGroup *Group

				for _, enemyGroup := range enemyGroups {
					if enemyGroup.Units <= 0 || enemyGroup.Attacker != nil || group.DamageDealt(enemyGroup) == 0 || group.DamageDealt(enemyGroup) < mostDamage {
						continue
					}
					if group.DamageDealt(enemyGroup) == mostDamage && targetGroup != nil {
						if enemyGroup.EffectivePower() < targetGroup.EffectivePower() {
							continue
						}
						if enemyGroup.EffectivePower() == targetGroup.EffectivePower() && enemyGroup.Initiative < targetGroup.Initiative {
							continue
						}
					}
					mostDamage = group.DamageDealt(enemyGroup)
					targetGroup = enemyGroup
				}
				if targetGroup != nil {
					group.Target = targetGroup
					targetGroup.Attacker = group
				}
			}
		}
	}
}

func (b Battlefield) Clean() {
	for army := range b {
		c := b[army][:0]
		for _, g := range b[army] {
			if g.Units > 0 {
				c = append(c, g)
			}
		}
		b[army] = c
	}
}

func (b Battlefield) Active() bool {
	for _, a := range b {
		if !a.Alive() {
			return false
		}
	}
	return true
}

func (b Battlefield) Result() (int, int) {
	var winner, units int

	for army, groups := range b {
		if groups.Alive() {
			winner = army

			for _, g := range groups {
				if g.Units > 0 {
					units += g.Units
				}
			}
		}
	}

	return winner, units
}

func (b Battlefield) TotalUnits() int {
	var sum int
	for _, groups := range b {
		for _, group := range groups {
			if group.Units > 0 {
				sum += group.Units
			}
		}
	}
	return sum
}

func intOrPanic(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func ImmuneSystemBoost(input []string) int {
	var boost int

	for {
		var stalemate bool
		battle, initiative := PrepareForBattle(input)

		battle[ArmyImmuneSystem].Boost(boost)

		for battle.Active() {
			before := battle.TotalUnits()

			battle.FindTargets()
			initiative.Attack()

			if battle.TotalUnits() == before {
				stalemate = true
				break
			}
			battle.Clean()
			initiative.Clean()
		}

		if !stalemate {
			winner, units := battle.Result()
			if winner == ArmyImmuneSystem {
				return units
			}
		}

		boost++
	}
}

func main() {
	input, err := os.ReadFile("day24/input.txt")
	if err != nil {
		panic(err)
	}
	armies := strings.Split(strings.TrimSpace(string(input)), "\n")
	fmt.Println(ImmuneSystemBoost(armies))
}
