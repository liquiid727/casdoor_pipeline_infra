# Admin Route Inventory

Generated from `web/src/ManagementPage.js`.

- Routes discovered: 65

## Home

| Route | Component | Source file | Spec module |
| --- | --- | --- | --- |
| `/` | `Dashboard` | `web/src/basic/Dashboard.js` | `overview` |
| `/apps` | `AppListPage` | `web/src/basic/AppListPage.js` | `overview` |
| `/shortcuts` | `ShortcutsPage` | `web/src/basic/ShortcutsPage.js` | `overview` |

## User Management

| Route | Component | Source file | Spec module |
| --- | --- | --- | --- |
| `/groups` | `GroupListPage` | `web/src/GroupListPage.js` | `user-management` |
| `/groups/:organizationName/:groupName` | `GroupEditPage` | `web/src/GroupEditPage.js` | `user-management` |
| `/invitations` | `InvitationListPage` | `web/src/InvitationListPage.js` | `user-management` |
| `/invitations/:organizationName/:invitationName` | `InvitationEditPage` | `web/src/InvitationEditPage.js` | `user-management` |
| `/organizations` | `OrganizationListPage` | `web/src/OrganizationListPage.js` | `user-management` |
| `/organizations/:organizationName` | `OrganizationEditPage` | `web/src/OrganizationEditPage.js` | `user-management` |
| `/organizations/:organizationName/users` | `UserListPage` | `web/src/UserListPage.js` | `user-management` |
| `/trees/:organizationName` | `GroupTreePage` | `web/src/GroupTreePage.js` | `user-management` |
| `/trees/:organizationName/:groupName` | `GroupTreePage` | `web/src/GroupTreePage.js` | `user-management` |
| `/users` | `UserListPage` | `web/src/UserListPage.js` | `user-management` |
| `/users/:organizationName/:userName` | `UserEditPage` | `web/src/UserEditPage.js` | `user-management` |

## Identity

| Route | Component | Source file | Spec module |
| --- | --- | --- | --- |
| `/applications` | `ApplicationListPage` | `web/src/ApplicationListPage.js` | `identity` |
| `/applications/:organizationName/:applicationName` | `ApplicationEditPage` | `web/src/ApplicationEditPage.js` | `identity` |
| `/certs` | `CertListPage` | `web/src/CertListPage.js` | `identity` |
| `/certs/:organizationName/:certName` | `CertEditPage` | `web/src/CertEditPage.js` | `identity` |
| `/keys` | `KeyListPage` | `web/src/KeyListPage.js` | `identity` |
| `/keys/:organizationName/:keyName` | `KeyEditPage` | `web/src/KeyEditPage.js` | `identity` |
| `/providers` | `ProviderListPage` | `web/src/ProviderListPage.js` | `identity` |
| `/providers/:organizationName/:providerName` | `ProviderEditPage` | `web/src/ProviderEditPage.js` | `identity` |
| `/resources` | `ResourceListPage` | `web/src/ResourceListPage.js` | `identity` |

## Authorization

| Route | Component | Source file | Spec module |
| --- | --- | --- | --- |
| `/adapters` | `AdapterListPage` | `web/src/AdapterListPage.js` | `authorization` |
| `/adapters/:organizationName/:adapterName` | `AdapterEditPage` | `web/src/AdapterEditPage.js` | `authorization` |
| `/enforcers` | `EnforcerListPage` | `web/src/EnforcerListPage.js` | `authorization` |
| `/enforcers/:organizationName/:enforcerName` | `EnforcerEditPage` | `web/src/EnforcerEditPage.js` | `authorization` |
| `/models` | `ModelListPage` | `web/src/ModelListPage.js` | `authorization` |
| `/models/:organizationName/:modelName` | `ModelEditPage` | `web/src/ModelEditPage.js` | `authorization` |
| `/permissions` | `PermissionListPage` | `web/src/PermissionListPage.js` | `authorization` |
| `/permissions/:organizationName/:permissionName` | `PermissionEditPage` | `web/src/PermissionEditPage.js` | `authorization` |
| `/roles` | `RoleListPage` | `web/src/RoleListPage.js` | `authorization` |
| `/roles/:organizationName/:roleName` | `RoleEditPage` | `web/src/RoleEditPage.js` | `authorization` |

## LLM AI

| Route | Component | Source file | Spec module |
| --- | --- | --- | --- |
| `/agents` | `AgentListPage` | `web/src/AgentListPage.js` | `gateway` |
| `/agents/:organizationName/:agentName` | `AgentEditPage` | `web/src/AgentEditPage.js` | `gateway` |
| `/entries` | `EntryListPage` | `web/src/EntryListPage.js` | `gateway` |
| `/entries/:organizationName/:entryName` | `EntryEditPage` | `web/src/EntryEditPage.js` | `gateway` |
| `/entries/:organizationName/:entryName/transcript` | `OpenClawSessionTranscriptPage` | `web/src/OpenClawSessionTranscriptPage.js` | `gateway` |
| `/rules` | `RuleListPage` | `web/src/RuleListPage.js` | `gateway` |
| `/rules/:organizationName/:ruleName` | `RuleEditPage` | `web/src/RuleEditPage.js` | `gateway` |
| `/server-store` | `ServerStorePage` | `web/src/ServerStorePage.js` | `gateway` |
| `/servers` | `ServerListPage` | `web/src/ServerListPage.js` | `gateway` |
| `/servers/:organizationName/:serverName` | `ServerEditPage` | `web/src/ServerEditPage.js` | `gateway` |
| `/sites` | `SiteListPage` | `web/src/SiteListPage.js` | `gateway` |
| `/sites/:organizationName/:siteName` | `SiteEditPage` | `web/src/SiteEditPage.js` | `gateway` |

## Auditing

| Route | Component | Source file | Spec module |
| --- | --- | --- | --- |
| `/records` | `RecordListPage` | `web/src/RecordListPage.js` | `auditing` |
| `/sessions` | `SessionListPage` | `web/src/SessionListPage.js` | `auditing` |
| `/tokens` | `TokenListPage` | `web/src/TokenListPage.js` | `auditing` |
| `/tokens/:tokenName` | `TokenEditPage` | `web/src/TokenEditPage.js` | `auditing` |
| `/verifications` | `VerificationListPage` | `web/src/VerificationListPage.js` | `auditing` |

## Admin

| Route | Component | Source file | Spec module |
| --- | --- | --- | --- |
| `/.well-known/openid-configuration` | `OdicDiscoveryPage` | `web/src/auth/OidcDiscoveryPage.js` | `admin-ops` |
| `/account` | `AccountPage` | `web/src/account/AccountPage.js` | `admin-ops` |
| `/forms` | `FormListPage` | `web/src/FormListPage.js` | `admin-ops` |
| `/forms/:formName` | `FormEditPage` | `web/src/FormEditPage.js` | `admin-ops` |
| `/ldap/:organizationName/:ldapId` | `LdapEditPage` | `web/src/LdapEditPage.js` | `admin-ops` |
| `/ldap/sync/:organizationName/:ldapId` | `LdapSyncPage` | `web/src/LdapSyncPage.js` | `admin-ops` |
| `/mfa/setup` | `MfaSetupPage` | `web/src/auth/MfaSetupPage.js` | `admin-ops` |
| `/syncers` | `SyncerListPage` | `web/src/SyncerListPage.js` | `admin-ops` |
| `/syncers/:organizationName/:syncerName` | `SyncerEditPage` | `web/src/SyncerEditPage.js` | `admin-ops` |
| `/sysinfo` | `SystemInfo` | `web/src/SystemInfo.js` | `admin-ops` |
| `/tickets` | `TicketListPage` | `web/src/TicketListPage.js` | `admin-ops` |
| `/tickets/:organizationName/:ticketName` | `TicketEditPage` | `web/src/TicketEditPage.js` | `admin-ops` |
| `/webhook-events` | `WebhookEventListPage` | `web/src/WebhookEventListPage.js` | `admin-ops` |
| `/webhooks` | `WebhookListPage` | `web/src/WebhookListPage.js` | `admin-ops` |
| `/webhooks/:webhookName` | `WebhookEditPage` | `web/src/WebhookEditPage.js` | `admin-ops` |
