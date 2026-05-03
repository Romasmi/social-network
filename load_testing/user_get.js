import http from "k6/http";
import { check, sleep } from "k6";
import { SharedArray } from "k6/data";

const data = new SharedArray("users", function () {
    return JSON.parse(open("./1000_random_users.json"));
});

export const options = {
    vus: 10,
    duration: '30s',
};

export default function () {
    const user = data[Math.floor(Math.random() * data.length)];
    const url = `${__ENV.BASE_URL || "http://api.localhost"}/user/get/${user.user_id}`;

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
