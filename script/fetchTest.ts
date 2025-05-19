import * as zlib from "zlib";
import { promisify } from "util";
import { exit } from "process";

const kvEndpoint = process.env.CAMPUS_DB_PUT_URL
const kvApiKey = process.env.CAMPUS_DB_PUT_SECRET

if (kvEndpoint === undefined || kvApiKey === undefined) {
  console.error("env CAMPUS_DB_PUT_URL or CAMPUS_DB_PUT_SECRET is undefined")
  exit(1)
}

const gunzip = promisify(zlib.gunzip);

async function fetchTest() {
  const resp = await fetch(`${kvEndpoint}/GXSupportCard`, {
    headers: {
      "Authorization": `Bearer ${kvApiKey}`,
    },
  })
  console.log(resp.status)
  const data = await gunzip(await resp.arrayBuffer())

  console.log(data.toString())
}

await fetchTest()
