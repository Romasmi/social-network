import http from "k6/http";
import { check, group, sleep } from "k6";

const apiURL = 'http://localhost:8000';
const countOfUsersToGet = 5;

/*
export const options = {
    stages: [
        { duration: '5s', target: 100 },
        { duration: '10s', target: 200 },
        { duration: '5s', target: 300 },
    ],
    thresholds: {
        http_req_failed: ['rate<0.05'],
    },
};
*/

export default function () {
    /**
     * Scenario: search users and check profile of 5 random found users
     */
    group('User Search Flow', () => {
        const users = searchUser();
        sleep(0.5);
        getRandomUsers(users);
    });
}

const getRandomUsers = (users) => {
    const maxFound = Math.min(countOfUsersToGet, users.length)
    for (let i = 0; i < maxFound;  i++) {
        getUser(users[Math.round(Math.random() * users.length - 1)].id)
        sleep(0.5);
    }
}

const getUser = (userId) => {
    const res = http.get(`${apiURL}/user/get/${userId}`, {
        headers: {
            Accept: "application/json",
            Authorization: "Bearer " + __ENV.TOKEN,
        },
    });

    check(res, {
        "status is 200": (r) => r.status === 200,
    });
}

const searchUser = () => {
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

    return res.json()
}
