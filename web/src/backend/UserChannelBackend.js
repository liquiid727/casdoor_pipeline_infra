import * as Setting from "../Setting";

export function getUserChannels(owner = "admin", rootOrganization = "", user = "", page = "", pageSize = "", field = "", value = "", sortField = "", sortOrder = "") {
  return fetch(`${Setting.ServerUrl}/api/get-user-channels?owner=${owner}&rootOrganization=${encodeURIComponent(rootOrganization)}&user=${encodeURIComponent(user)}&p=${page}&pageSize=${pageSize}&field=${field}&value=${value}&sortField=${sortField}&sortOrder=${sortOrder}`, {
    method: "GET",
    credentials: "include",
    headers: {
      "Accept-Language": Setting.getAcceptLanguage(),
    },
  }).then(res => res.json());
}

export function bindUserChannel(binding) {
  return fetch(`${Setting.ServerUrl}/api/bind-user-channel`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(Setting.deepCopy(binding)),
    headers: {
      "Accept-Language": Setting.getAcceptLanguage(),
    },
  }).then(res => res.json());
}

export function updateUserChannel(owner, name, binding) {
  return fetch(`${Setting.ServerUrl}/api/update-user-channel?id=${owner}/${encodeURIComponent(name)}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(Setting.deepCopy(binding)),
    headers: {
      "Accept-Language": Setting.getAcceptLanguage(),
    },
  }).then(res => res.json());
}

export function unbindUserChannel(binding) {
  return fetch(`${Setting.ServerUrl}/api/unbind-user-channel`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(Setting.deepCopy(binding)),
    headers: {
      "Accept-Language": Setting.getAcceptLanguage(),
    },
  }).then(res => res.json());
}
