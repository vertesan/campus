import { readFile, readFileSync } from "fs"
import path from "path"
import { Cidol, Csprt, Master, MemoryInspector, PCard, UsedDB, isNonNull } from "~/types"
import { MappedUsedDBTuple, UnionArrayToTuple } from "~/types/utils"

export async function getMemoryInspector(dirPath: string): Promise<MemoryInspector | null> {
  return await getAllJson([
    "ProduceCard",
    "ProduceExamEffect",
    "ProduceItem",
    "MemoryAbility",
    "ProduceSkill",
    "ProduceEffect",
  ], dirPath)
}

export async function getPCard(dirPath: string): Promise<PCard | null> {
  return await getAllJson([
    "ProduceCard",
    "ProduceExamEffect",
    "ProduceCardCustomize",
    "ProduceCardCustomizeRarityEvaluation",
    "ProduceCardGrowEffect",
    "ProduceDescriptionProduceCardGrowEffect",
    "ProduceCardStatusEnchant",
    "ProduceExamTrigger",
  ], dirPath)
}

export async function getMaster(dirPath: string): Promise<Master | null> {
  return await getAllJson([
    "Character",
    "ProduceEffectIcon",
    "Produce",
    "ExamInitialDeck",
    "ProduceDescriptionProduceEffect",
    "ProduceDescriptionExamEffect",
    "CharacterTrueEndBonus",
    "HomeEnter",
    "NoticeListAll",
    "PvpRateGet",
    "PvpRateConfig",
    "PvpRateCommonProduceCard",
    "ExamSetting",
    "ProduceExamBattleScoreConfig",
    "ProduceCard",
    "ProduceItem",
    "ProduceExamGimmickEffectGroup",
    "StoryEvent",
    "CharacterDetail",
    "Achievement",
    "AchievementProgress",
    "EventLabel",
    "ProduceExamEffect",
    "ResultGradePattern",
    "GuildReaction",
    "ProduceDescriptionLabel",
  ], dirPath)
}

export async function getCsprt(dirPath: string): Promise<Csprt | null> {
  return await getAllJson([
    "SupportCard",
    "ProduceCard",
    "ProduceItem",
    "ProduceEventSupportCard",
    "ProduceStepEventDetail",
    "ProduceEffect",
    "SupportCardProduceSkillLevelDance",
    "SupportCardProduceSkillLevelVocal",
    "SupportCardProduceSkillLevelVisual",
    "SupportCardProduceSkillLevelAssist",
    "ProduceSkill",
    "ProduceTrigger",
    "ProduceExamEffect",
    "ProduceCardCustomize",
    "ProduceCardCustomizeRarityEvaluation",
    "ProduceCardGrowEffect",
    "ProduceDescriptionProduceCardGrowEffect",
    "ProduceCardStatusEnchant",
    "ProduceExamTrigger",
  ], dirPath)
}

export async function getCidol(dirPath: string): Promise<Cidol | null> {
  return await getAllJson([
    "IdolCard",
    "IdolCardSkin",
    "IdolCardLevelLimit",
    "IdolCardLevelLimitProduceSkill",
    "IdolCardLevelLimitStatusUp",
    "IdolCardPotential",
    "IdolCardPotentialProduceSkill",
    "ProduceCard",
    "ProduceItem",
    "ProduceSkill",
    "ProduceEffect",
    "ProduceStepAuditionDifficulty",
    "ProduceExamBattleNpcGroup",
    "ProduceExamBattleConfig",
    "ProduceExamBattleScoreConfig",
    "ProduceExamGimmickEffectGroup",
    "ProduceExamEffect",
    "ProduceGroup",
    "Produce",
  ], dirPath)
}

async function getAllJson<T extends (keyof UsedDB)[]>(keys: UnionArrayToTuple<T>, dirPath: string): Promise<MappedUsedDBTuple<T> | null> {
  const results = await Promise.all(keys.map(key => {
    return getJson(key, dirPath)
  }))
  if (!isNonNull(results)) {
    return null
  }
  return results as MappedUsedDBTuple<T>
}

async function getJson<T extends keyof UsedDB>(key: T, dirPath: string): Promise<UsedDB[T] | null> {
  const val = JSON.parse(readFileSync(path.join(dirPath, key + ".json"), { encoding: "utf-8" }))
  if (val) {
    return val as UsedDB[T]
  }
  console.error(`${key} is null in KV`)
  return null
}
