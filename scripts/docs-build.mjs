// Regenerates the embedded Redoc page (internal/api/docs/static) from the
// OpenAPI spec at docs/openapi.yaml and vendors the exact Redoc bundle version
// the page references. Run via: make docs-build
//
// We deliberately do NOT serve redocly `build-docs`' SSR output: that page
// relies on React hydration, which throws the acknowledged upstream hydration
// errors (#418/#423) and is fragile when self-hosted. Instead we use it only to
// resolve the spec to a self-contained JSON and to learn which bundle version
// to vendor, then emit a minimal page that calls Redoc.init(spec) — the
// canonical, hydration-free embedding mode. The result is fully offline, has no
// console errors, and needs no SRI/crossorigin attributes.
import { execSync } from "node:child_process";
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";

const outDir = "internal/api/docs/static";
const htmlPath = `${outDir}/index.html`;
const bundlePath = `${outDir}/redoc.standalone.js`;

mkdirSync(outDir, { recursive: true });

execSync(
  'npx --yes @redocly/cli build-docs docs/openapi.yaml -o /tmp/build-docs.html ' +
    '--title "Flood Response API" --disableGoogleFont',
  { stdio: "inherit" },
);

const buildHtml = readFileSync("/tmp/build-docs.html", "utf8");

const versionMatch = buildHtml.match(
  /https:\/\/cdn\.redocly\.com\/redoc\/(v[0-9.]+)\/bundles\/redoc\.standalone\.js/,
);
if (!versionMatch) {
  throw new Error("no redoc.standalone.js CDN reference found in build-docs output");
}

const stateMatch = buildHtml.match(/const __redoc_state = (\{.*\});/);
if (!stateMatch) {
  throw new Error("could not extract __redoc_state from build-docs output");
}
const spec = JSON.parse(stateMatch[1]).spec.data;

const html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Flood Response API</title>
  <style>body{margin:0;padding:0}</style>
</head>
<body>
  <div id="redoc"></div>
  <script src="/docs/redoc.standalone.js"></script>
  <script>
    var spec = ${JSON.stringify(spec)};
    Redoc.init(spec, {
      title: "Flood Response API",
      disableGoogleFont: true,
      hideDownloadButton: true
    }, document.getElementById("redoc"));
  </script>
</body>
</html>
`;

writeFileSync(htmlPath, html);

const response = await fetch(
  `https://cdn.redocly.com/redoc/${versionMatch[1]}/bundles/redoc.standalone.js`,
);
if (!response.ok) {
  throw new Error(`failed to download redoc bundle: HTTP ${response.status}`);
}
writeFileSync(bundlePath, Buffer.from(await response.arrayBuffer()));

console.log(
  `docs-build: wrote ${htmlPath} and vendored redoc ${versionMatch[1]} in ${bundlePath}`,
);