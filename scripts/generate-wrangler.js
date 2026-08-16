const fs = require("fs");

const dbName = process.env.DATABASE_NAME;
const dbId = process.env.DATABASE_ID;
const kvId = process.env.KV_ID;
const corsAllowOrigins = process.env.CORS_ALLOW_ORIGINS;

if (!dbName || !dbId) {
  console.error("DATABASE_NAME and DATABASE_ID are required");
  process.exit(1);
}

if (!kvId) {
  console.error("KV_ID is required");
  process.exit(1);
}

const template = fs.readFileSync("wrangler.env.jsonc", "utf-8");
const generated = template
  .replace(/<DATABASE_NAME>/g, dbName)
  .replace(/<DATABASE_ID>/g, dbId)
  .replace(/<KV_ID>/g, kvId)
  .replace(/<CORS_ALLOW_ORIGINS>/g, corsAllowOrigins);

fs.writeFileSync("wrangler.jsonc", generated);

console.log("Generated wrangler.jsonc");
