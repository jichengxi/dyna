const httpsRE = /^https:\/\//
export function createProxy(proxy: any) {
  let proxyList = proxy as Array<Array<string>>
  const ret: any = {}
  console.log("proxyList", proxyList)
  for (const [prefix, target] of proxyList) {
    const isHttps = httpsRE.test(target)

    // https://github.com/http-party/node-http-proxy#options
    ret[prefix] = {
      target: target,
      changeOrigin: true,
      ws: true,
      rewrite: (path: string) => path.replace(new RegExp(`^${prefix}`), ''),
      // https is require secure=false
      ...(isHttps ? { secure: false } : {}),
    }
  }
  return ret
}
