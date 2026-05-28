import * as Setting from "../Setting";

export function getChannelRelations(owner, parentOrganization = "", page = "", pageSize = "", field = "", value = "", sortField = "", sortOrder = "") {
  return fetch(`${Setting.ServerUrl}/api/get-channel-relations?owner=${owner}&parentOrganization=${encodeURIComponent(parentOrganization)}&p=${page}&pageSize=${pageSize}&field=${field}&value=${value}&sortField=${sortField}&sortOrder=${sortOrder}`, {
    method: "GET",
    credentials: "include",
    headers: {
      "Accept-Language": Setting.getAcceptLanguage(),
    },
  }).then(res => res.json());
}

export function addChannelRelation(relation) {
  return fetch(`${Setting.ServerUrl}/api/add-channel-relation`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(Setting.deepCopy(relation)),
    headers: {
      "Accept-Language": Setting.getAcceptLanguage(),
    },
  }).then(res => res.json());
}

export function updateChannelRelation(owner, name, relation) {
  return fetch(`${Setting.ServerUrl}/api/update-channel-relation?id=${owner}/${encodeURIComponent(name)}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(Setting.deepCopy(relation)),
    headers: {
      "Accept-Language": Setting.getAcceptLanguage(),
    },
  }).then(res => res.json());
}

export function deleteChannelRelation(relation) {
  return fetch(`${Setting.ServerUrl}/api/delete-channel-relation`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(Setting.deepCopy(relation)),
    headers: {
      "Accept-Language": Setting.getAcceptLanguage(),
    },
  }).then(res => res.json());
}
