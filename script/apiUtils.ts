type KeyValueTuple<T, K extends keyof T> = [K, T[K]]

export function filterItems<T>(
  instances: T[],
  instanceKey: keyof T,
  matchValue: string | string[],
  options?: {
    sortRules?: [keyof T, boolean],
    limitRules?: KeyValueTuple<T, keyof T>
  },
): T[] {
  const sortRules = options?.sortRules
  const limitRules = options?.limitRules
  let filtered = instances
  if (limitRules) {
    filtered = filtered.filter(instance => instance[limitRules[0]] <= limitRules[1])
  }
  if (typeof matchValue === "string") {
    filtered = filtered.filter(instance => instance[instanceKey] === matchValue)
  } else {
    filtered = matchValue.map(val => {
      return filtered.filter(instance => instance[instanceKey] === val)
    }).flat()
  }

  if (sortRules) {
    const sortKey = sortRules[0]
    const ascending = sortRules[1]
    filtered.sort((a, b) =>
      ascending
        ? (+a[sortKey]) - (+b[sortKey])
        : (+b[sortKey]) - (+a[sortKey])
    )
  }
  return filtered
}
