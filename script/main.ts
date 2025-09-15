import { exit } from "process";
import { getXIdolCard } from "./cidol";
import { getCidol, getCsprt, getMaster, getMemoryInspector, getPCard } from "./kv";
import * as zlib from "zlib";
import { promisify } from "util";
import { getXSupportCard } from "./csprt";
import { getXCustProduceCards } from "./pcard";
import { getXMemoryInspector } from "./memoryInspector";
import path from "path";
import { getXMaster } from "./master";

const kvEndpoint = process.env.CAMPUS_DB_PUT_URL
const kvApiKey = process.env.CAMPUS_DB_PUT_SECRET
const cacheDir = process.env.CAMPUS_CACHE_DIR

if (kvEndpoint === undefined || kvApiKey === undefined || cacheDir === undefined) {
  console.error("env CAMPUS_DB_PUT_URL or CAMPUS_DB_PUT_SECRET or CAMPUS_CACHE_DIR is undefined")
  exit(1)
}

const gzip = promisify(zlib.gzip);

async function putToKv(key: string, rawJSONString: string) {
  const compressed = await gzip(rawJSONString)
  const formData = new FormData()
  formData.append("metadata", "{}")
  formData.append("value", new Blob([compressed]))
  const resp = await fetch(`${kvEndpoint}/${key}`, {
    method: "put",
    headers: {
      "Authorization": `Bearer ${kvApiKey}`,
    },
    body: formData,
  })
  if (resp.status != 200) {
    console.error("failed putting data to kv: ", resp.status, resp.statusText)
    console.error(await resp.text())
    exit(1)
  } else {
    console.log(`DB "${key}" is successfully put to kv`)
  }
}

async function updateDB<
  T extends (_: string) => Promise<any>,
  U extends (_: any) => any,
>(
  key: string,
  getDbFunc: T,
  getObjectFunc: U,
) {
  const db = await getDbFunc(path.join(cacheDir!, "masterJson"))
  if (!db) {
    console.error(`failed reading ${key} json files`)
    exit(1)
  }
  const xobject = getObjectFunc(db!)
  putToKv(key, JSON.stringify(xobject))
}

async function run() {
  await updateDB("GXIdolCard", getCidol, getXIdolCard)
  await updateDB("GXSupportCard", getCsprt, getXSupportCard)
  await updateDB("GXProduceCard", getPCard, getXCustProduceCards)
  await updateDB("GXMemory", getMemoryInspector, getXMemoryInspector)
  await updateDB("GXMaster", getMaster, getXMaster)
}
await run()
