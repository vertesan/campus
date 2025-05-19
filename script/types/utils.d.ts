import { UsedDB } from "."

export type MappedUsedDBTuple<T extends (keyof UsedDB)[]> = { [P in keyof T]: UsedDB[T[P]] }
export type UnionArrayToTuple<T extends (keyof UsedDB)[]> = { [P in keyof T]: T[P] }
export type UnArray<T> = T extends Array<infer U> ? U : T
