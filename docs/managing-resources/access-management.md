---
title: Access management
description: A guide to managing access to resources in Beneath
menu:
  docs:
    parent: managing-resources
    weight: 400
weight: 400
---

The gist of access management in Beneath is as follows: [Users]({{< ref "/docs/managing-resources/resources.md#users" >}}) and [services]({{< ref "/docs/managing-resources/resources.md#services" >}}) can be granted permissions to other resources, such as streams, projects and organizations. Different resources have different permissions, like `read` and `write` for streams.

You issue and use [secrets]({{< ref "/docs/managing-resources/resources.md#secrets" >}}) to authenticate as a *user* or a *service* from your code, notebook, command-line, or similar. *User* secrets are useful during development to connect to Beneath from your code. *Service* secrets should be used in production systems or publicly-exposed code to more strictly limit access permissions and monitor usage.

## Examples

### Authenticating in Beneath CLI

See example in [the CLI installation guide]({{< ref "/docs/managing-resources/cli.md#Installation" >}}).

### Inviting a user to an organization

> **Note:** When you *invite* a user to your organization, you take over the full billing responsibility for the user's activity. If you just want to grant the user access to view stuff in your organization (but not pay their bills), see [Granting a user access to an organization]({{< relref "#granting-a-user-access-to-an-organization" >}}).

First, an organization admin should send an invitation to the user (the flags designate the invitees new permissions):

```bash
beneath organization invite-member ORGANIZATION_NAME USERNAME --view --admin
```

Second, the invited user should accept the invitation:

```bash
beneath organization accept-invite ORGANIZATION_NAME
```

### Granting a user access to an organization

> **Note:** When you *grant* a user access to your organization, they remain responsible for paying the bills for their own usage on Beneath. If you also want to pay for their usage on Beneath, see [Inviting a user to an organization]({{< relref "#inviting-a-user-to-an-organization" >}}).

Run the following command (change the `--create` or `--admin` flags to configure permissions):

```bash
beneath organization update-permissions ORGANIZATION_NAME USERNAME --view true --create true --admin true
```

The user doesn't have to be a part of the organization in advance.

### Granting a user access to a project and its streams

In Beneath, *user* access to streams is managed at the project-level. You cannot grant a *user* access to only one stream (however, if you need a secret with permissions for just a single stream, use a *service*). 

Run the following command (change the `--create` or `--admin` flags to configure permissions):

```bash
beneath project update-permissions ORGANIZATION_NAME/PROJECT_NAME USERNAME --view true --create true --admin true
```

The user doesn't have to be a part of the same organization as you or the project to get access.

### Creating a secret for a user

You can issue and copy a new secret from the ["Secrets" tab of your profile page](https://beneath.dev/-/redirects/secrets) in the Beneath Terminal.

There are two types of secrets:
- **Command-line secrets:** These contain all your access permissions
- **Read-only secrets:** These are limited to `view` permissions on the resources you has access to

**Never share your user secrets!** You should not share or expose user secrets nor use them in production systems. Use a service secret if you need a secret in a production system or if you need to expose a secret (e.g. in a shared notebook or in your frontend code).

### Creating a service and granting it access to a stream

Services are useful when deploying or publishing code that reads or writes to Beneath. You can control their access permissions, and monitor and limit their usage. Read ["Services"]({{< ref "/docs/managing-resources/resources.md#services" >}}) for more details.

First, create a new service in your organization (customize their quotas to set monthly limits):

```bash
beneath service create ORGANIZATION_NAME/NEW_SERVICE_NAME --read-quota-mb 100 --write-quota-mb 100
```

Second, give the service the permissions it needs (remove `--write` if it doesn't need write permissions):

```bash
beneath service update-permission ORGANIZATION_NAME/SERVICE_NAME ORGANIZATION_NAME/PROJECT_NAME/STREAM_NAME --read --write 
```

> **Note:** Services do not automatically get read permissions for public streams (unlike users). You must set them manually.

You're now ready to issue a secret for the service.

### Creating a secret for a service

If you haven't yet created a service, follow the example given above.

Run the following command:

```bash
beneath service issue-secret ORGANIZATION_NAME/SERVICE_NAME --description "YOUR SERVICE DESCRIPTION"
```

You can now use the secret to connect to Beneath from your code. Most client libraries will automatically use your secret if you set it in the `BENEATH_SECRET` environment variable (see the documentation for your client library for other ways of passing the secret).

**Think carefully before sharing sharing service secrets!** If you need to expose a secret publicly (e.g. in your front-end code or in a notebook), make sure it belongs to a service with sensible usage quotas and only `read` permissions. In all other cases, keep your secret very safe and do not check it into Git.