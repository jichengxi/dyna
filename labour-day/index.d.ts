declare module '*.json' {
    const value: any;
    export default value;
}

declare interface Resp {
    code: number
    message: string
    data: []|{}
}
