import * as f from 'fs'
import * as p from 'path'
import * as d from 'dotenv'

export function wrapperEnv(envOptions: Record<string, string>) {
  if (!envOptions) return {}
  const ret = {}

  for (const key in envOptions) {
    let val = envOptions[key]
    // if (['true', 'false'].includes(val)) {
    //   val = val === 'true'
    // }
    // if (['VITE_PORT'].includes(key)) {
    //   val = +val
    // }
    if (key === 'VITE_PROXY' && val) {
      try {
        val = val.replace(/'/g, '"')
      } catch (error) {
        val = ''
      }
    }
    // ret[key] = val
    process.env[key] = val
    // if (typeof key === 'string') {
    //   process.env[key] = val
    // } else if (typeof key === 'object') {
    //   process.env[key] = JSON.stringify(val)
    // }
  }
  // return ret
}

/**
 * 获取当前环境下生效的配置文件名
 */
function getConfFiles() {
  const script = process.env.npm_lifecycle_script
  const reg = new RegExp('--mode ([a-z_\\d]+)')
  const result = reg.exec(script as string)
  console.log("result: ", result)
  if (result) {
    const mode = result[1]
    return ['.env', '.env.local', `.env.${mode}`]
  }
  return ['.env', '.env.local', '.env.production']
}

export function getEnvConfig(match = 'VITE_APP_GLOB_', confFiles = getConfFiles()) {
  let envConfig = {}
  confFiles.forEach((item) => {
    try {
      if (f.existsSync(p.resolve(process.cwd() + "/env", item))) {
        const env = d.parse(f.readFileSync(p.resolve(process.cwd() + "/env", item)))
        envConfig = { ...envConfig, ...env }
      }
    } catch (e) {
      console.error(`Error in parsing ${item}`, e)
    }
  })
  const reg = new RegExp(`^(${match})`)
  Object.keys(envConfig).forEach((key) => {
    if (!reg.test(key)) {
      Reflect.deleteProperty(envConfig, key)
    }
  })
  return envConfig
}

export function getRootPath(...dir: string[]) {
  return p.resolve(process.cwd(), ...dir)
}
