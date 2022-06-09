export function resolveToken(authorization: string) {
  /**
   * * jwt token
   * * Bearer + token
   * ! 认证方案: Bearer
   */
  const reqTokenSplit = authorization.split(' ')
  console.log("authorization", authorization)
  if (reqTokenSplit.length === 2) {
    return reqTokenSplit[1]
  }
  return ''
}
