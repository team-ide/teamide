import http from '@/server/http';

let power = {
    data(param) {
        return http.post('api/power/data', param);
    },
};


export default power;