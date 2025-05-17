import api from "../api";
import {EmailCodeRequest} from "@/api/models/request/emailCodeRequest.ts";
import {EmailCodeResponse} from "@/api/models/response/emailCodeResponse.ts";
import {AxiosResponse} from "axios"
import {VerifyCodeRequest} from "@/api/models/request/verifyCodeRequest.ts";
import {VerifyCodeResponse} from "@/api/models/response/verifyCodeResponse.ts";
import {TelegramLoginRequest} from "@/api/models/request/telegramLoginRequest.ts";
import {TelegramLoginResponse} from "@/api/models/response/telegramLoginResponse.ts";

export default class AuthService {
    static async requestCode(request: EmailCodeRequest): Promise<AxiosResponse<EmailCodeResponse>> {
        return api.post<EmailCodeResponse>('/api/v1/auth/email/request-code', request);
    }

    static async verifyCode(request: VerifyCodeRequest): Promise<AxiosResponse<VerifyCodeResponse>> {
        return api.post<VerifyCodeResponse>('/api/v1/auth/email/verify-code', request)
    }

    static async loginWithTelegram(request: TelegramLoginRequest): Promise<AxiosResponse<TelegramLoginResponse>> {
        return api.post<TelegramLoginResponse>('api/v1/auth/telegram/login', request)
    }
}
