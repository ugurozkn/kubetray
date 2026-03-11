const os = require("os");
const fs = require("fs");
const path = require("path");
const https = require("https");
const { execSync } = require("child_process");

const VERSION = require("./package.json").version;
const REPO = "ugurozkn/kubetray";

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();

  const osMap = { darwin: "darwin", linux: "linux" };
  const archMap = { x64: "amd64", arm64: "arm64" };

  const mappedOS = osMap[platform];
  const mappedArch = archMap[arch];

  if (!mappedOS || !mappedArch) {
    throw new Error(`Unsupported platform: ${platform}/${arch}`);
  }

  return { os: mappedOS, arch: mappedArch };
}

function download(url) {
  return new Promise((resolve, reject) => {
    https.get(url, (res) => {
      if (res.statusCode === 302 || res.statusCode === 301) {
        return download(res.headers.location).then(resolve).catch(reject);
      }
      if (res.statusCode !== 200) {
        return reject(new Error(`Download failed: ${res.statusCode}`));
      }
      const chunks = [];
      res.on("data", (chunk) => chunks.push(chunk));
      res.on("end", () => resolve(Buffer.concat(chunks)));
      res.on("error", reject);
    }).on("error", reject);
  });
}

async function main() {
  const { os: pOS, arch: pArch } = getPlatform();
  const name = `kubetray_${pOS}_${pArch}`;
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${name}.tar.gz`;

  console.log(`Downloading kubetray v${VERSION} for ${pOS}/${pArch}...`);

  const tarball = await download(url);
  const tarPath = path.join(__dirname, "kubetray.tar.gz");
  fs.writeFileSync(tarPath, tarball);

  execSync(`tar -xzf kubetray.tar.gz`, { cwd: __dirname });
  fs.unlinkSync(tarPath);
  fs.chmodSync(path.join(__dirname, "kubetray"), 0o755);

  console.log("kubetray installed successfully.");
}

main().catch((err) => {
  console.error("Failed to install kubetray:", err.message);
  process.exit(1);
});
