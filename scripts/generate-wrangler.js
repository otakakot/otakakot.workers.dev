const fs = require("fs");

const name = process.env.DATABASE_NAME;
const id = process.env.DATABASE_ID;

if (!name || !id) {
  console.error("DATABASE_NAME and DATABASE_ID are required");
  process.exit(1);
}

const template = fs.readFileSync("wrangler.jsonc", "utf-8");
const generated = template
  .replace(/<DATABASE_NAME>/g, name)
  .replace(/<DATABASE_ID>/g, id);

fs.writeFileSync("wrangler.generated.jsonc", generated);
console.log("Generated wrangler.generated.jsonc");
