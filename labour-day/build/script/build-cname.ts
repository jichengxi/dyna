import chalk from 'chalk'
import { writeFileSync } from 'fs-extra'
import { OUTPUT_DIR } from '../constant'
import { getEnvConfig, getRootPath } from '../utils'

export function runBuildCNAME() {
  const envConfig = getEnvConfig()
  console.log("envConfig", envConfig)
  if (!envConfig['VITE_APP_GLOB_CNAME']) return
  try {
    writeFileSync(getRootPath(`${OUTPUT_DIR}/CNAME`), envConfig['VITE_APP_GLOB_CNAME'])
  } catch (error) {
    console.log(chalk.red('CNAME file failed to package:\n' + error))
  }
}
