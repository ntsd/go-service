import urlencode from "https://jslib.k6.io/form-urlencoded/3.0.0/index.js";
import http from "k6/http";
import { b64encode } from "k6/encoding";
import { check, fail } from "k6";

const api = "http://localhost:8080/v1";

export const options = {
  vus: 10,
  duration: "15s",
};

function oauth2() {
  const client_id = randomString(32);
  const client_secret = randomString(32);
  const body = {
    client_id,
    client_secret,
  };

  const res = http.post(`${api}/oauth/clients`, JSON.stringify(body), {
    headers: {
      "Content-Type": "application/json",
    },
  });
  if (!check(res, { "status was 200": (r) => r.status == 200 })) {
    fail(`status code was *not* 200 but ${res.status}`);
  }

  const basicAuth = b64encode(`${client_id}:${client_secret}`);
  const res2 = http.post(
    `${api}/oauth/token`,
    urlencode({ grant_type: "client_credentials" }),
    {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
        Authorization: `Basic ${basicAuth}`,
      },
    }
  );
  if (!check(res2, { "status was 200": (r) => r.status == 200 })) {
    fail(`status code was *not* 200 but ${res2.status}`);
  }

  let data = {
    access_token: JSON.parse(res2.body).access_token,
  };

  data.userId = loadCreateUser(data);

  return data;
}

export function setup() {
  return oauth2();
}

function loadUnauthorized(data) {
  let res = http.get(api);
  check(res, { "status was 401": (r) => r.status == 401 });
}

function loadCreateUser(data) {
  const res = http.post(
    `${api}/users`,
    JSON.stringify({
      email: `${randomString(20)}@example.com`,
      name: "test",
    }),
    {
      headers: {
        "Content-Type": "application/json",
        authorization: `Bearer ${data.access_token}`,
      },
    }
  );
  check(res, { "status was 200": (r) => r.status == 200 });
  return JSON.parse(res.body).id;
}

function loadGetAllUsers(data) {
  const res = http.get(`${api}/users`, {
    headers: {
      "Content-Type": "application/json",
      authorization: `Bearer ${data.access_token}`,
    },
  });
  check(res, { "status was 200": (r) => r.status == 200 });
}

function loadGetOneUser(data) {
  const res = http.get(`${api}/users/${data.userId}`, {
    headers: {
      "Content-Type": "application/json",
      authorization: `Bearer ${data.access_token}`,
    },
  });
  check(res, { "status was 200": (r) => r.status == 200 });
}

export default function (data) {
  // loadUnauthorized(data)
  // loadCreateUser(data)
  // loadGetAllUsers(data);
  loadGetOneUser(data);
}

function randomString(length) {
  const characters = "0123456789abcdefghijklmnopqrstuvwxyz";
  let result = "";
  for (let i = length; i > 0; i--) {
    result += characters[Math.floor(Math.random() * characters.length)];
  }
  return result;
}
