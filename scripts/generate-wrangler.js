const fs = require("fs");

const dbName = process.env.DATABASE_NAME;
const dbId = process.env.DATABASE_ID;
const kvId = process.env.KV_ID;

if (!dbName || !dbId) {
  console.error("DATABASE_NAME and DATABASE_ID are required");
  process.exit(1);
}

if (!kvId) {
  console.error("KV_ID is required");
  process.exit(1);
}

const template = fs.readFileSync("wrangler.jsonc", "utf-8");
const generated = template
  .replace(/<DATABASE_NAME>/g, dbName)
  .replace(/<DATABASE_ID>/g, dbId)
  .replace(/<KV_ID>/g, kvId);

fs.writeFileSync("wrangler.generated.jsonc", generated);

console.log("Generated wrangler.generated.jsonc");
