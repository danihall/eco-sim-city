import * as fs from "node:fs";

const regex_src = /src\//;
const dist = "dist/";

fs.promises
  .readdir("./src/assets", {
    recursive: true,
    withFileTypes: true,
  })
  .then((result) => Promise.all(result.filter(_PNGFile).map(_copyToDist)))
  .then(console.log)
  .catch(console.error);

/**
 * @param {fs.Dirent} entry
 */
function _PNGFile(entry) {
  return entry.isFile() && entry.name.endsWith(".png");
}

/**
 * @param {fs.Dirent} file
 */
function _copyToDist(file) {
  const src = `${file.parentPath}/${file.name}`;
  const dest = src.replace(regex_src, dist);

  return fs.promises.copyFile(src, dest).then(() => _success(file.name, dest));
}

/**
 * @param {string} file_name
 * @param {string} destination
 */
function _success(file_name, destination) {
  return `copied "${file_name}" to "${destination}"`;
}
