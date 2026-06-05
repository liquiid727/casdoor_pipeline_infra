#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const yaml = require(path.resolve(__dirname, "../../web/node_modules/js-yaml"));

const repoRoot = path.resolve(__dirname, "../..");
const mode = process.argv[2] || "all";

const files = {
  router: path.join(repoRoot, "routers/router.go"),
  swagger: path.join(repoRoot, "swagger/swagger.yml"),
  management: path.join(repoRoot, "web/src/ManagementPage.js"),
  backendDir: path.join(repoRoot, "web/src/backend"),
  referenceBackendDir: path.join(repoRoot, "docs/reference/backend-api"),
  referenceFrontendDir: path.join(repoRoot, "docs/reference/frontend"),
};

function read(filePath) {
  return fs.readFileSync(filePath, "utf8");
}

function ensureDir(dirPath) {
  fs.mkdirSync(dirPath, {recursive: true});
}

function write(filePath, content) {
  ensureDir(path.dirname(filePath));
  fs.writeFileSync(filePath, content);
}

function titleCase(text) {
  return text
    .replace(/[_-]+/g, " ")
    .replace(/\s+/g, " ")
    .trim()
    .split(" ")
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

function humanizeAction(action) {
  if (!action) {
    return "Controller-defined action";
  }

  const pairs = [
    ["Get", "Get "],
    ["Add", "Add "],
    ["Update", "Update "],
    ["Delete", "Delete "],
    ["Upload", "Upload "],
    ["Sync", "Sync "],
    ["Run", "Run "],
    ["Send", "Send "],
    ["Handle", "Handle "],
    ["Grant", "Grant "],
    ["Revoke", "Revoke "],
    ["Refresh", "Refresh "],
    ["Check", "Check "],
    ["Verify", "Verify "],
    ["Set", "Set "],
    ["Cancel", "Cancel "],
  ];

  let value = action;
  for (const [from, to] of pairs) {
    value = value.replace(new RegExp(`^${from}`), to);
  }
  value = value
    .replace(/([a-z0-9])([A-Z])/g, "$1 $2")
    .replace(/\s+/g, " ")
    .trim();

  return value.charAt(0).toUpperCase() + value.slice(1);
}

function uniqueSorted(values) {
  return [...new Set(values)].sort((a, b) => a.localeCompare(b));
}

function inferDomain(pathname, action) {
  const matchers = [
    [["/.well-known", "/cas/"], "protocol"],
    [["/api/login", "/api/signup", "/api/logout", "/api/sso-logout", "/api/callback", "/api/device-auth", "/api/cancel-device-auth", "/api/device-auth-complete", "/api/native-sso-complete", "/api/kerberos-login", "/api/get-saml-login", "/api/acs", "/api/saml/", "/api/webauthn/", "/api/mfa/", "/api/grant-consent", "/api/revoke-consent", "/api/faceid-signin-begin"], "auth"],
    [["/api/get-account", "/api/userinfo", "/api/user", "/api/unlink"], "account"],
    [["/api/get-organizations", "/api/get-organization", "/api/update-organization", "/api/add-organization", "/api/delete-organization", "/api/get-default-application", "/api/get-organization-names"], "organization"],
    [["/api/get-groups", "/api/get-group", "/api/update-group", "/api/add-group", "/api/delete-group", "/api/upload-groups"], "group"],
    [["/api/get-global-users", "/api/get-users", "/api/get-sorted-users", "/api/get-user-count", "/api/get-user", "/api/update-user", "/api/add-user", "/api/delete-user", "/api/upload-users", "/api/remove-user-from-group", "/api/verify-identification", "/api/impersonate-user", "/api/exit-impersonate-user"], "user"],
    [["/api/get-invitations", "/api/get-invitation", "/api/get-invitation-info", "/api/update-invitation", "/api/add-invitation", "/api/delete-invitation", "/api/verify-invitation", "/api/send-invitation"], "invitation"],
    [["/api/get-applications", "/api/get-application", "/api/get-user-application", "/api/get-organization-applications", "/api/update-application", "/api/add-application", "/api/delete-application"], "application"],
    [["/api/get-providers", "/api/get-provider", "/api/get-global-providers", "/api/update-provider", "/api/add-provider", "/api/delete-provider"], "provider"],
    [["/api/get-resources", "/api/get-resource", "/api/update-resource", "/api/add-resource", "/api/delete-resource", "/api/upload-resource"], "resource"],
    [["/api/get-certs", "/api/get-global-certs", "/api/get-cert", "/api/update-cert", "/api/add-cert", "/api/delete-cert", "/api/update-cert-domain-expire"], "cert"],
    [["/api/get-keys", "/api/get-global-keys", "/api/get-key", "/api/update-key", "/api/add-key", "/api/delete-key"], "key"],
    [["/api/get-roles", "/api/get-role", "/api/update-role", "/api/add-role", "/api/delete-role", "/api/upload-roles"], "role"],
    [["/api/get-permissions", "/api/get-permissions-by-submitter", "/api/get-permissions-by-role", "/api/get-permission", "/api/update-permission", "/api/add-permission", "/api/delete-permission", "/api/upload-permissions"], "permission"],
    [["/api/get-models", "/api/get-model", "/api/update-model", "/api/add-model", "/api/delete-model"], "model"],
    [["/api/get-adapters", "/api/get-adapter", "/api/update-adapter", "/api/add-adapter", "/api/delete-adapter", "/api/get-policies", "/api/get-filtered-policies", "/api/update-policy", "/api/add-policy", "/api/remove-policy"], "adapter"],
    [["/api/get-enforcers", "/api/get-enforcer", "/api/update-enforcer", "/api/add-enforcer", "/api/delete-enforcer", "/api/enforce", "/api/batch-enforce", "/api/get-all-objects", "/api/get-all-actions", "/api/get-all-roles", "/api/run-casbin-command", "/api/refresh-engines"], "enforcer"],
    [["/api/get-agents", "/api/get-agent", "/api/update-agent", "/api/add-agent", "/api/delete-agent"], "agent"],
    [["/api/get-servers", "/api/get-online-servers", "/api/scan", "/api/sync-intranet-servers", "/api/get-server", "/api/update-server", "/api/sync-mcp-tool", "/api/add-server", "/api/delete-server", "/api/server/", "/api/get-mcp-access-token"], "server"],
    [["/api/get-entries", "/api/get-entry", "/api/get-openclaw-session-graph", "/api/get-openclaw-session-transcript", "/api/update-entry", "/api/add-entry", "/api/delete-entry"], "entry"],
    [["/api/get-sites", "/api/get-global-sites", "/api/get-site", "/api/update-site", "/api/add-site", "/api/delete-site"], "site"],
    [["/api/get-rules", "/api/get-rule", "/api/add-rule", "/api/update-rule", "/api/delete-rule"], "rule"],
    [["/api/get-sessions", "/api/get-session", "/api/update-session", "/api/add-session", "/api/delete-session", "/api/is-session-duplicated"], "session"],
    [["/api/get-tokens", "/api/get-token", "/api/update-token", "/api/add-token", "/api/delete-token"], "token"],
    [["/api/get-records", "/api/get-records-filter", "/api/add-record"], "record"],
    [["/api/get-verifications", "/api/get-email-and-phone", "/api/send-verification-code", "/api/verify-code", "/api/verify-captcha", "/api/reset-email-or-phone", "/api/get-captcha"], "verification"],
    [["/api/get-system-info", "/api/get-version-info", "/api/health", "/api/get-prometheus-info", "/api/metrics", "/api/v1/traces", "/api/v1/metrics", "/api/v1/logs"], "system"],
    [["/api/get-forms", "/api/get-global-forms", "/api/get-form", "/api/update-form", "/api/add-form", "/api/delete-form"], "form"],
    [["/api/get-syncers", "/api/get-syncer", "/api/update-syncer", "/api/add-syncer", "/api/delete-syncer", "/api/run-syncer", "/api/test-syncer-db"], "syncer"],
    [["/api/get-webhooks", "/api/get-webhook", "/api/update-webhook", "/api/add-webhook", "/api/delete-webhook", "/api/get-webhook-events", "/api/get-webhook-event-detail", "/api/replay-webhook-event", "/api/delete-webhook-event", "/api/webhook", "/api/get-webhook-event"], "webhook"],
    [["/api/get-tickets", "/api/get-ticket", "/api/update-ticket", "/api/add-ticket", "/api/delete-ticket", "/api/add-ticket-message"], "ticket"],
    [["/api/get-ldap-users", "/api/get-ldaps", "/api/get-ldap", "/api/add-ldap", "/api/update-ldap", "/api/delete-ldap", "/api/sync-ldap-users", "/scim/"], "ldap"],
    [["/api/send-email", "/api/send-sms", "/api/send-notification"], "service"],
    [["/api/oauth/register", "/api/login/oauth/access_token", "/api/login/oauth/refresh_token", "/api/login/oauth/introspect"], "protocol"],
    [["/api/mcp"], "mcp"],
  ];

  for (const [prefixes, domain] of matchers) {
    if (prefixes.some(prefix => pathname.startsWith(prefix))) {
      return domain;
    }
  }

  if (action && action.toLowerCase().includes("mcp")) {
    return "mcp";
  }

  return "system";
}

function inferFrontendModule(domain) {
  const mapping = {
    organization: "user-management",
    group: "user-management",
    user: "user-management",
    invitation: "user-management",
    application: "identity",
    provider: "identity",
    resource: "identity",
    cert: "identity",
    key: "identity",
    role: "authorization",
    permission: "authorization",
    model: "authorization",
    adapter: "authorization",
    enforcer: "authorization",
    agent: "gateway",
    server: "gateway",
    entry: "gateway",
    site: "gateway",
    rule: "gateway",
    mcp: "gateway",
    session: "auditing",
    token: "auditing",
    record: "auditing",
    verification: "auditing",
    form: "admin-ops",
    syncer: "admin-ops",
    webhook: "admin-ops",
    ticket: "admin-ops",
    ldap: "admin-ops",
    system: "admin-ops",
  };

  return mapping[domain] || "out-of-scope";
}

function inferAuthMode(pathname, domain) {
  if (pathname.startsWith("/.well-known") || pathname.startsWith("/cas/")) {
    return "Public protocol surface";
  }
  if (pathname.startsWith("/scim/")) {
    return "Controller-defined admin credential";
  }
  if (["auth", "account", "protocol"].includes(domain)) {
    return "Public or login-session flow";
  }
  return "Session admin surface";
}

function inferPrimaryUse(pathname, action) {
  const actionText = humanizeAction(action);
  if (pathname.startsWith("/.well-known")) {
    return "Publish discovery and key metadata";
  }
  if (pathname.startsWith("/cas/")) {
    return "Serve CAS protocol validation";
  }
  if (pathname.startsWith("/scim/")) {
    return "Serve SCIM and directory compatibility";
  }
  return actionText;
}

function domainToTag(domain) {
  return `${titleCase(domain)} API`;
}

function parseRouter() {
  const source = read(files.router);
  const regex = /web\.Router\("([^"]+)",\s*&([A-Za-z0-9_.]+)\{\},\s*"([^"]+)"\)/g;
  const routeRows = [];
  const pathMap = new Map();
  let match;

  while ((match = regex.exec(source)) !== null) {
    const pathname = match[1];
    const controllerRef = match[2];
    const spec = match[3];
    const [methodsPart, action = ""] = spec.split(":");
    const methods = uniqueSorted(methodsPart.split(",").map(item => item.trim()).filter(Boolean));
    const controller = controllerRef.split(".").pop();
    const domain = inferDomain(pathname, action);
    const frontendModule = inferFrontendModule(domain);
    const authMode = inferAuthMode(pathname, domain);
    const primaryUse = inferPrimaryUse(pathname, action);

    routeRows.push({
      path: pathname,
      methods,
      controller,
      action,
      controllerAction: `${controller}.${action}`,
      domain,
      frontendModule,
      authMode,
      primaryUse,
    });

    const existing = pathMap.get(pathname) || {
      path: pathname,
      methods: new Set(),
      controllerActions: new Set(),
      domains: new Set(),
    };

    methods.forEach(method => existing.methods.add(method));
    existing.controllerActions.add(`${controller}.${action}`);
    existing.domains.add(domain);
    pathMap.set(pathname, existing);
  }

  const uniquePaths = [...pathMap.values()].map(item => ({
    path: item.path,
    methods: uniqueSorted([...item.methods]),
    controllerActions: uniqueSorted([...item.controllerActions]),
    domains: uniqueSorted([...item.domains]),
  })).sort((a, b) => a.path.localeCompare(b.path));

  return {
    routeRows: routeRows.sort((a, b) => a.path.localeCompare(b.path) || a.controllerAction.localeCompare(b.controllerAction)),
    uniquePaths,
  };
}

function loadSwagger() {
  const doc = yaml.load(read(files.swagger));
  const paths = doc.paths || {};
  return {doc, paths};
}

function buildGeneratedOpenAPI(routerData, swaggerData) {
  const clone = JSON.parse(JSON.stringify(swaggerData.doc));
  const generatedPaths = {};
  const routeByPath = new Map(routerData.routeRows.map(route => [route.path, route]));

  for (const pathEntry of routerData.uniquePaths) {
    const existing = swaggerData.paths[pathEntry.path] ? JSON.parse(JSON.stringify(swaggerData.paths[pathEntry.path])) : {};
    const liveMethods = new Set(pathEntry.methods.map(item => item.toLowerCase()));
    const nextPath = {};

    for (const [key, value] of Object.entries(existing)) {
      if (key.startsWith("x-")) {
        nextPath[key] = value;
      } else if (liveMethods.has(key)) {
        nextPath[key] = value;
      }
    }

    if (!pathEntry.methods.includes("*")) {
      for (const method of pathEntry.methods) {
        const opKey = method.toLowerCase();
        if (!nextPath[opKey]) {
          const route = routeByPath.get(pathEntry.path);
          nextPath[opKey] = {
            tags: [domainToTag(route.domain)],
            operationId: `${route.controller}.${route.action}.${method.toLowerCase()}`,
            description: "Generated from live router definition; controller annotation is missing from the published Swagger artifact.",
            responses: {
              "200": {
                description: "Controller-defined response",
              },
            },
          };
        }
      }
    }

    nextPath["x-live-router-methods"] = pathEntry.methods;
    nextPath["x-live-controller-actions"] = pathEntry.controllerActions;
    nextPath["x-live-domains"] = pathEntry.domains;
    generatedPaths[pathEntry.path] = nextPath;
  }

  clone.swagger = clone.swagger || "2.0";
  clone.info = clone.info || {};
  clone.info.title = "Pipeline Auth Live Router API Reference";
  clone.info.description = "Generated from live router definitions and merged with the published Swagger artifact when matching operations exist.";
  clone.paths = generatedPaths;
  clone["x-generated-note"] = "Do not edit manually. Regenerate from scripts/docs/generate-reference-docs.js.";

  return yaml.dump(clone, {
    lineWidth: 120,
    noRefs: true,
    sortKeys: false,
  });
}

function parseManagementImports() {
  const source = read(files.management);
  const regex = /^import\s+([A-Za-z0-9_]+)\s+from\s+"([^"]+)";$/gm;
  const imports = new Map();
  let match;
  while ((match = regex.exec(source)) !== null) {
    imports.set(match[1], match[2]);
  }
  return imports;
}

function inferMenuGroup(routePath) {
  if (routePath === "/" || routePath.includes("/shortcuts") || routePath.includes("/apps")) {
    return "Home";
  }
  if (routePath.includes("/organizations") || routePath.includes("/trees") || routePath.includes("/groups") || routePath.includes("/users") || routePath.includes("/invitations")) {
    return "User Management";
  }
  if (routePath.includes("/applications") || routePath.includes("/providers") || routePath.includes("/resources") || routePath.includes("/certs") || routePath.includes("/keys")) {
    return "Identity";
  }
  if (routePath.includes("/roles") || routePath.includes("/permissions") || routePath.includes("/models") || routePath.includes("/adapters") || routePath.includes("/enforcers")) {
    return "Authorization";
  }
  if (routePath.includes("/agents") || routePath.includes("/servers") || routePath.includes("/server-store") || routePath.includes("/entries") || routePath.includes("/sites") || routePath.includes("/rules")) {
    return "LLM AI";
  }
  if (routePath.includes("/sessions") || routePath.includes("/records") || routePath.includes("/tokens") || routePath.includes("/verifications")) {
    return "Auditing";
  }
  return "Admin";
}

function inferSpecModuleFromRoute(routePath) {
  const mapping = {
    "Home": "overview",
    "User Management": "user-management",
    "Identity": "identity",
    "Authorization": "authorization",
    "LLM AI": "gateway",
    "Auditing": "auditing",
    "Admin": "admin-ops",
  };

  return mapping[inferMenuGroup(routePath)] || "overview";
}

function parseManagementRoutes() {
  const source = read(files.management);
  const imports = parseManagementImports();
  const regex = /<Route exact path="([^"]+)" render=\{\(props\) =>[\s\S]*?<([A-Za-z0-9_]+)\b/g;
  const routes = [];
  let match;

  while ((match = regex.exec(source)) !== null) {
    const routePath = match[1];
    const component = match[2];
    routes.push({
      path: routePath,
      component,
      importPath: imports.get(component) || null,
      menuGroup: inferMenuGroup(routePath),
      specModule: inferSpecModuleFromRoute(routePath),
    });
  }

  return routes;
}

function parseBackendFunctions() {
  const endpointsByModule = new Map();
  const functionToEndpoint = new Map();
  const filesInDir = fs.readdirSync(files.backendDir).filter(name => name.endsWith(".js")).sort();

  for (const fileName of filesInDir) {
    const moduleName = fileName.replace(".js", "");
    const moduleAlias = moduleName;
    const source = read(path.join(files.backendDir, fileName));
    const helperBlocks = new Map();
    const allFunctionMatches = [...source.matchAll(/(?:export\s+)?function\s+([A-Za-z0-9_]+)\([^)]*\)\s*\{([\s\S]*?)(?=\n(?:export\s+)?function|\n$)/g)];

    for (const match of allFunctionMatches) {
      helperBlocks.set(match[1], match[2]);
    }

    const matches = [...source.matchAll(/export function\s+([A-Za-z0-9_]+)\([^)]*\)\s*\{([\s\S]*?)(?=\nexport function|\n$)/g)];
    const functions = [];

    for (const match of matches) {
      const fnName = match[1];
      const block = match[2];
      let endpoint = "";
      const endpointMatch = block.match(/(\/(?:api|\.well-known|cas|scim)[^"'`\s?)]*)/);
      if (endpointMatch) {
        endpoint = endpointMatch[1];
      } else {
        const dashboardMatch = block.match(/fetchDashboardApi\("([^"]+)"/);
        if (dashboardMatch) {
          endpoint = `/api/${dashboardMatch[1]}`;
        } else {
          const helperCallMatch = block.match(/([A-Za-z0-9_]+Url)\(/) || block.match(/return\s+([A-Za-z0-9_]+)\(/);
          if (helperCallMatch) {
            const helperBlock = helperBlocks.get(helperCallMatch[1]) || "";
            const helperEndpointMatch = helperBlock.match(/(\/(?:api|\.well-known|cas|scim)[^"'`\s?)]*)/);
            if (helperEndpointMatch) {
              endpoint = helperEndpointMatch[1];
            }
          }
        }
      }
      functions.push({name: fnName, endpoint});
      functionToEndpoint.set(`${moduleAlias}.${fnName}`, endpoint);
    }

    endpointsByModule.set(moduleAlias, functions);
  }

  return {endpointsByModule, functionToEndpoint};
}

function resolveComponentFile(importPath) {
  if (!importPath) {
    return null;
  }

  const candidates = [
    path.resolve(path.dirname(files.management), `${importPath}.js`),
    path.resolve(path.dirname(files.management), importPath, "index.js"),
  ];

  return candidates.find(candidate => fs.existsSync(candidate)) || null;
}

function buildPageApiMatrix(managementRoutes, backendMeta) {
  const rows = [];

  for (const route of managementRoutes) {
    const componentFile = resolveComponentFile(route.importPath);
    const wrappers = new Set();
    const endpoints = new Set();

    if (componentFile) {
      const source = read(componentFile);
      const importMatches = [...source.matchAll(/import\s+\*\s+as\s+([A-Za-z0-9_]+)\s+from\s+"([^"]*backend\/[A-Za-z0-9_]+)";/g)];
      const backendImports = new Map();

      for (const match of importMatches) {
        backendImports.set(match[1], path.basename(match[2]));
      }

      for (const [alias, moduleAlias] of backendImports.entries()) {
        const callRegex = new RegExp(`${alias}\\.([A-Za-z0-9_]+)\\(`, "g");
        let callMatch;
        while ((callMatch = callRegex.exec(source)) !== null) {
          const fnName = callMatch[1];
          const wrapperKey = `${moduleAlias}.${fnName}`;
          wrappers.add(wrapperKey);
          const endpoint = backendMeta.functionToEndpoint.get(wrapperKey);
          if (endpoint) {
            endpoints.add(endpoint);
          }
        }
      }
    }

    rows.push({
      ...route,
      componentFile,
      wrappers: uniqueSorted([...wrappers]),
      endpoints: uniqueSorted([...endpoints]),
    });
  }

  return rows.sort((a, b) => a.specModule.localeCompare(b.specModule) || a.path.localeCompare(b.path));
}

function renderBackendRoutesMarkdown(routerData) {
  const totalDeclarations = routerData.routeRows.length;
  const uniquePaths = routerData.uniquePaths.length;
  const byDomain = new Map();

  for (const route of routerData.routeRows) {
    const list = byDomain.get(route.domain) || [];
    list.push(route);
    byDomain.set(route.domain, list);
  }

  const lines = [
    "# Backend Route Reference",
    "",
    "Generated from `routers/router.go`. This document is the human-readable companion to `openapi.generated.yaml`.",
    "",
    `- Live route declarations: ${totalDeclarations}`,
    `- Live unique paths: ${uniquePaths}`,
    "",
    "Each row reflects router truth. Auth mode is an observed classification for frontend/spec work and must be treated as controller-defined unless a module spec narrows it further.",
    "",
  ];

  for (const domain of uniqueSorted([...byDomain.keys()])) {
    lines.push(`## ${titleCase(domain)}`);
    lines.push("");
    lines.push("| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |");
    lines.push("| --- | --- | --- | --- | --- | --- |");
    for (const route of byDomain.get(domain)) {
      lines.push(`| ${route.methods.join(", ")} | \`${route.path}\` | \`${route.controllerAction}\` | ${route.authMode} | ${route.primaryUse} | ${route.frontendModule} |`);
    }
    lines.push("");
  }

  return lines.join("\n");
}

function renderAdminRoutesMarkdown(managementRoutes) {
  const byGroup = new Map();
  for (const route of managementRoutes) {
    const list = byGroup.get(route.menuGroup) || [];
    list.push(route);
    byGroup.set(route.menuGroup, list);
  }

  const lines = [
    "# Admin Route Inventory",
    "",
    "Generated from `web/src/ManagementPage.js`.",
    "",
    `- Routes discovered: ${managementRoutes.length}`,
    "",
  ];

  for (const group of ["Home", "User Management", "Identity", "Authorization", "LLM AI", "Auditing", "Admin"]) {
    if (!byGroup.has(group)) {
      continue;
    }
    lines.push(`## ${group}`);
    lines.push("");
    lines.push("| Route | Component | Source file | Spec module |");
    lines.push("| --- | --- | --- | --- |");
    for (const route of byGroup.get(group).sort((a, b) => a.path.localeCompare(b.path))) {
      const relativeSource = route.importPath ? path.posix.normalize(path.join("web/src", route.importPath.replace(/^\.\//, ""))) + ".js" : "";
      lines.push(`| \`${route.path}\` | \`${route.component}\` | ${relativeSource ? `\`${relativeSource}\`` : ""} | \`${route.specModule}\` |`);
    }
    lines.push("");
  }

  return lines.join("\n");
}

function renderPageApiMatrixMarkdown(matrixRows) {
  const lines = [
    "# Page to API Matrix",
    "",
    "Generated from `web/src/ManagementPage.js`, component imports, and `web/src/backend/*` wrappers.",
    "",
    `- Pages analyzed: ${matrixRows.length}`,
    "",
    "| Spec module | Component | Route(s) | Backend wrappers | API endpoints |",
    "| --- | --- | --- | --- | --- |",
  ];

  const grouped = new Map();
  for (const row of matrixRows) {
    const key = row.component;
    const existing = grouped.get(key) || {
      specModule: row.specModule,
      component: row.component,
      routes: new Set(),
      wrappers: new Set(),
      endpoints: new Set(),
    };
    existing.routes.add(row.path);
    row.wrappers.forEach(item => existing.wrappers.add(item));
    row.endpoints.forEach(item => existing.endpoints.add(item));
    grouped.set(key, existing);
  }

  for (const row of [...grouped.values()].sort((a, b) => a.specModule.localeCompare(b.specModule) || a.component.localeCompare(b.component))) {
    const routes = uniqueSorted([...row.routes]).map(item => `\`${item}\``).join("<br>");
    const wrappers = uniqueSorted([...row.wrappers]).map(item => `\`${item}\``).join("<br>");
    const endpoints = uniqueSorted([...row.endpoints]).map(item => `\`${item}\``).join("<br>");
    lines.push(`| \`${row.specModule}\` | \`${row.component}\` | ${routes} | ${wrappers || ""} | ${endpoints || ""} |`);
  }

  lines.push("");
  return lines.join("\n");
}

function renderDriftReportMarkdown(routerData, swaggerData, matrixRows, backendMeta) {
  const livePaths = new Set(routerData.uniquePaths.map(item => item.path));
  const swaggerPaths = new Set(Object.keys(swaggerData.paths));
  const frontendEndpoints = new Set();

  for (const row of matrixRows) {
    row.endpoints.forEach(endpoint => frontendEndpoints.add(endpoint));
  }
  for (const functions of backendMeta.endpointsByModule.values()) {
    functions.forEach(item => {
      if (item.endpoint) {
        frontendEndpoints.add(item.endpoint);
      }
    });
  }

  const liveMissingInSwagger = uniqueSorted([...livePaths].filter(item => !swaggerPaths.has(item)));
  const swaggerMissingInLive = uniqueSorted([...swaggerPaths].filter(item => !livePaths.has(item)));
  const frontendMissingInLive = uniqueSorted([...frontendEndpoints].filter(item => !livePaths.has(item)));
  const frontendMissingInSwagger = uniqueSorted([...frontendEndpoints].filter(item => !swaggerPaths.has(item)));

  const lines = [
    "# API Drift Report",
    "",
    "Generated by comparing live router paths, the published Swagger artifact, and frontend backend-wrapper usage.",
    "",
    `- Live unique router paths: ${livePaths.size}`,
    `- Published Swagger paths: ${swaggerPaths.size}`,
    `- Frontend wrapper endpoints: ${frontendEndpoints.size}`,
    "",
    `- Live paths missing from Swagger: ${liveMissingInSwagger.length}`,
    `- Swagger paths missing from live router: ${swaggerMissingInLive.length}`,
    `- Frontend endpoints missing from live router: ${frontendMissingInLive.length}`,
    `- Frontend endpoints missing from Swagger: ${frontendMissingInSwagger.length}`,
    "",
    "## Live paths missing from Swagger",
    "",
  ];

  liveMissingInSwagger.forEach(item => lines.push(`- \`${item}\``));
  if (liveMissingInSwagger.length === 0) {
    lines.push("- None");
  }

  lines.push("");
  lines.push("## Swagger paths missing from live router");
  lines.push("");
  swaggerMissingInLive.forEach(item => lines.push(`- \`${item}\``));
  if (swaggerMissingInLive.length === 0) {
    lines.push("- None");
  }

  lines.push("");
  lines.push("## Frontend endpoints missing from live router");
  lines.push("");
  frontendMissingInLive.forEach(item => lines.push(`- \`${item}\``));
  if (frontendMissingInLive.length === 0) {
    lines.push("- None");
  }

  lines.push("");
  lines.push("## Frontend endpoints missing from Swagger");
  lines.push("");
  frontendMissingInSwagger.forEach(item => lines.push(`- \`${item}\``));
  if (frontendMissingInSwagger.length === 0) {
    lines.push("- None");
  }

  lines.push("");
  return lines.join("\n");
}

function generateBackendReference(routerData, swaggerData) {
  write(path.join(files.referenceBackendDir, "openapi.generated.yaml"), buildGeneratedOpenAPI(routerData, swaggerData));
  write(path.join(files.referenceBackendDir, "routes.generated.md"), renderBackendRoutesMarkdown(routerData));
}

function generateFrontendReference(managementRoutes, matrixRows) {
  write(path.join(files.referenceFrontendDir, "admin-routes.generated.md"), renderAdminRoutesMarkdown(managementRoutes));
  write(path.join(files.referenceFrontendDir, "page-api-matrix.generated.md"), renderPageApiMatrixMarkdown(matrixRows));
}

function generateDrift(routerData, swaggerData, matrixRows, backendMeta) {
  write(path.join(files.referenceBackendDir, "drift-report.generated.md"), renderDriftReportMarkdown(routerData, swaggerData, matrixRows, backendMeta));
}

function main() {
  const routerData = parseRouter();
  const swaggerData = loadSwagger();
  const managementRoutes = parseManagementRoutes();
  const backendMeta = parseBackendFunctions();
  const matrixRows = buildPageApiMatrix(managementRoutes, backendMeta);

  if (mode === "backend" || mode === "all") {
    generateBackendReference(routerData, swaggerData);
  }
  if (mode === "frontend" || mode === "all") {
    generateFrontendReference(managementRoutes, matrixRows);
  }
  if (mode === "drift" || mode === "all") {
    generateDrift(routerData, swaggerData, matrixRows, backendMeta);
  }
}

main();
