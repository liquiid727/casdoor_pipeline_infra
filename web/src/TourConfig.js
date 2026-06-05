import React from "react";
import * as Setting from "./Setting";

export const TourObj = {
  home: [
    {
      title: "Welcome to Pipeline Auth",
      description: "You can learn more in the repository documentation.",
      cover: (
        <img
          alt="pipeline-auth.png"
          src={`${Setting.StaticBaseUrl}/img/logo.png`}
        />
      ),
    },
    {
      title: "Statistic cards",
      description: "Here are four statistic cards for user information.",
      id: "statistic",
    },
    {
      title: "Import users",
      description: "You can add new users or update existing users by uploading a XLSX file of user information.",
      id: "echarts-chart",
    },
  ],
  webhooks: [
    {
      title: "Webhook List",
      description: "Event systems allow you to build integrations, which subscribe to certain events on Pipeline Auth. When one of those events is triggered, a POST JSON payload is sent to the configured URL. Events include signup, login, logout, and user updates.",
    },
  ],
  syncers: [
    {
      title: "Syncer List",
      description: "Pipeline Auth stores users in the user table and provides syncers to help you import user data from external systems.",
    },
  ],
  sysinfo: [
    {
      title: "CPU Usage",
      description: "You can see the CPU usage in real time.",
      id: "cpu-card",
    },
    {
      title: "Memory Usage",
      description: "You can see the Memory usage in real time.",
      id: "memory-card",
    },
    {
      title: "API Latency",
      description: "You can see the usage statistics of each API latency in real time.",
      id: "latency-card",
    },
    {
      title: "API Throughput",
      description: "You can see the usage statistics of each API throughput in real time.",
      id: "throughput-card",
    },
    {
      title: "About Pipeline Auth",
      description: "You can get more Pipeline Auth information in this card.",
      id: "about-card",
    },
  ],
  subscriptions: [
    {
      title: "Subscription List",
      description: "Subscription helps to manage user's selected plan that make easy to control application's features access.",
    },
  ],
  pricings: [
    {
      title: "Price List",
      description: "Pipeline Auth can be used as a subscription management system via plan, pricing and subscription.",
    },
  ],
  plans: [
    {
      title: "Plan List",
      description: "Plans describe application features with their own names and prices. Plan features can be gated by roles and permissions.",
    },
  ],
  payments: [
    {
      title: "Payment List",
      description: "After the payment is successful, you can see the transaction information of the products in Payment, such as organization, user, purchase time, product name, etc.",
    },
  ],
  products: [
    {
      title: "Session List",
      description: "You can add the product (or service) you want to sell. The following will tell you how to add a product.",
    },
  ],
  sessions: [
    {
      title: "Session List",
      description: "You can get Session ID in this list.",
    },
  ],
  tokens: [
    {
      title: "Token List",
      description: "Pipeline Auth is based on OAuth. You can inspect users' access tokens in this list.",
    },
  ],
  enforcers: [
    {
      title: "Enforcer List",
      description: "In addition to enforcement APIs, Pipeline Auth also provides interfaces that help external applications obtain permission policy information.",
    },
  ],
  adapters: [
    {
      title: "Adapter List",
      description: "Pipeline Auth supports using the UI to connect adapters and manage policy rules.",
    },
  ],
  models: [
    {
      title: "Model List",
      description: "Model defines your permission policy structure, and how requests should match these permission policies and their effects. Then you can user model in Permission.",
    },
  ],
  permissions: [
    {
      title: "Permission List",
      description: "All users associated with a single organization are shared between that organization's applications. Use permissions when you need finer-grained access control.",
    },
    {
      title: "Permission Add",
      description: "In the Pipeline Auth Web UI, you can add a model for your organization and define policies in the permission configuration.",
      id: "add-button",
    },
    {
      title: "Permission Upload",
      description: "With Casbin Online Editor, you can generate model and policy files suitable for your usage scenarios and import them into Pipeline Auth.",
      id: "upload-button",
    },
  ],
  roles: [
    {
      title: "Role List",
      description: "Each user may have multiple roles. You can see the user's roles on the user's profile.",
    },
  ],
  resources: [
    {
      title: "Resource List",
      description: "You can upload resources in Pipeline Auth. Before uploading resources, configure a storage provider.",
    },
    {
      title: "Upload Resource",
      description: "Users can upload resources such as files and images to the previously configured cloud storage.",
      id: "upload-button",
    },
  ],
  providers: [
    {
      title: "Provider List",
      description: "We have 6 kinds of providers:OAuth providers、SMS Providers、Email Providers、Storage Providers、Payment Provider、Captcha Provider.",
    },
    {
      title: "Provider Add",
      description: "You must add the provider to application, then you can use the provider in your application",
      id: "add-button",
    },
  ],
  organizations: [
    {
      title: "Organization List",
      description: "Organization is the basic unit of Pipeline Auth, which manages users and applications. If a user signs in to an organization, they can access all of that organization's applications without signing in again.",
    },
  ],
  groups: [
    {
      title: "Group List",
      description: "In the groups list pages, you can see all the groups in organizations.",
    },
  ],
  users: [
    {
      title: "User List",
      description: "As an authentication platform, Pipeline Auth is able to manage users.",
    },
    {
      title: "Import users",
      description: "You can add new users or update existing users by uploading a XLSX file of user information.",
      id: "upload-button",
    },
  ],
  applications: [
    {
      title: "Application List",
      description: "If you want to use Pipeline Auth to provide login service for your web apps, you can add them as applications. Users can access all applications in their organizations without logging in twice.",
    },
  ],
};

export const TourUrlList = ["home", "organizations", "groups", "users", "applications", "providers", "resources", "roles", "permissions", "models", "adapters", "enforcers", "tokens", "sessions", "products", "payments", "plans", "pricings", "subscriptions", "sysinfo", "syncers", "webhooks"];

export function getNextUrl(pathName = window.location.pathname) {
  return TourUrlList[TourUrlList.indexOf(pathName.replace("/", "")) + 1] || "";
}

let orgIsTourVisible = true;

export function setOrgIsTourVisible(visible) {
  orgIsTourVisible = visible;
  if (orgIsTourVisible === false) {
    setIsTourVisible(false);
  }
}

export function setIsTourVisible(visible) {
  localStorage.setItem("isTourVisible", visible);
  window.dispatchEvent(new Event("storageTourChanged"));
}

export function setTourLogo(tourLogoSrc) {
  if (tourLogoSrc !== "") {
    TourObj["home"][0]["cover"] = (<img alt="casdoor.png" src={tourLogoSrc} />);
  }
}

export function getTourVisible() {
  return localStorage.getItem("isTourVisible") !== "false";
}

export function getNextButtonChild(nextPathName) {
  return nextPathName !== "" ?
    `Go to "${nextPathName.charAt(0).toUpperCase()}${nextPathName.slice(1)} List"`
    : "Finish";
}

export function getSteps() {
  const path = window.location.pathname.replace("/", "");
  const res = TourObj[path];
  if (res === undefined) {
    return [];
  } else {
    return res;
  }
}
