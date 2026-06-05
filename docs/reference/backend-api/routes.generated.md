# Backend Route Reference

Generated from `routers/router.go`. This document is the human-readable companion to `openapi.generated.yaml`.

- Live route declarations: 273
- Live unique paths: 272

Each row reflects router truth. Auth mode is an observed classification for frontend/spec work and must be treated as controller-defined unless a module spec narrows it further.

## Account

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| GET | `/api/get-account` | `ApiController.GetAccount` | Public or login-session flow | Get Account | out-of-scope |
| POST | `/api/unlink` | `ApiController.Unlink` | Public or login-session flow | Unlink | out-of-scope |
| GET | `/api/user` | `ApiController.GetUserinfo2` | Public or login-session flow | Get Userinfo2 | out-of-scope |
| GET | `/api/userinfo` | `ApiController.GetUserinfo` | Public or login-session flow | Get Userinfo | out-of-scope |

## Adapter

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-adapter` | `ApiController.AddAdapter` | Session admin surface | Add Adapter | authorization |
| POST | `/api/add-policy` | `ApiController.AddPolicy` | Session admin surface | Add Policy | authorization |
| POST | `/api/delete-adapter` | `ApiController.DeleteAdapter` | Session admin surface | Delete Adapter | authorization |
| GET | `/api/get-adapter` | `ApiController.GetAdapter` | Session admin surface | Get Adapter | authorization |
| GET | `/api/get-adapters` | `ApiController.GetAdapters` | Session admin surface | Get Adapters | authorization |
| POST | `/api/get-filtered-policies` | `ApiController.GetFilteredPolicies` | Session admin surface | Get Filtered Policies | authorization |
| GET | `/api/get-policies` | `ApiController.GetPolicies` | Session admin surface | Get Policies | authorization |
| POST | `/api/remove-policy` | `ApiController.RemovePolicy` | Session admin surface | Remove Policy | authorization |
| POST | `/api/update-adapter` | `ApiController.UpdateAdapter` | Session admin surface | Update Adapter | authorization |
| POST | `/api/update-policy` | `ApiController.UpdatePolicy` | Session admin surface | Update Policy | authorization |

## Agent

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-agent` | `ApiController.AddAgent` | Session admin surface | Add Agent | gateway |
| POST | `/api/delete-agent` | `ApiController.DeleteAgent` | Session admin surface | Delete Agent | gateway |
| GET | `/api/get-agent` | `ApiController.GetAgent` | Session admin surface | Get Agent | gateway |
| GET | `/api/get-agents` | `ApiController.GetAgents` | Session admin surface | Get Agents | gateway |
| POST | `/api/update-agent` | `ApiController.UpdateAgent` | Session admin surface | Update Agent | gateway |

## Application

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-application` | `ApiController.AddApplication` | Session admin surface | Add Application | identity |
| POST | `/api/delete-application` | `ApiController.DeleteApplication` | Session admin surface | Delete Application | identity |
| GET | `/api/get-application` | `ApiController.GetApplication` | Session admin surface | Get Application | identity |
| GET | `/api/get-applications` | `ApiController.GetApplications` | Session admin surface | Get Applications | identity |
| POST | `/api/update-application` | `ApiController.UpdateApplication` | Session admin surface | Update Application | identity |

## Auth

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/acs` | `ApiController.HandleSamlLogin` | Public or login-session flow | Handle Saml Login | out-of-scope |
| POST | `/api/callback` | `ApiController.Callback` | Public or login-session flow | Callback | out-of-scope |
| POST | `/api/cancel-device-auth` | `ApiController.CancelDeviceAuth` | Public or login-session flow | Cancel Device Auth | out-of-scope |
| POST | `/api/device-auth` | `ApiController.DeviceAuth` | Public or login-session flow | Device Auth | out-of-scope |
| POST | `/api/device-auth-complete` | `ApiController.DeviceAuthComplete` | Public or login-session flow | Device Auth Complete | out-of-scope |
| GET | `/api/faceid-signin-begin` | `ApiController.FaceIDSigninBegin` | Public or login-session flow | Face IDSignin Begin | out-of-scope |
| GET | `/api/get-saml-login` | `ApiController.GetSamlLogin` | Public or login-session flow | Get Saml Login | out-of-scope |
| POST | `/api/grant-consent` | `ApiController.GrantConsent` | Public or login-session flow | Grant Consent | out-of-scope |
| GET | `/api/kerberos-login` | `ApiController.KerberosLogin` | Public or login-session flow | Kerberos Login | out-of-scope |
| POST | `/api/login` | `ApiController.Login` | Public or login-session flow | Login | out-of-scope |
| POST | `/api/login/oauth/access_token` | `ApiController.GetOAuthToken` | Public or login-session flow | Get OAuth Token | out-of-scope |
| POST | `/api/login/oauth/introspect` | `ApiController.IntrospectToken` | Public or login-session flow | Introspect Token | out-of-scope |
| POST | `/api/login/oauth/refresh_token` | `ApiController.RefreshToken` | Public or login-session flow | Refresh Token | out-of-scope |
| GET, POST | `/api/logout` | `ApiController.Logout` | Public or login-session flow | Logout | out-of-scope |
| POST | `/api/mfa/setup/enable` | `ApiController.MfaSetupEnable` | Public or login-session flow | Mfa Setup Enable | out-of-scope |
| POST | `/api/mfa/setup/initiate` | `ApiController.MfaSetupInitiate` | Public or login-session flow | Mfa Setup Initiate | out-of-scope |
| POST | `/api/mfa/setup/verify` | `ApiController.MfaSetupVerify` | Public or login-session flow | Mfa Setup Verify | out-of-scope |
| POST | `/api/native-sso-complete` | `ApiController.NativeSsoComplete` | Public or login-session flow | Native Sso Complete | out-of-scope |
| POST | `/api/revoke-consent` | `ApiController.RevokeConsent` | Public or login-session flow | Revoke Consent | out-of-scope |
| GET | `/api/saml/metadata` | `ApiController.GetSamlMeta` | Public or login-session flow | Get Saml Meta | out-of-scope |
| * | `/api/saml/redirect/:owner/:application` | `ApiController.HandleSamlRedirect` | Public or login-session flow | Handle Saml Redirect | out-of-scope |
| POST | `/api/signup` | `ApiController.Signup` | Public or login-session flow | Signup | out-of-scope |
| GET, POST | `/api/sso-logout` | `ApiController.SsoLogout` | Public or login-session flow | Sso Logout | out-of-scope |
| GET | `/api/webauthn/signin/begin` | `ApiController.WebAuthnSigninBegin` | Public or login-session flow | Web Authn Signin Begin | out-of-scope |
| POST | `/api/webauthn/signin/finish` | `ApiController.WebAuthnSigninFinish` | Public or login-session flow | Web Authn Signin Finish | out-of-scope |
| GET | `/api/webauthn/signup/begin` | `ApiController.WebAuthnSignupBegin` | Public or login-session flow | Web Authn Signup Begin | out-of-scope |
| POST | `/api/webauthn/signup/finish` | `ApiController.WebAuthnSignupFinish` | Public or login-session flow | Web Authn Signup Finish | out-of-scope |

## Cert

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-cert` | `ApiController.AddCert` | Session admin surface | Add Cert | identity |
| POST | `/api/delete-cert` | `ApiController.DeleteCert` | Session admin surface | Delete Cert | identity |
| GET | `/api/get-cert` | `ApiController.GetCert` | Session admin surface | Get Cert | identity |
| GET | `/api/get-certs` | `ApiController.GetCerts` | Session admin surface | Get Certs | identity |
| GET | `/api/get-global-certs` | `ApiController.GetGlobalCerts` | Session admin surface | Get Global Certs | identity |
| POST | `/api/update-cert` | `ApiController.UpdateCert` | Session admin surface | Update Cert | identity |
| POST | `/api/update-cert-domain-expire` | `ApiController.UpdateCertDomainExpire` | Session admin surface | Update Cert Domain Expire | identity |

## Enforcer

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-enforcer` | `ApiController.AddEnforcer` | Session admin surface | Add Enforcer | authorization |
| POST | `/api/batch-enforce` | `ApiController.BatchEnforce` | Session admin surface | Batch Enforce | authorization |
| POST | `/api/delete-enforcer` | `ApiController.DeleteEnforcer` | Session admin surface | Delete Enforcer | authorization |
| POST | `/api/enforce` | `ApiController.Enforce` | Session admin surface | Enforce | authorization |
| GET | `/api/get-all-actions` | `ApiController.GetAllActions` | Session admin surface | Get All Actions | authorization |
| GET | `/api/get-all-objects` | `ApiController.GetAllObjects` | Session admin surface | Get All Objects | authorization |
| GET | `/api/get-all-roles` | `ApiController.GetAllRoles` | Session admin surface | Get All Roles | authorization |
| GET | `/api/get-enforcer` | `ApiController.GetEnforcer` | Session admin surface | Get Enforcer | authorization |
| GET | `/api/get-enforcers` | `ApiController.GetEnforcers` | Session admin surface | Get Enforcers | authorization |
| POST | `/api/refresh-engines` | `ApiController.RefreshEngines` | Session admin surface | Refresh Engines | authorization |
| GET | `/api/run-casbin-command` | `ApiController.RunCasbinCommand` | Session admin surface | Run Casbin Command | authorization |
| POST | `/api/update-enforcer` | `ApiController.UpdateEnforcer` | Session admin surface | Update Enforcer | authorization |

## Entry

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-entry` | `ApiController.AddEntry` | Session admin surface | Add Entry | gateway |
| POST | `/api/delete-entry` | `ApiController.DeleteEntry` | Session admin surface | Delete Entry | gateway |
| GET | `/api/get-entries` | `ApiController.GetEntries` | Session admin surface | Get Entries | gateway |
| GET | `/api/get-entry` | `ApiController.GetEntry` | Session admin surface | Get Entry | gateway |
| GET | `/api/get-openclaw-session-graph` | `ApiController.GetOpenClawSessionGraph` | Session admin surface | Get Open Claw Session Graph | gateway |
| GET | `/api/get-openclaw-session-transcript` | `ApiController.GetOpenClawSessionTranscript` | Session admin surface | Get Open Claw Session Transcript | gateway |
| POST | `/api/update-entry` | `ApiController.UpdateEntry` | Session admin surface | Update Entry | gateway |

## Form

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-form` | `ApiController.AddForm` | Session admin surface | Add Form | admin-ops |
| POST | `/api/delete-form` | `ApiController.DeleteForm` | Session admin surface | Delete Form | admin-ops |
| GET | `/api/get-form` | `ApiController.GetForm` | Session admin surface | Get Form | admin-ops |
| GET | `/api/get-forms` | `ApiController.GetForms` | Session admin surface | Get Forms | admin-ops |
| GET | `/api/get-global-forms` | `ApiController.GetGlobalForms` | Session admin surface | Get Global Forms | admin-ops |
| POST | `/api/update-form` | `ApiController.UpdateForm` | Session admin surface | Update Form | admin-ops |

## Group

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-group` | `ApiController.AddGroup` | Session admin surface | Add Group | user-management |
| POST | `/api/delete-group` | `ApiController.DeleteGroup` | Session admin surface | Delete Group | user-management |
| GET | `/api/get-group` | `ApiController.GetGroup` | Session admin surface | Get Group | user-management |
| GET | `/api/get-groups` | `ApiController.GetGroups` | Session admin surface | Get Groups | user-management |
| POST | `/api/update-group` | `ApiController.UpdateGroup` | Session admin surface | Update Group | user-management |
| POST | `/api/upload-groups` | `ApiController.UploadGroups` | Session admin surface | Upload Groups | user-management |

## Invitation

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-invitation` | `ApiController.AddInvitation` | Session admin surface | Add Invitation | user-management |
| POST | `/api/delete-invitation` | `ApiController.DeleteInvitation` | Session admin surface | Delete Invitation | user-management |
| GET | `/api/get-invitation` | `ApiController.GetInvitation` | Session admin surface | Get Invitation | user-management |
| GET | `/api/get-invitation-info` | `ApiController.GetInvitationCodeInfo` | Session admin surface | Get Invitation Code Info | user-management |
| GET | `/api/get-invitations` | `ApiController.GetInvitations` | Session admin surface | Get Invitations | user-management |
| POST | `/api/send-invitation` | `ApiController.SendInvitation` | Session admin surface | Send Invitation | user-management |
| POST | `/api/update-invitation` | `ApiController.UpdateInvitation` | Session admin surface | Update Invitation | user-management |
| GET | `/api/verify-invitation` | `ApiController.VerifyInvitation` | Session admin surface | Verify Invitation | user-management |

## Key

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-key` | `ApiController.AddKey` | Session admin surface | Add Key | identity |
| POST | `/api/delete-key` | `ApiController.DeleteKey` | Session admin surface | Delete Key | identity |
| GET | `/api/get-global-keys` | `ApiController.GetGlobalKeys` | Session admin surface | Get Global Keys | identity |
| GET | `/api/get-key` | `ApiController.GetKey` | Session admin surface | Get Key | identity |
| GET | `/api/get-keys` | `ApiController.GetKeys` | Session admin surface | Get Keys | identity |
| POST | `/api/update-key` | `ApiController.UpdateKey` | Session admin surface | Update Key | identity |

## Ldap

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-ldap` | `ApiController.AddLdap` | Session admin surface | Add Ldap | admin-ops |
| POST | `/api/delete-ldap` | `ApiController.DeleteLdap` | Session admin surface | Delete Ldap | admin-ops |
| GET | `/api/get-ldap` | `ApiController.GetLdap` | Session admin surface | Get Ldap | admin-ops |
| GET | `/api/get-ldap-users` | `ApiController.GetLdapUsers` | Session admin surface | Get Ldap Users | admin-ops |
| GET | `/api/get-ldaps` | `ApiController.GetLdaps` | Session admin surface | Get Ldaps | admin-ops |
| POST | `/api/sync-ldap-users` | `ApiController.SyncLdapUsers` | Session admin surface | Sync Ldap Users | admin-ops |
| POST | `/api/update-ldap` | `ApiController.UpdateLdap` | Session admin surface | Update Ldap | admin-ops |
| * | `/scim/*` | `RootController.HandleScim` | Controller-defined admin credential | Serve SCIM and directory compatibility | admin-ops |

## Mcp

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/mcp` | `McpController.HandleMcp` | Session admin surface | Handle Mcp | gateway |

## Model

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-model` | `ApiController.AddModel` | Session admin surface | Add Model | authorization |
| POST | `/api/delete-model` | `ApiController.DeleteModel` | Session admin surface | Delete Model | authorization |
| GET | `/api/get-model` | `ApiController.GetModel` | Session admin surface | Get Model | authorization |
| GET | `/api/get-models` | `ApiController.GetModels` | Session admin surface | Get Models | authorization |
| POST | `/api/update-model` | `ApiController.UpdateModel` | Session admin surface | Update Model | authorization |

## Organization

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-organization` | `ApiController.AddOrganization` | Session admin surface | Add Organization | user-management |
| POST | `/api/delete-organization` | `ApiController.DeleteOrganization` | Session admin surface | Delete Organization | user-management |
| GET | `/api/get-default-application` | `ApiController.GetDefaultApplication` | Session admin surface | Get Default Application | user-management |
| GET | `/api/get-organization` | `ApiController.GetOrganization` | Session admin surface | Get Organization | user-management |
| GET | `/api/get-organization-applications` | `ApiController.GetOrganizationApplications` | Session admin surface | Get Organization Applications | user-management |
| GET | `/api/get-organization-names` | `ApiController.GetOrganizationNames` | Session admin surface | Get Organization Names | user-management |
| GET | `/api/get-organizations` | `ApiController.GetOrganizations` | Session admin surface | Get Organizations | user-management |
| POST | `/api/update-organization` | `ApiController.UpdateOrganization` | Session admin surface | Update Organization | user-management |

## Permission

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-permission` | `ApiController.AddPermission` | Session admin surface | Add Permission | authorization |
| POST | `/api/delete-permission` | `ApiController.DeletePermission` | Session admin surface | Delete Permission | authorization |
| GET | `/api/get-permission` | `ApiController.GetPermission` | Session admin surface | Get Permission | authorization |
| GET | `/api/get-permissions` | `ApiController.GetPermissions` | Session admin surface | Get Permissions | authorization |
| GET | `/api/get-permissions-by-role` | `ApiController.GetPermissionsByRole` | Session admin surface | Get Permissions By Role | authorization |
| GET | `/api/get-permissions-by-submitter` | `ApiController.GetPermissionsBySubmitter` | Session admin surface | Get Permissions By Submitter | authorization |
| POST | `/api/update-permission` | `ApiController.UpdatePermission` | Session admin surface | Update Permission | authorization |
| POST | `/api/upload-permissions` | `ApiController.UploadPermissions` | Session admin surface | Upload Permissions | authorization |

## Protocol

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| * | `/.well-known/:application/jwks` | `RootController.GetJwksByApplication` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| GET | `/.well-known/:application/oauth-authorization-server` | `RootController.GetOAuthServerMetadataByApplication` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| GET | `/.well-known/:application/oauth-protected-resource` | `RootController.GetOauthProtectedResourceMetadataByApplication` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| GET | `/.well-known/:application/openid-configuration` | `RootController.GetOidcDiscoveryByApplication` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| GET | `/.well-known/:application/webfinger` | `RootController.GetWebFingerByApplication` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| * | `/.well-known/jwks` | `RootController.GetJwks` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| GET | `/.well-known/oauth-authorization-server` | `RootController.GetOAuthServerMetadata` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| GET | `/.well-known/oauth-protected-resource` | `RootController.GetOauthProtectedResourceMetadata` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| GET | `/.well-known/openid-configuration` | `RootController.GetOidcDiscovery` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| GET | `/.well-known/webfinger` | `RootController.GetWebFinger` | Public protocol surface | Publish discovery and key metadata | out-of-scope |
| POST | `/api/oauth/register` | `ApiController.DynamicClientRegister` | Public or login-session flow | Dynamic Client Register | out-of-scope |
| GET | `/cas/:organization/:application/p3/proxyValidate` | `RootController.CasP3ProxyValidate` | Public protocol surface | Serve CAS protocol validation | out-of-scope |
| GET | `/cas/:organization/:application/p3/serviceValidate` | `RootController.CasP3ServiceValidate` | Public protocol surface | Serve CAS protocol validation | out-of-scope |
| GET | `/cas/:organization/:application/proxy` | `RootController.CasProxy` | Public protocol surface | Serve CAS protocol validation | out-of-scope |
| GET | `/cas/:organization/:application/proxyValidate` | `RootController.CasProxyValidate` | Public protocol surface | Serve CAS protocol validation | out-of-scope |
| POST | `/cas/:organization/:application/samlValidate` | `RootController.SamlValidate` | Public protocol surface | Serve CAS protocol validation | out-of-scope |
| GET | `/cas/:organization/:application/serviceValidate` | `RootController.CasServiceValidate` | Public protocol surface | Serve CAS protocol validation | out-of-scope |
| GET | `/cas/:organization/:application/validate` | `RootController.CasValidate` | Public protocol surface | Serve CAS protocol validation | out-of-scope |

## Provider

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-provider` | `ApiController.AddProvider` | Session admin surface | Add Provider | identity |
| POST | `/api/delete-provider` | `ApiController.DeleteProvider` | Session admin surface | Delete Provider | identity |
| GET | `/api/get-global-providers` | `ApiController.GetGlobalProviders` | Session admin surface | Get Global Providers | identity |
| GET | `/api/get-provider` | `ApiController.GetProvider` | Session admin surface | Get Provider | identity |
| GET | `/api/get-providers` | `ApiController.GetProviders` | Session admin surface | Get Providers | identity |
| POST | `/api/update-provider` | `ApiController.UpdateProvider` | Session admin surface | Update Provider | identity |

## Record

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-record` | `ApiController.AddRecord` | Session admin surface | Add Record | auditing |
| GET | `/api/get-records` | `ApiController.GetRecords` | Session admin surface | Get Records | auditing |
| POST | `/api/get-records-filter` | `ApiController.GetRecordsByFilter` | Session admin surface | Get Records By Filter | auditing |

## Resource

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-resource` | `ApiController.AddResource` | Session admin surface | Add Resource | identity |
| POST | `/api/delete-resource` | `ApiController.DeleteResource` | Session admin surface | Delete Resource | identity |
| GET | `/api/get-resource` | `ApiController.GetResource` | Session admin surface | Get Resource | identity |
| GET | `/api/get-resources` | `ApiController.GetResources` | Session admin surface | Get Resources | identity |
| POST | `/api/update-resource` | `ApiController.UpdateResource` | Session admin surface | Update Resource | identity |
| POST | `/api/upload-resource` | `ApiController.UploadResource` | Session admin surface | Upload Resource | identity |

## Role

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-role` | `ApiController.AddRole` | Session admin surface | Add Role | authorization |
| POST | `/api/delete-role` | `ApiController.DeleteRole` | Session admin surface | Delete Role | authorization |
| GET | `/api/get-role` | `ApiController.GetRole` | Session admin surface | Get Role | authorization |
| GET | `/api/get-roles` | `ApiController.GetRoles` | Session admin surface | Get Roles | authorization |
| POST | `/api/update-role` | `ApiController.UpdateRole` | Session admin surface | Update Role | authorization |
| POST | `/api/upload-roles` | `ApiController.UploadRoles` | Session admin surface | Upload Roles | authorization |

## Rule

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-rule` | `ApiController.AddRule` | Session admin surface | Add Rule | gateway |
| POST | `/api/delete-rule` | `ApiController.DeleteRule` | Session admin surface | Delete Rule | gateway |
| GET | `/api/get-rule` | `ApiController.GetRule` | Session admin surface | Get Rule | gateway |
| GET | `/api/get-rules` | `ApiController.GetRules` | Session admin surface | Get Rules | gateway |
| POST | `/api/update-rule` | `ApiController.UpdateRule` | Session admin surface | Update Rule | gateway |

## Server

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-server` | `ApiController.AddServer` | Session admin surface | Add Server | gateway |
| POST | `/api/delete-server` | `ApiController.DeleteServer` | Session admin surface | Delete Server | gateway |
| GET | `/api/get-mcp-access-token` | `ApiController.GetMcpAccessToken` | Session admin surface | Get Mcp Access Token | gateway |
| GET | `/api/get-online-servers` | `ApiController.GetOnlineServers` | Session admin surface | Get Online Servers | gateway |
| GET | `/api/get-server` | `ApiController.GetServer` | Session admin surface | Get Server | gateway |
| GET | `/api/get-servers` | `ApiController.GetServers` | Session admin surface | Get Servers | gateway |
| GET | `/api/scan` | `ApiController.Scan` | Session admin surface | Scan | gateway |
| GET | `/api/server/:owner/:name` | `ApiController.ProxyServer` | Session admin surface | Proxy Server | gateway |
| POST | `/api/server/:owner/:name` | `ApiController.ProxyServer` | Session admin surface | Proxy Server | gateway |
| POST | `/api/sync-intranet-servers` | `ApiController.SyncIntranetServers` | Session admin surface | Sync Intranet Servers | gateway |
| POST | `/api/sync-mcp-tool` | `ApiController.SyncMcpTool` | Session admin surface | Sync Mcp Tool | gateway |
| POST | `/api/update-server` | `ApiController.UpdateServer` | Session admin surface | Update Server | gateway |

## Service

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/send-email` | `ApiController.SendEmail` | Session admin surface | Send Email | out-of-scope |
| POST | `/api/send-notification` | `ApiController.SendNotification` | Session admin surface | Send Notification | out-of-scope |
| POST | `/api/send-sms` | `ApiController.SendSms` | Session admin surface | Send Sms | out-of-scope |

## Session

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-session` | `ApiController.AddSession` | Session admin surface | Add Session | auditing |
| POST | `/api/delete-session` | `ApiController.DeleteSession` | Session admin surface | Delete Session | auditing |
| GET | `/api/get-session` | `ApiController.GetSingleSession` | Session admin surface | Get Single Session | auditing |
| GET | `/api/get-sessions` | `ApiController.GetSessions` | Session admin surface | Get Sessions | auditing |
| GET | `/api/is-session-duplicated` | `ApiController.IsSessionDuplicated` | Session admin surface | Is Session Duplicated | auditing |
| POST | `/api/update-session` | `ApiController.UpdateSession` | Session admin surface | Update Session | auditing |

## Site

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-site` | `ApiController.AddSite` | Session admin surface | Add Site | gateway |
| POST | `/api/delete-site` | `ApiController.DeleteSite` | Session admin surface | Delete Site | gateway |
| GET | `/api/get-global-sites` | `ApiController.GetGlobalSites` | Session admin surface | Get Global Sites | gateway |
| GET | `/api/get-site` | `ApiController.GetSite` | Session admin surface | Get Site | gateway |
| GET | `/api/get-sites` | `ApiController.GetSites` | Session admin surface | Get Sites | gateway |
| POST | `/api/update-site` | `ApiController.UpdateSite` | Session admin surface | Update Site | gateway |

## Syncer

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-syncer` | `ApiController.AddSyncer` | Session admin surface | Add Syncer | admin-ops |
| POST | `/api/delete-syncer` | `ApiController.DeleteSyncer` | Session admin surface | Delete Syncer | admin-ops |
| GET | `/api/get-syncer` | `ApiController.GetSyncer` | Session admin surface | Get Syncer | admin-ops |
| GET | `/api/get-syncers` | `ApiController.GetSyncers` | Session admin surface | Get Syncers | admin-ops |
| GET | `/api/run-syncer` | `ApiController.RunSyncer` | Session admin surface | Run Syncer | admin-ops |
| POST | `/api/test-syncer-db` | `ApiController.TestSyncerDb` | Session admin surface | Test Syncer Db | admin-ops |
| POST | `/api/update-syncer` | `ApiController.UpdateSyncer` | Session admin surface | Update Syncer | admin-ops |

## System

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/check-user-password` | `ApiController.CheckUserPassword` | Session admin surface | Check User Password | admin-ops |
| POST | `/api/delete-mfa` | `ApiController.DeleteMfa` | Session admin surface | Delete Mfa | admin-ops |
| GET | `/api/get-app-login` | `ApiController.GetApplicationLogin` | Session admin surface | Get Application Login | admin-ops |
| GET | `/api/get-dashboard` | `ApiController.GetDashboard` | Session admin surface | Get Dashboard | admin-ops |
| GET | `/api/get-dashboard-heatmap` | `ApiController.GetDashboardLoginHeatmap` | Session admin surface | Get Dashboard Login Heatmap | admin-ops |
| GET | `/api/get-dashboard-mfa` | `ApiController.GetDashboardMfaCoverage` | Session admin surface | Get Dashboard Mfa Coverage | admin-ops |
| GET | `/api/get-dashboard-providers` | `ApiController.GetDashboardProviderDistribution` | Session admin surface | Get Dashboard Provider Distribution | admin-ops |
| GET | `/api/get-prometheus-info` | `ApiController.GetPrometheusInfo` | Session admin surface | Get Prometheus Info | admin-ops |
| GET | `/api/get-qrcode` | `ApiController.GetQRCode` | Session admin surface | Get QRCode | admin-ops |
| GET | `/api/get-system-info` | `ApiController.GetSystemInfo` | Session admin surface | Get System Info | admin-ops |
| GET | `/api/get-version-info` | `ApiController.GetVersionInfo` | Session admin surface | Get Version Info | admin-ops |
| GET | `/api/health` | `ApiController.Health` | Session admin surface | Health | admin-ops |
| GET | `/api/metrics` | `ApiController.GetMetrics` | Session admin surface | Get Metrics | admin-ops |
| POST | `/api/set-password` | `ApiController.SetPassword` | Session admin surface | Set Password | admin-ops |
| POST | `/api/set-preferred-mfa` | `ApiController.SetPreferredMfa` | Session admin surface | Set Preferred Mfa | admin-ops |
| POST | `/api/v1/logs` | `ApiController.AddOtlpLogs` | Session admin surface | Add Otlp Logs | admin-ops |
| POST | `/api/v1/metrics` | `ApiController.AddOtlpMetrics` | Session admin surface | Add Otlp Metrics | admin-ops |
| POST | `/api/v1/traces` | `ApiController.AddOtlpTrace` | Session admin surface | Add Otlp Trace | admin-ops |

## Ticket

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-ticket` | `ApiController.AddTicket` | Session admin surface | Add Ticket | admin-ops |
| POST | `/api/add-ticket-message` | `ApiController.AddTicketMessage` | Session admin surface | Add Ticket Message | admin-ops |
| POST | `/api/delete-ticket` | `ApiController.DeleteTicket` | Session admin surface | Delete Ticket | admin-ops |
| GET | `/api/get-ticket` | `ApiController.GetTicket` | Session admin surface | Get Ticket | admin-ops |
| GET | `/api/get-tickets` | `ApiController.GetTickets` | Session admin surface | Get Tickets | admin-ops |
| POST | `/api/update-ticket` | `ApiController.UpdateTicket` | Session admin surface | Update Ticket | admin-ops |

## Token

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-token` | `ApiController.AddToken` | Session admin surface | Add Token | auditing |
| POST | `/api/delete-token` | `ApiController.DeleteToken` | Session admin surface | Delete Token | auditing |
| GET | `/api/get-token` | `ApiController.GetToken` | Session admin surface | Get Token | auditing |
| GET | `/api/get-tokens` | `ApiController.GetTokens` | Session admin surface | Get Tokens | auditing |
| POST | `/api/update-token` | `ApiController.UpdateToken` | Session admin surface | Update Token | auditing |

## User

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-user` | `ApiController.AddUser` | Session admin surface | Add User | user-management |
| POST | `/api/delete-user` | `ApiController.DeleteUser` | Session admin surface | Delete User | user-management |
| POST | `/api/exit-impersonate-user` | `ApiController.ExitImpersonateUser` | Session admin surface | Exit Impersonate User | user-management |
| GET | `/api/get-global-users` | `ApiController.GetGlobalUsers` | Session admin surface | Get Global Users | user-management |
| GET | `/api/get-sorted-users` | `ApiController.GetSortedUsers` | Session admin surface | Get Sorted Users | user-management |
| GET | `/api/get-user` | `ApiController.GetUser` | Session admin surface | Get User | user-management |
| GET | `/api/get-user-application` | `ApiController.GetUserApplication` | Session admin surface | Get User Application | user-management |
| GET | `/api/get-user-count` | `ApiController.GetUserCount` | Session admin surface | Get User Count | user-management |
| GET | `/api/get-users` | `ApiController.GetUsers` | Session admin surface | Get Users | user-management |
| POST | `/api/impersonate-user` | `ApiController.ImpersonateUser` | Session admin surface | Impersonate User | user-management |
| POST | `/api/remove-user-from-group` | `ApiController.RemoveUserFromGroup` | Session admin surface | Remove User From Group | user-management |
| POST | `/api/update-user` | `ApiController.UpdateUser` | Session admin surface | Update User | user-management |
| POST | `/api/upload-users` | `ApiController.UploadUsers` | Session admin surface | Upload Users | user-management |
| POST | `/api/verify-identification` | `ApiController.VerifyIdentification` | Session admin surface | Verify Identification | user-management |

## Verification

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| GET | `/api/get-captcha` | `ApiController.GetCaptcha` | Session admin surface | Get Captcha | auditing |
| GET | `/api/get-captcha-status` | `ApiController.GetCaptchaStatus` | Session admin surface | Get Captcha Status | auditing |
| GET | `/api/get-email-and-phone` | `ApiController.GetEmailAndPhone` | Session admin surface | Get Email And Phone | auditing |
| GET | `/api/get-verifications` | `ApiController.GetVerifications` | Session admin surface | Get Verifications | auditing |
| POST | `/api/reset-email-or-phone` | `ApiController.ResetEmailOrPhone` | Session admin surface | Reset Email Or Phone | auditing |
| POST | `/api/send-verification-code` | `ApiController.SendVerificationCode` | Session admin surface | Send Verification Code | auditing |
| POST | `/api/verify-captcha` | `ApiController.VerifyCaptcha` | Session admin surface | Verify Captcha | auditing |
| POST | `/api/verify-code` | `ApiController.VerifyCode` | Session admin surface | Verify Code | auditing |

## Webhook

| Method | Path | Controller/Action | Auth mode | Primary use | Frontend module |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/add-webhook` | `ApiController.AddWebhook` | Session admin surface | Add Webhook | admin-ops |
| POST | `/api/delete-webhook` | `ApiController.DeleteWebhook` | Session admin surface | Delete Webhook | admin-ops |
| POST | `/api/delete-webhook-event` | `ApiController.DeleteWebhookEvent` | Session admin surface | Delete Webhook Event | admin-ops |
| GET | `/api/get-webhook` | `ApiController.GetWebhook` | Session admin surface | Get Webhook | admin-ops |
| GET | `/api/get-webhook-event` | `ApiController.GetWebhookEventType` | Session admin surface | Get Webhook Event Type | admin-ops |
| GET | `/api/get-webhook-event-detail` | `ApiController.GetWebhookEvent` | Session admin surface | Get Webhook Event | admin-ops |
| GET | `/api/get-webhook-events` | `ApiController.GetWebhookEvents` | Session admin surface | Get Webhook Events | admin-ops |
| GET | `/api/get-webhooks` | `ApiController.GetWebhooks` | Session admin surface | Get Webhooks | admin-ops |
| POST | `/api/replay-webhook-event` | `ApiController.ReplayWebhookEvent` | Session admin surface | Replay Webhook Event | admin-ops |
| POST | `/api/update-webhook` | `ApiController.UpdateWebhook` | Session admin surface | Update Webhook | admin-ops |
| * | `/api/webhook` | `ApiController.HandleOfficialAccountEvent` | Session admin surface | Handle Official Account Event | admin-ops |
