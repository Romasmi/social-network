import http from "k6/http";
import { sleep } from "k6";
import { check } from "k6";

export default function () {
    const url = "http://localhost:8000/user/search?first_name=%D0%90%D0%B1%D1%80%D0%B0&last_name=%D0%A2%D0%B8%D0%BC";

    const params = {
        headers: {
            Accept: "application/json",
            Authorization: "Bearer " + __ENV.TOKEN,
        },
    };

    const res = http.get(url, params);

    check(res, {
        "status is 200": (r) => r.status === 200,
    });

    sleep(1);
}
