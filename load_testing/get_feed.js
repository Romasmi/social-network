import http from "k6/http";
import { sleep } from "k6";
import { check } from "k6";

export default function () {
    const url = `${__ENV.BASE_URL || "http://api.localhost"}/post/feed?offset=0&limit=20`;

    const params = {
        headers: {
            "Accept": "application/json",
            "Authorization": "Bearer " + __ENV.TOKEN,
        },
    };

    const res = http.get(url, params);

    check(res, {
        "status is 200": (r) => r.status === 200,
    });

    sleep(1);
}