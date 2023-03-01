import axios from "axios";

const baseURL = "http://localhost:8000/v1/affiliates/99b6b56f-e692-425d-b1f7-16150acdaa90"
const config = { headers: { "x-api-key": "50145a3f-97cc-4fd7-be53-5b978b430936" } }

export const SendFile = (file) => {
    const data = new FormData();
    data.append('attachment', file);
    return axios.post(`${baseURL}/batches`, data, config)
}

export const GetFiles = () => {
    return axios.get(`${baseURL}/batches`, config)
}

export const GetTransactions = () => {
    return axios.get(`${baseURL}/transactions`, config)
}
