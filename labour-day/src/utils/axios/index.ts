import axios, {AxiosInstance, AxiosRequestConfig, AxiosRequestHeaders} from "axios";
import {getToken} from "./token";
import {useRouter} from "vue-router";

const axiosInstance = (options: AxiosRequestConfig) => {
    const defBaseURL = import.meta.env.VITE_APP_GLOB_BASE_API
    let instance = axios.create({
        baseURL: options.baseURL || defBaseURL,
        timeout: options.timeout || 1000 * 10,
    })
    setAxiosInstance(instance)
    return instance
}

function setAxiosInstance(instance: AxiosInstance) {
    // 请求拦截器
    instance.interceptors.request.use(
        resConfig => {
            // 防止缓存，给get请求加上时间戳
            if (resConfig.method === 'get') {
                resConfig.params = {...resConfig.params, t: new Date().getTime()}
            }

            // 处理不需要token的请求
            const WITHOUT_TOKEN_API = [{url: '/auth/login', method: 'POST'}]
            const isWithoutToken = WITHOUT_TOKEN_API.some((item: { url: string, method: string }) => item.url === resConfig.url && item.method === (resConfig.method as string).toUpperCase())
            if (isWithoutToken) {
                return resConfig
            }

            const token = getToken()
            if (token) {
                /**
                 * * jwt token
                 * ! 认证方案: Bearer
                 */
                if (!(resConfig.headers as AxiosRequestHeaders).Authorization) (resConfig.headers as AxiosRequestHeaders).Authorization = 'Bearer ' + token
                return resConfig
            }

            /**
             * * 未登录或者token过期的情况下
             * * 跳转登录页重新登录，携带当前路由及参数，登录成功会回到原来的页面
             */
                // const { currentRoute } = router
            const router = useRouter()
            router.replace({
                path: '/login',
                query: {...router.currentRoute.value.query, redirect: router.currentRoute.value.path},
            }).catch(() => {
                return Promise.reject({code: '-1', message: '未登录'})
            })
        },
        (err: any) => err,
    )

    instance.interceptors.response.use(
        // 因为我们接口的数据都在res.data下，所以我们直接返回res.data
        (response) => {
            // return Promise.resolve<Resp>(response?.data)
            return response
        },
        (error: any) => {
            switch (error.response?.status) {
                case 400:
                    error.message = '请求错误(400)';
                    break;
                case 401:
                    error.message = '未授权(401)';
                    break;
                case 403:
                    error.message = '拒绝访问(403)';
                    break;
                case 404:
                    error.message = '请求出错(404)';
                    break;
                case 408:
                    error.message = '请求超时(408)';
                    break;
                case 500:
                    error.message = '服务器错误(500)';
                    break;
                case 501:
                    error.message = '服务未实现(501)';
                    break;
                case 502:
                    error.message = '网络错误(502)';
                    break;
                case 503:
                    error.message = '服务不可用(503)';
                    break;
                case 504:
                    error.message = '网络超时(504)';
                    break;
                case 505:
                    error.message = 'HTTP版本不受支持(505)';
                    break;
                default:
                    error.message = `连接出错(${error.response?.status})!`;
            }
            console.log("error.message", error.message);
            let {code, message} = error.response?.data
            return Promise.reject({code, message})
        },
    )
}

export const request = axiosInstance({})



