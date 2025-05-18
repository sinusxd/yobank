import api from "../api";
import { AxiosResponse } from "axios";

export interface Rate {
    currency: string;
    value: number;
    date: string;
}

export default class RateService {

    static async getRate(currency: string): Promise<AxiosResponse<Rate>> {
       return  api.get<Rate>(`/api/v1/rates/${currency}`);
    }

    static async getRatesHistory(currency: string, from: string, to: string): Promise<AxiosResponse<Rate[]>> {
        return  api.get<Rate[]>(`/api/v1/rates/${currency}/history`, {
            params: { from, to }
        });
    }
}
