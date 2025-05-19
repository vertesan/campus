import { MemoryInspector, XMemoryInspector } from "~/types"
import { getExamEffects, getSingleXProduceCard } from "./pcard"

export function getXMemoryInspector([
  ProduceCard,
  ProduceExamEffect,
  ProduceItem,
  MemoryAbility,
  ProduceSkill,
  ProduceEffect,
]: MemoryInspector): XMemoryInspector {
  const examEffects = getExamEffects(ProduceExamEffect)

  const xProduceCards = ProduceCard
    .filter(x => x.upgradeCount < 2)
    .reduce<XMemoryInspector['produceCards']>((acc, cur) => {
      acc[`${cur.id}-${cur.upgradeCount}`] = getSingleXProduceCard(cur, examEffects)
      return acc
    }, {})

  const produceItems = ProduceItem.reduce<XMemoryInspector['produceItems']>((acc, cur) => {
    acc[cur.id] = cur
    return acc
  }, {})

  const xMemoryAbilities = MemoryAbility.reduce<XMemoryInspector['memoryAbilities']>((acc, cur) => {
    const skill = ProduceSkill.find(x => x.id === cur.skillId)
    if (!skill) return acc
    const effects = []
    for (let i = 1; i <= 3; i++) {
      const produceEffectId = skill[`produceEffectId${i.toString() as '1' | '2' | '3'}`]
      if (produceEffectId === '') continue
      const effect = ProduceEffect.find(x => x.id === produceEffectId)
      if (!effect) continue
      effects.push(effect)
    }
    const xMAbility = {
      ...cur,
      skill: {
        ...skill,
        produceEffects: effects,
      },
    }
    acc[cur.id] = xMAbility
    return acc
  }, {})

  return {
    produceCards: xProduceCards,
    produceItems: produceItems,
    memoryAbilities: xMemoryAbilities,
  }
}