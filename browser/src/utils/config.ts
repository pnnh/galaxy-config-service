interface viteEnv {
  BASE_URL: string,
  MODE: string,
  DEV: boolean,
  PROD: boolean,
  SSR: boolean
}

const importMeta = (import.meta as any).env as viteEnv
let jsEnv = ''

if (importMeta.DEV) {
  jsEnv = 'dev'
} else {
  jsEnv = 'release'
}

export function isDev (): boolean {
  return jsEnv === 'dev'
}

export function isRelease (): boolean {
  return jsEnv === 'release'
}

// export function getApiUrl (): string {
//   if (isDev()) {
//     return 'http://127.0.0.1:8080'
//   }
//   return 'https://sfx.xyz'
// }

export function getHost (): string {
  console.log('getHost', isDev())
  if (isDev()) {
    return 'http://127.0.0.1:3000'
  }
  return 'https://res.sfx.xyz'
}

const ApiUrl = {
  home: '/',
  account: {
    personal: '/account/personal',
    edit: '/account/edit'
  }
}

export { ApiUrl }
